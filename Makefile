server:
	go build ./cmd/chessterd

test:
	go test ./chesster

clean:
	rm chessterd
