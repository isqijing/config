SET CGO_ENABLED=0
SET GOARCH=amd64
:: build windows
go build -p 4 -o build/main_config.exe servers/main/main.go

:: build linux
SET GOOS=linux
go build -p 4 -o build/main_config.bin servers/main/main.go
exit
