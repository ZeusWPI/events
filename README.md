# Events

A tool to help you plan a successful event!

## Getting started

### Quickstart

- Download the version for `golang`, `nodejs` and `pnpm`. The versions can be found inside the [asdf tool versions file](./.tool-versions) (Run `asdf install` if you're using asdf)
- Install make
- Run `make setup` to download the backend tools `Air`, `Goose`, `Sqlc` and frontend dependencies
- Copy the example [environment file](./.env.example) to `.env` in the backend and change the values if needed (necessary for a clean slate run)
- Configure the config file inside the [config directory](./config/) if needed (not necessary if you're only using the makefile)
- Migrate the database `make migrate`
- Run the project `make watch`

At the time of writing it supports the following features

### Full Explanation

The backend is written in Golang, the frontend in React + Typescript.
Versions for both can be found inside the [asdf tool versions file](./.tool-versions).

Workflows are used to ensure code quality.
You can run them locally before each commit by installing the githook `git config --local core.hooksPath .githooks/`. It requires `Golangci-lint`, the version can again be found inside the [asdf tool versions file](./.tool-versions).

A Makefile is used to simplify some workflows.

`make setup` downloads 2 additional golang tools.

- `goose` manages the migrations
- `sqlc` generates statically typed golang code from SQL queries

It also downloads all the frontend dependencies. To manually install them run `cd ui && pnpm install`.

`make migrate` will automatically start the postgres database and attempt to migrate it to the newest version.
If you don't want it to use the makefile or you have a separate database you can run `go run migrate.go` which uses the database connection information defined inside the config file.

`make watch` starts the entire docker stack. It supports HMR (hot module reloading) for both the backend and frontend. It only follows the logs of the backend and frontend.
If you want to see more logs you can use the command `docker compose up backend frontend` and add any additional container that you want to see the logs of.
For example to include the database logs use the command `docker compose up backend frontend db`.

> [!NOTE]
> A restart is required after adding or removing dependencies

## Server Setup

It is recommended to run the application in a docker container.

It needs the following additional resources:

- Postgres

> !NOTE
> Make sure to set the environment to `production` and populate the `production.yml` file.

## Useful flows

### Add a new typed query (sqlc)

1) Add your new query to db/queries/{target}.sql
2) Run `make query`
3) Enjoy your statically typed query

### Adding a migration

1) Run `make create-migration`

> [!NOTE]
> Nix users using devshell need to run `goose -dir ./db/migrations postgres create my_migration_name sql`

2) Edit the newly made migration that can be found in the `db/migrations` folder
3) Update the queries in the `db/queries` accordingly
4) Run `make query` to generate the new table structs
