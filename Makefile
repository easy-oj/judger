all: build

build:
	go build -a -v -o output/judger

fmt:
	find ./ -name "*.go" | grep -v "/vendor/" | xargs gofmt -w

clean:
	rm -rf output
