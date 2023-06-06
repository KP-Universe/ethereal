// Copyright © 2017-2019 Weald Technology Trading
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/wealdtech/ethereal/v2/cli"
	"github.com/wealdtech/ethereal/v2/util"
)

var (
	accountKeysAddress    string
	accountKeysPassphrase string
	accountKeysPrivateKey string
)

// accountKeysCmd represents the account keys command.
var accountKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Display keys for a given address",
	Long: `Display public and private keys for a given address.  For example:

    ethereal account keys --address=0x5FfC014343cd971B7eb70732021E26C35B744cc4 --passphrase=secret

Note that this will only work for filesystem-based keystores.  Hardware wallets never reveal their keys so these cannot be obtained.

In quiet mode this will return 0 if the account was successfully decoded, otherwise 1.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.Assert((accountKeysAddress != "" && accountKeysPassphrase != "") || accountKeysPrivateKey != "", quiet, "--privatekey or both of --address and --passphrase are required")

		var key *ecdsa.PrivateKey
		if accountKeysPrivateKey != "" {
			key, err = crypto.HexToECDSA(strings.TrimPrefix(accountKeysPrivateKey, "0x"))
			cli.ErrCheck(err, quiet, "Invalid private key")
		} else {
			address := common.HexToAddress(accountKeysAddress)
			key, err = util.PrivateKeyForAccount(c.ChainID(), address, accountKeysPassphrase)
		}
		cli.ErrCheck(err, quiet, "Failed to access account")
		if quiet {
			os.Exit(exitSuccess)
		}

		fmt.Printf("Private key:\t\t0x%032x\n", key.D)
		fmt.Printf("Public key:\t\t0x%s\n", hex.EncodeToString(crypto.FromECDSAPub(&key.PublicKey)))
		fmt.Printf("Ethereum address:\t%s\n", crypto.PubkeyToAddress(key.PublicKey).Hex())
	},
}

func init() {
	offlineCmds["account:keys"] = true
	accountCmd.AddCommand(accountKeysCmd)
	accountKeysCmd.Flags().StringVar(&accountKeysAddress, "address", "", "address for account keys")
	accountKeysCmd.Flags().StringVar(&accountKeysPassphrase, "passphrase", "", "passphrase for account keys")
	accountKeysCmd.Flags().StringVar(&accountKeysPrivateKey, "privatekey", "", "private key for account keys")
}
