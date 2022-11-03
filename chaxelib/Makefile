build-mac:
	go build -o ./dist/chaxelib

build-window:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./dist/chaxelib.exe

build-self:
	sudo go build -o /usr/local/bin/chaxelib 