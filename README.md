# go-data-inspector

## Initial Setup

```bash
cd backend
go mod init github.com/ReadableCode/go-data-inspector
```

## Create the main.go file

```bash
cd backend
touch main.go
```

## Run the main.go file

```bash
cd backend
go mod tidy
go run main.go --file ..\data\data.csv
go run main.go --file ..\data\data.csv --filter "Age>=30"
go run main.go --file ..\data\data.csv --filter "Age>=30" --sort "Age" --desc
go run main.go --file ..\data\data.csv --interactive
```
