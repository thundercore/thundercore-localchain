## Thunder Local Chain

For original Ethereum readme, please [check here](README_ETH.md)

## Building from source

Clone the repository to a directory on your local,

    # git clone https://github.com/thundercore/go-ethereum.git

Make sure you have Go compiler installed,

on mac,

    # brew install go

on Ubuntu,

    # sudo apt-get install -y build-essential golang

on Fedora,

    # sudo dnf install -y golang

Build the local chain binary,

    # make thunder

Build the whole suite,

    # make all

## Start local chain

### Create an account
Find a folder and run "thunderlocal account new " from shell and set a password when asked,

    # thunderlocal account new --datadir ./datadir
    Your new account is locked with a password. Please give a password. Do not forget this password.
    Passphrase:
    Repeat passphrase:
    Address: {0af454242c456d1fc25c1d74a56a00a816ec336b}

This command will create an account and account information is saved in ./datadir/keystore/UTC--XXX

### Create a genesis json file
We can prefund the account created above. Please note the "alloc" block. And name the file as, for example thunder.json

    {
        "config": {
            "chainId": 3606,
            "homesteadBlock": 0,
            "eip150Block": 0,
            "eip155Block": 0,
            "eip158Block": 0,
            "byzantiumBlock": 0,
            "thunder": {
            }
        },
        "nonce": "0x0",
        "gasLimit": "0x989680",
        "difficulty": "0x1",
        "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "coinbase": "0x0000000000000000000000000000000000000000",
        "alloc": {
            "0x0af454242c456d1fc25c1d74a56a00a816ec336b" : { "balance": "1000000000000000000000000000000000000000" }
        },
        "number": "0x0",
        "gasUsed": "0x0",
        "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
    }

Please note the "thunder" block in the json file, which will use thunder consensus engine when starting a local chain. Alternatively, we can also create such genesis json file by [using setup wizard, puppeth](#puppeth).

### <a name="initchain"></a>Initialize the chain

    # thunderlocal init thunder.json --datadir ./datadir

### Start the chain

    # thunderlocal --datadir ./datadir  --networkid 3606 --rpc --rpcport "8545" --rpccorsdomain "*" --rpcapi "db,eth,net,miner,web3,personal" --nodiscover  --mine

Now, you are good to go. Try to send some transactions or deploy smart contract.

### Interactive Javascript console

    # thunderlocal attach http://127.0.0.1:8545

### <a name="puppeth"></a>Use puppeth to initialize genesis json config

"puppeth" is built when you build the whole suite by running "make all". Start the wizard from command line,

    # puppeth

    Please specify a network name to administer (no spaces or hyphens, please)
    > thunder

    What would you like to do? (default = stats)
    1. Show network stats
    2. Configure new genesis
    3. Track new remote server
    4. Deploy network components
    > 2

    Which consensus engine to use? (default = thunder)
    1. Ethash - proof-of-work
    2. Clique - proof-of-authority
    3. Thunder - thundercore
    > 3

    Which accounts should be pre-funded? (advisable at least one)
    > 0x8ddf4f3e475f5b5f10ec5bf3452f94830e9d2ce9

    Specify your chain/network ID if you want an explicit one (default = random)
    > 3606

    What would you like to do? (default = stats)
    1. Show network stats
    2. Manage existing genesis
    3. Track new remote server
    4. Deploy network components
    > 2

    1. Modify existing fork rules
    2. Export genesis configuration
    3. Remove genesis configuration
    > 2

    Which file to save the genesis into? (default = ThunderLocal.json)
    > thunder.json

Then, you can initialize then start chain as described [above](#initchain)

## License

  Copyright 2018 Thunder Token Inc., The ThunderCore™ Authors
  This file comprises an original work of authorship that may make use of, or
  interface with another work licensed under a GNU or third party license, but
  which is not otherwise based on said another work.

  To the extent that portions of this file contains source code that is subject
  to the terms of the GNU or third party license, the minimal corresponding source
  code for those portions can be freely redistributed and/or modified under the
  terms of the respective license, either of GNU Lesser General Public License version 3
  or (at your option) any later version.

  The remaining code for the ThunderCore™ network application is not a contribution
  to be incorporated into said another work.  Rather, it is open source and licensed
  from Thunder Token Inc. to you, the recipient, to copy, modify and distribute the
  original or modified work without a fee, subject to reciprocity and recipient’s
  (i) promise and covenant not to sue Thunder Token Inc., its assigns, successors,
  affiliates and subsidiaries (hereinafter “Thunder Token”) on claims arising from
  any of their use of recipient’s code, if any; (ii) promise and ongoing commitment
  to not unfairly compete against or interfere with Thunder Token’s business or commercial
  relationships; and (iii) promise and ongoing commitment to not challenge the validity,
  enforceability, title, or ownership (by Thunder Token) of any intellectual property
  rights arising from or relating to the ThunderCore™ network application.  Further, you,
  the recipient, agree to and must do the following: (1) give prominent notice and
  attribution to Thunder Token Inc. and the ThunderCore™ Authors for their work on the
  original work and include any appropriate copyright, trademark, patent notices,
  (2) accompany the original or modified work with a copy of this notice (TT license v1.0
  or, at your option, any later version) in its entirety or a link directing the user to
  the same, (3) accompany the modified work with a prominent notice indicating that it
  has been modified and that it was based off of the original work; and (4) convey or
  otherwise make freely available the source code corresponding to the modified work
  under the same conditions and restrictions on the exercise of rights granted or
  affirmed under this license.

  Your copying, reverse-engineering, debugging, modifying, or distributing the original
  or modified work constitutes assent and agreement to these terms.  You may not use this
  file in any way except in compliance with the terms of this license.

  The code is distributed AS-IS in the hope that it will be useful, but WITHOUT ANY WARRANTY;
  without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE or
  TITLE or of non-infringement.  Thunder Token Inc. and any contributors to the software shall
  not be liable for any direct, indirect, incidental, special, punitive, exemplary, or
  consequential damages (including, without limitation, procurement of substitute goods or
  services, loss of use, data or profits or business interruption) however caused and under
  any theory of liability, whether in contract, strict liability, or tort (including negligence)
  or otherwise arising in any way out of the use of or inability to use the software, even if
  advised of the possibility of such damage.  The foregoing limitations of liability shall apply
  even if deemed to fail of their essential purpose.  The software may only be distributed under
  these terms and this disclaimer.

  This license does not grant permission to use the trade names, trademarks, service marks, or
  product names of ThunderCore™ or of Thunder Token Inc., except as required for reasonable and
  customary use in describing the origin of the work and reproducing the content of this file.

  Thunder Token Inc. and The ThunderCore™ Authors may publish revised and/or new versions of
  this TT license from time to time.

  You should have received a copy of the specific GNU license along with this file,
  the ThunderCore™ library, or the go-ethereum library.  If not, then see, e.g.,
  <https://www.gnu.org/licenses/lgpl-3.0.en.html> and/or <http://www.gnu.org/licenses/>.
