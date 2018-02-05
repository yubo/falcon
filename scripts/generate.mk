.PHONY: clean parse doc deploy run

PBFILES=transfer/transfer.pb.go transfer/transfer.pb.gw.go \
	lib/tsdb/tsdb.pb.go \
	service/service.pb.go service/service.pb.gw.go \
	alarm/alarm.pb.go alarm/alarm.pb.gw.go 

DOCFILES=transfer/transfer.swagger.json \
	service/service.swagger.json \
	agent/agent.swagger.json \
	alarm/alarm.swagger.json

SUBMODULES=ctrl agent transfer service alarm
YACCFILES=parse/parse.go $(SUBMODULES:%=%/config/parse.go)

all: $(PBFILES) $(DOCFILES) $(YACCFILES)

%.go: %.y
	goyacc -o $@ $<

%.pb.go: %.proto
	protoc -I/usr/local/include -I. \
	  -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --go_out=plugins=grpc:$(GOPATH)/src \
	  $<


%.pb.gw.go: %.proto
	protoc -I/usr/local/include -I. \
	  -I$(GOPATH)/src \
	  -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --grpc-gateway_out=logtostderr=true:. \
	  $<

%.swagger.json: %.proto
	protoc -I/usr/local/include -I. \
	  -I$(GOPATH)/src \
	  -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --swagger_out=logtostderr=true:. \
	  $<
