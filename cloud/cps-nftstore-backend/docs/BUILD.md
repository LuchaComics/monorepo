# Building Instructions

1. Download and install the following onto your computer before proceeding:
  * [git](https://git-scm.com/downloads)
  * [Golang](https://go.dev/dl/)
  * [Taskfile](https://taskfile.dev/installation/)
  * [Node](https://nodejs.org/en/download/package-manager)
  * [Docker](https://www.docker.com/products/docker-desktop/)

2. Go to your `~/go/src/github.com` folder and clone this monorepo. Then go to the backend project.

    ```shell
    cd ~/go/src/github.com
    mkdir LuchaComics
    cd LuchaComics
    git clone git@github.com:LuchaComics/monorepo.git
    cd ./monorepo/cloud/cps-nftstore-backend
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
