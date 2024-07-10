## How to:

Run docker 1st
```
docker compose up
```

Install swagger
```
go install github.com/swaggo/swag/cmd/swag@latest
```

Run swagger, since using gorm.DeletedAt
```
swag init --parseDependency --parseInternal  
```

Run lint
```
 golangci-lint run
```

Run go
```
go run main.go
```

Run unit test
```
go test ./... -v -cover
```