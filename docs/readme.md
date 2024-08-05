# polyrepo.pro API

This is the API for the polyrepo.pro project.

## Commands

| Command                              | Description                            |
| ------------------------------------ | -------------------------------------- |
| [make db/generate](#make-dbgenerate) | Generates the Prisma client            |
| [make db/migrate](#make-dbmigrate)   | Migrates the database                  |
| [make db/spy](#make-dbspy)           | Generate schema documentation and ERD. |
| [make run](#make-run)                | Runs the API                           |

### `make db/generate`

Generates the Prisma client library from the [prisma/schema.prisma](prisma/schema.prisma) file and outputs it to [prisma/db](prisma/db).

### `make db/migrate`

Performs database migrations to the database specified in the [.prisma.env](.prisma.env) file.

### `make run`

Runs the API.

### `make db/spy`

Generates schema documentation and ERD.

This will output to the [tmp/schemaspy](tmp/schemaspy) directory. You can then open the index.html file in your browser to view the documentation.

![alt text](<Google Chrome-000020@2x.png>)
