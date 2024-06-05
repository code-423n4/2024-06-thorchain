########################################################################################
# Environment Checks
########################################################################################

CHECK_ENV:=$(shell ./scripts/check-env.sh)
ifneq ($(CHECK_ENV),)
$(error Check environment dependencies.)
endif

########################################################################################
# Config
########################################################################################

.PHONY: build test tools export healthcheck run-mocknet build-mocknet stop-mocknet halt-mocknet ps-mocknet reset-mocknet logs-mocknet openapi

# pull branch name from CI if unset and available
ifdef CI_COMMIT_BRANCH
BRANCH?=${CI_COMMIT_BRANCH}
BUILDTAG?=${CI_COMMIT_BRANCH}
endif

# image build settings
COMMIT?=$(shell git log -1 --format='%H' 2>/dev/null)
BRANCH?=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GITREF?=$(shell git rev-parse --short HEAD 2>/dev/null)
BUILDTAG?=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)

# compiler flags
VERSION:=$(shell cat version)
TAG?=mocknet
ldflags = -X gitlab.com/thorchain/thornode/constants.Version=$(VERSION) \
      -X gitlab.com/thorchain/thornode/constants.GitCommit=$(COMMIT) \
      -X github.com/cosmos/cosmos-sdk/version.Name=THORChain \
      -X github.com/cosmos/cosmos-sdk/version.AppName=thornode \
      -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
      -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
      -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(TAG) \
      -buildid=

# golang settings
TEST_DIR?="./..."
BUILD_FLAGS := -ldflags '$(ldflags)' -tags ${TAG} -trimpath
TEST_BUILD_FLAGS := -parallel=1 -tags=mocknet
GOBIN?=${GOPATH}/bin
BINARIES=./cmd/thornode ./cmd/bifrost ./tools/generate

# docker tty args are disabled in CI
ifndef CI
DOCKER_TTY_ARGS=-it
endif

########################################################################################
# Targets
########################################################################################

# ------------------------------ Generate ------------------------------

generate: go-generate openapi protob-docker

go-generate:
	@go install golang.org/x/tools/cmd/stringer@v0.15.0
	@go generate ./...


protob:
	@./scripts/protocgen.sh

protob-docker:
	@docker run --rm -v $(shell pwd):/app -w /app \
		registry.gitlab.com/thorchain/thornode:builder-v5@sha256:33a2345924750d45a7e9060f429f66ddc1888482e6ef73e2358a23024dc41f01 \
		make protob

openapi:
	@docker run --rm \
		--user $(shell id -u):$(shell id -g) \
		-v $$PWD/openapi:/mnt \
		openapitools/openapi-generator-cli:v6.0.0@sha256:310bd0353c11863c0e51e5cb46035c9e0778d4b9c6fe6a7fc8307b3b41997a35 \
		generate -i /mnt/openapi.yaml -g go -o /mnt/gen
	@rm openapi/gen/go.mod openapi/gen/go.sum
	@find ./openapi/gen -type f | xargs sed -i '/^[- ]*API version.*$(shell cat version)/d;/APIClient.*$(shell cat version)/d'
	@find ./openapi/gen -type f | grep model | xargs sed -i 's/MarshalJSON(/MarshalJSON_deprecated(/'

# ------------------------------ Docs ------------------------------

docs-init:
	@cargo install mdbook --version 0.4.37
	@cargo install mdbook-admonish --version 1.14.0
	@cargo install mdbook-catppuccin --version 2.1.0
	@cargo install mdbook-katex --version 0.5.10
	@cargo install mdbook-embed --version 0.2.0
	@cd docs && mdbook-catppuccin install
	@cd docs && mdbook-admonish install --css-dir theme

docs-generate: docs-init
	@cd docs && mdbook build -d ../public

docs-dev: docs-init
	@cd docs && mdbook serve -d ../public --open

# ------------------------------ Build ------------------------------

build:
	go build ${BUILD_FLAGS} ${BINARIES}

install:
	go install ${BUILD_FLAGS} ${BINARIES}

tools:
	go install -tags ${TAG} ./tools/generate
	go install -tags ${TAG} ./tools/pubkey2address
	go install -tags ${TAG} ./tools/p2p-check

# ------------------------------ Gitlab CI ------------------------------

gitlab-trigger-ci:
	@./scripts/gitlab-trigger-ci.sh

# ------------------------------ Housekeeping ------------------------------

format:
	@git ls-files '*.go' | grep -v -e '^docs/' | xargs gofumpt -w

lint:
	@./scripts/lint.sh
	@./scripts/trunk check --no-fix --upstream origin/develop

lint-ci:
	@./scripts/lint.sh
	@./scripts/trunk-ci.sh

# ------------------------------ Unit Tests ------------------------------

