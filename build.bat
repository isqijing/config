SET CGO_ENABLED=0
SET GOARCH=amd64
:: build windows
go build -p 4 -o build/main_config.exe servers/main/main.go

:: build linux if you need
:: SET GOOS=linux
:: go build -p 4 -o build/main_config.bin servers/main/main.go

:: build mac if you need
:: SET GOOS=darwin
:: go build -o ./build/main_config servers/main/main.go

:: copy template.txt dynamic.json
xcopy /y template.txt build\

xcopy /y dynamic.json build\

mkdir build\config
mkdir build\output
mkdir build\proto\output\proto

mkdir build\webserver
xcopy /y/e webserver  build\webserver\

xcopy /y clean.*  build\
xcopy /y README.* build\

exit
