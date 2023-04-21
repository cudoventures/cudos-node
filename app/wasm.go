package app

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmutils "github.com/CosmWasm/wasmd/x/wasm/client/utils"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvm "github.com/CosmWasm/wasmvm"
	nftCustomBindings "github.com/CudoVentures/cudos-node/x/nft/custom-bindings"
	nftkeeper "github.com/CudoVentures/cudos-node/x/nft/keeper"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	snapshot "github.com/cosmos/cosmos-sdk/snapshots/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	protoio "github.com/gogo/protobuf/io"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func GetCustomMsgEncodersOptions() []wasmkeeper.Option {
	nftEncodingOptions := wasmkeeper.WithMessageEncoders(nftEncoders())
	return []wasm.Option{nftEncodingOptions}
}

func GetCustomMsgQueryOptions(keeper nftkeeper.Keeper) []wasmkeeper.Option {
	nftQueryOptions := wasmkeeper.WithQueryPlugins(nftQueryPlugins(keeper))
	return []wasm.Option{nftQueryOptions}
}

func nftEncoders() *wasmkeeper.MessageEncoders {
	return &wasmkeeper.MessageEncoders{
		Custom: nftCustomBindings.EncodeNftMessage(),
	}
}

// nftQueryPlugins needs to be registered in test setup to handle custom query callbacks
func nftQueryPlugins(keeper nftkeeper.Keeper) *wasmkeeper.QueryPlugins {
	return &wasmkeeper.QueryPlugins{
		Custom: nftCustomBindings.PerformCustomNftQuery(keeper),
	}
}

var _ snapshot.ExtensionSnapshotter = &wasmSnapshotter{}

// snapshotFormat format 1 is just gzipped wasm byte code for each item payload. No protobuf envelope, no metadata.
const snapshotFormat = 1

type wasmSnapshotter struct {
	k      *wasmkeeper.Keeper
	cms    sdk.MultiStore
	wasmVM wasmtypes.WasmerEngine
}

// WasmSnapshotter was introduced in wasmd v0.27.0 to include contract wasm codes in the snapshot.
// wasmd v0.25.0 -> v0.27.0 requires ibc-go update from v2.2.0 -> v3.0.0
// it is a breaking change and requires software upgrade
// temporarily introduced here manually to support state-sync until next cudos-noded upgrade
func NewWasmSnapshotter(cms sdk.MultiStore, keeper *wasmkeeper.Keeper, wasmVM wasmtypes.WasmerEngine) *wasmSnapshotter {
	return &wasmSnapshotter{
		k:      keeper,
		cms:    cms,
		wasmVM: wasmVM,
	}
}

func (ws *wasmSnapshotter) SnapshotName() string {
	return wasmtypes.ModuleName
}

func (ws *wasmSnapshotter) SnapshotFormat() uint32 {
	return snapshotFormat
}

func (ws *wasmSnapshotter) SupportedFormats() []uint32 {
	// If we support older formats, add them here and handle them in Restore
	return []uint32{snapshotFormat}
}

func (ws *wasmSnapshotter) Snapshot(height uint64, protoWriter protoio.Writer) error {
	cacheMS, err := ws.cms.CacheMultiStoreWithVersion(int64(height))
	if err != nil {
		return err
	}

	ctx := sdk.NewContext(cacheMS, tmproto.Header{}, false, log.NewNopLogger())
	seenBefore := make(map[string]bool)
	var rerr error

	ws.k.IterateCodeInfos(ctx, func(id uint64, info wasmtypes.CodeInfo) bool {
		// Many code ids may point to the same code hash... only sync it once
		hexHash := hex.EncodeToString(info.CodeHash)
		// if seenBefore, just skip this one and move to the next
		if seenBefore[hexHash] {
			return false
		}
		seenBefore[hexHash] = true

		// load code and abort on error
		wasmBytes, err := ws.k.GetByteCode(ctx, id)
		if err != nil {
			rerr = err
			return true
		}

		compressedWasm, err := wasmutils.GzipIt(wasmBytes)
		if err != nil {
			rerr = err
			return true
		}

		if err = snapshot.WriteExtensionItem(protoWriter, compressedWasm); err != nil {
			rerr = err
			return true
		}

		return false
	})

	return rerr
}

