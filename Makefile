.PHONY: run build-lambda deploy test clean fw-build fw-upload fw-monitor fw-clean

run:
	go run ./cmd/local --output output.jpg

build-lambda:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bootstrap ./cmd/lambda

deploy: build-lambda
	lambroll deploy

test:
	go test ./...

clean:
	rm -f output.jpg output_dark.jpg bootstrap function.zip

fw-build:
	cd firmware && pio run

fw-upload:
	cd firmware && pio run --target upload

fw-monitor:
	cd firmware && pio device monitor

fw-clean:
	cd firmware && pio run --target clean
