package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/similadayo/cosmosbank/x/bank/types"
)

type QueryServer struct {
	Keeper
	types.UnimplementedQueryServer
}

func NewQueryServerImpl(keeper Keeper) QueryServer {
	return QueryServer{Keeper: keeper}
}

func (k QueryServer) Balance(ctx context.Context, req *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	balance := k.GetBalance(sdkCtx, addr)
	return &types.QueryBalanceResponse{Balance: balance.String()}, nil
}
