all: dependencies server api test

server: api
	go build ./cmd/chessterd

PROTOBUFS := $(wildcard api/protobuf/*.proto)
PROTOBUFS_PB := $(PROTOBUFS:.proto=.pb.go)
GENERATED := $(patsubst api/protobuf/%,api/%,$(PROTOBUFS_PB))
api: $(GENERATED)

api/%.pb.go: api/protobuf/%.proto
	protoc -I=api/protobuf/ --go_out=api/ $<

dependencies:
	go get ./chesster

test:
	go test ./chesster

clean:
	rm chessterd