func (ws *wasmSnapshotter) Restore(
	height uint64, format uint32, protoReader protoio.Reader,
) (snapshot.SnapshotItem, error) {
	if format != snapshotFormat {
		return snapshot.SnapshotItem{}, snapshot.ErrUnknownFormat
	}

	// keep the last item here... if we break, it will either be empty (if we hit io.EOF)
	// or contain the last item (if we hit payload == nil)
	var item snapshot.SnapshotItem
	for {
		item = snapshot.SnapshotItem{}
		if err := protoReader.ReadMsg(&item); err == io.EOF {
			break
		} else if err != nil {
			return snapshot.SnapshotItem{}, sdkerrors.Wrap(err, "invalid protobuf message")
		}

		// if it is not another ExtensionPayload message, then it is not for us.
		// we should return it an let the manager handle this one
		payload := item.GetExtensionPayload()
		if payload == nil {
			break
		}

		wasmCode, err := uncompressGzip(payload.Payload, uint64(wasmtypes.MaxWasmSize))
		if err != nil {
			return snapshot.SnapshotItem{}, sdkerrors.Wrap(wasmtypes.ErrCreateFailed, err.Error())
		}

		if _, err = ws.wasmVM.Create(wasmCode); err != nil {
			return snapshot.SnapshotItem{}, sdkerrors.Wrap(wasmtypes.ErrCreateFailed, err.Error())
		}
	}

	ctx := sdk.NewContext(ws.cms, tmproto.Header{Height: int64(height)}, false, log.NewNopLogger())
	return item, ws.k.InitializePinnedCodes(ctx)
}

// magic bytes to identify gzip.
// See https://www.ietf.org/rfc/rfc1952.txt
// and https://github.com/golang/go/blob/master/src/net/http/sniff.go#L186
var gzipIdent = []byte("\x1F\x8B\x08")

// uncompress returns gzip uncompressed content or given src when not gzip.
func uncompressGzip(src []byte, limit uint64) ([]byte, error) {
	switch n := uint64(len(src)); {
	case n < 3:
		return src, nil
	case n > limit:
		return nil, wasmtypes.ErrLimit
	}
	if !bytes.Equal(gzipIdent, src[0:3]) {
		return src, nil
	}
	zr, err := gzip.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}
	zr.Multistream(false)
	defer zr.Close()
	return ioutil.ReadAll(newLimitedReader(zr, int64(limit)))
}

// newLimitedReader returns a Reader that reads from r
// but stops with types.ErrLimit after n bytes.
// The underlying implementation is a *io.LimitedReader.
func newLimitedReader(r io.Reader, n int64) io.Reader {
	return &limitedReader{r: &io.LimitedReader{R: r, N: n}}
}

type limitedReader struct {
	r *io.LimitedReader
}

func (l *limitedReader) Read(p []byte) (n int, err error) {
	if l.r.N <= 0 {
		return 0, wasmtypes.ErrLimit
	}
	return l.r.Read(p)
}

const (
	supportedFeatures   = "iterator,staking,stargate"
	contractMemoryLimit = uint32(32)
)

func NewWasmVM(homePath string, appOpts servertypes.AppOptions) (*wasmvm.VM, error) {
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		return nil, fmt.Errorf("error while reading wasm config: %w", err)
	}

	return wasmvm.NewVM(
		filepath.Join(homePath, "wasm", "wasm"),
		supportedFeatures,
		contractMemoryLimit,
		wasmConfig.ContractDebugMode,
		wasmConfig.MemoryCacheSize,
	)
}
