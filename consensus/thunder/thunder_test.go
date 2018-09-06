// Copyright 2018 Thunder Token Inc., The ThunderCore™ Authors
// This file comprises an original work of authorship that may make use of, or
// interface with another work licensed under a GNU or third party license, but
// which is not otherwise based on said another work.

// To the extent that portions of this file contains source code that is subject
// to the terms of the GNU or third party license, the minimal corresponding source
// code for those portions can be freely redistributed and/or modified under the
// terms of the respective license, either of GNU Lesser General Public License version 3
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

package thunder

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/assert"
)

func TestAuthor(t *testing.T) {
	assert := assert.New(t)

	thunder := New(new(params.ThunderConfig))

	coinbase, err := thunder.Author(nil)

	assert.Equal(coinbase, zeroCoinbaseAddress, "error")
	assert.Equal(err, nil, "error")
}

func TestCalcDifficulty(t *testing.T) {
	assert := assert.New(t)

	var difficulty *big.Int

	thunder := New(new(params.ThunderConfig))

	difficulty = thunder.CalcDifficulty(nil, 0, nil)

	assert.Equal(difficulty.Int64(), int64(0), "error")
}

func makeThunderTestChain() *core.BlockChain {

	var (
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		db      = ethdb.NewMemDatabase()
	)

	gspec := &core.Genesis{
		Config:   params.TestThunderChainConfig,
		Alloc:    core.GenesisAlloc{addr1: {Balance: big.NewInt(1000000)}},
		GasLimit: blockGasLimit,
	}
	genesis := gspec.MustCommit(db)

	blockchain, _ := core.NewBlockChain(db, nil, gspec.Config, New(new(params.ThunderConfig)), vm.Config{})

	blocks := make(types.Blocks, 1)
	blocks[0] = genesis

	blockchain.InsertChain(blocks)

	return blockchain
}

func makeNewHeader(blockchain *core.BlockChain) *types.Header {

	currBlock := blockchain.CurrentBlock()
	header := &types.Header{
		Number:     new(big.Int).Add(currBlock.Number(), common.Big1),
		ParentHash: currBlock.Hash(),
	}

	return header
}

func TestThunderPrepare(t *testing.T) {

	assert := assert.New(t)

	blockchain := makeThunderTestChain()
	header := makeNewHeader(blockchain)

	thunder := New(new(params.ThunderConfig))
	err := thunder.Prepare(blockchain, header)
	assert.Equal(err, nil)
	assert.Equal(header.UncleHash, zeroUncleHash)
	assert.Equal(header.Coinbase, zeroCoinbaseAddress)
	assert.Equal(header.Difficulty, unityDifficulty)
	assert.Equal(header.Extra, zeroExtraData)
	assert.Equal(header.MixDigest, zeroMixDigest)
	assert.Equal(header.Nonce, types.BlockNonce{0, 0, 0, 0, 0, 0, 0, 0})
	assert.Equal(header.GasLimit, blockGasLimit)
}

func TestThunderFinalize(t *testing.T) {

	assert := assert.New(t)

	blockchain := makeThunderTestChain()
	header := makeNewHeader(blockchain)

	thunder := New(new(params.ThunderConfig))
	db := ethdb.NewMemDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db))

	block, err := thunder.Finalize(blockchain, header, state, nil, nil, nil)
	assert.Equal(err, nil)
	assert.Equal(block.UncleHash(), zeroUncleHash)
	assert.Equal(block.Coinbase(), zeroCoinbaseAddress)
	assert.Equal(block.Difficulty(), unityDifficulty)
	assert.Equal(block.Extra(), zeroExtraData)
	assert.Equal(block.MixDigest(), zeroMixDigest)
	assert.Equal(block.Nonce(), uint64(0))

}

func TestVerifyHeader(t *testing.T) {

	assert := assert.New(t)

	blockchain := makeThunderTestChain()
	header := makeNewHeader(blockchain)
	number := header.Number

	thunder := New(new(params.ThunderConfig))
	thunder.Prepare(blockchain, header)

	err := thunder.VerifyHeader(blockchain, header, false)
	assert.Equal(err, nil)

	err = thunder.VerifyHeader(blockchain, header, true)
	assert.Equal(err, nil)

	header.Number = nil
	assert.Equal(thunder.VerifyHeader(blockchain, header, false), errUnknownBlock)
	header.Number = number

	header.Number = big.NewInt(100)
	assert.Equal(thunder.VerifyHeader(blockchain, header, false), consensus.ErrInvalidNumber)
	header.Number = number

	header.GasLimit = blockGasLimit + 1
	assert.Errorf(thunder.VerifyHeader(blockchain, header, false),
		fmt.Sprintf("invalid gasLimit: have %v, max %v", header.GasLimit, blockGasLimit))
	header.GasLimit = blockGasLimit

	header.GasUsed = blockGasLimit + 1
	assert.Errorf(thunder.VerifyHeader(blockchain, header, false),
		fmt.Sprintf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit))
}

func TestSeal(t *testing.T) {

	assert := assert.New(t)

	resultsCh := make(chan *types.Block)

	blockchain := makeThunderTestChain()
	header := makeNewHeader(blockchain)

	thunder := New(new(params.ThunderConfig))
	thunder.Prepare(blockchain, header)

	block := types.NewBlock(header, nil, nil, nil)

	start := time.Now()
	err := thunder.Seal(blockchain, block, resultsCh, make(chan struct{}))
	assert.Equal(err, nil)

	<-resultsCh

	// Block interval is 1 second
	elapsed := time.Since(start)
	assert.True(elapsed.Seconds() >= blockInterval.Seconds())
}
