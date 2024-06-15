SHELL := /bin/bash # Use bash syntax

.EXPORT_ALL_VARIABLES:
APP=prismo
APP_COMMIT:=$(shell git rev-parse HEAD)
APP_EXECUTABLE="./out/$(APP)"
CURRENT_DIR=$(shell pwd)
ALL_PACKAGES=$(shell go list ./... | grep -v /vendor)
ALL_TESTABLE_PACKAGES=$(shell go list ./... | grep -v /vendor | grep -v /nonfunctional | grep -v /mocks)
SOURCE_DIRS=$(shell go list ./... | grep -v /vendor | grep -v /out | cut -d "/" -f4 | uniq)
MIGRATION_DB_HOST=$(shell cat application.yml | grep -i ^DB_HOST | cut -d " " -f2)
MIGRATION_DB_NAME=$(shell cat application.yml | grep -i ^DB_NAME | cut -d " " -f2)
MIGRATION_TEST_DB_NAME=$(shell cat application.yml | grep -i ^TEST_DB_NAME | cut -d " " -f2)
MIGRATION_DB_USER=$(shell cat application.yml | grep -i ^DB_USER | cut -d " " -f2)
GO111MODULE:=on
COVERAGE_MIN=60

all: check-quality clean test

copy-config:
	@echo "Copying configs/application.yml.sample to application.yml ..."
	@cp configs/application.yml.sample application.yml
	@echo "Done."

setup:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint
	go get -u github.com/fzipp/gocyclo

setup-local:
	@echo "Copying commit-hook script from repo to .git folder"
	@cp scripts/commit-msg .git/hooks/commit-msg
	@chmod +x .git/hooks/commit-msg

check-quality: lint fmt vet

build: lint fmt vet
	mkdir -p out/
	go build -o $(APP_EXECUTABLE) -ldflags "-X main.commit=$(APP_COMMIT)" cmd/*.go

clean:
	rm -rf out/
	rm -f coverage*.out

testdb.create:
	@echo "creating test database $(MIGRATION_TEST_DB_NAME)..."
	@createdb -Opostgres -h $(MIGRATION_DB_HOST) -U $(MIGRATION_DB_USER) -Eutf8 $(MIGRATION_TEST_DB_NAME)

testdb.drop:
	@echo "dropping test database $(MIGRATION_TEST_DB_NAME)..."
	@dropdb --if-exists -h $(MIGRATION_DB_HOST) -U $(MIGRATION_DB_USER) $(MIGRATION_TEST_DB_NAME)

testdb.migrate: build
	@IS_TEST=true \
	${APP_EXECUTABLE} migrate

test: copy-ci-test-config testdb.drop testdb.create testdb.migrate
	@IS_TEST=true \
	go test -p 1 -count=1 -cover -coverprofile .cover.out $(ALL_TESTABLE_PACKAGES)
	go tool cover -func=.cover.out | tail -1 | awk '{print "total coverage: "$$3}'
	go tool cover -html=.cover.out -o .testCoverage.html


db.create:
	@echo "creating database $(MIGRATION_DB_NAME)..."
	@createdb -Opostgres -h $(MIGRATION_DB_HOST) -U $(MIGRATION_DB_USER) -Eutf8 $(MIGRATION_DB_NAME)

db.drop:
	@echo "dropping database $(MIGRATION_DB_NAMEm)..."
	@dropdb --if-exists -h $(MIGRATION_DB_HOST) -U $(MIGRATION_DB_USER) $(MIGRATION_DB_NAME)

fmt:
	gofmt -l -s -w $(SOURCE_DIRS)

vet:
	@go vet ./...

lint:
	@if [[ `golint $(ALL_PACKAGES) | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } | wc -l | tr -d ' '` -ne 0 ]]; then \
          golint $(ALL_PACKAGES) | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; }; \
          exit 2; \
    fi;

migrate: build
	${APP_EXECUTABLE} migrate

rollback: build
	${APP_EXECUTABLE} rollback

run: build
	${APP_EXECUTABLE} server
