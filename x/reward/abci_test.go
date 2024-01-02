package reward_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/initia-labs/initia/x/reward/types"
)

func Test_BeginBlocker(t *testing.T) {
	app := createApp(t)

	// new block
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	ctx := app.BaseApp.NewContext(false, header)
	header.Time = app.RewardKeeper.GetLastReleaseTimestamp(ctx)

	// update params & mint coins for reward distribution
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	params := app.RewardKeeper.GetParams(ctx)
	rewardDenom := params.RewardDenom
	rewardAmount := math.NewInt(10_000_000)
	rewardCoins := sdk.NewCoins(sdk.NewCoin(rewardDenom, rewardAmount))
	err := app.BankKeeper.MintCoins(ctx, authtypes.Minter, rewardCoins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.Minter, types.ModuleName, rewardCoins)
	require.NoError(t, err)

	supply := app.BankKeeper.GetSupply(ctx, rewardDenom)

	params.ReleaseEnabled = true
	params.ReleaseRate = math.LegacyNewDecWithPrec(7, 2) // 7%
	params.DilutionPeriod = time.Hour * 24
	app.RewardKeeper.SetParams(ctx, params)

	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

	// new block after 24 hours
	header = tmproto.Header{Height: app.LastBlockHeight() + 1, Time: header.Time.Add(time.Hour * 24).Add(time.Second)}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	app.Commit()

	// check supply
	expectedReleasedAmount := math.LegacyNewDec(7).QuoInt64(100).MulInt(supply.Amount).QuoInt64(365).TruncateInt()
	checkBalance(t, app, authtypes.NewModuleAddress(types.ModuleName), rewardCoins.Sub(sdk.NewCoin(rewardDenom, expectedReleasedAmount)))

	// release rate should be half
	ctx = app.BaseApp.NewContext(true, header)
	require.Equal(t, math.LegacyNewDecWithPrec(35, 3), app.RewardKeeper.GetReleaseRate(ctx))
	require.Equal(t, header.Time, app.RewardKeeper.GetLastReleaseTimestamp(ctx))
	require.Equal(t, header.Time, app.RewardKeeper.GetLastDilutionTimestamp(ctx))
}

func Test_BeginBlockerNotEnabled(t *testing.T) {
	app := createApp(t)

	// new block
	header := tmproto.Header{Height: app.LastBlockHeight() + 1}
	ctx := app.BaseApp.NewContext(false, header)
	header.Time = app.RewardKeeper.GetLastReleaseTimestamp(ctx)

	// update params & mint coins for reward distribution
	app.BeginBlock(abci.RequestBeginBlock{Header: header})

	params := app.RewardKeeper.GetParams(ctx)
	rewardDenom := params.RewardDenom
	rewardAmount := math.NewInt(10_000_000)
	rewardCoins := sdk.NewCoins(sdk.NewCoin(rewardDenom, rewardAmount))
	err := app.BankKeeper.MintCoins(ctx, authtypes.Minter, rewardCoins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.Minter, types.ModuleName, rewardCoins)
	require.NoError(t, err)

	params.ReleaseEnabled = false
	params.ReleaseRate = math.LegacyNewDecWithPrec(7, 2) // 7%
	params.DilutionPeriod = time.Hour * 24
	app.RewardKeeper.SetParams(ctx, params)

	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

	// new block after 24 hours
	header = tmproto.Header{Height: app.LastBlockHeight() + 1, Time: header.Time.Add(time.Hour * 24)}
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	app.Commit()

	// check supply
	expectedReleasedAmount := math.ZeroInt()
	checkBalance(t, app, authtypes.NewModuleAddress(types.ModuleName), rewardCoins.Sub(sdk.NewCoin(rewardDenom, expectedReleasedAmount)))

	// only timestamps updated
	require.Equal(t, math.LegacyNewDecWithPrec(7, 2), app.RewardKeeper.GetParams(app.BaseApp.NewContext(true, tmproto.Header{})).ReleaseRate)
	require.Equal(t, header.Time, app.RewardKeeper.GetLastReleaseTimestamp(app.BaseApp.NewContext(true, tmproto.Header{})))
	require.Equal(t, header.Time, app.RewardKeeper.GetLastDilutionTimestamp(app.BaseApp.NewContext(true, tmproto.Header{})))
}
