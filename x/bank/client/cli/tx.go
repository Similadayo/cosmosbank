package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/similadayo/cosmosbank/x/bank/types"
	"github.com/spf13/cobra"
)

// CmdSend sends tokens from one address to another.
func CmdSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [from_address] [to_address] [amount]",
		Short: "Send tokens from one address to another",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromAddr := args[0]
			toAddr := args[1]
			amount := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgSend{
				FromAddress: fromAddr,
				ToAddress:   toAddr,
				Amount:      amount,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
