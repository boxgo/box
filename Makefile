TEST_DIR?="./pkg/..."
CONF_CI?="./testdata/ci.yaml"
CONF_LOCAL?="./testdata/local.yaml"

lint:
	golangci-lint run $(TEST_DIR)

test:
	BOX_BOOT_CONFIG=$(CONF_CI) go test -race -coverprofile=coverage.out -covermode=atomic $(TEST_DIR)

local:
	BOX_BOOT_CONFIG=$(CONF_LOCAL) go test -v -race $(TEST_DIR)

.IGNORE:
doc:
	pkill godoc
	godoc -http=":6060" -play&
	@echo wait 3 second
	sleep 3
	open http://127.0.0.1:6060/pkg/github.com/boxgo/box/
