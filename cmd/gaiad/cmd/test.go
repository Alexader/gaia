package cmd

import (
	"fmt"
	"strconv"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// NewParallelSendTxCmd returns a CLI command handler for creating a MsgSend transaction.
func NewParallelSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "send [from_key_or_address] [to_address] [amount]",
		Short: `Send funds from one account to another. Note, the'--from' flag is
ignored as it is implied from [from_key_or_address].`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return parallelRun(cmd)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	addParallelFlags(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func addParallelFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("type", "p", "transfer,interchain", "To send tx type, like interchain,transfer")
	cmd.Flags().IntP("duration", "d", 60, "Test duration in second")
	cmd.Flags().IntP("tps", "t", 6, "Test target tps")
	cmd.Flags().IntP("concurrent", "c", 1, "Concurrent thread number")
}

var (
	froms = []string{
		"cosmos1uxcv2ux5jyp9grrsvgte4udfkrqfs8j26w0sg3",
	}
	tos = []string{
		"cosmos16jdrmjfm2ygqldzgkuurgsnly8sqn4x5g8mtw7",
	}
)

func parallelRun(cmd *cobra.Command) error {
	flagSet := cmd.Flags()
	typ, err := flagSet.GetString("type")
	if err != nil {
		return err
	}
	concurrent, err := flagSet.GetInt("concurrent")
	if err != nil {
		return err
	}
	tps, err := flagSet.GetInt("tps")
	if err != nil {
		return err
	}
	duration, err := flagSet.GetInt("duration")
	if err != nil {
		return err
	}
	fmt.Printf("%s, %d, %d, %d\n", typ, concurrent, tps, duration)

	var wg sync.WaitGroup
	wg.Add(concurrent)
	for i := 0; i < concurrent; i++ {
		// set from and to address
		cmd.Flags().Set(flags.FlagFrom, froms[i])
		clientCtx, err := client.GetClientTxContext(cmd)
		if err != nil {
			return err
		}
		accNum, accSeq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, clientCtx.GetFromAddress())
		if err != nil {
			return err
		}
		fmt.Printf("account content is %d and %d\n", accNum, accSeq)

		toAddr, err := sdk.AccAddressFromBech32(tos[i])
		if err != nil {
			return err
		}

		coins, err := sdk.ParseCoinsNormalized("1validatortoken")
		if err != nil {
			return err
		}

		msg := types.NewMsgSend(clientCtx.GetFromAddress(), toAddr, coins)
		if err := msg.ValidateBasic(); err != nil {
			fmt.Printf("new tx error %s\n", err.Error())
		}

		go func(clientCtx client.Context) {
			defer wg.Done()
			for j := 0; j < tps; j++ {
				err := flagSet.Set(flags.FlagSequence, strconv.FormatUint(accSeq, 10))
				if err != nil {
					fmt.Printf("send tx error %s\n", err.Error())
					return
				}
				err = tx.GenerateOrBroadcastTxCLI(clientCtx, flagSet, msg)
				if err != nil {
					fmt.Printf("send tx error %s\n", err.Error())
				}
				accSeq++
			}
		}(clientCtx)
	}
	wg.Wait()
	return nil
}
