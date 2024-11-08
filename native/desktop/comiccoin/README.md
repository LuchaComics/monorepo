# 📚🪙 ComicCoin

**Project still under active development - use at your own risk**

ComicCoin is proof-of-authority based blockchain which allow users to transfer **coins** and **non-fungible tokens** among themselves.

The **comiccoin** repository contains code to run a fullnode on the ComicCoin blockchain network; in addition, allow you to import the package library into your application to build ComicCoin blockchain powered applications.

## ⭐️ Features

* TODO

## 👐 Installation

Follow these steps to setup the project locally on your development machine.

1. Go to your `$GOPATH` directory and clone our *monorepo*.

   ```shell
   cd $GOPATH/src/github.com
   mkdir LuchaComics
   cd ./LuchaComics
   git clone git@github.com:LuchaComics/monorepo.git
   ```

2. Go into our monorepo folder.

   ```shell
   cd ./LuchaComics/monorepo
   ```

3. Activate the golang workspace which is required.

    ```shell
    go work use ./native/desktop/comiccoin
    ```

4. Go into our comiccoin folder

    ```shell
    cd ./native/desktop/comiccoin
    ```

5. Install our dependencies.

   ```shell
   go mod tidy
   ```

6. Run `comiccoin` as a full node on the ComicCoin blockchain network.

   ```shell
   go run main.go daemon --datadir=./data/ComicCoin
   ```

   Notes:
   * `daemon` is the command used to run the full node.
   * `--datadir=./data/ComicCoin` means we want the node to save the entire blockchain in this location.
   * To learn more, please read the [**Documentation**](./docs).

7. Congrats, the node is running and connected as a peer to peer.


## 📕 Documentation

See the [**Documentation**](./docs) for more information.

## 🛠️ Building

See [Build Instructions](./docs/BUILD.md) for more information on building **CPS NFT Store (backend)** and working with the source code.

## 🤝 Contributing

Found a bug? Want a feature to improve the package? Please create an [issue](https://github.com/LuchaComics/monorepo/issues/new).

## 📝 License

This application is licensed under the [**GNU Affero General Public License v3.0**](https://opensource.org/license/agpl-v3). See [LICENSE](LICENSE) for more information.
