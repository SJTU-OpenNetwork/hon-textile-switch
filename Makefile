.PHONY:shadow
.PHONY:protos
.PHONY:shadow-arm
.PHONY:init
init:
	mkdir bin
	mkdir bin/arm
shadow:
	go build  -ldflags "-s -w" -o bin/shadow github.com/SJTU-OpenNetwork/hon-textile-switch
shadow-arm:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc go build  -ldflags "-linkmode external -extldflags -static -s -w" -o bin/shadow-arm github.com/SJTU-OpenNetwork/hon-textile-switch
protos:
	cd pb/protos; protoc --go_out=../. *.proto
