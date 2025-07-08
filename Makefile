build:
	go build -o bin/gifter .
run:
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make run FILE=example2.gif"; \
		exit 1; \
	fi
	@SFLAG=$$( [ -n "$(S)" ] && echo "-s $(S)" ); \
	go run . -w 80 -h 40 $$SFLAG $(FILE)
run-c:
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make run-c FILE=example2.gif [S=style]"; \
		exit 1; \
	fi
	@SFLAG=$$( [ -n "$(S)" ] && echo "-s $(S)" ); \
	go run . -w 80 -h 40 -c $$SFLAG $(FILE)
run-url:
	@SFLAG=$$( [ -n "$(S)" ] && echo "-s $(S)" ); \
	go run . -w 40 -h 40 -c $$SFLAG $(URL)
run-grahpic:
	go run . -w 200 -h 200 -m graphic example.gif
run-1:
	go run . -w 80 -h 40 -c -s shaded example.gif
run-2:
	go run . -w 80 -h 40 -c example2.gif
run-normal:
	go run . -w 200 -h 200 -s normal example2.gif
run-help:
	go run . -h
