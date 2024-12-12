# StocksBack

## What is stocks back?

an api for stocks usage

### User functionality

- Creating accounts `/signup`
- Sign in using password or secret key (in the future jwt) 
- Farming (getting solids based on the stock amount) every hour `/farm`
- Getting solids from stocks every day at 21 (server time)
- Buying stocks `/buy`
- Changing name and password `/change/name`, `/change/password`
- Getting user `/get`

### Backend functionality 

- Two database types, that could be easy changed
    1. [using file system (for testing)](/pkg/file_db/main.go)
    2. [PostgreSQL database](/pkg/db/main.go)
- [Graceful shutdown](/pkg/closer/main.go)
- [Custom cron usage](/pkg/cron/main.go)
- [Hasher](/pkg/hash/hash.go)
- [Custom logger](/pkg/logger/main.go)
- [Specific query expressions that can be used in both database types](/pkg/query/query.go)
- [User service for all user activities](/pkg/user_service/main.go)
- [Http handler and server](/http/server/)
- Good structured headers, requests, responses
    1. [Headers](/config/headers/headers.go)
    2. [Requests](/config/requests/requests.go)
    3. [Responses](/config/responses/responses.go)
- Timeout server and service mode

## Setup program 

### Getting project

```bash
git clone https://github.com/VandiKond/StocksBack.git
```

### Config setup

[Config example](/config/config.yaml)

Replace it with your data

### Running program

```bash
go run cmd/main.go
```

You can edit [main file](/cmd/main.go)

Examples

1. With timeout
    ```go
    package main

    import (
        "context"
        "os/signal"
        "syscall"
        "time"

        "github.com/VandiKond/StocksBack/internal/application"
        "github.com/VandiKond/StocksBack/pkg/db"
    )

    func main() {
        // Creating a new application with a hour timeout
        app := application.New(time.Hour, db.Constructor{})

        // Adding graceful  shutdown
        ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
        defer stop()

        // Running the app
        app.Run(ctx)
    }
    ```
2. Service mode
    ```go 
    package main

    import (
        "context"
        "os/signal"
        "syscall"

        "github.com/VandiKond/StocksBack/internal/application"
        "github.com/VandiKond/StocksBack/pkg/db"
    )

    func main() {
        // Creating a new service application
        app := application.NewService(db.Constructor{})

        // Adding graceful  shutdown
        ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
        defer stop()

        // Running the app
        app.Run(ctx)
    }
    ```

## License 

[LICENSE](LICENSE)


