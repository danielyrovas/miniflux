APP          := 'miniflux'
VERSION      := `git describe --tags --abbrev=0`
COMMIT       := `git rev-parse --short HEAD`
BUILD_DATE   := `date +%FT%T%z`
LD_FLAGS     := "-s -w -X 'miniflux.app/version.Version=" + VERSION +"' -X 'miniflux.app/version.Commit="+ COMMIT +"' -X 'miniflux.app/version.BuildDate=" + BUILD_DATE +"'"
PKG_LIST     := `go list ./... | grep -v /vendor/`
DB_URL       := 'postgres://postgres:postgres@localhost/miniflux_test?sslmode=disable'

list:
  just --list
echo:
	echo {{LD_FLAGS}}

build:
  go build -buildmode=pie -ldflags="{{LD_FLAGS}}" -o {{APP}} main.go

run:
  @ LOG_DATE_TIME=1 DEBUG=1 RUN_MIGRATIONS=1 go run main.go

run-create-admin:
  @ LOG_DATE_TIME=1 DEBUG=1 RUN_MIGRATIONS=1 CREATE_ADMIN=1 ADMIN_USERNAME=admin ADMIN_PASSWORD=admin123 go run main.go

test:
  go test -cover -race -count=1 ./...

lint:
  golint -set_exit_status "{{PKG_LIST}}"

