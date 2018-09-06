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
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	// TODO: set/change this to a better limit.
	// Currently it is large since we don't know how this will play with fast path recovery.
	allowedFutureBlockTime = 365 * 24 * 3600 * time.Second
	// We are using fixed gas limit for test net. This limit is same as that in genesis.go
	blockGasLimit = uint64(10000000)

	blockInterval = 1 * time.Second
)

var (
	zeroUncleHash = types.EmptyUncleHash
	// Currently to zero
	zeroCoinbaseAddress = common.Address{}
	unityDifficulty     = big.NewInt(1)
	// Used by ethhash for DAO header. Clique uses it for it's own thing.
	// We don't need it in thunder.
	zeroExtraData = make([]byte, 0)
	// ethhash uses MixDigest for something related to PoW.
	// clique consensus doesn't use it and checks it to be zero.
	// So we likely don't need it at this point.
	zeroMixDigest = common.Hash{}
	// Nonce is zero since we are not using PoW algorithm.
	zeroNonce = [8]byte{0, 0, 0, 0, 0, 0, 0, 0}

	// Various error messages to mark blocks invalid. These should be private to
	// prevent engine specific errors from being referenced in the remainder of the
	// codebase, inherently breaking if the engine is swapped out. Please put common
	// error types into the consensus package.
	errUnknownBlock = errors.New("block number is nil")

	errSealOperationOnGenesisBlock = errors.New(
		"verifySeal/Seal operations on genesis block not permitted")

	// If child block's timestamp < parent block's timestamp
	errBackwardBlockTime = errors.New("block timestamp less than parent's timestamp")

	// Errors for unused fields not set to zero/empty

	errNonEmptyUncleHash = errors.New("non empty uncle hash")
	errNonEmptyAddress   = errors.New("non empty coinbase address")
	// Thunder PoS doesn't have difficulty as in PoW.
	errNonZeroDifficulty = errors.New("non zero difficulty")
	errNonEmptyExtra     = errors.New("non empty extra")
	errNonEmptyMixDigest = errors.New("non-zero mix digest")
	errNonZeroNonce      = errors.New("non-zero nonce")
)

// Thunder is the proof-of-stake consensus engine.
type Thunder struct {
	config *params.ThunderConfig // Consensus engine configuration parameters
}

// New creates a Thunder proof-of-stake consensus engine.
func New(config *params.ThunderConfig) *Thunder {
	return &Thunder{config: config}
}

//////////////////////////////////
// consensus.Engine implementation
//////////////////////////////////

// Author implements consensus.Engine.
func (thunder *Thunder) Author(header *types.Header) (common.Address, error) {
	return zeroCoinbaseAddress, nil
}

// CalcDifficulty implements consensus.Engine
// CalcDifficulty is used for difficulty adjustment in PoW algorithms. Not relevant in Thunder.
func (thunder *Thunder) CalcDifficulty(chain consensus.ChainReader, time uint64,
	parent *types.Header) *big.Int {
	return big.NewInt(0)

}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (thunder *Thunder) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header,
	seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		for i, header := range headers {
			err := thunder.VerifyHeader(chain, header, seals[i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

func verifyHeaderUnusedFieldsAreZero(header *types.Header) error {
	// Ensure that the block doesn't contain any uncles
	if header.UncleHash != zeroUncleHash {
		return errNonEmptyUncleHash
	}
	if header.Coinbase != zeroCoinbaseAddress {
		return errNonEmptyAddress
	}
	// Ensure that the block's difficulty is zero.
	if header.Difficulty.Cmp(unityDifficulty) != 0 {
		return errNonZeroDifficulty
	}
	if bytes.Compare(header.Extra, zeroExtraData) != 0 {
		return errNonEmptyExtra
	}
	// Ensure that the mix digest is zero. Provisioned for fork protection in Ethereum.
	if header.MixDigest != zeroMixDigest {
		return errNonEmptyMixDigest
	}
	if !bytes.Equal(header.Nonce[:], zeroNonce[:]) {
		return errNonZeroNonce
	}
	return nil
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (thunder *Thunder) VerifyHeader(chain consensus.ChainReader, header *types.Header,
	seal bool) error {
	if header.Number == nil {
		return errUnknownBlock
	}
	// If the block is already in local chain, no need to verify again.
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// Verify that the block number is parent's +1
	if diff := new(big.Int).Sub(header.Number, parent.Number); diff.Cmp(big.NewInt(1)) != 0 {
		return consensus.ErrInvalidNumber
	}
	// Don't waste time checking blocks from the future
	if header.Time.Cmp(big.NewInt(time.Now().Add(allowedFutureBlockTime).Unix())) > 0 {
		return consensus.ErrFutureBlock
	}
	if header.Time.Cmp(parent.Time) < 0 {
		return errBackwardBlockTime
	}
	if err := verifyHeaderUnusedFieldsAreZero(header); err != nil {
		return err
	}
	if header.GasLimit != blockGasLimit {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit,
			blockGasLimit)
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed,
			header.GasLimit)
	}
	// Verify the engine specific seal securing the block
	if seal {
		if err := thunder.VerifySeal(chain, header); err != nil {
			return err
		}
	}
	// Note: Removed validations regarding hard-forks.
	// We may want to put add a check here when we plan for forks.
	return nil
}

// VerifyUncles implements consensus.Engine, always returning an error for any
// uncles as this consensus mechanism doesn't permit uncles.
func (thunder *Thunder) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errors.New("uncles not allowed")
	}
	return nil
}

