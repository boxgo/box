lint:
	golangci-lint run pkg/...

test:
	go test -v -tags no_config_init ./pkg/...

cover:
	go test -v -tags no_config_init -race -coverprofile=coverage.txt -covermode=atomic ./pkg/...

.IGNORE:
doc:
	pkill godoc
	godoc -http=":6060"&
	@echo wait 3 second
	sleep 3
	open http://127.0.0.1:6060/pkg/github.com/boxgo/box/
