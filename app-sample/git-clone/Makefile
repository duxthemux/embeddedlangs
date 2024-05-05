IMG = gopherconbr-2024-emb-languages/git-clone

setup:
	docker buildx create --use --name multi-arch-builder
build:
	
	CGO=0 GOOS=linux GOARCH=arm64 go build -ldflags "-w -s -extldflags=-static " -o bin/main_arm64 .
	CGO=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -extldflags=-static " -o bin/main_amd64 .
	
	docker buildx build -t $(IMG) --platform=linux/amd64,linux/arm64 --load .
