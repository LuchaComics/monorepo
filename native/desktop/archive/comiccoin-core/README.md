# 📚🪙 ComicCoin - Core

**Project still under active development - use at your own risk**

ComicCoin Core is a full ComicCoin client and builds the backbone of the network. It offers high levels of security, privacy, and stability. However it takes a lot of space and memory.

## ⭐️ Features

* Control over your coins - This wallet gives you full control over your comiccoins. This means no third party can freeze or lose your funds. You are however still responsible for securing and backing up your wallet.

* Complete transparency - The source code for this wallet is public and the build process is deterministic. This means any developer in the world can audit the code and make sure the final software isn't hiding any secrets.

* Vulnerable environment - This wallet can be loaded on computers which are vulnerable to malware. Securing your computer, using a strong passphrase, moving most of your funds to cold storage, or enabling two-factor authentication can make it harder to steal your comiccoins.

* No fees - Free to transfer your tokens and coins among individuals without incurring any fees.

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
    go work use ./native/desktop/comiccoin-nftassetstore
    ```

4. Go into our `comiccoin-nftassetstore` folder

    ```shell
    cd ./native/desktop/comiccoin-nftassetstore
    ```

5. Install our dependencies.

   ```shell
   go mod tidy
   ```

6. Start the GUI application (in developer mode). Please note that you will need the `comiccoin` full node running in **proof of authority** mode turned on in addition the `comiccoin-nftassetstore` CLI application running as well. To run in live development mode, run `wails dev` in the project directory. This will run a Vite development server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect to this in your browser, and you can call your Go code from devtools:

   ```shell
   wails dev
   ```

## 🛠️ Building

See **Build Instructions (TODO)** for more information on building **ComicCoin Registry** GUI application and working with the source code.

## 🤝 Contributing

Found a bug? Want a feature to improve the package? Please create an [issue](https://github.com/LuchaComics/monorepo/issues/new).

## 📝 License

This application is licensed under the [**GNU Affero General Public License v3.0**](https://opensource.org/license/agpl-v3). See [LICENSE](LICENSE) for more information.

## ⚙️ Tech Stack

**Client:** React, Bulma.css

**Server:** Golang, Wails, IPFS, leveldb
