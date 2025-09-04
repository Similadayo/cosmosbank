package bank

import (
	"encoding/json"

	"cosmossdk.io/api/tendermint/abci"
	client "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/similadayo/cosmosbank/x/bank/client/cli"
	"github.com/similadayo/cosmosbank/x/bank/keeper"
	"github.com/similadayo/cosmosbank/x/bank/types"
	"github.com/spf13/cobra"
)

type AppModule struct{}

func (AppModule) Name() string {
	return banktypes.ModuleName
}

func (AppModule) RegisteerLegacyAminoCodec(cdc *codec.LegacyAmino) {
	banktypes.RegisterLegacyAminoCodec(cdc)
}

func (AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(banktypes.DefaultGenesisState())
}

func (AppModule) validateGenesisCodec(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data banktypes.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return err
	}

	return banktypes.ValidateGenesis(&data)
}

// AppModule implements an application module for the bank module.
type AppModuleBasic struct {
	AppModule
	keeper keeper.Keeper
}

func NewAppModuleBasic(keeper keeper.Keeper) AppModuleBasic {
	return AppModuleBasic{
		AppModule: AppModule{},
		keeper:    keeper,
	}
}

func (am AppModuleBasic) Name() string {
	return banktypes.ModuleName
}

func (am AppModuleBasic) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

func (am AppModuleBasic) InitGenesis(ctx client.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState banktypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModuleBasic) ExportGenesis(ctx client.Context, cdc codec.JSONCodec) json.RawMessage {
	genesisState := banktypes.DefaultGenesisState()
	return cdc.MustMarshalJSON(genesisState)
}

func (am AppModuleBasic) ConsensusVersion() uint64 {
	return 1
}

func (am AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

func (am AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.NewQueryCmd()
}
