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


## Other Notes

* Probably need to use GO's RPC library, probably most difficult

* Use HTTP for sending requests and setting up server 

* How do I load a file in?

* Maybe some mechanism to figure out how many coins you have before you send



