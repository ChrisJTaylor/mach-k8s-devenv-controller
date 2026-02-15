default_group := "dev"

# show all available tasks
_default:
  @just --list

# initialise kubebuilder project
[group("kubebuilder")]
initialise-project name domain="machinology":
  kubebuilder init --domain {{domain}}.dev --repo github.com/{{domain}}/{{name}}

# create the api/crd
[group("kubebuilder")]
create-api kind:
  kubebuilder create api --group {{default_group}} --version v1alpha1 --kind {{kind}}

# create manifests
[group("kubebuilder")]
manifests:
  make manifests

# generate deepcopies
[group("kubebuilder")]
generate:
  make generate


# init module
[group("golang")]
[no-cd]
init-module name="":
  go mod init {{name}}

# restore dependencies
[group("golang")]
[no-cd]
restore:
  go mod tidy

# build solution
[group("golang")]
[no-cd]
build:
  mkdir -p bin
  go build -o bin/devenv-controller ./cmd/main.go

# run tests 
[group("golang")]
test: restore
  make test

# run all tests with coverage
[group("golang")]
test-with-coverage:
  go test -cover -v ./internal/... || true

# generate coverage report
[group("golang")]
generate-coverage:
  go test -coverprofile=coverage.out ./internal/... || true
  go tool cover -html=coverage.out

# run all tests including e2e 
[group("golang")]
test-all:
  go test -v ./internal/...
  go test -v ./test/e2e/...

# run benchmarks
[group("golang")]
benchmark:
  go test -bench=. ./internal/...

# run tests with race detector
[group("golang")]
test-with-race-detector:
  go test -race ./...

# run module
[group("golang")]
run:
  go run ./cmd/main.go

# Watch for file changes and run tests automatically (requires `watchexec`)
[group("golang")]
watch-tests:
  watchexec -c -e go "just test"


# run linter
[group("quality")]
lint:
  golangci-lint run

# run formatter
[group("quality")]
format:
  golangci-lint fmt
