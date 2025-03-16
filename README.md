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

## Run the Go program

```bash
cd backend
go mod tidy
# run the app with options
go run . # to host the web version
go run . --cli --file ..\data\data.csv
go run . --cli --file ..\data\data.csv --filter "Age>=30"
go run . --cli --file ..\data\data.csv --filter "Age>=30" --sort "Age" --desc
go run . --cli --file ..\data\data.csv --interactive
```

## Compile the Go program

```bash
cd backend
go mod tidy
go build
# run the app with options
.\go-data-inspector
.\go-data-inspector --cli --file ..\data\data.csv
.\go-data-inspector --cli --file ..\data\data.csv --filter "Age>=30"
.\go-data-inspector --cli --file ..\data\data.csv --filter "Age>=30" --sort "Age" --desc
.\go-data-inspector --cli --file ..\data\data.csv --interactive
```
