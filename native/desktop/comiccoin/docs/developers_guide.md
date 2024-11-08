# Developers Guide to Getting Started

## Prerequisites

1. Before we start, please confirm you are in the correct directory:

   ```shell
   cd $GOPATH/src/github.com/LuchaComics/monorepo/native/desktop/comiccoin
   ```

2. Also please confirm you have activated the `golang` workspace for this project:

   ```shell
   go work use ./native/desktop/comiccoin
   ```

3. **IMPORTANT:** In this tutorial, whenever we open a new terminal window, please make sure *Prerequisites (1) and (2)* have been done.

## Part 1:  Setup the ComicCoin Proof of Authority & Genesis Block
### (A) Initializing the Proof of Authority peer service identity.

1. Run this command to setup the private-public keys necessary for server identity.

   ```shell
   go run main.go init --datadir=./data/peer1/ComicCoin --id=peer1
   ```

2. When you you finish successfully you should see something like this (but not exact values as the values will be different every time your run the command):

   ```text
   msg="Blockchain node intitialized and ready" "peer identity"=QmcmBzpxfCe7dR4GhegYwP8JKgY1rEtYqkjMeQTFfcue6X "full address"=/ip4/127.0.0.1/tcp/26642/p2p/QmcmBzpxfCe7dR4GhegYwP8JKgY1rEtYqkjMeQTFfcue6X
   ```

3. From the output we must do the following:

   * `QmcmBzpxfCe7dR4GhegYwP8JKgY1rEtYqkjMeQTFfcue6X` is a p2p address we use in our network to indicate that we are the **proof of authority**.
   * In your `comiccoin/config/constants/constants.go` folder, please update `ComicCoinBootstrapPeers` value to be set this value.


4. For convenience in these tutorials, save the output of this into an environment variable.

   ```shell
   export COMICCOIN_POF_ADDRESS=QmcmBzpxfCe7dR4GhegYwP8JKgY1rEtYqkjMeQTFfcue6X
   export COMICCOIN_BOOTSTRAP_PEERS=/ip4/127.0.0.1/tcp/26642/p2p/QmcmBzpxfCe7dR4GhegYwP8JKgY1rEtYqkjMeQTFfcue6X
   ```

### (B) Initializing the ComicCoin Blockchain.
1. Generate a random 64 character long password which we will use as the password for our `coinbase` account. **Please store this password somewhere safe and do not lose it.**

   ```shell
   openssl rand -hex 64
   ```

2. For convenience, save the output of this into an environment variable. Please replace `...` with the output from `openssl`.

    ```shell
    export COMICCOIN_COINBASE_PASSWORD=...
    ```

2. Run this command to initialize the ComicCoin blockchain for the very first time. Command will create the genesis block and the `coinbase` account.

   ```shell
   go run main.go blockchain init \
   --datadir=./data/peer1/ComicCoin \
   --coinbase-password=$COMICCOIN_COINBASE_PASSWORD \
   --coinbase-password-repeated=$COMICCOIN_COINBASE_PASSWORD;
   ```

3. You have successfully initialized the blockchain! For our example, let's say the `coinbase` account was created with the following address `0x704185d9511b11bb03fd492ae306423343377a09`. For convenience, save the output of this into an environment variable.

    ```shell
    export COMICCOIN_COINBASE_ADDRESS=0x704185d9511b11bb03fd492ae306423343377a09;
    ```

## Part 2: Start a Peer-to-Peer Network

1. Generate the second and third server's identity to simulate a distributed network.

   ```shell
   go run main.go init --datadir=./data/peer2/ComicCoin --id=peer2
   go run main.go init --datadir=./data/peer3/ComicCoin --id=peer3
   ```

