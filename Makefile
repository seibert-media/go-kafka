
run:
	docker network create kafka || echo 'network already exists'
	docker-compose up -d
	docker-compose logs -f

ksqlcli:
	docker run -ti \
	--net=kafka \
	confluentinc/cp-ksql-cli:5.1.0 \
	http://ksql-server:8088

prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/maxbrunsfeld/counterfeiter
	go get -u github.com/google/addlicense
	go get -u github.com/golang/dep/

ginkgo:
	go get github.com/onsi/ginkgo/ginkgo
	ginkgo -r -race -cover

precommit: ensure generate format test check addlicense
	@echo "ready to commit"

ensure:
	dep ensure

generate:
	go get github.com/maxbrunsfeld/counterfeiter
	find . -name mocks -type d | xargs rm -rf -
	go generate ./...

format:
	@go get golang.org/x/tools/cmd/goimports
	@find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	@find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)

check: lint vet errcheck

lint:
	@go get github.com/golang/lint/golint
	@golint -min_confidence 1 $(shell go list ./... | grep -v /vendor/)

vet:
	@go vet $(shell go list ./... | grep -v /vendor/)

errcheck:
	@go get github.com/kisielk/errcheck
	@errcheck -ignore '(Close|Write|Fprint)' $(shell go list ./... | grep -v /vendor/)

addlicense:
	go get github.com/google/addlicense
	addlicense -c "//SEIBERT/MEDIA GmbH" -y 2018 -l bsd ./schema/*.go ./consumer/*.go ./persistent/*.go
