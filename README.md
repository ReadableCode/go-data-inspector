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
# run the app with options
go run main.go # to host the web version
go run main.go --cli --file ..\data\data.csv
go run main.go --cli --file ..\data\data.csv --filter "Age>=30"
go run main.go --cli --file ..\data\data.csv --filter "Age>=30" --sort "Age" --desc
go run main.go --cli --file ..\data\data.csv --interactive
```
