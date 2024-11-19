# 📚🪙 ComicCoin - Registry

**Project still under active development - use at your own risk**

The purpose of this GUI application is provide the ComicCoin authority the ability to create NFT metadata (along with the digital assets associated with the NFT) and submit it to the IPFS network to be shared; in addition, submit the NFT metadata/assets to [`comiccoin-nftstorage`](../comiccoin-nftstorage) server to provide exclusive hosting of the said assets.

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
    go work use ./native/desktop/comiccoin-nftstorage
    ```

4. Go into our `comiccoin-nftstorage` folder

    ```shell
    cd ./native/desktop/comiccoin-nftstorage
    ```

5. Install our dependencies.

   ```shell
   go mod tidy
   ```

6. Start the GUI application (in developer mode). Please note that you will need the **Remote Address** and **API Key** of the [NFT Asset Store servers](../comiccoin-nftstorage) you have concurrently running for this GUI application.

   ```shell
   wails dev
   ```

7. If the GUI application loads up then you have successfully started running the **ComicCoin Registry**. You may now populate the NFT metadata/assets of the ComicCoin blockchain network!

## 🛠️ Building

See **Build Instructions (TODO)** for more information on building **ComicCoin Registry** GUI application and working with the source code.

## 🤝 Contributing

Found a bug? Want a feature to improve the package? Please create an [issue](https://github.com/LuchaComics/monorepo/issues/new).

## 📝 License

This application is licensed under the [**GNU Affero General Public License v3.0**](https://opensource.org/license/agpl-v3). See [LICENSE](LICENSE) for more information.

## ⚙️ Tech Stack

**Client:** React, Bulma.css

**Server:** Golang, Wails, IPFS, leveldb
