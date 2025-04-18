# Events

**E**venementen **V**erzamelen en **E**valueren, **N**otificeren van **T**odo's en **S**tatussen

As the name suggests it's a tool to help remember all the required steps for a successful event!

## Goal

It strives to create as little manual work as possible. \
While this is not possible for every check, the vast majority of them will be fully automatic. \
In the end the following checks will (hopefully) be supported:

- Announcements
- A well written website event page
- Posters
- DSA website entry
- Reservations (if it's not taking place in the Kelder)
- Mentioned in a presentation for every bachelor year
- Mentioned in an email
- Custom checks

With each check having it's own deadline and reminder notifications (emails) for the organisers.

## Backend

The backend is written in `Golang`.
It uses (not a complete list):

- [Fiber](https://pkg.go.dev/github.com/gofiber/fiber/v2)
- [Validator](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [Sqlc](https://pkg.go.dev/github.com/kyleconroy/sqlc)
- [Zap](https://pkg.go.dev/go.uber.org/zap)

## Frontend

The frontend is located in [ui](./ui) and written in `Typescript`.
It uses:

- [React](https://react.dev/)
- [Tanstack Router](https://tanstack.com/router/latest)
- [Tanstack Query](https://tanstack.com/query/latest)
- [Shadcn](https://ui.shadcn.com/)
- [Tailwind](https://tailwindcss.com/)
- [Zod](https://zod.dev/)

## Production

A docker container gets build every time main gets updated.

1. Set `APP_ENV=production` in a `.env` file and mount it in the container to `.env`
2. Create and configure a `production.yml` file and mount it to `/config/production.yml`
3. Make sure an external database is running
4. Run the container, the server listens to port 4000

## Development

### Prerequisites

- Install the required versions of `Golang` and `Nodejs`. They can be found in the [asdf tool versions file](.tool-versions).
- (Optional) Install the pre-commit hooks
  - Install the required version of `Golangci-lint`. It can be found in the [asdf tool versions file](.tool-versions)
  - Install the pre-commit hooks `git config --local core.hooksPath .githooks/`.
- Install sqlc `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`.
- Install goose `go install github.com/pressly/goose/v3/cmd/goose@latest`.
- If you're using asdf you might have to run `asdf reshim golang`.

### Run the application

Configure the environment variables in both the [backend](.env.example) and [frontend](ui/.env.example). \
Migrate the database by starting the database `make db` and running the migrations `make migrate`. \
Start the backend & frontend (both supporting HMR) `make watch`.

### Commits

Commit messages should adhere to the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) standards.
A [githook](.githooks/commit-msg) and [workflow](.github/workflows/commit-message.yml) is used to enforce it.

## Commands

All commands can be found in the [makefile](makefile).
A short list:

- `make build` - Build the application to a single docker container.
- `make watch` - Start the database, backend & frontend with HMR
- `make migrate` - Run all pending migrations.
- `make test` - Run all tests

## Useful flows

### Adding a new typed query (sqlc)

1. Add your new query to `db/queries/{target}.sql`` file.
2. Run `make sqlc`.
3. Use the new query in your code.

### Adding a migration

1. Run `make create-migration`.
2. Edit the newly made migration that can be found in the `db/migrations` folder.
3. Update the queries in the `db/queries` folder accordingly.
4. Run `make sqlc`
