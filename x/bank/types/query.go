package types

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func RegisterQueryHandlerClient(ctx client.Context, mux *runtime.ServeMux, client QueryClient) error {
	return nil
}
