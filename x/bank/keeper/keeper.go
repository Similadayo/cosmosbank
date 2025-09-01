package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/similadayo/cosmosbank/x/bank/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.Codec
}

func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey) Keeper {
	return Keeper{storeKey: storeKey, cdc: cdc}
}

func (k Keeper) SetBalance(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coins) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("balances/"))

	balance := &types.Balance{Coins: amount.String()}
	bz := k.cdc.MustMarshal(balance)
	store.Set(addr.Bytes(), bz)
}

func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("balances/"))
	bz := store.Get(addr.Bytes())
	if bz == nil {
		return sdk.NewCoins()
	}

	var balance types.Balance
	k.cdc.MustUnmarshal(bz, &balance)

	coins, err := sdk.ParseCoinsNormalized(balance.Coins)
	if err != nil {
		return sdk.NewCoins()
	}

	return coins
}

func (k Keeper) SendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) error {
	fromBalance := k.GetBalance(ctx, fromAddr)

	if !fromBalance.IsAllGTE(amount) {
		return sdkerrors.ErrInsufficientFunds
	}

	// Subtract from sender
	newFromBalance := fromBalance.Sub(amount...)
	k.SetBalance(ctx, fromAddr, newFromBalance)

	// Add to recipient
	toBalance := k.GetBalance(ctx, toAddr)
	newToBalance := toBalance.Add(amount...)
	k.SetBalance(ctx, toAddr, newToBalance)

	return nil
}
