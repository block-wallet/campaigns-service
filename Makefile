TAG=${shell git rev-parse --short=7 HEAD}
REGION=eu-central-1
ACCOUNT=472641090641.dkr.ecr.$(REGION).amazonaws.com
SERVICE=campaigns-service
ECR=$(ACCOUNT)/$(SERVICE)
GH_ACCESS_TOKEN?=
TAG_INTEGRATION_TESTS=$(TAG)-integration-tests
AWS_USERNAME=AWS

# GitHub commands
.PHONY: gh
gh:
	@[ "${GH_ACCESS_TOKEN}" ] && echo "GitHub token found" || ( echo "GitHub token is not set"; exit 1 )


#######################################################
############ formats, lint, test and build ############
#######################################################

.PHONY: fmt
fmt:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Formatting code..."
	@echo "----------------------------------------------------------------"
	gofmt -s -w ./.

.PHONY: lint
lint:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Linting code..."
	@echo "----------------------------------------------------------------"
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2 run -E gofmt --out-format=github-actions --timeout=10m ./...

.PHONY: test
test:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Testign the code..."
	@echo "----------------------------------------------------------------"
	go test ./...

.PHONY: generate
generate:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Generating GRPC code..."
	@echo "----------------------------------------------------------------"
	buf generate

.PHONY: mocks
mocks:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Generating Mocks..."
	@echo "----------------------------------------------------------------"
	mockery --dir=domain/campaigns-service --output=domain/campaigns-service/mocks --outpkg=campaignsservicemocks --all

.PHONY: build
build:
	@echo "----------------------------------------------------------------"
	@echo " :package: Building binaries..."
	@echo "----------------------------------------------------------------"
	go build -a -ldflags "-w -X 'main.Version=${shell git rev-parse --short=7 HEAD}'" -o build/campaignsservice main.go


################################################
############ docker builds and push ############
################################################

.PHONY: docker-login
docker-login:
	@echo "----------------------------------------------------------------"
	@echo " ️🏃 Login to $(ACCOUNT)..."
	@echo "----------------------------------------------------------------"
	aws ecr get-login-password --region $(REGION) | docker login --username $(AWS_USERNAME) --password-stdin $(ACCOUNT)

.PHONY: docker-build/service
docker-build/service:
	@echo "----------------------------------------------------------------"
	@echo " ⚙️  Building the service Docker image..."
	@echo "----------------------------------------------------------------"
	docker build . -f Dockerfile --no-cache --network=host -t $(ECR):$(TAG) --build-arg VERSION=$(TAG) --build-arg GH_ACCESS_TOKEN=$(GH_ACCESS_TOKEN)

.PHONY: docker-push/service
docker-push/service:
	@echo "----------------------------------------------------------------"
	@echo " ⬆️ Pushing service image $(TAG)..."
	@echo "----------------------------------------------------------------"
	docker push $(ECR):$(TAG)


.PHONY: docker-publish/integration-tests
docker-publish/integration-tests: docker-build/integration-tests docker-push/integration-tests


.PHONY: docker-publish/service
docker-publish/service: docker-build/service docker-push/service

.PHONY: aws/login
aws/login:
	aws ecr get-login-password --region eu-central-1 | docker login --username AWS --password-stdin $(ACCOUNT)


#######################################################
###################### local run ######################
#######################################################

.PHONY: run
run: build
	@echo "----------------------------------------------------------------"
	@echo " ️🏃 Running..."
	@echo " :package: Building binaries..."
	@echo "----------------------------------------------------------------"
	./build/campaignsservice serve

.PHONY: dev/run
dev/run:
	@echo "----------------------------------------------------------------"
	@echo " Building and running the local dev environment using docker-compose..."
	@echo "----------------------------------------------------------------"
	docker-compose -f docker-compose.yml up --build

.PHONY: dev/up
dev/up:
	@echo "----------------------------------------------------------------"
	@echo " Running local dev environment using docker-compose..."
	@echo "----------------------------------------------------------------"
	docker-compose -f docker-compose.yml up

.PHONY: dev/down
dev/down:
	@echo "----------------------------------------------------------------"
	@echo " Stopping local dev environment..."
	@echo "----------------------------------------------------------------"
	docker-compose -f docker-compose.yml down

.PHONY: service/up
service/up:
	@echo "----------------------------------------------------------------"
	@echo " Running local campaigns-service using docker..."
	@echo "----------------------------------------------------------------"
	docker-compose -f docker-compose.yml up campaigns-service

.PHONY: service/down
service/down:
	@echo "----------------------------------------------------------------"
	@echo " Stopping local campaigns-service..."
	@echo "----------------------------------------------------------------"
	docker-compose -f docker-compose.yml down campaigns-service

.PHONY: db/up
db/up:
	@echo "----------------------------------------------------------------"
	@echo " Running local postgres database..."
	@echo "----------------------------------------------------------------"
	docker-compose -f docker-compose.yml up -d postgres-database

.PHONY: db/down
db/down:
	@echo "----------------------------------------------------------------"
	@echo " Stopping local postgres database..."
	@echo "----------------------------------------------------------------"
	docker stop postgres-campaigns-service