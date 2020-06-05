package simapp

import (
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/modules/incubator/nft"
)

const appName = "SimApp"

var (
	// DefaultCLIHome default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.simapp")

	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome = os.ExpandEnv("$HOME/.simapp")

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		supply.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		params.AppModuleBasic{},
		nft.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{}
)

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

// Verify app interface at compile time
var _ simapp.App = (*SimApp)(nil)

// SimApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type SimApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	AccountKeeper auth.AccountKeeper
	BankKeeper    bank.Keeper
	SupplyKeeper  supply.Keeper
	StakingKeeper staking.Keeper
	ParamsKeeper  params.Keeper
	NFTKeeper     nft.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// NewSimApp returns a reference to an initialized SimApp.
func NewSimApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *SimApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, params.StoreKey, nft.StoreKey)
	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	app := &SimApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
		subspaces:      make(map[string]params.Subspace),
	}

	// init params keeper and subspaces
	app.ParamsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.ParamsKeeper.Subspace(staking.DefaultParamspace)

	// add keepers
	app.AccountKeeper = auth.NewAccountKeeper(
		app.cdc, keys[auth.StoreKey], app.subspaces[auth.ModuleName], auth.ProtoBaseAccount,
	)
	app.BankKeeper = bank.NewBaseKeeper(
		app.AccountKeeper, app.subspaces[bank.ModuleName], app.BlacklistedAccAddrs(),
	)
	app.SupplyKeeper = supply.NewKeeper(
		app.cdc, keys[supply.StoreKey], app.AccountKeeper, app.BankKeeper, maccPerms,
	)
	stakingKeeper := staking.NewKeeper(
		app.cdc, keys[staking.StoreKey], app.SupplyKeeper, app.subspaces[staking.ModuleName],
	)
	app.NFTKeeper = nft.NewKeeper(app.cdc, keys[nft.StoreKey])

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(),
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.SupplyKeeper),
		nft.NewAppModule(app.NFTKeeper, app.AccountKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderEndBlockers(staking.ModuleName)

	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		auth.ModuleName, staking.ModuleName, bank.ModuleName,
		supply.ModuleName, genutil.ModuleName, nft.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.SupplyKeeper),
		nft.NewAppModule(app.NFTKeeper, app.AccountKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.AccountKeeper, app.SupplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}
	return app
}

// application updates every begin block
func (app *SimApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// application updates every end block
func (app *SimApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// application update at chain initialization
func (app *SimApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

// load a particular height
func (app *SimApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *SimApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlacklistedAccAddrs returns all the app's module account addresses black listed for receiving tokens.
func (app *SimApp) BlacklistedAccAddrs() map[string]bool {
	blacklistedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blacklistedAddrs[supply.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blacklistedAddrs
}

// Codec returns simapp's codec
func (app *SimApp) Codec() *codec.Codec {
	return app.cdc
}

// GetKey returns the KVStoreKey for the provided store key
func (app *SimApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key
func (app *SimApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// SimulationManager implements the SimulationApp interface
func (app *SimApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}
