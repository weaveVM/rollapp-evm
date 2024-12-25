package drs4

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	rollappparamskeeper "github.com/dymensionxyz/dymension-rdk/x/rollappparams/keeper"
	evmkeeper "github.com/evmos/evmos/v12/x/evm/keeper"

	drs3 "github.com/dymensionxyz/rollapp-evm/app/upgrades/drs-3"
)

func CreateUpgradeHandler(
	rpKeeper rollappparamskeeper.Keeper,
	evmKeeper *evmkeeper.Keeper,
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if rpKeeper.Version(ctx) < 3 {
			// first run drs-3 migration
			if err := drs3.HandleUpgrade(ctx, rpKeeper, evmKeeper); err != nil {
				return nil, err
			}
		}
		// upgrade drs to 4
		if err := rpKeeper.SetVersion(ctx, uint32(4)); err != nil {
			return nil, err
		}
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}