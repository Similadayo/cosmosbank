package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/similadayo/cosmosbank/x/bank/types"
	"github.com/spf13/cobra"
)

func CmdQueryBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-balance [address]",
		Short: "Query the balance of an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			address := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Balance(context.Background(), &types.QueryBalanceRequest{
				Address: address,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
