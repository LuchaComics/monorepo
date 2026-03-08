# How to Import a MongoDB Dump into a New Local Dev Instance

A quick guide for importing an existing `dump/cps_db` directory into a fresh local development environment.

> For the full export-from-production + import workflow, see [HOW_TO_EXPORT_MONGODB_AND_IMPORT_LOCALLY.md](HOW_TO_EXPORT_MONGODB_AND_IMPORT_LOCALLY.md).

## Prerequisites

- **Docker** installed and running
- **MongoDB Database Tools** installed locally (for Option A):
  ```bash
  brew install mongodb-database-tools mongosh
  ```
- A `./dump/cps_db/` directory containing `.bson` and `.metadata.json` files

## Step 1: Start the Dev Environment

From the `cloud/cps-backend/` directory:

```bash
task start
```

Or manually:

```bash
docker-compose -p cps -f dev.docker-compose.yml up -d
```

This starts a 3-node MongoDB replica set (`rs0`):

| Container  | Port  | Role      |
|------------|-------|-----------|
| `cps_db1`  | 27017 | Primary   |
| `cps_db2`  | 27018 | Secondary |
| `cps_db3`  | 27019 | Secondary |

Wait ~10 seconds for the replica set to initialize before proceeding.

## Step 2: Import the Dump

### Option A: From your host machine (recommended)

Requires `mongorestore` installed locally.

```bash
mongorestore \
  --uri="mongodb://localhost:27017/?replicaSet=rs0" \
  --db=cps_db \
  --drop \
  ./dump/cps_db
```

### Option B: From inside the Docker container

If you don't have MongoDB tools installed locally:

```bash
# Copy the dump into the primary container
docker cp ./dump/cps_db cps_db1:/tmp/dump

# Run mongorestore inside the container
docker exec cps_db1 mongorestore \
  --db=cps_db \
  --drop \
  /tmp/dump

# Clean up
docker exec cps_db1 rm -rf /tmp/dump
```

**Note:** The `mongo:7.0` image includes `mongorestore`. If you get "command not found", use Option A instead.

## Step 3: Create an Admin Account

After importing a production dump, you won't know the passwords for the existing users. The app has a built-in mechanism to auto-create a root admin on startup — we just need to configure it with an email that **does not already exist** in the imported dump.

1. Edit `cloud/cps-backend/.env` and set the following:

   ```dotenv
   CPS_BACKEND_INITIAL_ADMIN_EMAIL=admin@cpsapp.ca
   CPS_BACKEND_INITIAL_ADMIN_PASSWORD=123password
   CPS_BACKEND_INITIAL_ADMIN_ORG_NAME=CPS
   ```

   Replace the email and password with your desired credentials. The email **must not** match any user already in the imported dump.

2. Restart the app so it picks up the new env vars:

   ```bash
   task end && task start
   ```

   Or manually:

   ```bash
   docker-compose -p cps -f dev.docker-compose.yml restart app
   ```

3. On startup, the app checks if a user with `CPS_BACKEND_INITIAL_ADMIN_EMAIL` exists in the database. Since the email is new, it will auto-create a root admin account and an associated root Store.

4. Log in at [http://localhost:8000](http://localhost:8000) with the email and password you configured.

> **Note:** If the email already exists in the imported dump, the app skips creation and logs `"Root user already exists, skipping creation."` — pick an email that is not in the dump.

## Step 4: Verify the Import

```bash
mongosh "mongodb://localhost:27017/?replicaSet=rs0"
```

```javascript
use cps_db
show collections

// Check document counts
db.getCollectionNames().forEach(function(c) {
    print(c + ": " + db.getCollection(c).countDocuments() + " documents");
})
```

## Step 5: Access the App

| Service         | URL                    |
|-----------------|------------------------|
| Backend API     | http://localhost:8000   |
| Mongo Express   | http://localhost:8081   |

## Troubleshooting

### Replica set not ready

```
Failed: error connecting to db server: no reachable servers
```

The replica set needs a few seconds to initialize after startup. Check status:

```bash
mongosh "mongodb://localhost:27017/?replicaSet=rs0" --eval "rs.status()"
```

Verify all containers are running:

```bash
docker ps | grep cps_db
```

### "not primary" error

If `mongorestore` reports the node is not primary, connect directly to port 27017:

```bash
mongorestore \
  --host=localhost \
  --port=27017 \
  --db=cps_db \
  --drop \
  ./dump/cps_db
```

### `mongorestore: command not found`

Install MongoDB Database Tools:

```bash
brew install mongodb-database-tools
```
