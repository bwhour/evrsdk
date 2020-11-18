package keeper_test

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Evrynetlabs/evrsdk/simapp"
	sdk "github.com/Evrynetlabs/evrsdk/types"
	"github.com/Evrynetlabs/evrsdk/x/auth"
)

func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})

	app.AccountKeeper.SetParams(ctx, auth.DefaultParams())
	app.BankKeeper.SetSendEnabled(ctx, true)

	return app, ctx
}
