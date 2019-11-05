// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package state

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

type DumpAccount struct {
	Balance  string            `json:"balance"`
	Nonce    uint64            `json:"nonce"`
	Root     string            `json:"root"`
	CodeHash string            `json:"codeHash"`
	Code     string            `json:"code"`
	Storage  map[string]string `json:"storage"`
}

type Dump struct {
	Root     string                 `json:"root"`
	Accounts map[string]DumpAccount `json:"accounts"`
}

func (self *StateDB) RawDump() Dump {
	dump := Dump{
		Root:     fmt.Sprintf("%x", self.trie.Hash()),
		Accounts: make(map[string]DumpAccount),
	}

	it := trie.NewIterator(self.trie.NodeIterator(nil))
	for it.Next() {
		addr := self.trie.GetKey(it.Key)
		var data Account
		if err := rlp.DecodeBytes(it.Value, &data); err != nil {
			panic(err)
		}

		obj := newObject(nil, common.BytesToAddress(addr), data)
		account := DumpAccount{
			Balance:  data.Balance.String(),
			Nonce:    data.Nonce,
			Root:     common.Bytes2Hex(data.Root[:]),
			CodeHash: common.Bytes2Hex(data.CodeHash),
			Code:     common.Bytes2Hex(obj.Code(self.db)),
			Storage:  make(map[string]string),
		}
		storageIt := trie.NewIterator(obj.getTrie(self.db).NodeIterator(nil))
		for storageIt.Next() {
			account.Storage[common.Bytes2Hex(self.trie.GetKey(storageIt.Key))] = common.Bytes2Hex(storageIt.Value)
		}
		dump.Accounts[common.Bytes2Hex(addr)] = account
	}
	return dump
}

func (self *StateDB) Dump() []byte {
	json, err := json.MarshalIndent(self.RawDump(), "", "    ")
	if err != nil {
		fmt.Println("dump err", err)
	}

	return json
}

var ( // Thunder Pre-Compiled Contracts
	commElectionTPCHash = sha256.Sum256([]byte("Thunder_CommitteeElection"))
	// CommElectionTPCAddress is 0x30d87bd4D1769437880c64A543bB649a693EB348
	CommElectionTPCAddress = common.BytesToAddress(commElectionTPCHash[:20])

	vaultTPCHash = sha256.Sum256([]byte("Thunder_Vault"))
	// VaultTPCAddress is 0xEC45c94322EaFEEB2Cf441Cd1aB9e81E58901a08
	VaultTPCAddress = common.BytesToAddress(vaultTPCHash[:20])

	randomTPCHash = sha256.Sum256([]byte("Thunder_Random"))
	// RandomTPCAddress is 0x8cC9C2e145d3AA946502964B1B69CE3cD066A9C7
	RandomTPCAddress = common.BytesToAddress(randomTPCHash[:20])

	// from thunder/tests/utils/accounts.py
	genesisAddress    = common.HexToAddress("0x9A78d67096bA0c7C1bCdc0a8742649Bc399119c0")
	txStressAddress   = common.HexToAddress("0x4Bc87B58CfD96A4627a76C3dA5A8A26486eE7Fc9")
	genesisWebAddress = common.HexToAddress("0xbb8718BE30d331A9D98E74C0FE92391dc2b437c3")
	srcAddress        = common.HexToAddress("0xCD1191CAe116bDCBB24657c15C10aDfdb506aD85")
	destAddress       = common.HexToAddress("0xb58972114Bf1624165ed8dF5eF755F8927dA4730")
	txnFeeAddress     = common.HexToAddress("0xc4F3c85Bb93F33A485344959CF03002B63D7c4E3")
	foundationAddress = common.HexToAddress("0x0000000000000000000000000000001234567989")
	zeroAddress       = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

func shouldIgnoreAddress(addr []byte) bool {
	if IsPrecompiledContract(addr) {
		return true
	}
	if bytes.Equal(addr, genesisAddress[:]) ||
		bytes.Equal(addr, txStressAddress[:]) ||
		bytes.Equal(addr, genesisWebAddress[:]) ||
		bytes.Equal(addr, srcAddress[:]) ||
		bytes.Equal(addr, destAddress[:]) ||
		bytes.Equal(addr, txnFeeAddress[:]) ||
		bytes.Equal(addr, foundationAddress[:]) ||
		bytes.Equal(addr, zeroAddress[:]) {
		return true
	}
	return false
}

func IsPrecompiledContract(addr []byte) bool {
	if bytes.Equal(addr, CommElectionTPCAddress[:]) ||
		bytes.Equal(addr, VaultTPCAddress[:]) ||
		bytes.Equal(addr, RandomTPCAddress[:]) {
		return true
	}
	return false
}

func (self *StateDB) RawDumpNonContracts() Dump {
	dump := Dump{
		Root:     fmt.Sprintf("%x", self.trie.Hash()),
		Accounts: make(map[string]DumpAccount),
	}

	it := trie.NewIterator(self.trie.NodeIterator(nil))
	for it.Next() {
		addr := self.trie.GetKey(it.Key)

		if shouldIgnoreAddress(addr) {
			continue
		}

		var data Account
		if err := rlp.DecodeBytes(it.Value, &data); err != nil {
			panic(err)
		}

		if !bytes.Equal(data.CodeHash, emptyCodeHash) {
			continue
		}

		obj := newObject(nil, common.BytesToAddress(addr), data)
		account := DumpAccount{
			Balance:  data.Balance.String(),
			Nonce:    data.Nonce,
			Root:     common.Bytes2Hex(data.Root[:]),
			CodeHash: common.Bytes2Hex(data.CodeHash),
			Code:     common.Bytes2Hex(obj.Code(self.db)),
			Storage:  make(map[string]string),
		}
		storageIt := trie.NewIterator(obj.getTrie(self.db).NodeIterator(nil))
		for storageIt.Next() {
			account.Storage[common.Bytes2Hex(self.trie.GetKey(storageIt.Key))] = common.Bytes2Hex(storageIt.Value)
		}
		dump.Accounts[common.Bytes2Hex(addr)] = account
	}
	return dump
}

func (self *StateDB) DumpNonContracts() []byte {
	json, err := json.MarshalIndent(self.RawDumpNonContracts(), "", "    ")
	if err != nil {
		fmt.Println("dump err", err)
	}

	return json
}
