# ComicCoin NFT Asset Store - 👐 Installation

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

6. Run the following to generate a secure HMAC code.

    ```shell
    openssl rand -hex 64
    ```

7. For convenience take the output from the `openssl` command and create an environment variable by replacing the `...` value.

    ```shell
    export COMICCOIN_NFTSTORAGE_APP_HMAC_SECRET='...';
    ```

8. Generate an API key.

    ```shell
    go run main.go genapikey --hmac-secret=$COMICCOIN_NFTSTORAGE_APP_HMAC_SECRET
    ```

9. Great! Now take the console outputs and create the environment variables by replacing the `...` values below:

    ```shell
    export COMICCOIN_NFTSTORAGE_APP_SECRET_KEY='...';
    export COMICCOIN_NFTSTORAGE_API_KEY='...';
    ```

10. Before we start the application, we will need to run concurrently in the background of your operating system the application called [IPFS](https://ipfs.tech). Please visit their website and start an IPFS node instance on your machine. Once running proceed to next step 11 and do not proceed if you cannot setup an IPFS node locally on your machine.

10. Start the app.

   ```shell
   go run main.go daemon --listen-http-address=127.0.0.1:8080;
   ```
