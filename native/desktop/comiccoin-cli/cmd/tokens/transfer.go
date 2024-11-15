package tokens

import (
	"github.com/spf13/cobra"
)

// Command line argument flags
var (
	flagSenderAccountAddress          string
	flagSenderAccountPassword         string
	flagSenderAccountPasswordRepeated string
	flagRecipientAddress              string
	flagQuantity                      uint64
	flagData                          string

	flagDataDirectory     string
	flagChainID           uint16
	flagAuthorityAddress  string
	flagNFTStorageAddress string
)

func TransferTokensCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "transfer",
		Short: "Submit a (pending) transaction to the ComicCoin blockchain network to transfer tokens from your account to another account",
		Run: func(cmd *cobra.Command, args []string) {
			doRunTransferTokensCommand()
		},
	}

	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", preferences.DataDirectory, "The data directory to save to")
	cmd.Flags().Uint16Var(&flagChainID, "chain-id", preferences.ChainID, "The blockchain to sync with")
	cmd.Flags().StringVar(&flagAuthorityAddress, "authority-address", preferences.AuthorityAddress, "The BlockChain authority address to connect to")
	cmd.Flags().StringVar(&flagNFTStorageAddress, "nftstorage-address", preferences.NFTStorageAddress, "The NFT storage service adress to connect to")

	cmd.Flags().StringVar(&flagSenderAccountAddress, "sender-account-address", "", "The address of the account we will use in our token transfer")
	cmd.MarkFlagRequired("sender-account-address")

	cmd.Flags().StringVar(&flagSenderAccountPassword, "sender-account-password", "", "The password to unlock the account which will transfer the token")
	cmd.MarkFlagRequired("sender-account-password")

	cmd.Flags().Uint64Var(&flagQuantity, "value", 0, "The amount of tokens to send")
	cmd.MarkFlagRequired("value")

	cmd.Flags().StringVar(&flagData, "data", "", "Optional data to include with this transaction")

	cmd.Flags().StringVar(&flagRecipientAddress, "recipient-address", "", "The address of the account whom will receive this token")
	cmd.MarkFlagRequired("recipient-address")

	return cmd
}

func doRunTransferTokensCommand() {

}