test-coverage: test-network-specific
	@go test ${TEST_BUILD_FLAGS} -v -coverprofile=coverage.txt -covermode count ${TEST_DIR}
	sed -i '/\.pb\.go:/d' coverage.txt

coverage-report: test-coverage
	@go tool cover -html=coverage.txt

test-coverage-sum: test-network-specific
	@go run gotest.tools/gotestsum --junitfile report.xml --format testname -- ${TEST_BUILD_FLAGS} -v -coverprofile=coverage.txt -covermode count ${TEST_DIR}
	sed -i '/\.pb\.go:/d' coverage.txt
	@GOFLAGS='${TEST_BUILD_FLAGS}' go run github.com/boumenot/gocover-cobertura < coverage.txt > coverage.xml
	@go tool cover -func=coverage.txt
	@go tool cover -html=coverage.txt -o coverage.html

test: test-network-specific
	@CGO_ENABLED=0 go test ${TEST_BUILD_FLAGS} ${TEST_DIR}

test-network-specific:
	@CGO_ENABLED=0 go test -tags stagenet ./common
	@CGO_ENABLED=0 go test -tags mainnet ./common ./bifrost/pkg/chainclients/utxo/...
	@CGO_ENABLED=0 go test -tags mocknet ./common ./bifrost/pkg/chainclients/utxo/...

test-race:
	@go test -race ${TEST_BUILD_FLAGS} ${TEST_DIR}

# ------------------------------ Regression Tests ------------------------------

test-regression: build-test-regression
	@docker run --rm ${DOCKER_TTY_ARGS} \
		-e DEBUG -e RUN -e EXPORT -e TIME_FACTOR -e PARALLELISM -e FAIL_FAST \
		-e UID=$(shell id -u) -e GID=$(shell id -g) \
		-p 1317:1317 -p 26657:26657 \
		-v $(shell pwd)/test/regression/mnt:/mnt \
		-v $(shell pwd)/test/regression/suites:/app/test/regression/suites \
		-v $(shell pwd)/test/regression/templates:/app/test/regression/templates \
		-w /app thornode-regtest sh -c 'make _test-regression'

build-test-regression:
	@DOCKER_BUILDKIT=1 docker build . \
		-t thornode-regtest \
		-f ci/Dockerfile.regtest \
		--build-arg COMMIT=$(COMMIT)

test-regression-coverage:
	@go tool cover -html=test/regression/mnt/coverage/coverage.txt

# internal target used in docker build
_build-test-regression:
	@go install -ldflags '$(ldflags)' -tags=mocknet,regtest ./cmd/thornode
	@go build -ldflags '$(ldflags)' -cover -tags=mocknet,regtest -o /regtest/cover-thornode ./cmd/thornode
	@go build -ldflags '$(ldflags)' -tags mocknet -o /regtest/regtest ./test/regression/cmd

# internal target used in test run
_test-regression:
	@rm -rf /mnt/coverage && mkdir -p /mnt/coverage
	@cd test/regression && /regtest/regtest
	@go tool covdata textfmt -i /mnt/coverage -o /mnt/coverage/coverage.txt
	@grep -v -E -e archive.go -e 'v[0-9]+.go' -e openapi/gen /mnt/coverage/coverage.txt > /mnt/coverage/coverage-filtered.txt
	@go tool cover -func /mnt/coverage/coverage-filtered.txt > /mnt/coverage/func-coverage.txt
	@awk '/^total:/ {print "Regression Coverage: " $$3}' /mnt/coverage/func-coverage.txt
	@chown -R ${UID}:${GID} /mnt

# ------------------------------ Simulation Tests ------------------------------

test-simulation: build-test-simulation build-mocknet reset-mocknet
	@docker run --rm ${DOCKER_TTY_ARGS} \
		-e PARALLELISM --network host -w /app \
		thornode-simtest sh -c 'make _test-simulation'

build-test-simulation:
	@DOCKER_BUILDKIT=1 docker build . \
		-t thornode-simtest \
		-f ci/Dockerfile.simtest \
		--build-arg COMMIT=$(COMMIT) \

# internal target used in docker build
_build-test-simulation:
	@go build -ldflags '$(ldflags)' -tags mocknet -o /simtest/simtest ./test/simulation/cmd

# internal target used in test run
_test-simulation:
	@cd test/simulation && /simtest/simtest

# ------------------------------ Single Node Mocknet ------------------------------

cli-mocknet:
	@docker compose -f build/docker/docker-compose.yml run --rm cli

