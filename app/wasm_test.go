package app

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/snapshots"
	sdk "github.com/cosmos/cosmos-sdk/types"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func TestWasmSnapshotter(t *testing.T) {
	reflectContract, err := ioutil.ReadFile("../cosmwasm-testing/testdata/reflect.wasm")
	require.NoError(t, err)

	burnerContract, err := ioutil.ReadFile("../cosmwasm-testing/testdata/burner.wasm")
	require.NoError(t, err)

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	for _, tc := range []struct {
		desc          string
		contractFiles [][]byte
	}{
		{
			desc:          "single contract",
			contractFiles: [][]byte{reflectContract},
		},
		{
			desc:          "multiple contract",
			contractFiles: [][]byte{reflectContract, burnerContract, reflectContract},
		},
		{
			desc:          "duplicate contract",
			contractFiles: [][]byte{reflectContract, reflectContract},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			// setup source app
			srcApp := newTestApp(t)
			srcApp.Commit()
			srcApp.BeginBlock(abcitypes.RequestBeginBlock{Header: tmproto.Header{
				Height: srcApp.LastBlockHeight() + 1,
			}})
			msgStoreCodeHandler := srcApp.MsgServiceRouter().Handler(&wasm.MsgStoreCode{})

			// store contracts to srcApp and get state hash
			srcAppStateHash := make(map[uint64][]byte, len(tc.contractFiles))
			for i, c := range tc.contractFiles {
				ctx := srcApp.BaseApp.NewContext(false, tmproto.Header{
					Height: srcApp.LastBlockHeight() + 1,
				})

				res, err := msgStoreCodeHandler(ctx, &wasm.MsgStoreCode{
					Sender:       addr.String(),
					WASMByteCode: c,
				})
				require.NoError(t, err)

				var msgStoreCodeRes wasm.MsgStoreCodeResponse
				err = msgStoreCodeRes.Unmarshal(res.Data)
				require.NoError(t, err)
				require.Equal(t, uint64(i+1), msgStoreCodeRes.CodeID)

				hash := sha256.Sum256(c)
				srcAppStateHash[msgStoreCodeRes.CodeID] = hash[:]
			}

			// commit state and create srcApp snapshot
			srcApp.Commit()
			snapshot, err := srcApp.SnapshotManager().Create(uint64(srcApp.LastBlockHeight()))
			require.NoError(t, err)
			require.NotNil(t, snapshot)

			// setup destination app and restore srcApp snapshot
			dstApp := newTestApp(t)

			err = dstApp.SnapshotManager().Restore(*snapshot)
			require.NoError(t, err)

			for i := uint32(0); i < snapshot.Chunks; i++ {
				chunkBz, err := srcApp.SnapshotManager().LoadChunk(snapshot.Height, snapshot.Format, i)
				require.NoError(t, err)

				end, err := dstApp.SnapshotManager().RestoreChunk(chunkBz)
				require.NoError(t, err)
				if end {
					break
				}
			}

			// get dstApp state hash after restore
			ctx := dstApp.BaseApp.NewContext(false, tmproto.Header{
				Height: dstApp.LastBlockHeight() + 1,
			})

			queryClient := wasmtypes.NewQueryClient(&baseapp.QueryServiceTestHelper{
				GRPCQueryRouter: dstApp.GRPCQueryRouter(),
				Ctx:             ctx,
			})

			contracts, err := queryClient.Codes(ctx, &wasmtypes.QueryCodesRequest{})
			require.NoError(t, err)

			dstAppStateHash := make(map[uint64][]byte, len(tc.contractFiles))
			for _, c := range contracts.CodeInfos {
				contractFile, err := queryClient.Code(ctx, &wasmtypes.QueryCodeRequest{CodeId: c.CodeID})
				require.NoError(t, err)

				hash := sha256.Sum256(contractFile.Data)
				require.Equal(t, hash[:], c.DataHash.Bytes())

				dstAppStateHash[c.CodeID] = hash[:]
			}

			// assert state hashes
			require.Equal(t, srcAppStateHash, dstAppStateHash)
		})
	}
}

