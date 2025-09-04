package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	math "cosmossdk.io/math"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/similadayo/cosmosbank/x/bank/keeper"
	"github.com/similadayo/cosmosbank/x/bank/types"
)

const denom = "ucosmos"

// -----------------------------
// Test Suite Setup
// -----------------------------

type BankTestSuite struct {
	suite.Suite
	ctx    sdk.Context
	keeper keeper.Keeper
	server keeper.MsgServer
	addr1  sdk.AccAddress
	addr2  sdk.AccAddress
}

func (s *BankTestSuite) SetupTest() {
	// in-memory DB + store
	db := dbm.NewMemDB()
	logger := log.NewNopLogger()
	metrics := metrics.NewMetrics(nil)
	cms := store.NewCommitMultiStore(db, logger, metrics)

	key := storetypes.NewKVStoreKey(types.StoreKey)
	cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	err := cms.LoadLatestVersion()
	s.Require().NoError(err)

	reg := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(reg)

	s.keeper = keeper.NewKeeper(cdc, key)
	s.ctx = sdk.NewContext(cms, tmproto.Header{Height: 1}, false, log.NewNopLogger())
	s.server = keeper.NewMsgServerImpl(s.keeper)

	// sample addrs
	s.addr1 = sdk.AccAddress([]byte("addr1_______________"))
	s.addr2 = sdk.AccAddress([]byte("addr2_______________"))
}

// -----------------------------
// Keeper Tests
// -----------------------------

func (s *BankTestSuite) TestKeeper_SetAndGetBalance() {
	tests := []struct {
		name     string
		addr     sdk.AccAddress
		amount   sdk.Coins
		expected string
	}{
		{
			name:     "set and get non-empty balance",
			addr:     s.addr1,
			amount:   sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(100))),
			expected: "100" + denom,
		},
		{
			name:     "get empty balance",
			addr:     s.addr2,
			amount:   sdk.NewCoins(),
			expected: "",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.keeper.SetBalance(s.ctx, tc.addr, tc.amount)
			balance := s.keeper.GetBalance(s.ctx, tc.addr)
			s.Equal(tc.expected, balance.String())
		})
	}
}

func (s *BankTestSuite) TestKeeper_SendCoins() {
	tests := []struct {
		name       string
		initAmount sdk.Coins
		sendAmount sdk.Coins
		expectErr  bool
		expected1  string
		expected2  string
	}{
		{
			name:       "successful send",
			initAmount: sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(100))),
			sendAmount: sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(40))),
			expectErr:  false,
			expected1:  "60" + denom,
			expected2:  "40" + denom,
		},
		{
			name:       "insufficient funds",
			initAmount: sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(30))),
			sendAmount: sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(50))),
			expectErr:  true,
			expected1:  "30" + denom,
			expected2:  "",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			s.keeper.SetBalance(s.ctx, s.addr1, tc.initAmount)
			s.keeper.SetBalance(s.ctx, s.addr2, sdk.NewCoins())

			err := s.keeper.SendCoins(s.ctx, s.addr1, s.addr2, tc.sendAmount)

			if tc.expectErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}

			b1 := s.keeper.GetBalance(s.ctx, s.addr1)
			b2 := s.keeper.GetBalance(s.ctx, s.addr2)

			s.Equal(tc.expected1, b1.String())
			s.Equal(tc.expected2, b2.String())
		})
	}
}

// -----------------------------
// MsgServer Tests
// -----------------------------

func (s *BankTestSuite) TestMsgServer_Send_TableDriven() {
	tests := []struct {
		name       string
		amount     string
		expectErr  bool
		expectFrom int64
		expectTo   int64
	}{
		{
			name:       "valid transfer",
			amount:     "50",
			expectErr:  false,
			expectFrom: 50, // 100 - 50
			expectTo:   50,
		},
		{
			name:       "zero amount",
			amount:     "0",
			expectErr:  true,
			expectFrom: 100, // unchanged
			expectTo:   0,
		},
		{
			name:       "invalid amount string",
			amount:     "abc",
			expectErr:  true,
			expectFrom: 100, // unchanged
			expectTo:   0,
		},
		{
			name:       "insufficient funds",
			amount:     "200",
			expectErr:  true,
			expectFrom: 100, // unchanged
			expectTo:   0,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			// reset balances before each case
			s.keeper.SetBalance(s.ctx, s.addr1, sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(100))))
			s.keeper.SetBalance(s.ctx, s.addr2, sdk.NewCoins())

			_, err := s.server.Send(s.ctx, &types.MsgSend{
				FromAddress: s.addr1.String(),
				ToAddress:   s.addr2.String(),
				Amount:      tc.amount,
			})

			if tc.expectErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}

			// check balances regardless of error
			fromBal := s.keeper.GetBalance(s.ctx, s.addr1)
			toBal := s.keeper.GetBalance(s.ctx, s.addr2)

			s.Equal(tc.expectFrom, fromBal.AmountOf(denom).Int64(), "from balance mismatch")
			s.Equal(tc.expectTo, toBal.AmountOf(denom).Int64(), "to balance mismatch")
		})
	}
}

// -----------------------------
// Run Suite
// -----------------------------

func TestBankTestSuite(t *testing.T) {
	suite.Run(t, new(BankTestSuite))
}