// VerifySeal implements consensus.Engine.
// In thunder protocol, we don't store signed proposals in the block.
func (thunder *Thunder) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errSealOperationOnGenesisBlock
	}
	return nil
}

// All header fields which are not relevant in Thunder protocol are set to predefined zero values.
func setHeaderUnusedFieldsToZero(header *types.Header) {
	header.UncleHash = zeroUncleHash
	header.Coinbase = zeroCoinbaseAddress
	header.Difficulty = unityDifficulty
	header.Extra = zeroExtraData
	header.MixDigest = zeroMixDigest
	header.Nonce = zeroNonce
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (thunder *Thunder) Prepare(chain consensus.ChainReader, header *types.Header) error {
	setHeaderUnusedFieldsToZero(header)
	number := header.Number.Uint64()

	// Ensure the timestamp has the correct delay
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	header.Time = new(big.Int).Add(parent.Time, new(big.Int).SetUint64(0))
	if header.Time.Int64() < time.Now().Unix() {
		header.Time = big.NewInt(time.Now().Unix())
	}
	header.GasLimit = blockGasLimit
	return nil
}

// Finalize implements consensus.Engine, ensuring no uncles are set, nor block
// rewards given, and returns the final block.
func (thunder *Thunder) Finalize(chain consensus.ChainReader, header *types.Header,
	state *state.StateDB, txs []*types.Transaction, uncles []*types.Header,
	receipts []*types.Receipt) (*types.Block, error) {
	// No block rewards in Thunder, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	setHeaderUnusedFieldsToZero(header)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts), nil
}

// Seal implements consensus.Engine.
func (thunder *Thunder) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block,
	stop <-chan struct{}) error {
	header := block.Header()

	// Sealing the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errSealOperationOnGenesisBlock
	}

	go func() {
		select {
		case <-stop:
			return
		case <-time.After(blockInterval):
		}

		select {
		case results <- block.WithSeal(header):
		default:
			log.Warn("Sealing result is not read by miner", "sealhash", thunder.SealHash(header))
		}
	}()
	return nil
}

// Close implements consensus.Engine. It's a noop for Thunder as there is are no background threads.
func (thunder *Thunder) Close() error {
	return nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (thunder *Thunder) SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
	})
	hasher.Sum(hash[:0])
	return hash
}

type API struct {
	chain   consensus.ChainReader
	thunder *Thunder
}

// GetBlockInterval retrieves block interval.
func (api *API) GetBlockInterval() float64 {
	return blockInterval.Seconds()
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the signer voting.
func (thunder *Thunder) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "thunder",
		Version:   "0.1",
		Service:   &API{chain: chain, thunder: thunder},
		Public:    false,
	}}
}
