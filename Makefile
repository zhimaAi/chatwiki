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

.PHONY:message_service
message_service:
	go version
	cd cmd/message_service&&go mod tidy
	set GOARCH=amd64&&set GOOS=linux&&go build -o build/message_service -ldflags "-s -w" cmd/message_service/main.go
	cd build&&git add message_service&&git update-index --chmod=+x message_service&&git ls-files --stage message_service

.PHONY:message_service_mac
message_service_mac:
	go version
	cd cmd/message_service&&go mod tidy
	GOARCH=amd64 GOOS=linux go build -o build/message_service -ldflags "-s -w" cmd/message_service/main.go
	cd build&&chmod a+x message_service&&ls -l message_service

.PHONY:crawler
crawler:
	cd cmd/crawler&&go mod tidy
	set GOARCH=amd64&&set GOOS=linux&&go build -o build/crawler -ldflags "-s -w" cmd/crawler/main.go cmd/crawler/process_page.go
	cd build&&git add crawler&&git update-index --chmod=+x crawler&&git ls-files --stage crawler

.PHONY:crawler_mac
crawler_mac:
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

.PHONY:make_all
make_all:
	make chatwiki
	make message_service
	make crawler
	make client_side_build

.PHONY:make_all_mac
make_all_mac:
	make chatwiki_mac
	make message_service_mac
	make client_side_build_mac
	make crawler_mac