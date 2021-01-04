# Simplest list

An extremelly simple todo list app, with random task pick.

## Backend
The backend is built using Go.

## Build
To build the server run:

```sh
go build main.go
```

## Run
You can directly run the server from the code with:
```sh
go run main.go
```

or if you built it you can just run the executable:
```sh
./main
```

### Parameters
There are a few configuration parameters you can pass to the server when running it:
* --port: the port to run the server at (default: 8080)
* --cert: the path to the TLS certificate file (optional)
* --pk:   the path to the TLS private key (optional, mandatory if cert is defined)

### HTTPS
To run the server through HTTPS you need to pass in the `cert` and `pk` parameters.