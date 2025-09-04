package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "bank",
		Short:                      "Bank transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(CmdSend())
	return txCmd
}

func NewQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        "bank",
		Short:                      "Bank query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(CmdQueryBalance())
	return queryCmd
}
