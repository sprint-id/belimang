# Beli Mang

Requirement: [Project Beli Mang](https://openidea-projectsprint.notion.site/BeliMang-7979300c7ce54dbf8ecd0088806eff14)


## Database Migration

Database migration must use [golang-migrate](https://github.com/golang-migrate/migrate) as a tool to manage database migration

- **Short Tutorial:**
    - Direct your terminal to your project folder first
    - Initiate folder
        
        ```bash
        mkdir db/migrations
        
        ```
        
    - Create migration
        
        ```bash
        migrate create -ext sql -dir db/migrations add_user_table
        
        ```
        
        This command will create two new files named `add_user_table.up.sql` and `add_user_table.down.sql` inside the `db/migrations` folder
        
        - `.up.sql` can be filled with database queries to create / delete / change the table
        - `.down.sql` can be filled with database queries to perform a `rollback` or return to the state before the table from `.up.sql` was created
    - Execute migration
        
        ```bash
        migrate -database "postgres://postgres:password@host:5432/postgres?sslmode=disable" -path ./db/migrations -verbose up
        
        ```
        
    - Rollback migration (one migration)
        
        ```bash
        migrate -database "postgres://postgres:password@host:5432/postgres?sslmode=disable" -path db/migrations -verbose down
        
        ```
        
    - Rollback migration (all migration)
        
        ```bash
        migrate -database "postgres://postgres:password@host:5432/postgres?sslmode=disable" -path db/migrations -verbose drop
        ```


## Run & Build Beli Mang

Run for debugging

```
make run
```

Build app

```
make build
```



