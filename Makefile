.PHONY: all
all: build test

.PHONY: protos
protos: model/limits.pb.go

%.pb.go: %.proto
	protoc --proto_path=$(@D) --proto_path=./vendor --go_out=plugins=grpc:$(@D) $<

.PHONY: dep
dep:
	dep ensure --vendor-only

.PHONY: dep-up
dep-up:
	dep ensure

.PHONY: checkfmt
checkfmt:
	find . ! \( -path ./vendor -prune \) ! \( -path ./.git -prune \) -name '*.go' -exec gofmt -l {} +

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o srv cmd/server/main.go
	go build -o client cmd/client/main.go

.PHONY: build
run: build
	./srv
