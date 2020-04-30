.PHONY:shadow
.PHONY:protos
.PHONY:shadow-arm
.PHONY:init
init:
	mkdir bin
	mkdir bin/arm
shadow:
	go build  -ldflags "-s -w" -o bin/shadow github.com/SJTU-OpenNetwork/hon-textile-switch
protos:
	cd pb/protos; protoc --go_out=../. *.proto
