build:
	go build -o build/calculate src/*.go

rebuild:
	rm build/calculate
	go build -o build/calculate src/*.go

test:
	go test -v -cover src/*
	
coverage:
	go test -v -coverprofile=cover.out src/*.go
	go tool cover -html=cover.out
	rm cover.out

clean:
	rm -rf build
