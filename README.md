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

## Run the Go program without compiling it separately

```bash
cd backend
go mod tidy
# run the web app
go run .
# run the cli app
go run . --cli --file ..\data\data.csv
go run . --cli --file ..\data\data.csv --filter "Age>=30"
go run . --cli --file ..\data\data.csv --filter "Age>=30" --sort "Age" --desc
go run . --cli --file ..\data\data.csv --interactive
```

## Compile the Go program and run it

```bash
cd backend
go mod tidy
go build
# run the web app on windows
.\go-data-inspector.exe
# run the web app on linux
chmod +x go-data-inspector
./go-data-inspector
# run the cli app on windows
.\go-data-inspector --cli --file ..\data\data.csv
.\go-data-inspector --cli --file ..\data\data.csv --filter "Age>=30"
.\go-data-inspector --cli --file ..\data\data.csv --filter "Age>=30" --sort "Age" --desc
.\go-data-inspector --cli --file ..\data\data.csv --interactive
# run the cli app on linux
chmod +x go-data-inspector
./go-data-inspector --cli --file ../data/data.csv
./go-data-inspector --cli --file ../data/data.csv --filter "Age>=30"
./go-data-inspector --cli --file ../data/data.csv --filter "Age>=30" --sort "Age" --desc
./go-data-inspector --cli --file ../data/data.csv --interactive
```
