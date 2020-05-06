@echo off
echo %1

set os=%1
set ChassisConfDir=D:\\gopro\\src\\github.com\\go-chassis\\go-chassis\\examples\\rpc\\server\\conf
echo %ChassisConfDir%
set CHASSIS_HOME=D:\\gopro\\src\\github.com\\go-chassis\\go-chassis\\examples\\rpc\\server
echo %CHASSIS_HOME%

@echo off
if /i "%os%"=="linux" (
   SET CGO_ENABLED=0
   SET GOOS=linux
   SET GOARCH=amd64
   go  build  -o rpcserver  main.go
   rpcserver
)  else (
   SET CGO_ENABLED=0
   SET GOOS=windows
   SET GOARCH=amd64
   go  build  -o rpcserver.exe  main.go
   rpcserver.exe
)

echo success