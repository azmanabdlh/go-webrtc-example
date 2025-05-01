

make bundle:
	yarn bundle

run:	
	go run *.go

build:	
	yarn bundle
	go build -o build *.go