BUILD_FILES = src/main.go
TEST_FILES = src/calculator/*

build:
	go build -o build/calculate $(BUILD_FILES)

rebuild:
	rm build/calculate
	go build -o build/calculate $(BUILD_FILES)

test:
	go test -v -cover $(TEST_FILES)
	
coverage:
	go test -v -coverprofile=cover.out $(TEST_FILES)
	go tool cover -html=cover.out
	rm cover.out

clean:
	rm -rf build

.PHONY: build rebuild test coverage clean