gen-proto:
	protoc -I=. --go_out=./gen ./proto/woop-socket-message.proto
start-colima:
	colima start --memory 4