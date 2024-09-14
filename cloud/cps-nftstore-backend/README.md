# Collectible Protection Services NFT Store

**Project still under active development - use at your own risk**

The purpose of this project is to provide NFT minting for the **Collectible Protection Services Submission Token** NFT.

## Features

* Manage NFT Collections

* Manage NFTs

* Mint NFTs via **"Collectible Protection Service Submissions"** smart contract which dynamic and has no tokens limit

* Automatically handles file uploads to IPFS to get `cid` values for pinning.

* Automatically handles `IPNS` management for each NFT collection used.

## Installation

1. Download and install the following onto your computer before proceeding:
  * [git](https://git-scm.com/downloads)
  * [Golang](https://go.dev/dl/)
  * [Taskfile](https://taskfile.dev/installation/)
  * [Node](https://nodejs.org/en/download/package-manager)
  * [Docker](https://www.docker.com/products/docker-desktop/)

2. Go to your `~/go/src/github.com` folder and clone this monorepo.

    ```shell
    cd ~/go/src/github.com
    mkdir LuchaComics
    cd LuchaComics
    git clone git@github.com:LuchaComics/monorepo.git
    ```

3. Duplicate the provided sample environment variables of the project.

    ```shell
    cp .env.sample .env
    ```

4. Please open your `.env` file and modify it to your specifications.

5. Start the backend server by running the following in your console.

    ```shell
    docker-compose -p cps_nftstore_backend -f docker-compose.yml up
    ```

6. Congratulations, the backend is running, now please go to the [`cps-nftstore-frontend`](../../web/cps-nftstore-frontend) repository and startup the frontend to access this backend via the web-browser.

## Documentation

See the [**Documentation**](./docs) folder.

## Contributing

Found a bug? Want a feature to improve the package? Please create an [issue](https://github.com/LuchaComics/monorepo/issues/new).

## License

This application is licensed under the [**GNU Affero General Public License v3.0**](https://opensource.org/license/agpl-v3). See [LICENSE](LICENSE) for more information.
