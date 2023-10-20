🚧 sandbox for a seaport rfq node

implementing interfaces based on [valorem's rfq-over-seaport](https://github.com/valorem-labs-inc/trade-interfaces)

using [connect-rpc](https://connectrpc.com/) & buf for proto tooling:

* [schemas over here](https://buf.build/sambarnes/seaport/file/main:proto/seaport/v1/rfq.proto)
* [interactive postman-like studio](https://buf.build/studio/sambarnes/seaport/main?serviceDialog=open)

## development

get the deps: `go get`

install some tooling: `make install-dev-tooling`

run the server: `go run main.go`

unary endpoints can be used thru grpc or plain json:
```shell
$ curl --header "Content-Type: application/json" --data '{}' localhost:8090/proto.seaport.v1.AuthService/Nonce

{"nonce":"0xdadb0d"}
```

will need grpcurl or similar to tail a stream:
```shell
$ grpcurl -protoset <(buf build -o -) -plaintext -d '{}' localhost:8090 proto.seaport.v1.RFQService/WebTaker

# TODO: actually do this example
```
