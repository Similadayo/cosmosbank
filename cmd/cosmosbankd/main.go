package main

import (
	"io"
	"os"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"

	"github.com/similadayo/cosmosbank/app"
	"github.com/similadayo/cosmosbank/x/bank/client/cli"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cosmosbankd",
		Short: "CosmosBank Daemon (node & CLI)",
	}

	// Add persistent flags (chain-id, home, etc.)
	rootCmd.PersistentFlags().String("home", app.DefaultNodeHome, "node's home directory")

	// Add server commands (start, tendermint, etc.)
	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, nil)

	// Add tx/query commands from your bank module
	rootCmd.AddCommand(
		cli.NewTxCmd(),
		cli.NewQueryCmd(),
	)

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome, app.AppName); err != nil {
		os.Exit(1)
	}

}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	return app.NewCosmosBankApp(logger, db)
}

func appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	cosmosBankApp := app.NewCosmosBankApp(logger, db)
	return cosmosBankApp.ExportApp(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
