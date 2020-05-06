@echo off
echo %1

set os=%1
rem set ChassisConfDir=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\example\\rest\\server\\conf
set ChassisConfDir=..\\conf
echo %ChassisConfDir%
set CHASSIS_HOME=.

@echo off
if /i "%os%"=="linux" (
   rm -Rf ygx_grpc_client
   SET CGO_ENABLED=0
   SET GOOS=linux
   SET GOARCH=amd64
   go  build  -o ygx_grpc_client  main.go
   ygx_grpc_client
)  else (
   rm -Rf ygx_grpc_client.exe
   SET CGO_ENABLED=0
   SET GOOS=windows
   SET GOARCH=amd64
   go  build  -o ygx_grpc_client.exe  main.go
   ygx_grpc_client.exe
)

echo success