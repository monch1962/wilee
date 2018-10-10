clean:
	rm bin/wilee
build:
lambda:
	env GOOS=linux go build -ldflags="-s -w" -o bin/wilee wilee/main.go
linux:
	env GOOS=linux go build -ldflags="-s -w" -o bin/wilee wilee/main.go
linux-arm:
	env GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o bin/wilee wilee/main.go
local:
	go build -ldflags="-s -w" -o bin/wilee wilee/main.go
mac:
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/wilee wilee/main.go
windows:
	env GOOS=windows go build -ldflags="-s -w" -o bin/wilee.exe wilee/main.go
