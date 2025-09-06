package app

import (
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"

	"github.com/similadayo/cosmosbank/x/bank"
	"github.com/similadayo/cosmosbank/x/bank/keeper"
	"github.com/similadayo/cosmosbank/x/bank/types"
)

type CosmosBankApp struct {
	*baseapp.BaseApp

	cdc        codec.Codec
	BankKeeper keeper.Keeper
	mm         *module.Manager
}

func NewCosmosBankApp(logger log.Logger, db dbm.DB) *CosmosBankApp {
	// encoding
	encoding := MakeEncodingConfig()
	cdc := encoding.Marshaler

	bApp := baseapp.NewBaseApp("cosmosbank", logger, db, encoding.TxConfig.TxDecoder())
	bApp.SetVersion("0.1.0")

	// keys
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	// keeper
	bankKeeper := keeper.NewKeeper(cdc, storeKey)

	// module
	appModule := bank.NewAppModule(bankKeeper)
	mm := module.NewManager(appModule)

	// register services
	configurator := module.NewConfigurator(cdc, bApp.MsgServiceRouter(), bApp.GRPCQueryRouter())
	mm.RegisterServices(configurator)

	app := &CosmosBankApp{
		BaseApp:    bApp,
		cdc:        cdc,
		BankKeeper: bankKeeper,
		mm:         mm,
	}

	// mount store
	app.MountStore(storeKey, storetypes.StoreTypeIAVL)

	return app
}

// --- encoding ---
type EncodingConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          client.TxConfig
}

func MakeEncodingConfig() EncodingConfig {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
	}
}
