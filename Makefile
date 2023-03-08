SHELL=/bin/sh
include Makefile.*
TAG=${shell git rev-parse --short=5 HEAD}
AWS_USERNAME=AWS
REGION=eu-central-1
ACCOUNT=472641090641.dkr.ecr.$(REGION).amazonaws.com
SERVICE="go-template"
PROD_NAMESPACE=$(SERVICE)
DEV_NAMESPACE=$(SERVICE)-dev
ECR=$(ACCOUNT)/$(SERVICE)
ENV?=
LATEST=false

################################################################################
## Protos Go src & Doc generation
################################################################################
.PHONY: generate
generate:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Creating (gRPC & endpoints) src and docs from protos..."
	@echo "----------------------------------------------------------------"
	$(PROTOTOOL) generate

################################################################################
## Binary build & execution
################################################################################
.PHONY: install
install:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Installing protoc-gen-go..."
	@echo "----------------------------------------------------------------"
	$(GOINSTALL) github.com/golang/protobuf/protoc-gen-go
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Installing protoc-gen-grpc-gateway..."
	@echo "----------------------------------------------------------------"
	$(GOINSTALL) github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Installing protoc-gen-doc..."
	@echo "----------------------------------------------------------------"
	$(GOINSTALL) github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
	$(GOMOD) tidy

.PHONY: fmt
fmt:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Formatting code..."
	@echo "----------------------------------------------------------------"
	$(GO) fmt ./...
	$(GOIMPORTS) -w .
	$(GOGET) -d github.com/uber/prototool/cmd/prototool
	$(PROTOTOOL) format -w
	$(GOMOD) tidy

.PHONY: lint
lint:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Linting code..."
	@echo "----------------------------------------------------------------"
	$(GOLINT) run
	$(GOGET) -d github.com/uber/prototool/cmd/prototool
	$(PROTOTOOL) lint
	$(GOMOD) tidy

.PHONY: test
test:
	@echo "----------------------------------------------------------------"
	@echo " ✅  Testing code..."
	@echo "----------------------------------------------------------------"
	$(GO) test ./... -v -coverprofile=coverage.out

.PHONY: test-report
test-report:
	@echo "----------------------------------------------------------------"
	@echo " ✅  Testing code with report..."
	@echo "----------------------------------------------------------------"
	$(GO) get -u github.com/jstemmer/go-junit-report
	$(GO) test ./... -v -coverprofile=coverage.out 2>&1 | go-junit-report -set-exit-code > report.xml

.PHONY: coverage
coverage:
	@echo "----------------------------------------------------------------"
	@echo " 📊  Checking coverage..."
	@echo "----------------------------------------------------------------"
	$(GOTOOL) cover -html=coverage.out -o coverage.html
	$(GOTOOL) cover -func=coverage.out

.PHONY: godoc
godoc:
	@echo "----------------------------------------------------------------"
	@echo " 📄 Serving Go documentation..."
	@echo "----------------------------------------------------------------"
	@echo "Serving documentation at ~> http://localhost:9090"
	$(GODOC) -http=:9090 > /dev/null

.PHONY: compile
compile:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Compiling code..."
	@echo "----------------------------------------------------------------"
	$(GOBUILD) ./...
	$(PROTOTOOL) compile

.PHONY: deps
deps:
	@echo "----------------------------------------------------------------"
	@echo " ⬇️  Downloading dependencies..."
	@echo "----------------------------------------------------------------"
	$(GOGET) ./...

.PHONY: build
build: deps fmt
	@echo "----------------------------------------------------------------"
	@echo " 📦 Building binary..."
	@echo "----------------------------------------------------------------"
	$(GOBUILD) -a -ldflags "-w -X 'main.Version=$(TAG)'" -tags 'netgo osusergo' -o ethservice main.go

.PHONY: run
run: build
	@echo "----------------------------------------------------------------"
	@echo " ️🏃 Running..."
	@echo "----------------------------------------------------------------"
	./ethservice serve

.PHONY: all
all: generate lint build

################################################################################


################################################################################
## Docker commands
################################################################################
.PHONY: docker-login
docker-login:
	@echo "----------------------------------------------------------------"
	@echo " ️🏃 Login to $(ACCOUNT)..."
	@echo "----------------------------------------------------------------"
	aws ecr get-login-password --region $(REGION) | docker login --username $(AWS_USERNAME) --password-stdin $(ACCOUNT)

.PHONY: docker-build
docker-build: 
	@echo "----------------------------------------------------------------"
	@echo " 🔨 Building image $(TAG)..."
	@echo "----------------------------------------------------------------"
	docker pull $(ECR):latest || true
	docker build . --network=host --cache-from $(ECR):latest -t $(ECR):$(TAG) --build-arg	VERSION=$(TAG)

.PHONY: docker-push
docker-push: 
	@echo "----------------------------------------------------------------"
	@echo " ⬆️ Pushing image $(TAG)..."
	@echo "----------------------------------------------------------------"
	docker push $(ECR):$(TAG)
ifeq ($(LATEST), true)
	docker tag $(ECR):$(TAG) $(ECR):latest
	docker push $(ECR):latest
endif

.PHONY: docker-publish
docker-publish: docker-build docker-push

################################################################################
## Deploy commands
################################################################################

.PHONY: deploy
deploy:
	@echo "----------------------------------------------------------------"
	@echo " 🖥️ Deploying to $(ENV)..."
	@echo "----------------------------------------------------------------"
ifeq ($(ENV), prod)
	kubectl --namespace $(PROD_NAMESPACE) apply -f k8s/prod/main.yaml
else ifeq ($(ENV), dev)
		kubectl --namespace $(DEV_NAMESPACE) apply -f k8s/dev/main.yaml
else
	@echo "Invalid ENV"
	exit 1
endif
