
all: build scan-state.svg

build:
	go build
	go test

test:
	go test

scan-state.svg: scan-state.dot
	dot -Tsvg scan-state.dot >scan-state.svg
