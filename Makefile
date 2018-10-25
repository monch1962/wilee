clean:
	rm bin/wilee
build:
lambda:
	env GOOS=linux go build -ldflags="-s -w" -o bin/wilee wilee.go
linux:
	env GOOS=linux go build -ldflags="-s -w" -o bin/wilee wilee.go
linux-arm:
	env GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o bin/wilee wilee.go
local:
	go build -ldflags="-s -w" -o bin/wilee wilee.go
mac:
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/wilee wilee.go
windows:
	env GOOS=windows go build -ldflags="-s -w" -o bin/wilee.exe wilee.go
