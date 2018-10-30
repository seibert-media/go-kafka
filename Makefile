
prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/maxbrunsfeld/counterfeiter
	go get -u github.com/google/addlicense
	go get -u github.com/golang/dep/

precommit: ensure format generate test check
	@echo "ready to commit"

ensure:
	dep ensure

ginkgo:
	go get github.com/onsi/ginkgo/ginkgo
	ginkgo -r -race -cover

generate:
	go get github.com/maxbrunsfeld/counterfeiter
	find . -name mocks -type d | xargs rm -rf -
	go generate ./...

addlicense:
	go get github.com/google/addlicense
	addlicense -c "//SEIBERT/MEDIA GmbH" -y 2018 -l bsd ./schema/*.go ./consumer/*.go ./persistent/*.go

test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)

check: format lint vet errcheck

format:
	@go get golang.org/x/tools/cmd/goimports
	@find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	@find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

vet:
	@go vet $(shell go list ./... | grep -v /vendor/)

lint:
	@go get github.com/golang/lint/golint
	@golint -min_confidence 1 $(shell go list ./... | grep -v /vendor/)

errcheck:
	@go get github.com/kisielk/errcheck
	@errcheck -ignore '(Close|Write|Fprintf)' $(shell go list ./... | grep -v /vendor/)

