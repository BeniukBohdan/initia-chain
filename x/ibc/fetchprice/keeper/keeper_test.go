package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/initia-labs/initia/x/ibc/fetchprice/keeper"
	"github.com/initia-labs/initia/x/ibc/fetchprice/types"
	ibctesting "github.com/initia-labs/initia/x/ibc/testing"

	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v8/types"
)

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
}

func (suite *KeeperTestSuite) SetupTest() {
	// to bypass authority check
	keeper.IsTesting = true

	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

func NewFetchPricePath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = types.PortID
	path.EndpointB.ChannelConfig.PortID = icqtypes.PortID
	path.EndpointA.ChannelConfig.Version = types.Version
	path.EndpointB.ChannelConfig.Version = types.Version

	return path
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
