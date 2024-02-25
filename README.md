# Producer/Consumer Nodes

## Requirements

1) Clients uses RPC protocol in GO to send request

2) User sends request to server for file

3) User storing file will then send the file back

4) User sending the request will then make a request to send tokens/coins
 

## Assumptions

1) Each consumer/producer has their own IP address

2) Producer sets up local HTTP server

3) Consumer can fetch document from producer's local HTTP server


## Running

Generating gRPC files:

``` bash

$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    fileshare/file_share.proto 
```

GO Version: 1.21.4

```bash

$ go build

$ ./client

```

## gRPC API
* RequestFileStoreIP: Give file, file size, price, get back IP address

* StoreFile: Give file, 

* RequestFileIP [Consumer to Market]: Give file name, get back IP address

* RequestFile [Consumer to Producer]: Give file name, get back file

* Store Transaction [Producer to Market]

## HTTP Functionality

Server should only look for two things:

* Route /requestFile with a GET Request, parameter of `filename`, a string that represents name of file

 Demo below. The client requests a file in the server's local directory.

 


https://github.com/kahshiuhtang/PeerNodes/assets/78182536/321af8e3-0a5f-4731-9544-f67b3c4418e8


* Route /sendCoin with a POST Request, data contains a JSON object with one field named `amount`

## Other Notes

* Probably need to use GO's RPC library, probably most difficult

* Use HTTP for sending requests and setting up server 

* How do I load a file in?

* Maybe some mechanism to figure out how many coins you have before you send




