package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/initia-labs/initia/x/move/types"
	vmtypes "github.com/initia-labs/initiavm/types"
)

// BaseDenom - base denom of native move dex
func (k Keeper) BaseDenom(ctx context.Context) string {
	return k.GetParams(ctx).BaseDenom
}

// BaseMinGasPrice - min gas price in base denom unit
func (k Keeper) BaseMinGasPrice(ctx context.Context) sdk.Dec {
	return k.GetParams(ctx).BaseMinGasPrice
}

// ArbitraryEnabled - arbitrary enabled flag
func (k Keeper) ArbitraryEnabled(ctx context.Context) (bool, error) {
	return NewCodeKeeper(&k).GetAllowArbitrary(ctx)
}

// AllowedPublishers - allowed publishers
func (k Keeper) AllowedPublishers(ctx context.Context) ([]vmtypes.AccountAddress, error) {
	return NewCodeKeeper(&k).GetAllowedPublishers(ctx)
}

// SetArbitraryEnabled - update arbitrary enabled flag
func (k Keeper) SetArbitraryEnabled(ctx context.Context, arbitraryEnabled bool) error {
	return NewCodeKeeper(&k).SetAllowArbitrary(ctx, arbitraryEnabled)
}

// SetAllowedPublishers - update allowed publishers
func (k Keeper) SetAllowedPublishers(ctx context.Context, allowedPublishers []vmtypes.AccountAddress) error {
	return NewCodeKeeper(&k).SetAllowedPublishers(ctx, allowedPublishers)
}

// ContractSharedRevenueRatio - percentage of fees distributed to developers
func (k Keeper) ContractSharedRevenueRatio(ctx context.Context) sdk.Dec {
	return k.GetParams(ctx).ContractSharedRevenueRatio
}

// SetParams sets the x/move module parameters.
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	if err := k.SetRawParams(ctx, params.ToRaw()); err != nil {
		return err
	}

	return NewCodeKeeper(&k).SetAllowArbitrary(ctx, params.ArbitraryEnabled)
}

// GetParams returns the x/move module parameters.
func (k Keeper) GetParams(ctx context.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		panic("params not found")
	}

	rawParams := types.RawParams{}
	k.cdc.MustUnmarshal(bz, &rawParams)

	allowArbitrary, allowedPublishers, err := NewCodeKeeper(&k).GetParams(ctx)
	if err != nil {
		panic(err)
	}

	_allowedPublishers := make([]string, len(allowedPublishers))
	for i, addr := range allowedPublishers {
		_allowedPublishers[i] = addr.String()
	}

	return rawParams.ToParams(allowArbitrary, _allowedPublishers)
}

// SetRawParams stores raw params to store.
func (k Keeper) SetRawParams(ctx context.Context, params types.RawParams) error {
	store := ctx.KVStore(k.storeKey)
	if bz, err := k.cdc.Marshal(&params); err != nil {
		return err
	} else {
		store.Set(types.ParamsKey, bz)
	}

	return nil
}
