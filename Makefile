.PHONY: run build-lambda deploy test clean

run:
	go run ./cmd/local --output output.jpg

build-lambda:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bootstrap ./cmd/lambda

deploy: build-lambda
	lambroll deploy

test:
	go test ./...

clean:
	rm -f output.jpg bootstrap function.zip
