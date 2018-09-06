// Copyright 2018 Thunder Token Inc., The ThunderCore™ Authors
// This file comprises an original work of authorship that may make use of, or
// interface with another work licensed under a GNU or third party license, but
// which is not otherwise based on said another work.

// To the extent that portions of this file contains source code that is subject
// to the terms of the GNU or third party license, the minimal corresponding source
// code for those portions can be freely redistributed and/or modified under the
// terms of the respective license, either of GNU General Public License version 3
// or (at your option) any later version.

// The remaining code for the ThunderCore™ network application is not a contribution
// to be incorporated into said another work.  Rather, it is open source and licensed
// from Thunder Token Inc. to you, the recipient, to copy, modify and distribute the
// original or modified work without a fee, subject to reciprocity and recipient’s
// (i) promise and covenant not to sue Thunder Token Inc., its assigns, successors,
// affiliates and subsidiaries (hereinafter “Thunder Token”) on claims arising from
// any of their use of recipient’s code, if any; (ii) promise and ongoing commitment
// to not unfairly compete against or interfere with Thunder Token’s business or commercial
// relationships; and (iii) promise and ongoing commitment to not challenge the validity,
// enforceability, title, or ownership (by Thunder Token) of any intellectual property
// rights arising from or relating to the ThunderCore™ network application.  Further, you,
// the recipient, agree to and must do the following: (1) give prominent notice and
// attribution to Thunder Token Inc. and the ThunderCore™ Authors for their work on the
// original work and include any appropriate copyright, trademark, patent notices,
// (2) accompany the original or modified work with a copy of this notice (TT license v1.0
// or, at your option, any later version) in its entirety or a link directing the user to
// the same, (3) accompany the modified work with a prominent notice indicating that it
// has been modified and that it was based off of the original work; and (4) convey or
// otherwise make freely available the source code corresponding to the modified work
// under the same conditions and restrictions on the exercise of rights granted or
// affirmed under this license.

// Your copying, reverse-engineering, debugging, modifying, or distributing the original
// or modified work constitutes assent and agreement to these terms.  You may not use this
// file in any way except in compliance with the terms of this license.

// The code is distributed AS-IS in the hope that it will be useful, but WITHOUT ANY WARRANTY;
// without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE or
// TITLE or of non-infringement.  Thunder Token Inc. and any contributors to the software shall
// not be liable for any direct, indirect, incidental, special, punitive, exemplary, or
// consequential damages (including, without limitation, procurement of substitute goods or
// services, loss of use, data or profits or business interruption) however caused and under
// any theory of liability, whether in contract, strict liability, or tort (including negligence)
// or otherwise arising in any way out of the use of or inability to use the software, even if
// advised of the possibility of such damage.  The foregoing limitations of liability shall apply
// even if deemed to fail of their essential purpose.  The software may only be distributed under
// these terms and this disclaimer.

// This license does not grant permission to use the trade names, trademarks, service marks, or
// product names of ThunderCore™ or of Thunder Token Inc., except as required for reasonable and
// customary use in describing the origin of the work and reproducing the content of this file.

// Thunder Token Inc. and The ThunderCore™ Authors may publish revised and/or new versions of
// this TT license from time to time.

