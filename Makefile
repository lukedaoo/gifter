build:
	go build -o bin/gifter .
run:
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make run FILE=example2.gif"; \
		exit 1; \
	fi
	go run . -w 80 -h 40 $(FILE)
run-c:
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make run FILE=example2.gif"; \
		exit 1; \
	fi
	go run . -w 80 -h 40 -c $(FILE)
run-1:
	go run . -w 80 -h 40 -c example.gif 
run-2:
	go run . -w 80 -h 40 -c example2.gif
run-help:
	go run . -h
