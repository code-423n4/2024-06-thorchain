package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	openapi "gitlab.com/thorchain/thornode/openapi/gen"

	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "THORChain transaction subcommands",
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetCmdSetNodeKeys())
	cmd.AddCommand(GetCmdSetVersion())
	cmd.AddCommand(GetCmdSetIPAddress())
	cmd.AddCommand(GetCmdBan())
	cmd.AddCommand(GetCmdMimir())
	cmd.AddCommand(GetCmdNodePauseChain())
	cmd.AddCommand(GetCmdNodeResumeChain())
	cmd.AddCommand(GetCmdDeposit())
	cmd.AddCommand(GetCmdSend())
	cmd.AddCommand(GetCmdObserveTxIns())
	cmd.AddCommand(GetCmdObserveTxOuts())
	for _, subCmd := range cmd.Commands() {
		flags.AddTxFlagsToCmd(subCmd)
	}
	return cmd
}

// GetCmdDeposit command to send a native transaction
func GetCmdDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [amount] [coin] [memo]",
		Short: "sends a deposit transaction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			amt, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount (must be an integer): %w", err)
			}

			asset, err := common.NewAsset(args[1])
			if err != nil {
				return fmt.Errorf("invalid asset: %w", err)
			}

			coin := common.NewCoin(asset, cosmos.NewUint(uint64(amt)))

			msg := types.NewMsgDeposit(common.Coins{coin}, args[2], clientCtx.GetFromAddress())
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

