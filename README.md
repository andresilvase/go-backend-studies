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
├── main.go                    # Selects and runs the current exercise
├── docker-compose.yml          # Local PostgreSQL 17 service
├── go.mod                     # Go module and dependencies
└── topics/
	├── syntax/
	│   └── syntax.go          # Go language fundamentals
	└── transactions/
		├── challenges/
		│   ├── 1.go           # Create related records atomically
		│   ├── 2.go           # Transfer money atomically
		│   ├── 3.go           # Simulate transaction failure and rollback
		│   ├── 4.go           # Transfer without transaction safeguards
		│   ├── 5.go           # Observe a non-repeatable read
		│   ├── 6.go           # Prevent a non-repeatable read
		│   ├── 7.go           # Observe a phantom read
		│   ├── 8.go           # Retry serializable transactions
		│   ├── 9.go           # Concurrent order intents
		│   ├── 10.go          # Create a user, wallet, and profile atomically
		│   ├── models/        # Shared challenge parameters
		│   └── shared/        # CSV-backed challenge helpers
		├── database/
		│   ├── connection.go  # PostgreSQL connection
		│   ├── queries/       # Application SQL queries
		│   ├── pgstore/       # sqlc-generated Go accessors
		│   └── migrations/    # Versioned database schema
		├── errors/
		│   └── custom_errors.go # Custom error types
		└── data/
			├── products.csv    # Product seed data
			└── users.csv       # User seed data
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
4. `004_create_products_table.sql` creates products.
5. `005_create_orders_table.sql` creates orders.
6. `006_order_items_table.sql` creates order items.
7. `007_product_inventory_table.sql` adds product stock.
8. `008_create_user_profile_table.sql` adds a profile related to a user and
	wallet.

## Transaction challenges

### Challenge 1 — create a user and wallet together

[`1.go`](topics/transactions/challenges/1.go) inserts a user, captures its
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

[`2.go`](topics/transactions/challenges/2.go) debits one wallet and credits
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

[`3.go`](topics/transactions/challenges/3.go) demonstrates transaction
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

### Challenge 4 — transfer without transaction safeguards

[`4.go`](topics/transactions/challenges/4.go) performs a money transfer
without using a transaction, exposing the risks of unprotected multi-step
operations.

It exercises:

- input validation (amount > 0, distinct wallets);
- executing debit and credit as separate, unprotected updates;
- checking `RowsAffected` to validate wallet existence and balance;
- simulating a failure between the debit and credit steps; and
- demonstrating data inconsistency when transaction safeguards are absent.

The key lesson is **the danger of unprotected updates**: a failure between the
debit and credit leaves the database in an inconsistent state (money removed
from one wallet but never added to another), illustrating why Challenge 3's
transactional approach is essential for money transfer operations.

### Challenge 5 — concurrent transactions and isolation levels

[`5.go`](topics/transactions/challenges/5.go) demonstrates concurrent
transactions running in parallel, using goroutines and channels to coordinate
their execution and reproduce a non-repeatable read under PostgreSQL's default
`READ COMMITTED` isolation level.

It exercises:

- spawning multiple goroutines with `sync.WaitGroup`;
- coordinating goroutine execution with channels;
- running concurrent database transactions;
- reading data from within an active transaction;
- observing how transaction isolation affects concurrent reads and writes; and
- detecting a non-repeatable read when the same transaction sees a newly
  committed value on its second query.

The key lesson is **isolation matters in concurrency**: how the database isolates
concurrent transactions determines what data each transaction observes, which is
critical for building correct multi-user applications.

### Challenge 6 — prevent non-repeatable reads

[`6.go`](topics/transactions/challenges/6.go) repeats Challenge 5 with
transaction A running at `REPEATABLE READ` isolation. Transaction B updates and
commits the wallet balance between transaction A's two reads, but transaction A
continues to see its original snapshot.

It exercises:

- selecting `sql.LevelRepeatableRead` with `sql.TxOptions`;
- coordinating two concurrent transactions with goroutines and channels;
- reading the same wallet before and after another transaction commits;
- comparing snapshot behavior with the default `READ COMMITTED` behavior; and
- preventing non-repeatable reads within a transaction.

The key lesson is **a stable snapshot produces repeatable reads**: changes
committed by another transaction are not visible until the repeatable-read
transaction finishes.

### Challenge 7 — observe a phantom read

[`7.go`](topics/transactions/challenges/7.go) runs the same predicate query
twice in transaction A while transaction B changes a wallet so that it newly
matches the predicate `balance > 1000` between those reads.

