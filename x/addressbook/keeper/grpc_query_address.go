package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/CudoVentures/cudos-node/x/addressbook/types"
)

func (k Keeper) AddressAll(c context.Context, req *types.QueryAllAddressRequest) (*types.QueryAllAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var addresss []types.Address
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	addressStore := prefix.NewStore(store, types.KeyPrefix(types.AddressKeyPrefix))

	pageRes, err := query.Paginate(addressStore, req.Pagination, func(key []byte, value []byte) error {
		var address types.Address
		if err := k.cdc.Unmarshal(value, &address); err != nil {
			return err
		}

		addresss = append(addresss, address)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAddressResponse{Address: addresss, Pagination: pageRes}, nil
}

func (k Keeper) Address(c context.Context, req *types.QueryGetAddressRequest) (*types.QueryGetAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetAddress(
		ctx,
		req.Creator,
		req.Network,
		req.Label,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetAddressResponse{Address: val}, nil
}
