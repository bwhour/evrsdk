package client

import (
	govclient "github.com/Evrynetlabs/evrsdk/x/gov/client"
	"github.com/Evrynetlabs/evrsdk/x/params/client/cli"
	"github.com/Evrynetlabs/evrsdk/x/params/client/rest"
)

// ProposalHandler handles param change proposals
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