It exercises:

- coordinating predicate reads and a concurrent update;
- counting rows that match a condition inside an active transaction;
- committing a change from a second transaction between two reads;
- observing the result set change under the default `READ COMMITTED` isolation
  level; and
- demonstrating the phantom-read isolation phenomenon.

The key lesson is **predicate results can change during a transaction**: at
`READ COMMITTED`, each statement receives a new snapshot and may see rows that
did not match an earlier query.

### Challenge 8 — serializable transactions with retry

[`8.go`](topics/transactions/challenges/8.go) runs two concurrent transactions
at `SERIALIZABLE` isolation, synchronizing both after they read the other
wallet's balance before attempting their updates.

It exercises:

- coordinating concurrent transactions with goroutines, wait groups, and buffered channels;
- using a barrier to reproduce a serialization conflict deterministically;
- cancelling the shared context when a transaction fails before reaching the barrier;
- identifying PostgreSQL serialization failures (`40001`) and deadlocks (`40P01`) as retryable errors;
- retrying transient failures with exponential backoff, up to five attempts; and
- collecting errors from both transactions while allowing each goroutine to finish cleanly.

The key lesson is **serializable isolation needs retry handling**: the database
can reject one of two conflicting transactions to preserve consistency, so the
application must retry transient failures instead of treating them as permanent
operation errors.

### Challenge 9 — concurrent order intents and inventory

[`9.go`](topics/transactions/challenges/9.go) generates order intents from the
CSV fixtures and attempts to complete them concurrently using serializable
transactions.

It exercises:

- loading users and products from CSV data;
- creating orders and order items in a transaction;
- decrementing product inventory atomically;
- coordinating many purchase goroutines with a wait group and channel;
- recognizing serialization failures and deadlocks; and
- retrying transient failures with exponential backoff and jitter.

The key lesson is **contention requires coordination**: concurrent purchases
must protect shared inventory and retry when PostgreSQL detects a conflict.

### Challenge 10 — create a user, wallet, and profile together

[`10.go`](topics/transactions/challenges/10.go) selects a user from the CSV
fixtures and creates that user, a random-balance wallet, and a profile in one
transaction.

It exercises:

- composing generated `sqlc` queries through a transaction;
- passing generated IDs from one insert to the next;
- creating related records in dependency order; and
- rolling back the whole operation when any step fails.

The key lesson is **related writes share one boundary**: the user, wallet, and
profile should either all be created or none should be persisted.

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

Create `topics/transactions/database/migrations/.env` with the following values:

```dotenv
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=transactions_lab
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
```

These values are used by Docker Compose, Tern, and the migration helper. The
database connection currently uses the same local PostgreSQL defaults directly
in `topics/transactions/database/connection.go`. The `.env` file is ignored by
Git.

### 5. Start PostgreSQL

From the repository root:

```bash
docker compose --env-file topics/transactions/database/migrations/.env up -d
```

### 6. Apply the schema migrations

```bash
make migrations
```

Check their status at any time with:

```bash
make status
```

### 7. Prepare data for the default exercise

`main.go` currently runs Challenge 10. It reads a random user from
`topics/transactions/data/users.csv`, then creates that user, a wallet, and a
profile. The migration-created tables are enough to run it; no manual seed
insert is required.

Challenges 1–8 use wallet IDs `1` and `2`. If you choose one of those exercises
on a fresh database, create two users and two wallets first. Both balances start
at `600` so the serializable challenge can exercise retry logic:

```bash
docker compose --env-file topics/transactions/database/migrations/.env exec postgres \
  psql -U postgres -d transactions_lab -c \
  "INSERT INTO users (name) VALUES ('Source User'), ('Target User'); INSERT INTO wallets (user_id, balance) VALUES (1, 600), (2, 600);"
```

### 8. Run the current challenge

```bash
go run .
```

The exact order varies because the transactions run concurrently, but the
output includes messages similar to:

```text
Connected to PostgreSQL.
Running Challenge 10...
```

You can inspect the result with:

```bash
docker compose --env-file topics/transactions/database/migrations/.env exec postgres \
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
var sourceWalletID int64 = 1
var targetWalletID int64 = 2
var amount int64 = 100

var chTwoParam mdl.TransferParams = mdl.TransferParams{
	Ctx:            &ctx,
	DB:             db,
	SourceWalletID: &sourceWalletID,
	TargetWalletID: &targetWalletID,
	Amount:         &amount,
}
if err := challenges.Two(chTwoParam); err != nil {
	log.Fatal(err)
}
```

