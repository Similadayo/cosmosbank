package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/log"
	store "cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/similadayo/cosmosbank/x/bank/keeper"
	banktypes "github.com/similadayo/cosmosbank/x/bank/types"
)

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context) {
	key := storetypes.NewKVStoreKey(banktypes.StoreKey)

	// create in-memory DB + store
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	if err := cms.LoadLatestVersion(); err != nil {
		t.Fatalf("failed to load store: %v", err)
	}

	reg := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(reg)

	k := keeper.NewKeeper(cdc, key)
	ctx := sdk.NewContext(cms, tmproto.Header{Height: 1}, false, log.NewNopLogger())

	return k, ctx
}

func TestSetAndGetBalance(t *testing.T) {
	k, ctx := setupKeeper(t)

	addr := sdk.AccAddress([]byte("testaddr____________"))
	amount := sdk.NewCoins(sdk.NewInt64Coin("ucosmos", 1000))

	k.SetBalance(ctx, addr, amount)
	balance := k.GetBalance(ctx, addr)

	require.Equal(t, amount, balance)
}

func TestSendCoins_EmptyAmount(t *testing.T) {
	k, ctx := setupKeeper(t)
	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	k.SetBalance(ctx, addr1, sdk.NewCoins(sdk.NewInt64Coin("ucosmos", 100)))

	err := k.SendCoins(ctx, addr1, addr2, sdk.NewCoins())
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty amount")
}

func TestMultipleDenoms(t *testing.T) {
	k, ctx := setupKeeper(t)
	addr1 := sdk.AccAddress([]byte("addr1_______________"))
	addr2 := sdk.AccAddress([]byte("addr2_______________"))

	k.SetBalance(ctx, addr1, sdk.NewCoins(
		sdk.NewInt64Coin("ucosmos", 500),
		sdk.NewInt64Coin("uatom", 300),
	))

	err := k.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewInt64Coin("uatom", 200)))
	require.NoError(t, err)

	balance1 := k.GetBalance(ctx, addr1)
	balance2 := k.GetBalance(ctx, addr2)

	require.Contains(t, balance1.String(), "uatom")
	require.Contains(t, balance2.String(), "uatom")
}
