# Peer Node

## Requirements

1) Clients uses gRPC protocol in GO to ask for address of file

2) User sends request to server for file

3) User storing file will then send the file back

4) User will then record the transaction by sending a message to the blockchain
 

## Assumptions

1) Each consumer/producer has their own IP address

2) Producer sets up local HTTP server

3) Consumer can fetch document from producer's local HTTP server


## Running

First generate the gRPC files for GO. Make sure you are in the root of the project and run the command below.

``` bash

$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    fileshare/file_share.proto 
```

GO Version: 1.21.4

```bash

$ go build

$ ./peer-node

```

## gRPC API

* RecordFileRequestTransaction: Tell blockchain of a completed transaction

* PlaceFileRequest: Ask market to tell you ALL possible locations where file is store 

* NotifyFileStore Tell market you will store a file for future access

* NotifyFileUnstore: Tell market you no longer have a specific file

* SendFile: Send a File

## HTTP Functionality

Server should only look for two things:

* Route /requestFile with a GET Request, parameter of `filename`, a string that represents name of file

* Route /storeFile with a GET Request, similar to the route /requestFile


## CLI interface

Requesting a file:

```bash
$ get [ip] [port] [filename]
```

Storing a file:

```bash
$ store [ip] [address] [filename]
```

Import a file:

```bash
$ import [filepath]
```

Listing all files stored for IPFS

```bash
$ list
```

Getting current peer node location

```bash
$ location
```

Exiting Program

```bash
$ exit
```




