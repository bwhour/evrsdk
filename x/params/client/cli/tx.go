package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Evrynetlabs/evrsdk/client/context"
	"github.com/Evrynetlabs/evrsdk/codec"
	sdk "github.com/Evrynetlabs/evrsdk/types"
	"github.com/Evrynetlabs/evrsdk/version"
	"github.com/Evrynetlabs/evrsdk/x/auth"
	"github.com/Evrynetlabs/evrsdk/x/auth/client/utils"
	govtypes "github.com/Evrynetlabs/evrsdk/x/gov/types"
	paramscutils "github.com/Evrynetlabs/evrsdk/x/params/client/utils"
	"github.com/Evrynetlabs/evrsdk/x/params/types"
)

// GetCmdSubmitProposal implements a command handler for submitting a parameter
// change proposal transaction.
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "param-change [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a parameter change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a parameter proposal along with an initial deposit.
The proposal details must be supplied via a JSON file. For values that contains
objects, only non-empty fields will be updated.

IMPORTANT: Currently parameter changes are evaluated but not validated, so it is
very important that any "value" change is valid (ie. correct type and within bounds)
for its respective parameter, eg. "MaxValidators" should be an integer and not a decimal.

Proper vetting of a parameter change proposal should prevent this from happening
(no deposits should occur during the governance process), but it should be noted
regardless.

Example:
$ %s tx gov submit-proposal param-change <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Staking Param Change",
  "description": "Update max validators",
  "changes": [
    {
      "subspace": "staking",
      "key": "MaxValidators",
      "value": 105
    }
  ],
  "deposit": [
    {
      "denom": "stake",
      "amount": "10000"
    }
  ]
}
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			proposal, err := paramscutils.ParseParamChangeProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewParameterChangeProposal(proposal.Title, proposal.Description, proposal.Changes.ToParamChanges())

			msg := govtypes.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
