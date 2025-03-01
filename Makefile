golb: format
	go build -o build/golb .

clean:
	rm build/golb

format:
	gofmt -w .

.PHONY: clean format
