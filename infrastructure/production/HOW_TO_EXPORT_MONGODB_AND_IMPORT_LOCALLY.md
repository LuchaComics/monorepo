# How to Export Production MongoDB and Import Locally

This guide walks through exporting the production MongoDB database (DigitalOcean Managed MongoDB) and importing it into your local Docker Compose development environment.

## Prerequisites

- **MongoDB Database Tools** installed locally (`mongodump`, `mongorestore`, `mongosh`)
  - macOS: `brew install mongodb-database-tools mongosh`
  - Ubuntu/Debian: Follow [MongoDB docs](https://www.mongodb.com/docs/database-tools/installation/)
- **Local development environment running** via `task dev` (from `cloud/backend/`)
- **Production connection string** from DigitalOcean (ask your team lead)

## Step 1: Export from Production

Run `mongodump` from your local machine against the DigitalOcean Managed MongoDB instance. No SSH is needed since managed databases expose a direct connection string.

```bash
mongodump \
  --uri="mongodb+srv://<username>:<password>@<host>/cps_db?authSource=admin&tls=true" \
  --out=./dump
```

If your connection string uses the standard `mongodb://` format (common with DigitalOcean):

```bash
mongodump \
  --uri="mongodb://<username>:<password>@<host>:<port>/cps_db?authSource=admin&tls=true" \
  --out=./dump
```

**Notes:**
- Replace `<username>`, `<password>`, `<host>`, and `<port>` with your actual production credentials.
- The `authSource=admin` parameter tells MongoDB to authenticate against the `admin` database.
- The `tls=true` parameter is required for DigitalOcean Managed MongoDB.
- If DigitalOcean provides a CA certificate, add `--tlsCAFile=<path-to-ca-cert.crt>`.

This creates a `./dump/cps_db/` directory containing BSON files for each collection.

### Export a single collection (optional)

If you only need one collection:

```bash
mongodump \
  --uri="mongodb+srv://<username>:<password>@<host>/cps_db?authSource=admin&tls=true" \
  --collection=<collection_name> \
  --out=./dump
```

## Step 2: Transfer the Dump

Since `mongodump` runs on your local machine and connects directly to the managed database, the dump files are already local. No file transfer is needed.

Verify the dump was created:

```bash
ls -la ./dump/cps_db/
```

You should see `.bson` and `.metadata.json` files for each collection.

## Step 3: Prepare Local Environment

### Start the development environment

From the `cloud/backend/` directory:

```bash
task dev
```

This starts the local 3-node MongoDB replica set:

| Container                   | Host        | Port  | Role              |
|-----------------------------|-------------|-------|-------------------|
| `cps_db1`    | localhost   | 27017 | Primary           |
| `cps_db2`    | localhost   | 27018 | Secondary         |
| `cps_db3`    | localhost   | 27019 | Secondary         |

Replica set name: `rs0`

### (Optional) Drop the existing local database

If you want a clean slate, drop the existing database first:

```bash
mongosh "mongodb://localhost:27017/?replicaSet=rs0" --eval "use cps_db" --eval "db.dropDatabase()"
```

Or in a single interactive session:

```bash
mongosh "mongodb://localhost:27017/?replicaSet=rs0"
```

```javascript
use cps_db
db.dropDatabase()
exit
```

## Step 4: Import into Local Docker

### Option A: Run `mongorestore` from your host machine

This is the simplest approach. Run `mongorestore` against the local replica set primary:

```bash
mongorestore \
  --uri="mongodb://localhost:27017/?replicaSet=rs0" \
  --db=cps_db \
  --drop \
  ./dump/cps_db
```

**Flags:**
- `--drop`: Drops each collection before importing (ensures a clean import).
- `--db=cps_db`: Target database name.

### Option B: Run `mongorestore` from inside the Docker container

If you don't have `mongorestore` installed locally, you can copy the dump into the container and restore from there:

```bash
# Copy the dump into the primary container
docker cp ./dump/cps_db cps_db1:/tmp/dump

# Run mongorestore inside the container
docker exec cps_db1 mongorestore \
  --db=cps_db \
  --drop \
  /tmp/dump

# Clean up the dump from the container
docker exec cps_db1 rm -rf /tmp/dump
```

**Note:** The container image may not include `mongorestore`. If you get a "command not found" error, install MongoDB Database Tools locally and use Option A.

## Step 5: Verify the Import

Connect to the local database and check that collections and documents are present:

```bash
mongosh "mongodb://localhost:27017/?replicaSet=rs0"
```

```javascript
use cps_db

// List all collections
show collections

// Check document counts for key collections
db.getCollectionNames().forEach(function(c) {
    print(c + ": " + db.getCollection(c).countDocuments() + " documents");
})
```

Compare collection counts against production to confirm the import is complete.

## Troubleshooting

### Authentication errors during export

```
Failed: error connecting to db server: authentication error
```

- Double-check your username, password, and host.
- Ensure `authSource=admin` is in the URI.
- Ensure `tls=true` is in the URI for DigitalOcean Managed MongoDB.
- If using a CA certificate, verify the path with `--tlsCAFile`.

### Replica set not ready during import

```
Failed: error connecting to db server: no reachable servers
```

- Wait a few seconds after `task dev` for the replica set to initialize.
- Check replica set status:
  ```bash
  mongosh "mongodb://localhost:27017/?replicaSet=rs0" --eval "rs.status()"
  ```
- Ensure all three containers are running:
  ```bash
  docker ps | grep cps_db
  ```

### Large dump files

For very large databases, the dump and restore can take a long time. You can speed up the restore with parallel collections:

```bash
mongorestore \
  --uri="mongodb://localhost:27017/?replicaSet=rs0" \
  --db=cps_db \
  --drop \
  --numParallelCollections=4 \
  ./dump/cps_db
```

### "not primary" error during restore

If `mongorestore` reports the node is not primary, explicitly connect to the primary:

```bash
mongorestore \
  --host=localhost \
  --port=27017 \
  --db=cps_db \
  --drop \
  ./dump/cps_db
```

### Index build failures

If indexes fail to build during restore, you can skip them and rebuild manually:

```bash
mongorestore \
  --uri="mongodb://localhost:27017/?replicaSet=rs0" \
  --db=cps_db \
  --drop \
  --noIndexRestore \
  ./dump/cps_db
```

Then rebuild indexes inside `mongosh`:

```javascript
use cps_db
db.getCollectionNames().forEach(function(c) {
    db.getCollection(c).reIndex();
})
```

## Security Reminders

- **Never commit database dumps to version control.** The `dump/` directory should be in your `.gitignore`.
- **Delete dumps after use.** Production data should not persist on developer machines longer than needed.
- **Sanitize sensitive data.** If working with PII or payment data, consider sanitizing the data after import:
  ```javascript
  // Example: anonymize email addresses
  use cps_db
  db.users.updateMany({}, [{ $set: { email: { $concat: ["user_", { $toString: "$public_id" }, "@example.com"] } } }])
  ```
- **Do not share dump files** over insecure channels (email, Slack, public links).
- **Rotate credentials** if production credentials are exposed.
