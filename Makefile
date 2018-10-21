
all: generate test format addlicense

test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)

ginkgo:
	go get github.com/onsi/ginkgo/ginkgo
	ginkgo -r -race -cover

format:
	go get golang.org/x/tools/cmd/goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

generate:
	go get github.com/maxbrunsfeld/counterfeiter
	go generate ./...

addlicense:
	go get github.com/google/addlicense
	addlicense -c "//SEIBERT/MEDIA GmbH" -y 2018 -l bsd ./schema/*.go ./consumer/*.go ./persistent/*.go

prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/maxbrunsfeld/counterfeiter
	go get -u github.com/google/addlicense
	go get -u github.com/golang/dep/
