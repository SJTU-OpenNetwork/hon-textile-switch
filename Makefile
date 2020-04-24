.PHONY:shadow
.PHONY:protos
.PHONY:shadow-arm
.PHONY:init
init:
    mkdir bin
    mkdir bin/arm
shadow:
	go build -o bin/shadow github.com/SJTU-OpenNetwork/hon-textile-switch
protos:
	cd pb/protos; protoc --go_out=../. *.proto
shadow-arm:
    go build -o bin/arm/shadow github.com/SJTU-OpenNetwork/hon-textile-switch
