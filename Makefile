push:
	go mod download && go mod vendor && git add . && git commit -m '$(m)'

build:
	go build  main.go

test:
	export RUN_ENVIRONMENT=test
	go test -race $(go list ./...)

cover:
	go test ./... -coverprofile=coverage.out
