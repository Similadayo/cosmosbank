package keeper

import (
	"context"

	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/similadayo/cosmosbank/x/bank/types"
)

const denom = "ucosmos"

type MsgServer struct {
	Keeper
	types.UnimplementedMsgServer
}

func NewMsgServerImpl(keeper Keeper) MsgServer {
	return MsgServer{Keeper: keeper}
}

func (k MsgServer) Send(goCtx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}

	to, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, err
	}

	amount := sdk.NewCoins(sdk.NewCoin(denom, math.NewIntFromUint64(msg.Amount)))

	fromBalance := k.GetBalance(ctx, from)
	if !fromBalance.IsAllGTE(amount) {
		return nil, sdkerrors.ErrInsufficientFunds
	}

	k.SetBalance(ctx, from, fromBalance.Sub(amount...))
	toBalance := k.GetBalance(ctx, to)
	k.SetBalance(ctx, to, toBalance.Add(amount...))

	return &types.MsgSendResponse{}, nil

}
