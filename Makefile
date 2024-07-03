.PHONY:default
default:
	echo "Please Specify The Packaged APP ... "

.PHONY:chatwiki
chatwiki:
	go version
	cd cmd/chatwiki&&go mod tidy
	set GOARCH=amd64&&set GOOS=linux&&go build -o build/chatwiki -ldflags "-s -w" cmd/chatwiki/main.go
	cd build&&git add chatwiki&&git update-index --chmod=+x chatwiki&&git ls-files --stage chatwiki

.PHONY:chatwiki_mac
chatwiki_mac:
	go version
	cd cmd/chatwiki&&go mod tidy
	GOARCH=amd64 GOOS=linux go build -o build/chatwiki -ldflags "-s -w" cmd/chatwiki/main.go
	cd build&&chmod a+x chatwiki&&ls -l chatwiki

.PHONY:crawler
crawler:
	go version
	cd cmd/crawler&&go mod tidy
	set GOARCH=amd64&&set GOOS=linux&&go build -o build/crawler -ldflags "-s -w" cmd/crawler/main.go cmd/crawler/process_page.go
	cd build&&git add crawler&&git update-index --chmod=+x crawler&&git ls-files --stage crawler

.PHONY:crawler_mac
crawler_mac:
	go version
	cd cmd/crawler&&go mod tidy
	GOARCH=amd64 GOOS=linux go build -o build/crawler -ldflags "-s -w" cmd/crawler/*.go
	cd build&&chmod a+x crawler&&ls -l crawler

.PHONY:client_side_build
client_side_build:
	go version
	cd cmd/client_side_build&&go mod tidy
	set GOARCH=amd64&&set GOOS=linux&&go build -o build/client_side_build -ldflags "-s -w" cmd/client_side_build/main.go
	cd build&&git add client_side_build&&git update-index --chmod=+x client_side_build&&git ls-files --stage client_side_build

.PHONY:client_side_build_mac
client_side_build_mac:
	go version
	cd cmd/client_side_build&&go mod tidy
	GOARCH=amd64 GOOS=linux go build -o build/client_side_build -ldflags "-s -w" cmd/client_side_build/main.go
	cd build&&chmod a+x client_side_build&&ls -l client_side_build

.PHONY:websocket
websocket:
	go version
	cd cmd/websocket&&go mod tidy
	set GOARCH=amd64&&set GOOS=linux&&go build -o build/websocket -ldflags "-s -w" cmd/websocket/main.go
	cd build&&git add websocket&&git update-index --chmod=+x websocket&&git ls-files --stage websocket

.PHONY:websocket_mac
websocket_mac:
	go version
	cd cmd/websocket&&go mod tidy
	GOARCH=amd64 GOOS=linux go build -o build/websocket -ldflags "-s -w" cmd/websocket/main.go
	cd build&&chmod a+x websocket&&ls -l websocket

.PHONY:make_all
make_all:
	make chatwiki
	make crawler
	make client_side_build
	make websocket

.PHONY:make_all_mac
make_all_mac:
	make chatwiki_mac
	make crawler_mac
	make client_side_build_mac
	make websocket_mac