To run **Challenge 3**, enable:

```go
var sourceWalletID int64 = 1
var targetWalletID int64 = 2
var simulatedFail = true
var amount int64 = 250

var chThreeParam mdl.TransferParams = mdl.TransferParams{
	Ctx:            &ctx,
	DB:             db,
	SourceWalletID: &sourceWalletID,
	TargetWalletID: &targetWalletID,
	SimulatedFail:  &simulatedFail,
	Amount:         &amount,
}
if err := challenges.Three(chThreeParam); err != nil {
	log.Fatal(err)
}
```

To run **Challenge 4**, enable:

```go
var sourceWalletID int64 = 1
var targetWalletID int64 = 2
var simulatedFail = true
var amount int64 = 250

var chFourParam mdl.TransferParams = mdl.TransferParams{
	Ctx:            &ctx,
	DB:             db,
	SourceWalletID: &sourceWalletID,
	TargetWalletID: &targetWalletID,
	SimulatedFail:  &simulatedFail,
	Amount:         &amount,
}
if err := challenges.Four(chFourParam); err != nil {
	log.Fatal(err)
}
```

To run **Challenge 5**, enable:

```go
var sourceWalletID int64 = 1
var targetWalletID int64 = 2
var amount int64 = 100

var chFiveParam mdl.TransferParams = mdl.TransferParams{
	Ctx:            &ctx,
	DB:             db,
	SourceWalletID: &sourceWalletID,
	TargetWalletID: &targetWalletID,
	Amount:         &amount,
}
var wg sync.WaitGroup
if err := challenges.Five(chFiveParam, &wg); err != nil {
	log.Fatal(err)
}
wg.Wait()
```

To run **Challenge 6**, enable:

```go
var sourceWalletID int64 = 1
var targetWalletID int64 = 2
var amount int64 = 250

var chSixParam mdl.TransferParams = mdl.TransferParams{
	Ctx:            &ctx,
	DB:             db,
	SourceWalletID: &sourceWalletID,
	TargetWalletID: &targetWalletID,
	Amount:         &amount,
}
var wg sync.WaitGroup
if err := challenges.Six(chSixParam, &wg); err != nil {
	log.Fatal(err)
}
wg.Wait()
```

To run **Challenge 7**, enable:

```go
var sourceWalletID int64 = 1
var targetWalletID int64 = 2
var amount int64 = 250

var chSevenParam mdl.TransferParams = mdl.TransferParams{
	Ctx:            &ctx,
	DB:             db,
	SourceWalletID: &sourceWalletID,
	TargetWalletID: &targetWalletID,
	Amount:         &amount,
}
var wg sync.WaitGroup
if err := challenges.Seven(chSevenParam, &wg); err != nil {
	log.Fatal(err)
}
wg.Wait()
```

To run **Challenge 8**, enable:

```go
var sourceWalletID int64 = 1
var targetWalletID int64 = 2
var amount int64 = 250

var chEightParam mdl.TransferParams = mdl.TransferParams{
	Ctx:            &ctx,
	DB:             db,
	SourceWalletID: &sourceWalletID,
	TargetWalletID: &targetWalletID,
	Amount:         &amount,
}
var wg sync.WaitGroup
if err := challenges.Eight(chEightParam, &wg); err != nil {
	log.Fatal(err)
}
wg.Wait()
```

To run **Challenge 9**, configure `main.go` with the order-intent parameters:

```go
const numberOfOrderIntents = 50

var wg sync.WaitGroup
params := mdl.ChallengeNineParams{
	Ctx:             ctx,
	DB:              db,
	WG:              &wg,
	NumberOfIntents: numberOfOrderIntents,
}
if err := challenges.Nine(params); err != nil {
	log.Fatal(err)
}
wg.Wait()
```

To run **Challenge 10**, configure `main.go` with:

```go
params := mdl.ChallengeTenParams{Ctx: ctx, DB: db}
if err := challenges.Ten(params); err != nil {
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
make status

# Roll the database schema back to version 0 (destructive)
make rollback-all

# Stop PostgreSQL but retain its data volume
docker compose --env-file topics/transactions/database/migrations/.env down

# Stop PostgreSQL and delete the local database volume
docker compose --env-file topics/transactions/database/migrations/.env down -v
```

The root `Makefile` also provides:

```bash
# Generate sqlc code after changing queries or schema
make sql

# Seed users and products from the CSV fixtures
make db_fed

# Roll back, migrate, and seed in one command
make reset
```
