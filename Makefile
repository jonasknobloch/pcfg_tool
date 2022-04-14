BINARY_NAME=pcfg_tool

build:
	go build -o ${BINARY_NAME} main.go

	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-amd64-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-amd64-linux main.go

run:
	go run main.go

clean:
	go clean
	rm ${BINARY_NAME}

deps:
	go get github.com/jonasknobloch/jinn@v0.6.0
	go get github.com/spf13/cobra@v1.4.0