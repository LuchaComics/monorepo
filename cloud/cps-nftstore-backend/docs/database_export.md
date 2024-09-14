# Export

1. Install the following applications in your mac. Special thanks to [this link](https://www.mongodb.com/docs/database-tools/installation/installation-macos/).

  ```shell
  brew tap mongodb/brew
  brew install mongodb-database-tools
```

2. Make sure you are in this repo's folder.

  ```shell
  cd ~/go/src/github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend
  ```

2. Start your `docker-compose` so the server is running.

3. Run the following [mongodump](https://www.mongodb.com/docs/database-tools/mongodump/) to backup the running `mongodb`:

  ```shell
  mongodump
  ```

4. If you look in the `dump` folder you can now see it!

  ```shell
  cd ~/go/src/github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend
  ```

# Restore

1. Sees above for steps 1 & 2.

2. Run the [mongorestore](https://www.mongodb.com/docs/database-tools/mongorestore/) to restore.

  ```shell
  mongorestore  dump/
  ```
