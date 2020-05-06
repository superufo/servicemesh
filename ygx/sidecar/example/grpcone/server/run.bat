@echo off
echo %1

set os=%1
rem set ChassisConfDir=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\example\\rest\\server\\conf
set ChassisConfDir=..\\conf
echo %ChassisConfDir%
set CHASSIS_HOME=.
rem set CHASSIS_HOME=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\example\\rest\\server

echo %CHASSIS_HOME%

@echo off
if /i "%os%"=="linux" (
   rm -Rf ygx_grpc_server
   SET CGO_ENABLED=0
   SET GOOS=linux
   SET GOARCH=amd64
   go  build  -o ygx_grpc_server  main.go
   ygx_grpc_server
)  else (
   rm -Rf ygx_grpc_server.exe
   SET CGO_ENABLED=0
   SET GOOS=windows
   SET GOARCH=amd64
   go  build  -o ygx_grpc_server.exe  main.go
   ygx_grpc_server.exe
)

echo success