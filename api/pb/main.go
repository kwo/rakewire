package pb

//go:generate protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/gengo/grpc-gateway/third_party/googleapis --go_out=Mgoogle/api/annotations.proto=github.com/gengo/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. --grpc-gateway_out=logtostderr=true:. api.proto
