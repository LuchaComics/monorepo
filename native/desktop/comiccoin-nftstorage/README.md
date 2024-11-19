# 📚🪙 ComicCoin - NFT & Assets Storage Service

**Project still under active development - use at your own risk**

The purpose of this server is to exclusively host NFT assets that belong to the ComicCoin blockchain network and share them via IPFS Gateway implementation.

## ⭐️ Features

* IPFS Gateway that exclusively hosts NFT assets that belong on the ComicCoin blockchain network.

* Authenticated `pinadd` API endpoint to allow upload access to only the ComicCoin authority.

* Connects with either a local or remote [IPFS node](https://ipfs.tech).

## 👐 Installation

See [Installation Instructions](./docs/installation.md) for more information on building **NFT Asset Store** and working with the source code.

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file.

`COMICCOIN_NFTSTORAGE_APP_HMAC_SECRET` - This is the applications HMAC secret after running the following in your console: `openssl rand -hex 64`.

`COMICCOIN_NFTSTORAGE_APP_SECRET_KEY` - This is the application secret that you generated after running `go run main.go apikey` in your console.

`COMICCOIN_NFTSTORAGE_API_KEY` - This is the API key that you generated after running `go run main.go apikey` in your console.

## CLI Usage

To start the NFT Asset Store server:

   ```shell
   go run main.go daemon;
   ```

To add and pin a digital asset:

   ```shell
   go run main.go pinadd --filepath=./README.md
   ```

To get a digital asset:

   ```shell
   go run main.go cat --cid=bafkreiew7pqyqoryi7ynwmtwv3rhilgr6hjc6hl364u7glrhrhiaya5poy
   ```

## HTTP JSON Usage

1. See the file [pinadd.go](./cmd/pinadd.go) on how to upload a digital asset from your Golang code and pin it on the IPFS network.

2. see the file [get.go](./cmd/get.go) on how to get the content of a digital asset from the IPFS network.

## 📕 Documentation

See the [**Documentation**](./docs) for more information.

## 🛠️ Building

See **Build Instructions (TODO)** for more information on building **NFT Asset Store** and working with the source code.

## 🤝 Contributing

Found a bug? Want a feature to improve the package? Please create an [issue](https://github.com/LuchaComics/monorepo/issues/new).

## 📝 License

This application is licensed under the [**GNU Affero General Public License v3.0**](https://opensource.org/license/agpl-v3). See [LICENSE](LICENSE) for more information.
