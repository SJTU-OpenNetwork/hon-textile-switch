.PHONY:shadow
.PHONY:protos
shadow:
	go build -o bin/shadow github.com/SJTU-OpenNetwork/hon-textile-switch
protos:
	cd pb/protos; protoc --go_out=../. *.proto
