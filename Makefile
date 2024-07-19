.PHONY: build server scheduler test certs swagger docker-build

########## Variables ##########
SRC_DIR=./
DST_DIR=$(SRC_DIR)/internal
BUILD_VERSION=$(shell date +'%y%m%d%H%M%S')
LDFLAGS=-X 'main.version=$(VERSION)'

VERSION=1.0.6

DOCKER_USERNAME=erdemozgen
IMAGE_NAME=blackdagger
FULL_IMAGE_NAME=$(DOCKER_USERNAME)/$(IMAGE_NAME):$(VERSION)
LATEST_IMAGE_NAME=$(DOCKER_USERNAME)/$(IMAGE_NAME):latest
BUILDER_NAME=mybuilder

DOCKER_CMD := docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 --build-arg VERSION=$(VERSION) --push --no-cache

DEV_CERT_SUBJ_CA="/C=TR/ST=ASIA/L=TOKYO/O=DEV/OU=blackdagger/CN=*.blackdagger.dev/emailAddress=ca@dev.com"
DEV_CERT_SUBJ_SERVER="/C=TR/ST=ASIA/L=TOKYO/O=DEV/OU=SERVER/CN=*.server.dev/emailAddress=server@dev.com"
DEV_CERT_SUBJ_CLIENT="/C=TR/ST=ASIA/L=TOKYO/O=DEV/OU=CLIENT/CN=*.client.dev/emailAddress=client@dev.com"
DEV_CERT_SUBJ_ALT="subjectAltName=DNS:localhost"

########## Main Targets ##########
main:
	go run . server

watch:
	nodemon --watch . --ext go,gohtml --verbose --signal SIGINT --exec 'make server'

test:
	@go test --race ./...

test-clean:
	@go clean -testcache
	@go test --race ./...

install-tools: install-nodemon install-swagger

swagger: clean-swagger gen-swagger

certs: cert-dir gencerts-ca gencerts-server gencerts-client gencert-check

build: build-ui build-dir go-lint build-bin

build-image:
ifeq ($(VERSION),)
	$(error "VERSION is null")
endif
	$(DOCKER_CMD) -t erdemozgen/blackdagger:$(VERSION) .
	$(DOCKER_CMD) -t erdemozgen/blackdagger:latest .

server: go-lint build-dir build-bin
	./bin/blackdagger server

https-server:
	@blackdagger_CERT_FILE=./cert/server-cert.pem \
		blackdagger_KEY_FILE=./cert/server-key.pem \
		go run . server

scheduler: go-lint build-dir build-bin
	./bin/blackdagger scheduler

########## Tools ##########

build-bin:
	go build -ldflags="$(LDFLAGS)" -o ./bin/blackdagger .

build-dir:
	@mkdir -p ./bin

build-ui:
	@cd ui; \
		yarn && yarn build

	@rm -f ./service/frontend/assets/*.js
	@rm -f ./service/frontend/assets/*.woff
	@rm -f ./service/frontend/assets/*.woff2

	@cp ui/dist/*.js ./service/frontend/assets/
	@cp ui/dist/*.woff ./service/frontend/assets/
	@cp ui/dist/*.woff2 ./service/frontend/assets/

go-lint:
	@golangci-lint run ./...

cert-dir:
	@mkdir -p ./cert

gencerts-ca:
	@openssl req -x509 -newkey rsa:4096 \
		-nodes -days 365 -keyout cert/ca-key.pem \
		-out cert/ca-cert.pem \
		-subj "$(DEV_CERT_SUBJ_CA)"

gencerts-server:
	@openssl req -newkey rsa:4096 -nodes -keyout cert/server-key.pem \
		-out cert/server-req.pem \
		-subj "$(DEV_CERT_SUBJ_SERVER)"

	@openssl x509 -req -in cert/server-req.pem -CA cert/ca-cert.pem -CAkey cert/ca-key.pem \
		-CAcreateserial -out cert/server-cert.pem \
		-extfile cert/openssl.conf

gencerts-client:
	@openssl req -newkey rsa:4096 -nodes -keyout cert/client-key.pem \
		-out cert/client-req.pem \
		-subj "$(DEV_CERT_SUBJ_CLIENT)"

	@openssl x509 -req -in cert/client-req.pem -days 60 -CA cert/ca-cert.pem \
		-CAkey cert/ca-key.pem -CAcreateserial -out cert/client-cert.pem \
		-extfile cert/openssl.conf

gencert-check:
	@openssl x509 -in cert/server-cert.pem -noout -text

clean-swagger:
	@echo "Cleaning files"
	@rm -rf service/frontend/restapi/models
	@rm -rf service/frontend/restapi/operations

gen-swagger:
	@echo "Validating swagger yaml"
	@swagger validate ./swagger.yaml
	@echo "Generating swagger server code from yaml"
	@swagger generate server -t service/frontend --server-package=restapi --exclude-main -f ./swagger.yaml -A blackdagger
	@echo "Running go mod tidy"
	@go mod tidy

install-nodemon:
	npm install -g nodemon

install-swagger:
	brew tap go-swagger/go-swagger
	brew install go-swagger

docker-build:
	# Enable Docker BuildKit
	@export DOCKER_BUILDKIT=1
	# Create and use a new buildx builder instance
	-@docker buildx create --name $(BUILDER_NAME) --use
	# Login to Docker Hub
	@docker login
	# Build the Docker image for multiple architectures and push to Docker Hub
	@docker buildx build --platform linux/amd64,linux/arm64 -t $(FULL_IMAGE_NAME) --push .
	@docker buildx build --platform linux/amd64,linux/arm64 -t $(LATEST_IMAGE_NAME) --push .
	# Clean up the builder instance
	-@docker buildx rm $(BUILDER_NAME)