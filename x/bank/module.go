package bank

import (
	"encoding/json"

	"cosmossdk.io/api/tendermint/abci"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	banktype "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/similadayo/cosmosbank/x/bank/client/cli"
	"github.com/similadayo/cosmosbank/x/bank/keeper"
	"github.com/similadayo/cosmosbank/x/bank/types"
	"github.com/spf13/cobra"
)

// AppModuleBasic defines the basic application module used by the bank module.
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return banktypes.ModuleName
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	banktypes.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(banktypes.DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data banktype.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return err
	}
	return banktype.ValidateGenesis(&data)
}

// AppModule implements an application module for the bank module.
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

func NewAppModule(keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

func (am AppModule) Name() string {
	return banktypes.ModuleName
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState banktypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	// init your keeper with genesisState if needed
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genesisState := banktypes.DefaultGenesisState()
	return cdc.MustMarshalJSON(genesisState)
}

func (am AppModule) ConsensusVersion() uint64 {
	return 1
}

func (am AppModule) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

func (am AppModule) GetQueryCmd() *cobra.Command {
	return cli.NewQueryCmd()
}

func (am AppModule) IsAppModule() {}

func (am AppModule) IsOnePerModuleType() {}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(
		clientCtx, // Cosmos SDK client.Context
		mux,
		types.NewQueryClient(clientCtx),
	); err != nil {
		panic(err)
	}
}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}
