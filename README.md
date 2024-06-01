To deploy, just run `docker-compose up`.

To run locally, copy-paste file `cfg/dev.go.example` with the name `cfg/dev.go` and change the config value appropriately.

The whole application features a simple CRUD for `Order` entity and only that. Seed the database by running the
statements in `db/seed.sql` file.
