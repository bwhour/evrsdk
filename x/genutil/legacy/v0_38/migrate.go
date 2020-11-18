package v038

import (
	"github.com/Evrynetlabs/evrsdk/codec"
	v036auth "github.com/Evrynetlabs/evrsdk/x/auth/legacy/v0_36"
	v038auth "github.com/Evrynetlabs/evrsdk/x/auth/legacy/v0_38"
	v036distr "github.com/Evrynetlabs/evrsdk/x/distribution/legacy/v0_36"
	v038distr "github.com/Evrynetlabs/evrsdk/x/distribution/legacy/v0_38"
	v036genaccounts "github.com/Evrynetlabs/evrsdk/x/genaccounts/legacy/v0_36"
	"github.com/Evrynetlabs/evrsdk/x/genutil/types"
	v036staking "github.com/Evrynetlabs/evrsdk/x/staking/legacy/v0_36"
	v038staking "github.com/Evrynetlabs/evrsdk/x/staking/legacy/v0_38"
)

// Migrate migrates exported state from v0.36/v0.37 to a v0.38 genesis state.
func Migrate(appState types.AppMap) types.AppMap {
	v036Codec := codec.New()
	codec.RegisterCrypto(v036Codec)

	v038Codec := codec.New()
	codec.RegisterCrypto(v038Codec)
	v038auth.RegisterCodec(v038Codec)

	if appState[v036genaccounts.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var authGenState v036auth.GenesisState
		v036Codec.MustUnmarshalJSON(appState[v036auth.ModuleName], &authGenState)

		var genAccountsGenState v036genaccounts.GenesisState
		v036Codec.MustUnmarshalJSON(appState[v036genaccounts.ModuleName], &genAccountsGenState)

		// delete deprecated genaccounts genesis state
		delete(appState, v036genaccounts.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v038auth.ModuleName] = v038Codec.MustMarshalJSON(v038auth.Migrate(authGenState, genAccountsGenState))
	}

	// migrate staking state
	if appState[v036staking.ModuleName] != nil {
		var stakingGenState v036staking.GenesisState
		v036Codec.MustUnmarshalJSON(appState[v036staking.ModuleName], &stakingGenState)

		delete(appState, v036staking.ModuleName) // delete old key in case the name changed
		appState[v038staking.ModuleName] = v038Codec.MustMarshalJSON(v038staking.Migrate(stakingGenState))
	}

	// migrate distribution state
	if appState[v036distr.ModuleName] != nil {
		var distrGenState v036distr.GenesisState
		v036Codec.MustUnmarshalJSON(appState[v036distr.ModuleName], &distrGenState)

		delete(appState, v036distr.ModuleName) // delete old key in case the name changed
		appState[v038distr.ModuleName] = v038Codec.MustMarshalJSON(v038distr.Migrate(distrGenState))
	}

	return appState
}
