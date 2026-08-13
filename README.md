# Go Backend Studies

A hands-on Go playground for learning backend fundamentals by writing small,
focused examples and then applying them to realistic database problems.

The repository currently explores:

- core Go syntax and built-in data structures;
- PostgreSQL access through `database/sql` and `pgx`;
- schema versioning with Tern migrations;
- atomic multi-step operations with SQL transactions; and
- common money-transfer safeguards such as preventing negative balances.

> This is a study repository, not a production service. The examples favor
> visibility and experimentation over abstractions and application structure.

## Repository map

```text
.
├── main.go                         # Selects and runs the current exercise
├── docker-compose.yml               # Local PostgreSQL 17 service
├── go.mod                          # Go module and dependencies
└── topics/
    ├── syntax/
    │   └── syntax.go                   # Go language fundamentals
    └── transactions/
        ├── Makefile                    # Migration shortcuts
        ├── challenges/
        │   ├── one.go                  # Create related records atomically
        │   ├── two.go                  # Transfer money atomically
        │   └── three.go                # Simulate transaction failure and rollback
        ├── database/
        │   ├── database.go              # PostgreSQL connection
        │   └── migrations/              # Versioned database schema
        └── errors/
            └── custom_errors.go         # Custom error types
```

## Topics

### `topics/syntax`

Small examples for becoming comfortable with the language before moving into
database-backed code. `syntax.Run()` demonstrates:

- classic and condition-based `for` loops;
- declaring, growing, iterating over, and combining slices;
- creating and updating maps;
- `range`; and
- formatted console output with `fmt`.

### `topics/transactions`

A PostgreSQL lab built around users and their wallets. It demonstrates opening
and validating a database connection, executing parameterized SQL, scanning
returned IDs, handling errors with context, and treating several writes as one
atomic unit.

The schema is introduced incrementally:

1. `001_create_users_table.sql` creates users.
2. `002_create_wallets_table.sql` creates one wallet per user through a unique
   foreign key.
3. `003_alter_wallets_table.sql` adds a database-level rule that prevents a
   negative balance.

## Transaction challenges

### Challenge 1 — create a user and wallet together

[`one.go`](topics/transactions/challenges/one.go) inserts a user, captures its
generated ID, and creates the user's wallet inside the same transaction.

It exercises:

- starting a transaction with `BeginTx`;
- passing `*sql.Tx` to helper functions;
- using `INSERT ... RETURNING id` with `QueryRowContext`;
- preserving referential integrity between users and wallets;
- committing only after every operation succeeds; and
- using a deferred rollback so partial work is discarded on failure.

The key lesson is **all or nothing**: a user should never be persisted without
the wallet that belongs to them.

### Challenge 2 — transfer money between wallets

[`two.go`](topics/transactions/challenges/two.go) debits one wallet and credits
another inside a single transaction.

It exercises:

- grouping a debit and credit into one atomic operation;
- parameterized `UPDATE` statements;
- enforcing sufficient funds in the debit query itself;
- checking `RowsAffected` to detect missing wallets or a rejected debit;
- rolling back when either half of the transfer fails; and
- reinforcing business rules with the database `CHECK` constraint.

The key lesson is **conservation of money**: the debit and credit must both
succeed, or neither change should remain.

### Challenge 3 — simulate a failure mid-transfer and verify rollback

[`three.go`](topics/transactions/challenges/three.go) demonstrates transaction
safety by simulating a failure between the withdrawal and deposit operations.

It exercises:

- implementing a parameterized challenge function with options;
- injecting a simulated failure point mid-transaction;
- verifying that when an error occurs between operations, the entire transaction
  rolls back and no partial changes persist;
- using deferred rollback as a safety net; and
- confirming that both wallets remain in their original state after a failed
  transfer.

The key lesson is **failure safety**: a partial or failed transfer leaves the
database in a consistent state, with no orphaned debits or credits.

## Run from a fresh clone

### Prerequisites

- [Go](https://go.dev/dl/) **1.25 or newer**
- [Docker](https://docs.docker.com/get-docker/) with Docker Compose
- [`tern`](https://github.com/jackc/tern), used to run the migrations
- `make`

### 1. Clone and enter the repository

```bash
git clone <repository-url>
cd go-backend-studies
```

### 2. Download the Go dependencies

```bash
go mod download
```

### 3. Install Tern

```bash
go install github.com/jackc/tern/v2@latest
```

Make sure the Go binary directory is on your `PATH` so the `tern` command can
be found.

### 4. Configure the local database

Create `topics/transactions/.env` with the following values:

```dotenv
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=transactions_lab
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
```

These values intentionally match the connection string currently used by
`topics/transactions/database/database.go`. The `.env` file is ignored by Git.

### 5. Start PostgreSQL

From the repository root:

```bash
docker compose --env-file topics/transactions/.env up -d
```

### 6. Apply the schema migrations

```bash
make -C topics/transactions run-migrations
```

Check their status at any time with:

```bash
make -C topics/transactions status
```

### 7. Prepare data for the default exercise

`main.go` currently runs Challenge 2 with wallet IDs `1` and `2`, transferring
`100` units. A new database therefore needs two users, two wallets, and a
balance in the source wallet:

```bash
docker compose --env-file topics/transactions/.env exec postgres \
  psql -U postgres -d transactions_lab -c \
  "INSERT INTO users (name) VALUES ('Source User'), ('Target User'); INSERT INTO wallets (user_id, balance) VALUES (1, 500), (2, 0);"
```

### 8. Run the current challenge

```bash
go run .
```

Expected output:

```text
Connected to PostgreSQL.
Transfering money from wallet 1 to wallet 2...
Transference completed successfully!
```

You can inspect the result with:

```bash
docker compose --env-file topics/transactions/.env exec postgres \
  psql -U postgres -d transactions_lab \
  -c "SELECT id, user_id, balance FROM wallets ORDER BY id;"
```

## Choosing an exercise

[`main.go`](main.go) is the repository's small exercise runner. Uncomment the
call you want to study and comment out the other active call.

To run **Challenge 1**, enable:

```go
if err := challenges.One(ctx, db); err != nil {
	log.Fatal(err)
}
```

To run **Challenge 2**, enable:

```go
if err := challenges.Two(ctx, db, 1, 2, 100); err != nil {
	log.Fatal(err)
}
```

To run the **syntax examples**, import the package and call `syntax.Run()`:

```go
import syntax "transactions-lab/topics/syntax"

func main() {
	syntax.Run()
}
```

## Useful commands

```bash
# Compile and run package checks
go test ./...

# Format all Go packages
go fmt ./...

# Show migration status
make -C topics/transactions status

# Roll the database schema back to version 0 (destructive)
make -C topics/transactions rollback-all

# Stop PostgreSQL but retain its data volume
docker compose --env-file topics/transactions/.env down

# Stop PostgreSQL and delete the local database volume
docker compose --env-file topics/transactions/.env down -v
```

## Ideas for further study

Natural next steps for this lab include adding table-driven tests, moving the
database URL fully into environment-based configuration, validating positive
transfer amounts, locking rows for concurrent transfers, and exposing the
operations through an HTTP API.
