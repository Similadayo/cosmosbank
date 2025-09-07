package main

import (
	"io"
	"os"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	"github.com/similadayo/cosmosbank/app"
	"github.com/similadayo/cosmosbank/x/bank/client/cli"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cosmosbankd",
		Short: "CosmosBank Daemon (node & CLI)",
	}

	// --- IMPORTANT: pass a ModuleInitFlags value (no-op if you don't need module flags) ---
	// This avoids the nil-function panic inside server.AddCommands.
	var addModuleInitFlags servertypes.ModuleInitFlags = func(startCmd *cobra.Command) {
		// no per-module start flags for now
	}

	// register standard server commands (start, export, etc.)
	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, addModuleInitFlags)

	// add module tx/query cli roots
	rootCmd.AddCommand(
		cli.NewTxCmd(),
		cli.NewQueryCmd(),
	)

	// execute
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome, app.AppName); err != nil {
		os.Exit(1)
	}
}

// newApp must match servertypes.AppCreator signature
func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	return app.NewCosmosBankApp(logger, db)
}

// appExport must match servertypes.AppExporter signature
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