2. Run this command to run the ComicCoin blockchain service (in P2P `host mode`) with the `ComicCoin Central Validation Authority` services running the background.

   ```shell
   go run main.go daemon \
   --datadir=./data/peer1/ComicCoin \
   --identitykey-id=peer1 \
   --listen-http-address=127.0.0.1:8000 \
   --listen-p2p-port=26642 \
   --consensus-protocol=PoA \
   --enable-miner=true \
   --poa-address=$COMICCOIN_COINBASE_ADDRESS \
   --poa-password=$COMICCOIN_COINBASE_PASSWORD;
   ```

3. Open a new `terminal` window, let us start a new node (in P2P `dial mode`):

   ```shell
   go run main.go daemon \
   --datadir=./data/peer2/ComicCoin \
   --identitykey-id=peer2 \
   --listen-http-address=127.0.0.1:8002 \
   --listen-p2p-port=26643 \
   --consensus-protocol=PoA \
   --bootstrap-peers=$COMICCOIN_BOOTSTRAP_PEERS;
   ```

   Developer Notes:
   * Notice the `--listen-p2p-port=26643` argument for our second node, while our first node had `--listen-p2p-port=26642`? The reason for this is because we are running two nodes on the same address, so we need to have unique port numbers for every node instance running.
   * Notice the same for `--listen-http-port` values being different.
   * Finally the argument `--bootstrap-peers` was set once the first node ran.

4. Finally open another new `terminal` window, and start our third node:

   ```shell
   go run main.go daemon \
   --datadir=./data/peer3/ComicCoin \
   --identitykey-id=peer3 \
   --listen-http-address=127.0.0.1:8004 \
   --listen-p2p-port=26644 \
   --consensus-protocol=PoA \
   --bootstrap-peers=$COMICCOIN_BOOTSTRAP_PEERS;
   ```

5. Congratulations, you have now setup a peer-to-peer network locally.

## Part 3: Transfer Coins

1. Open a new `terminal` window and to begin by setting the password for our new wallet with the label called `alice`. Also please change the value `...` with something else secure:

   ```shell
   export COMICCOIN_ALICE_WALLET_PASSWORD=...
   ```

and run the following command to create our `alice` wallet. Please replace the password value of `...` with your own password.

   ```shell
   go run main.go blockchain account new \
   --wallet-password=$COMICCOIN_ALICE_WALLET_PASSWORD \
   --wallet-password-repeated=$COMICCOIN_ALICE_WALLET_PASSWORD \
   --wallet-label=alice;
   ```

2. For our example, let's say outputted address was: `0xfaea8212cb24dd69b1cc60b40febe4e7015939fb`. So for convenience in these tutorials, save the output of this into an environment variable.

   ```shell
   export COMICCOIN_ALICE_WALLET_ADDRESS=0xfaea8212cb24dd69b1cc60b40febe4e7015939fb
   export COMICCOIN_ALICE_WALLET_PASSWORD=...
   ```

3. Transfer from `coinbase` wallet to the `alice` wallet.

    ```shell
    go run main.go blockchain coin transfer \
    --sender-account-address=$COMICCOIN_COINBASE_ADDRESS \
    --sender-account-password=$COMICCOIN_COINBASE_PASSWORD \
    --value=1 \
    --recipient-address=$COMICCOIN_ALICE_WALLET_ADDRESS;
    ```

4. If you get no errors then you have successfully transfered coins between accounts. In one of the terminals you should see the *miner* execute the request. Now run this command to confirm our transfer occured. You should see somesort of output.

   ```shell
   go run main.go blockchain account get --account-address=$COMICCOIN_ALICE_WALLET_ADDRESS;
   ```

## Part 4: Creating Non-Fungible Tokens (NFT)
### (A) Create the NFT metadata file
Every NFT is essentially a map of `token_id` values to a `metadata_uri` values. To begin we will need to create a **metdata file** and post to our ComicCoin NFT Assets store so we can get a `cid` value that we'll set for the **metadata file**; afterwords, we will need to mint a token with that newly generated `metadata_uri`.

1. Go to the [`comiccoin-nftassetstore`](../comiccoin-nftstore) repository and setup the project.

2. Go to the [`comiccoin-registry`](../comiccoin-registry) repository and setup the project.

TODO...
