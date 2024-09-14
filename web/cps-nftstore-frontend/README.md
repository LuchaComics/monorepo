# Collectible Protection Services NFT Store (Frontend)

**Project still under active development - use at your own risk**

The purpose of this project is to provide NFT collection management and NFT minting **frontend** for the backend. The backend project is found in the [`cps-nftstore-backend`](../../cloud/cps-nftstore-backend) repository.

## Features

* Manage NFT Collections

* Manage NFTs

* Mint NFTs via **"Collectible Protection Service Submissions"** smart contract which dynamic and has no tokens limit

* Automatically handles file uploads to IPFS to get `cid` values for pinning.

* Automatically handles `IPNS` management for each NFT collection used.

## Installation

1. Download and install the following onto your computer before proceeding:
  * [git](https://git-scm.com/downloads)
  * [Taskfile](https://taskfile.dev/installation/)
  * [Node](https://nodejs.org/en/download/package-manager)

2. Go to your `~/go/src/github.com` folder and clone this monorepo (if you haven't done that already).

    ```shell
    cd ~/go/src/github.com
    mkdir LuchaComics
    cd LuchaComics
    git clone git@github.com:LuchaComics/monorepo.git
    cd ./monorepo/web/cps-nftstore-frontend
    ```

3. Install our projects dependencies.

    ```shell
    npm install
    ```

4. Start the frontend server by running the following in your console.

    ```shell
    task start
    ```

5. Congratulations, the frontend is running, now please go to the [`cps-nftstore-backend`](../../cloud/cps-nftstore-backend) repository and startup the backend if you haven't already started it. You may now access the frontend via [http://127.0.0.1:8000](http://127.0.0.1:8000) in your favourite browser.

## Developers Notes:
* When you run `task start` then this loads up the development environment variables.
* If you don't want to use `Taskfile` then use the following alternative `npm run start:dev`.
* You can run `npm run start:qa"` to load up the quality assurance environment variables.
* You can run `npm run start:prod"` to load up the production environment variables.

## Documentation

See the [**Documentation**](./docs) folder.

## Contributing

Found a bug? Want a feature to improve the package? Please create an [issue](https://github.com/LuchaComics/monorepo/issues/new).

## License

This application is licensed under the [**GNU Affero General Public License v3.0**](https://opensource.org/license/agpl-v3). See [LICENSE](LICENSE) for more information.
