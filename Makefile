help:
	go

build:
	go build -o app

run:
	go run .

test:
	go test

bench_1000:
	COUNT=1000 ./prepare_bench.sh
	go test -bench . -benchmem
	rm ./file.txt

bench_100000:
	COUNT=100000 ./prepare_bench.sh
	go test -bench . -benchmem
	rm ./file.txt

bench_1000000:
	COUNT=1000000 ./prepare_bench.sh
	go test -bench . -benchmem
	rm ./file.txt