func newTestApp(t *testing.T) *App {
	homeDir := t.TempDir()
	snapshotDir := filepath.Join(homeDir, "snapshots")

	snapshotDb, err := sdk.NewLevelDB("metadata", snapshotDir)
	require.NoError(t, err)

	snapshotStore, err := snapshots.NewStore(snapshotDb, snapshotDir)
	require.NoError(t, err)

	testApp := New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		homeDir,
		0,
		MakeEncodingConfig(),
		server.NewDefaultContext().Viper,
		baseapp.SetSnapshotStore(snapshotStore),
	)

	genesisState := NewDefaultGenesisState(testApp.AppCodec())
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	require.NoError(t, err)

	testApp.InitChain(abcitypes.RequestInitChain{
		Validators:    []abcitypes.ValidatorUpdate{},
		AppStateBytes: stateBytes,
	})

	return testApp
}

func TestUncompressGzip(t *testing.T) {
	wasmRaw, err := ioutil.ReadFile("../cosmwasm-testing/testdata/hackatom.wasm")
	require.NoError(t, err)

	wasmGzipped, err := ioutil.ReadFile("../cosmwasm-testing/testdata/hackatom.wasm.gzip")
	require.NoError(t, err)

	const maxSize = 400_000

	for _, tc := range []struct {
		desc      string
		src       []byte
		expError  error
		expResult []byte
	}{
		{
			desc:      "handle wasm uncompressed",
			src:       wasmRaw,
			expResult: wasmRaw,
		},
		{
			desc:      "handle wasm compressed",
			src:       wasmGzipped,
			expResult: wasmRaw,
		},
		{
			desc:      "handle nil slice",
			src:       nil,
			expResult: nil,
		},
		{
			desc:      "handle short unidentified",
			src:       []byte{0x1, 0x2},
			expResult: []byte{0x1, 0x2},
		},
		{
			desc:     "handle input slice exceeding limit",
			src:      []byte(strings.Repeat("a", maxSize+1)),
			expError: wasmtypes.ErrLimit,
		},
		{
			desc:      "handle input slice at limit",
			src:       []byte(strings.Repeat("a", maxSize)),
			expResult: []byte(strings.Repeat("a", maxSize)),
		},
		{
			desc:     "handle gzip identifier only",
			src:      gzipIdent,
			expError: io.ErrUnexpectedEOF,
		},
		{
			desc:     "handle broken gzip",
			src:      append(gzipIdent, byte(0x1)),
			expError: io.ErrUnexpectedEOF,
		},
		{
			desc:     "handle incomplete gzip",
			src:      wasmGzipped[:len(wasmGzipped)-5],
			expError: io.ErrUnexpectedEOF,
		},
		{
			desc:      "handle limit gzip output",
			src:       asGzip(bytes.Repeat([]byte{0x1}, maxSize)),
			expResult: bytes.Repeat([]byte{0x1}, maxSize),
		},
		{
			desc:     "handle big gzip output",
			src:      asGzip(bytes.Repeat([]byte{0x1}, maxSize+1)),
			expError: wasmtypes.ErrLimit,
		},
		{
			desc:     "handle other big gzip output",
			src:      asGzip(bytes.Repeat([]byte{0x1}, 2*maxSize)),
			expError: wasmtypes.ErrLimit,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			r, err := uncompressGzip(tc.src, maxSize)
			require.True(t, errors.Is(tc.expError, err), "exp %v got %+v", tc.expError, err)
			if tc.expError != nil {
				return
			}
			require.Equal(t, tc.expResult, r)
		})
	}
}

func asGzip(src []byte) []byte {
	var buf bytes.Buffer
	zipper := gzip.NewWriter(&buf)
	if _, err := io.Copy(zipper, bytes.NewReader(src)); err != nil {
		panic(err)
	}

	if err := zipper.Close(); err != nil {
		panic(err)
	}

	return buf.Bytes()
}
