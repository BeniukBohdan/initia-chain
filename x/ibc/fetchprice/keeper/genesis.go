package keeper

import (
	"context"
	"fmt"

	"github.com/initia-labs/initia/x/ibc/fetchprice/types"
)

// InitGenesis initializes the ibc-transfer state and binds to PortID.
func (k Keeper) InitGenesis(ctx context.Context, state types.GenesisState) {
	if err := k.PortID.Set(ctx, state.PortId); err != nil {
		panic(err)
	}

	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if !k.IsBound(ctx, state.PortId) {
		// transfer module binds to the transfer port on InitChain
		// and claims the returned capability
		err := k.BindPort(ctx, state.PortId)
		if err != nil {
			panic(fmt.Sprintf("could not claim port capability: %v", err))
		}
	}

	if err := k.Params.Set(ctx, state.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis exports ibc-transfer module's portID and denom trace info into its genesis state.
func (k Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	portID, err := k.PortID.Get(ctx)
	if err != nil {
		panic(err)
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		PortId: portID,
		Params: params,
	}
}
