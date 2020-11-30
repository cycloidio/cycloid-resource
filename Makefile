build: build-in build-check build-out

build-in:
	@go build -o resource/in cmd/in/main.go

build-check:
	@go build -o resource/check cmd/check/main.go

build-out:
	@go build -o resource/out cmd/out/main.go