// GetCmdSend command to send funds
func GetCmdSend() *cobra.Command {
	return &cobra.Command{
		Use:   "send [to_address] [coins]",
		Short: "sends funds",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			toAddr, err := cosmos.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid address: %w", err)
			}

			coins, err := cosmos.ParseCoins(args[1])
			if err != nil {
				return fmt.Errorf("invalid coins: %w", err)
			}

			msg := types.NewMsgSend(clientCtx.GetFromAddress(), toAddr, coins)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

// GetCmdMimir command to change a mimir attribute
func GetCmdMimir() *cobra.Command {
	return &cobra.Command{
		Use:   "mimir [key] [value]",
		Short: "updates a mimir attribute (admin only)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			val, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid value (must be an integer): %w", err)
			}

			msg := types.NewMsgMimir(strings.ToUpper(args[0]), val, clientCtx.GetFromAddress())
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

// GetCmdNodePauseChain command to change node pause chain
func GetCmdNodePauseChain() *cobra.Command {
	return &cobra.Command{
		Use:   "pause-chain",
		Short: "globally pause chain (NOs only)",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgNodePauseChain(int64(1), clientCtx.GetFromAddress())
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

// GetCmdNodeResumeChain command to change node resume chain
func GetCmdNodeResumeChain() *cobra.Command {
	return &cobra.Command{
		Use:   "resume-chain",
		Short: "globally resume chain (NOs only)",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgNodePauseChain(int64(-1), clientCtx.GetFromAddress())
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

// GetCmdBan command to ban a node accounts
func GetCmdBan() *cobra.Command {
	return &cobra.Command{
		Use:   "ban [node address]",
		Short: "votes to ban a node address (caution: costs 0.1% of minimum bond)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addr, err := cosmos.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid node address: %w", err)
			}

			msg := types.NewMsgBan(addr, clientCtx.GetFromAddress())
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

// GetCmdSetIPAddress command to set a node accounts IP Address
func GetCmdSetIPAddress() *cobra.Command {
	return &cobra.Command{
		Use:   "set-ip-address [ip address]",
		Short: "update registered ip address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetIPAddress(args[0], clientCtx.GetFromAddress())
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

// GetCmdSetVersion command to set an admin config
func GetCmdSetVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "set-version",
		Short: "update registered version",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetVersion(constants.SWVersion.String(), clientCtx.GetFromAddress())
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

// GetCmdSetNodeKeys command to add a node keys
func GetCmdSetNodeKeys() *cobra.Command {
	return &cobra.Command{
		Use:   "set-node-keys  [secp256k1] [ed25519] [validator_consensus_pub_key]",
		Short: "set node keys, the account use to sign this tx has to be whitelist first",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			secp256k1Key, err := common.NewPubKey(args[0])
			if err != nil {
				return fmt.Errorf("fail to parse secp256k1 pub key ,err:%w", err)
			}
			ed25519Key, err := common.NewPubKey(args[1])
			if err != nil {
				return fmt.Errorf("fail to parse ed25519 pub key ,err:%w", err)
			}
			pk := common.NewPubKeySet(secp256k1Key, ed25519Key)
			validatorConsPubKey, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeConsPub, args[2])
			if err != nil {
				return fmt.Errorf("fail to parse validator consensus public key: %w", err)
			}
			validatorConsPubKeyStr, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeConsPub, validatorConsPubKey)
			if err != nil {
				return fmt.Errorf("fail to convert public key to string: %w", err)
			}
			msg := types.NewMsgSetNodeKeys(pk, validatorConsPubKeyStr, clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Manual Observations
////////////////////////////////////////////////////////////////////////////////////////

// GetCmdObserveTxIns command manually observes inbound transactions.
func GetCmdObserveTxIns() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "observe-tx-ins --txids [tx-id-1],[tx-id-2],[tx-id-3] --raw [json-array]",
		Short: "manually observe inbound transactions with either --txids or --raw flag",
		Args:  cobra.ExactArgs(0),
	}

	// setup flags
	cmd.Flags().String("raw-observations", "", "raw json array of txs to observe")
	cmd.Flags().String("txids", "", "comma separated list of tx ids to observe")
	cmd.Flags().String("thornode-api", "", "thornode api endpoint")

	cmd.RunE = observeTxs(false)

	return cmd
}

// GetCmdObserveTxOuts command manually observes outbound transactions.
func GetCmdObserveTxOuts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "observe-tx-outs --txids [tx-id-1],[tx-id-2],[tx-id-3] --raw [json-array]",
		Short: "manually observe outbound transactions with either --txids or --raw flag",
		Args:  cobra.ExactArgs(0),
	}

	// setup flags
	cmd.Flags().String("raw-observations", "", "raw json array of txs to observe")
	cmd.Flags().String("txids", "", "comma separated list of tx ids to observe")
	cmd.Flags().String("thornode-api", "", "thornode api endpoint")

	cmd.RunE = observeTxs(true)

	return cmd
}

// ------------------------------ internal ------------------------------

// observeTxs returns a command handler for observing inbound or outbound transactions.
// If outbound is false the returned closure will observe inbounds, otherwise outbounds.
func observeTxs(outbound bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientTxContext(cmd)
		if err != nil {
			return err
		}

		// get txids and raw observations from flags
		txids, err := cmd.Flags().GetString("txids")
		if err != nil {
			return err
		}
		rawObservations, err := cmd.Flags().GetString("raw-observations")
		if err != nil {
			return err
		}
		thorNodeAPI, err := cmd.Flags().GetString("thornode-api")
		if err != nil {
			return err
		}

		// if both flags are empty, return an error
		if txids == "" && rawObservations == "" {
			return fmt.Errorf("either --txids or --raw-observations flag must be set")
		}
		// if both flags are set, return an error
		if txids != "" && rawObservations != "" {
			return fmt.Errorf("only one of --txids or --raw-observations flag can be set")
		}
		// if txids is set, the thornoode api endpoint must be set
		if txids != "" && thorNodeAPI == "" {
			return fmt.Errorf("--txids requires --thornode-api")
		}

		// ensure node address is set
		nodeAddress := clientCtx.GetFromAddress().String()
		if nodeAddress == "" {
			return fmt.Errorf("--from must be set")
		}

		var observations []types.ObservedTx

		// if txids is set, retrieve highest count observation this node did not broadcast
		if txids != "" {
			txidList := strings.Split(txids, ",")
			for _, txid := range txidList {
				observation, err := findLackingObservation(txid, nodeAddress, thorNodeAPI)
				if err != nil {
					return fmt.Errorf("failed to find lacking observation: %w", err)
				}
				if observation == nil {
					continue
				}
				observations = append(observations, *observation)

				// small sleep to avoid spamming the node or hitting rate limit
				time.Sleep(100 * time.Millisecond)
			}
		}

		// if rawObservations is set, parse the raw json array into observations
		if rawObservations != "" {
			if err := json.Unmarshal([]byte(rawObservations), &observations); err != nil {
				return fmt.Errorf("failed to parse raw observations: %w", err)
			}
		}

		// allow up to 50 txs
		if len(observations) > 50 {
			return fmt.Errorf("cannot observe more than 50 transactions at once")
		}

		// abort if no observations
		if len(observations) == 0 {
			fmt.Println("node has broadcast all observation versions")
			return nil
		}

		// create the message
		var msg cosmos.Msg
		if outbound {
			msg = types.NewMsgObservedTxOut(observations, clientCtx.GetFromAddress())
		} else {
			msg = types.NewMsgObservedTxIn(observations, clientCtx.GetFromAddress())
		}

		// output the message to be broadcast
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		if err := enc.Encode(msg); err != nil {
			return fmt.Errorf("failed to encode message: %w", err)
		}

		return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
	}
}

// findLackingObservation retrieves the highest count observation this node did not
// broadcast for the provided transaction ID.
func findLackingObservation(txid, address, thornodeAPI string) (*types.ObservedTx, error) {
	// get tx details from thornode API
	url := fmt.Sprintf("%s/thorchain/tx/details/%s", thornodeAPI, txid)
	resp, err := http.Get(url) // trunk-ignore(golangci-lint/gosec): variable url ok
	if err != nil {
		return nil, fmt.Errorf("failed to get tx details: %w", err)
	}

	// parse the response
	var txDetails openapi.TxDetailsResponse
	if err = json.NewDecoder(resp.Body).Decode(&txDetails); err != nil {
		return nil, fmt.Errorf("failed to parse tx details: %w", err)
	}

	// short address for concise output
	shortAddress := address[len(address)-4:]

	// find the highest count observation this node did not broadcast
	var observation *types.ObservedTx
	highestCount := 0
	for _, tx := range txDetails.Txs {
		// determine if we have already taken part in this observation
		signed := false
		for _, signer := range tx.Signers {
			if signer == address {
				signed = true
				break
			}
		}

		// skip if we already signed
		if signed {
			fmt.Printf("[%s] %s observed tx with %d others\n", txid, shortAddress, len(tx.Signers))
			continue
		}

		// determine if this observation has a higher count
		if len(tx.Signers) > highestCount {
			highestCount = len(tx.Signers)
			observation, err = extractOpenAPIObservedTx(tx)
			if err != nil {
				return nil, fmt.Errorf("failed to extract observed tx: %w", err)
			}
		}
	}

	if observation != nil {
		fmt.Printf("[%s] %s will observe tx with %d signers\n", txid, shortAddress, highestCount)
	}

	return observation, nil
}

func extractOpenAPIObservedTx(otx openapi.ObservedTx) (*types.ObservedTx, error) {
	tx := &types.ObservedTx{}
	var err error

	tx.Tx.ID, err = common.NewTxID(*otx.Tx.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tx id: %w", err)
	}

	tx.Tx.Chain, err = common.NewChain(*otx.Tx.Chain)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chain: %w", err)
	}

	tx.Tx.FromAddress, err = common.NewAddress(*otx.Tx.FromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse from address: %w", err)
	}

	tx.Tx.ToAddress, err = common.NewAddress(*otx.Tx.ToAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse to address: %w", err)
	}

	tx.Tx.Coins = make(common.Coins, len(otx.Tx.Coins))
	for i, coin := range otx.Tx.Coins {
		var coinAsset common.Asset
		coinAsset, err = common.NewAsset(coin.Asset)
		if err != nil {
			return nil, fmt.Errorf("failed to parse coin asset: %w", err)
		}
		coinAmount := cosmos.NewUintFromString(coin.Amount)

		tx.Tx.Coins[i] = common.NewCoin(coinAsset, coinAmount)
		if coin.Decimals != nil {
			tx.Tx.Coins[i].Decimals = *coin.Decimals
		}
	}

	tx.Tx.Gas = common.Gas(make(common.Coins, len(otx.Tx.Gas)))
	for i, coin := range otx.Tx.Gas {
		var coinAsset common.Asset
		coinAsset, err = common.NewAsset(coin.Asset)
		if err != nil {
			return nil, fmt.Errorf("failed to parse gas asset: %w", err)
		}
		coinAmount := cosmos.NewUintFromString(coin.Amount)

		tx.Tx.Gas[i] = common.NewCoin(coinAsset, coinAmount)
		if coin.Decimals != nil {
			tx.Tx.Gas[i].Decimals = *coin.Decimals
		}
	}

	tx.Tx.Memo = *otx.Tx.Memo

	if otx.Aggregator != nil {
		tx.Aggregator = *otx.Aggregator
	}
	if otx.AggregatorTarget != nil {
		tx.AggregatorTarget = *otx.AggregatorTarget
	}
	if otx.AggregatorTargetLimit != nil {
		target := cosmos.NewUintFromString(*otx.AggregatorTargetLimit)
		tx.AggregatorTargetLimit = &target
	}

	tx.BlockHeight = *otx.ExternalObservedHeight
	tx.FinaliseHeight = *otx.ExternalConfirmationDelayHeight

	tx.ObservedPubKey, err = common.NewPubKey(*otx.ObservedPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse observed public key: %w", err)
	}

	return tx, nil
}
