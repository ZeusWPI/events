# Events

Tool to help remember all the required steps for a successful event!

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

## Development

### Prerequisites

- Install the required versions of `Golang` and `Typescript`. Can be found in the [asdf tool versions file](.tool-version).
- Install pre-commit hooks `git config --local core.hooksPath .githooks/`.
- Install sqlc `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`.
- Install goose `go install github.com/pressly/goose/v3/cmd/goose@latest`.

### Run the application

Configure the environment variables in both the [backend](.env.example) and [frontend](ui/.env.example).
Migrate the database by starting the database `make db` and running the migrations `make migrate`.
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
