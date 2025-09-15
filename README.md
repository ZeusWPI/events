# Events

A tool to help you plan and manage a Zeus WPI event!

## Features

The app automatically fetches all event and board data by listening to a GitHub webhook installed in the [Zeus website repository](https://github.com/ZeusWPI/zeus.ugent.be).
As a fallback, it periodically refetches the data.

### Organizers

You can assign organizers to each event.
Organizers are grouped by academic year and an organizer corresponds to a board member for that academic year.

### Checks

Each event includes checks to track the progress of its organization.
There are two types of checks: **manual** and **automatic**.

#### Manual Checks

Manual checks must be manually added and marked as done for each event.

#### Automatic Checks

Automatic checks are updated automatically. Currently, they include:

- Is it added to the DSA website?
- Is there an announcement written?
- Has the event been covered in a mail?
- Are there any posters for the event?

### Announcements

You can write announcements for each event and schedule them to be send at a later date.

### Mails

Similar to announcements, you can write and schedule emails.
For each mail you can select with events are covered.

### Posters

Each event can have 2 posters

- `Big` -> Poster meant to print out and hang up to advertise for the event
- `Scc` -> Poster displayed on cammie chat screen

The posters are automatically synced with the [visueel repository in Gitmate](https://git.zeus.gent/ZeusWPI/visueel).

If a poster is not yet in the `visueel` repository then a pull request will be created to add it after the event is finished.
If the `visueel` repository has a poster and events has none or a different one then the events poster will be deleted and replaced by the one in the `visueel` repository.

More information can be found in the [visueel repository](https://git.zeus.gent/ZeusWPI/visueel)

### Powerpoints

You can generate a PowerPoint covering one or more events.
If available it includes event posters and generates QR codes to the event's webpage.

### Public API

There's a public API to get some basic events data.
Click [here](https://events.zeus.gent/api/v1/docs) to go to the swagger.

---

## Development

In development, no roles or permissions are required to access the application.
However, only board members will be visible as organizers.

Some modules require external API keys. The will fail to launch without them and by result the application will fail to start
You can find each API key in the [example env file](./.env.example).
If you don't have an API key you should remove the startup of the relevant module inside the [main file](./cmd/api/main.go).

A more in depth explanation of this application can be found in the [internal README.md](/internal/README.md).

### Quickstart

1. Install all tools listed in the [asdf tool versions file](./.tool-versions).  
   If using `asdf`, run `asdf install`.
2. Install `make`.
3. Run `make setup` to install:
   - Backend tools: `Goose`, `Sqlc`
   - Frontend dependencies
4. Install Git hooks for code quality:

   ```bash
   git config --local core.hooksPath .githooks/
   ```

5) Copy the example environment file

   ```bash
   cp .env.example .env
   ```

   Update values as needed.
6) (Optional) Configure the [development config file](./config/development.yml) if you're not using the makefile.
7) Migrate the database:

   ```bash
   make migrate
   ```

8) Start the project

   ```bash
   make watch
   ```

- **Backend**: <http://localhost:4000>
- **Frontend**: <http://localhost:3000>

### Full Explanation

- **Backend**: Golang
- **Frontend**: React + Typescript

Workflows are used to ensure code quality.
You can run them manually before each commit by installing the githooks

```bash
git config --local core.hooksPath
```

They rely on [golangci-lint], which is included in the [asdf tool versions file](./.tool-versions).

A Makefile is used to simplify most tasks.

`make setup`

Installs:

- [goose](https://github.com/pressly/goose) manages the migrations
- [sqlc](https://github.com/sqlc-dev/sqlc) generates statically typed golang code from SQL queries
- Frontend dependencies (manually: `cd ui $$ pnpm install`)

`make migrate`

Starts a postgres container and applies the migrations.
If you want to use your own database you can run the migrations with

```bash
go run migrate.go
```

This uses the connection values specified in your config.

`make watch`

Starts the full Docker stack with hot module reloading (HMR) for both backend and frontend.
It follows logs for the backend and frontend by default.

To view logs from other services (like the database):

```bash
docker compose up backend frontend db
```

> [!NOTE]
> A restart is required after adding or removing dependencies

### Useful make targets

**Adding a migration**

1) Run `make create-migration`

> [!NOTE]
> Nix users using devshell need to run `goose -dir ./db/migrations postgres create my_migration_name sql`

2) Edit the newly made migration that can be found in the `db/migrations` folder
3) Update the queries in the `db/queries` accordingly
4) Run `make query` to generate the new table structs

**Adding a new typed query (Sqlc)**

1) Add your new query to `db/queries/{target}.sql`
2) Run

    ```bash
    make query
    ```

**Check for dead code**

1. Run `make dead`

**Generate swagger docs**

This will also format the swagger comments

1. Run `make swagger`

### Running without docker

> [!NOTE]
> Docker is strongly recommended for development

1) Install [Air](https://github.com/air-verse/air) for HMR

    ```bash
    go install github.com/air-verse/air@latest
    ```

2) Update your config files as needed
3) Start backend:

    ```bash
    air .
    ```

4) Start frontend

    ```bash
    cd ui && pnpm run dev
    ```

## Production

It is recommended to run the application using Docker.

**Requirements**

- Postgres
- Minio

> [!NOTE]
> Set the environment to `production` and populate the [production config file](./config/production.yml)

This repository automatically builds and publishes a docker container.
The container will run the migrations before starting the webserver.