// You should have received a copy of the specific GNU license along with this file,
// the ThunderCore™ library, or the go-ethereum library.  If not, then see, e.g.,
// <https://www.gnu.org/licenses/lgpl-3.0.en.html> and/or <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// makeGenesis creates a new genesis struct based on some user input.
func (w *wizard) makeGenesis() {
	// Construct a default genesis block
	genesis := &core.Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		GasLimit:   4700000,
		Difficulty: big.NewInt(524288),
		Alloc:      make(core.GenesisAlloc),
		Config: &params.ChainConfig{
			HomesteadBlock: big.NewInt(1),
			EIP150Block:    big.NewInt(2),
			EIP155Block:    big.NewInt(3),
			EIP158Block:    big.NewInt(3),
			ByzantiumBlock: big.NewInt(4),
		},
	}
	// Figure out which consensus engine to choose
	fmt.Println()
	fmt.Println("Which consensus engine to use? (default = thunder)")
	fmt.Println(" 1. Ethash - proof-of-work")
	fmt.Println(" 2. Clique - proof-of-authority")
	fmt.Println(" 3. Thunder - thundercore")

	choice := w.read()
	switch {
	case choice == "1":
		// In case of ethash, we're pretty much done
		genesis.Config.Ethash = new(params.EthashConfig)
		genesis.ExtraData = make([]byte, 32)

	case choice == "2":
		// In the case of clique, configure the consensus parameters
		genesis.Difficulty = big.NewInt(1)
		genesis.Config.Clique = &params.CliqueConfig{
			Period: 15,
			Epoch:  30000,
		}
		fmt.Println()
		fmt.Println("How many seconds should blocks take? (default = 15)")
		genesis.Config.Clique.Period = uint64(w.readDefaultInt(15))

		// We also need the initial list of signers
		fmt.Println()
		fmt.Println("Which accounts are allowed to seal? (mandatory at least one)")

		var signers []common.Address
		for {
			if address := w.readAddress(); address != nil {
				signers = append(signers, *address)
				continue
			}
			if len(signers) > 0 {
				break
			}
		}
		// Sort the signers and embed into the extra-data section
		for i := 0; i < len(signers); i++ {
			for j := i + 1; j < len(signers); j++ {
				if bytes.Compare(signers[i][:], signers[j][:]) > 0 {
					signers[i], signers[j] = signers[j], signers[i]
				}
			}
		}
		genesis.ExtraData = make([]byte, 32+len(signers)*common.AddressLength+65)
		for i, signer := range signers {
			copy(genesis.ExtraData[32+i*common.AddressLength:], signer[:])
		}

	case choice == "" || choice == "3":
		genesis.Difficulty = big.NewInt(1)
		genesis.GasLimit = 10000000
		genesis.Config.Thunder = &params.ThunderConfig{}
	default:
		log.Crit("Invalid consensus engine choice", "choice", choice)
	}
	// Consensus all set, just ask for initial funds and go
	fmt.Println()
	fmt.Println("Which accounts should be pre-funded? (advisable at least one)")
	for {
		// Read the address of the account to fund
		if address := w.readAddress(); address != nil {
			genesis.Alloc[*address] = core.GenesisAccount{
				Balance: new(big.Int).Lsh(big.NewInt(1), 256-7), // 2^256 / 128 (allow many pre-funds without balance overflows)
			}
			continue
		}
		break
	}
	// Add a batch of precompile balances to avoid them getting deleted
	for i := int64(0); i < 256; i++ {
		genesis.Alloc[common.BigToAddress(big.NewInt(i))] = core.GenesisAccount{Balance: big.NewInt(1)}
	}
	// Query the user for some custom extras
	fmt.Println()
	fmt.Println("Specify your chain/network ID if you want an explicit one (default = random)")
	genesis.Config.ChainID = new(big.Int).SetUint64(uint64(w.readDefaultInt(rand.Intn(65536))))

	// All done, store the genesis and flush to disk
	log.Info("Configured new genesis block")

	w.conf.Genesis = genesis
	w.conf.flush()
}

// manageGenesis permits the modification of chain configuration parameters in
// a genesis config and the export of the entire genesis spec.
func (w *wizard) manageGenesis() {
	// Figure out whether to modify or export the genesis
	fmt.Println()
	fmt.Println(" 1. Modify existing fork rules")
	fmt.Println(" 2. Export genesis configuration")
	fmt.Println(" 3. Remove genesis configuration")

	choice := w.read()
	switch {
	case choice == "1":
		// Fork rule updating requested, iterate over each fork
		fmt.Println()
		fmt.Printf("Which block should Homestead come into effect? (default = %v)\n", w.conf.Genesis.Config.HomesteadBlock)
		w.conf.Genesis.Config.HomesteadBlock = w.readDefaultBigInt(w.conf.Genesis.Config.HomesteadBlock)

		fmt.Println()
		fmt.Printf("Which block should EIP150 come into effect? (default = %v)\n", w.conf.Genesis.Config.EIP150Block)
		w.conf.Genesis.Config.EIP150Block = w.readDefaultBigInt(w.conf.Genesis.Config.EIP150Block)

		fmt.Println()
		fmt.Printf("Which block should EIP155 come into effect? (default = %v)\n", w.conf.Genesis.Config.EIP155Block)
		w.conf.Genesis.Config.EIP155Block = w.readDefaultBigInt(w.conf.Genesis.Config.EIP155Block)

		fmt.Println()
		fmt.Printf("Which block should EIP158 come into effect? (default = %v)\n", w.conf.Genesis.Config.EIP158Block)
		w.conf.Genesis.Config.EIP158Block = w.readDefaultBigInt(w.conf.Genesis.Config.EIP158Block)

		fmt.Println()
		fmt.Printf("Which block should Byzantium come into effect? (default = %v)\n", w.conf.Genesis.Config.ByzantiumBlock)
		w.conf.Genesis.Config.ByzantiumBlock = w.readDefaultBigInt(w.conf.Genesis.Config.ByzantiumBlock)

		out, _ := json.MarshalIndent(w.conf.Genesis.Config, "", "  ")
		fmt.Printf("Chain configuration updated:\n\n%s\n", out)

	case choice == "2":
		// Save whatever genesis configuration we currently have
		fmt.Println()
		fmt.Printf("Which file to save the genesis into? (default = %s.json)\n", w.network)
		out, _ := json.MarshalIndent(w.conf.Genesis, "", "  ")
		if err := ioutil.WriteFile(w.readDefaultString(fmt.Sprintf("%s.json", w.network)), out, 0644); err != nil {
			log.Error("Failed to save genesis file", "err", err)
		}
		log.Info("Exported existing genesis block")

	case choice == "3":
		// Make sure we don't have any services running
		if len(w.conf.servers()) > 0 {
			log.Error("Genesis reset requires all services and servers torn down")
			return
		}
		log.Info("Genesis block destroyed")

		w.conf.Genesis = nil
		w.conf.flush()

	default:
		log.Error("That's not something I can do")
	}
}
