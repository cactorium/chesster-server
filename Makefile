all: dependencies server api test

server: api
	go build ./cmd/chessterd

PROTOBUFS := $(wildcard api/*.proto)
api: $(PROTOBUFS:.proto=.pb.go)

api/%.pb.go: api/%.proto
	protoc -I=api/ --go_out=api/ $<

dependencies:
	go get ./chesster

test:
	go test ./chesster

clean:
	rm chessterd
