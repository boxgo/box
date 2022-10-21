TEST_DIR?="./pkg/..."
TEST_RUN=".*"
CONF_CI?="./testdata/ci.yaml"
CONF_LOCAL?="./testdata/local.yaml"

lint:
	golangci-lint run $(TEST_DIR)

test:
	BOX_BOOT_CONFIG=$(CONF_CI) go test -run ${TEST_RUN} -race -coverprofile=coverage.out -covermode=atomic $(TEST_DIR)

bench:
	BOX_BOOT_CONFIG=$(CONF_CI) go test -run ${TEST_RUN} -race -bench=. -benchmem $(TEST_DIR)

test_local:
	BOX_BOOT_CONFIG=$(CONF_LOCAL) go test -v -run ${TEST_RUN} -race $(TEST_DIR)

bench_local:
	BOX_BOOT_CONFIG=$(CONF_LOCAL) go test -v -run ${TEST_RUN} -race -bench=. -benchmem $(TEST_DIR)

.IGNORE:
doc:
	pkill godoc
	godoc -http=":6060" -play&
	@echo wait 3 second
	sleep 3
	open http://127.0.0.1:6060/pkg/github.com/boxgo/box/
