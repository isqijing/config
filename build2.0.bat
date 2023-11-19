SET CGO_ENABLED=0
SET GOARCH=amd64
:: build windows
go build -p 4 -o build2.0/main_config.exe servers/main2.0/main.go

:: build linux if you need
:: SET GOOS=linux
:: go build -p 4 -o build/main_config.bin servers/main/main.go

:: build mac if you need
:: SET GOOS=darwin
:: go build -o ./build/main_config servers/main/main.go

:: copy template.txt dynamic.json
xcopy /y template2.0.txt build2.0\

xcopy /y dynamic.json build2.0\

mkdir build2.0\output2.0


mkdir build2.0\webserver
xcopy /y/e webserver  build2.0\webserver\

xcopy /y clean2.0.*  build2.0\
xcopy /y README.* build2.0\

exit
