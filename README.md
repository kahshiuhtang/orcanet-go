# Peer Node

## Basic Functionality

* Retrieving
    1) Find addresses corresponding with a specific file hash inside the DHT
    2) Connect through HTTP with the cheapest option
    3) Recieve chunks from the HTTP server
        1) If the file is small (< 4KB), the entire transaction is done in one go
        2) Otherwise, files are sent in chunks
            1) A signed transaction must be sent to the sender
    4) Close connection
    5) Add file to imported folder inside files directory

* Storing
    1) Hash the file content that you are trying to store
    2) Send a request to store the file in to the DHT
    3) Make sure file is inside the stored folder inside files directory
    4) Accept/Decline any requests for said file, or set the default behavior
    5) Send all microtransactions to blockchain once everything is done


## Assumptions

1) Each consumer/producer has their own IP address

2) Producer sets up local HTTP server

3) Consumer can fetch document from producer's local HTTP server


## Running

First generate the gRPC files for GO. Make sure you are in the root of the project and run the command below.

``` bash

$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/fileshare/file_share.proto 
```

GO Version: 1.21.4

```bash

$ make all

```

## CLI interface

Requesting a file:

```bash
$ get [ip] [port] [filename]
```

Storing a file:

```bash
$ store [ip] [address] [filename]
```

Putting a key inside DHT

```bash

$ putKey [key] [value]

```

Getting a key inside DHT

```bash

$ getKey [key]

```

Import a file:

```bash
$ import [filepath]
```

Complete pipeline for getting a file from DHT:

```bash

$ fileGet [fileHash] 

```

Send a certain amount of coin to an address

```bash

$ send [amount] [ip] [amount]

```

Hash a file

```bash

$ hash [filename]

```

Listing all files stored for IPFS

```bash
$ list
```

Getting current peer node location

```bash
$ location
```

Testing network speeds

```bash
$ network
```

Exiting Program

```bash
$ exit
```

#### File System:

* There is a folder called <i>files</i>. This is where all the files that are available to the user is stored

* Any file directly stored inside <i>files</i> folder is considered <i>uploaded</i> to the client.

* Any file that has been requested by the user is stored in the <i>files/requested</i> folder.

* Any file that is available to be requested for by anyone on the network is in <i>files/stored</i>.

* Technically, you can import the files manually if you drag them inside the desired folder. There is currently no protection against this.

* The <i>transactions</i> folder stores all of the transactions that have been processed and stored.

#### Notes:

* Files that are on the network should be in the files folder. This can be done manually or by using the CLI

* Inside the config file, set your public key and private key location. If you don't want to, the CLI will generate a key-pair for you.

* Only .txt, .json and .mp4 file formats are currently supported.


## HTTP Functionality

Server should only look for two things:

* Route /requestFile/ with a GET Request, parameter of `filename`, a string that represents name of file

* Route /storeFile/ with a GET Request, similar to the route /requestFile

* Route /sendTransaction with a POST Request, must send the transaction and a signed version of the transaction


## gRPC protocol

Currently in a state of flux, will be update when anything changes