run-mocknet:
	@BTC_MASTER_ADDR="bcrt1qf4l5dlqhaujgkxxqmug4stfvmvt58vx2h44c39" \
		BCH_MASTER_ADDR="qpxh73huzlhjfzcccr03zkpd9nd3wsasegmrreet72" \
		DOGE_MASTER_ADDR="mnaioCtEGdw6bd6rWJ13Mbre1kN5rPa2Mo" \
		LTC_MASTER_ADDR="rltc1qf4l5dlqhaujgkxxqmug4stfvmvt58vx2fc03xm" \
		docker compose -f build/docker/docker-compose.yml --profile mocknet up -d

stop-mocknet:
	@docker compose -f build/docker/docker-compose.yml --profile mocknet --profile midgard down -v

# Halt the Mocknet without erasing the blockchain history, so it can be resumed later.
halt-mocknet:
	@docker compose -f build/docker/docker-compose.yml --profile mocknet --profile midgard down

build-mocknet:
	@docker compose -f build/docker/docker-compose.yml --profile mocknet --profile midgard build \
		--build-arg COMMIT=$(COMMIT)

bootstrap-mocknet:
	@docker run --rm ${DOCKER_TTY_ARGS} \
		-e PARALLELISM -e BOOTSTRAP_ONLY=true --network host -w /app \
		thornode-simtest sh -c 'make _test-simulation'

ps-mocknet:
	@docker compose -f build/docker/docker-compose.yml --profile mocknet --profile midgard images
	@docker compose -f build/docker/docker-compose.yml --profile mocknet --profile midgard ps

logs-mocknet:
	@docker compose -f build/docker/docker-compose.yml --profile mocknet logs -f thornode bifrost

reset-mocknet: stop-mocknet run-mocknet

# ------------------------------ Mocknet EVM Fork ------------------------------

reset-mocknet-fork-%: stop-mocknet
	@./tools/evm/run-mocknet-fork.sh $*

# ------------------------------ Multi Node Mocknet ------------------------------

run-mocknet-cluster:
	@docker compose -f build/docker/docker-compose.yml --profile mocknet-cluster --profile midgard up -d

stop-mocknet-cluster:
	@docker compose -f build/docker/docker-compose.yml --profile mocknet-cluster --profile midgard down -v

halt-mocknet-cluster:
	@docker compose -f build/docker/docker-compose.yml --profile mocknet-cluster --profile midgard down

build-mocknet-cluster:
	@docker compose -f build/docker/docker-compose.yml --profile mocknet-cluster --profile midgard build

ps-mocknet-cluster:
	@docker compose -f build/docker/docker-compose.yml --profile mocknet-cluster --profile midgard images
	@docker compose -f build/docker/docker-compose.yml --profile mocknet-cluster --profile midgard ps

reset-mocknet-cluster: stop-mocknet-cluster build-mocknet-cluster run-mocknet-cluster

# ------------------------------ Test Sync ------------------------------

test-sync-mainnet:
	@BUILDTAG=mainnet BRANCH=mainnet $(MAKE) docker-gitlab-build
	@docker run --rm -e CHAIN_ID=thorchain-mainnet-v1 -e NET=mainnet registry.gitlab.com/thorchain/thornode:mainnet

test-sync-stagenet:
	@BUILDTAG=stagenet BRANCH=stagenet $(MAKE) docker-gitlab-build
	@docker run --rm -e CHAIN_ID=thorchain-stagenet-v2 -e NET=stagenet registry.gitlab.com/thorchain/thornode:stagenet

# ------------------------------ Docker Build ------------------------------

docker-gitlab-login:
	docker login -u ${CI_REGISTRY_USER} -p ${CI_REGISTRY_PASSWORD} ${CI_REGISTRY}

docker-gitlab-push:
	./build/docker/semver_tags.sh registry.gitlab.com/thorchain/thornode ${BRANCH} $(shell cat version) \
		| xargs -n1 | grep registry | xargs -n1 docker push
	docker push registry.gitlab.com/thorchain/thornode:${GITREF}

docker-gitlab-build:
	docker build . \
		-f build/docker/Dockerfile \
		$(shell sh ./build/docker/semver_tags.sh registry.gitlab.com/thorchain/thornode ${BRANCH} $(shell cat version)) \
		-t registry.gitlab.com/thorchain/thornode:${GITREF} \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg TAG=$(BUILDTAG)

########################################################################################
# Tools
########################################################################################

thorscan-build:
	@docker build tools/thorscan -f tools/thorscan/Dockerfile \
		-t registry.gitlab.com/thorchain/thornode:thorscan-${GITREF} \
		-t registry.gitlab.com/thorchain/thornode:thorscan

thorscan-gitlab-push:
	@docker login -u ${CI_REGISTRY_USER} -p ${CI_REGISTRY_PASSWORD} ${CI_REGISTRY}
	@docker push registry.gitlab.com/thorchain/thornode:thorscan-${GITREF}
	@docker push registry.gitlab.com/thorchain/thornode:thorscan
