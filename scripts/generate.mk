.PHONY: clean parse doc deploy run

PBFILES=modules/transfer/transfer.pb.go \
	modules/transfer/transfer.pb.gw.go \
	modules/service/service.pb.go \
	modules/service/service.pb.gw.go \
	modules/alarm/alarm.pb.go \
	modules/alarm/alarm.pb.gw.go \
	lib/tsdb/tsdb.pb.go 

DOCFILES=modules/transfer/transfer.swagger.json \
	modules/service/service.swagger.json \
	modules/agent/agent.swagger.json \
	modules/alarm/alarm.swagger.json

SUBMODULES=ctrl agent transfer service alarm
#YACCFILES=parse/parse.go $(SUBMODULES:%=%/config/parse.go)

all: $(PBFILES) $(DOCFILES) modules/alarm/alarm.pb.gw.go

#%.go: %.y
#	goyacc -o $@ $<

%.pb.go: %.proto
	protoc -I/usr/local/include -I. \
	  -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --go_out=plugins=grpc:$(GOPATH)/src \
	  --grpc-gateway_out=logtostderr=true:$(GOPATH)/src \
	  --swagger_out=logtostderr=true:$(GOPATH)/src \
	  $<


%.pb.gw.go: %.proto
	protoc -I/usr/local/include -I. \
	  -I$(GOPATH)/src \
	  -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --grpc-gateway_out=logtostderr=true:$(GOPATH)/src \
	  $<

%.swagger.json: %.proto
	protoc -I/usr/local/include -I. \
	  -I$(GOPATH)/src \
	  -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --swagger_out=logtostderr=true:. \
	  $<
