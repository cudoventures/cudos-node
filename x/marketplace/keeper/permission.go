package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) isAdmin(ctx sdk.Context, sender string) error {
	for _, admin := range k.GetParams(ctx).Admins {
		if strings.ToLower(admin) == strings.ToLower(sender) {
			return nil
		}
	}

	return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "Insufficient permissions. '%s' is not admin.", sender)
}

func (k msgServer) transferAdminPermission(ctx sdk.Context, currentAdmin, newAdmin string) error {
	params := k.GetParams(ctx)
	for idx, admin := range params.Admins {
		if strings.ToLower(admin) == strings.ToLower(currentAdmin) {
			params.Admins[idx] = newAdmin
			k.SetParams(ctx, params)
			return nil
		}
	}

	return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "Insufficient permissions. '%s' is not admin.", currentAdmin)
}
