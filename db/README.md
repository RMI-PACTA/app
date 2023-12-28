# Database

This project uses [PostgreSQL](https://www.postgresql.org/) for its
database, and almost all of the code for it lives in `/db/sqldb`.

## Development

The database can be started locally by running:

```
bazel run //scripts:run_db
```

This will set it up in such a way that it can recieve connections from other
parts of the stack running locally.

To connect to the locally running database, run:

```
bazel run //scripts:db_shell
```

To run tests in this package (which will spin up their own databases, no need to
have the DB running), run:

```
bazel test //db/...
```

## About

### Migrations

To support safe versioning of your schema, we use migrations to describe any
and all mutations to your database schema. Any time you want to alter your
database's schema, create two files, one describing the transaction that you
want to apply, and the second describing how you would roll it back.

The naming convention is:

```
[Monotonically Increasing Number]_[brief description].[up|down].sql
```

For example:

```
db/sqldb/migrations/0002_account_table.down.sql
db/sqldb/migrations/0002_account_table.up.sql
```

There's a migration test that validates up and down migrations are true inverses
of one another, by comparing the initial database state to the database state if
the rollup and the rollback mutations are done in sequence.

Warning: Once a mutation has been checked in, it should not be altered. Instead,
a new migration should be created with the required changes. This approach
allows for robust data handling and data migrations that can have thorough
integration tests and backward compatability tests.

### Goldens

Two golden files are expected to be checked in alongside your code in
the `golden` repo: a dump of your schema, and a human readable version
of your schema. These are purely for code review purposes - it allows
a reviewer to know what mutations your migrations are proposing, and
validates that the migrations have in fact passed the basic tests for
stability and sanity embedded in the golden regeneration tests. You can
run these at any time (idempotently) with:

```
bazel run //scripts:regen_db_goldens
```

### Testing

Testing is demonstrated in the `_test.go` files. Clean versions of
the database are stood up to run each test, and since you're testing
against a local version of postgres, these tests are really
integration tests, and you can use them to not only validate your
business logic, but how you anticipate postgres responding to various
situations (foreign key constraint violations, etc).
