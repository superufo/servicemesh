@echo off
echo %1

set os=%1
set ChassisConfDir=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\usercenter\\user\\grpc-server\\conf
echo %ChassisConfDir%
set CHASSIS_HOME=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\usercenter\\user\\grpc-server
echo %CHASSIS_HOME%

@echo off
if /i "%os%"=="linux" (
   rm -Rf ygx_usercenter_user_grpc-server
   SET CGO_ENABLED=0
   SET GOOS=linux
   SET GOARCH=amd64
   go  build  -o ygx_usercenter_user_grpc-server  main.go
   ygx_usercenter_demo_rest_grpc-server
)  else (
   rm -Rf ygx_usercenter_user_grpc-server.exe
   SET CGO_ENABLED=0
   SET GOOS=windows
   SET GOARCH=amd64
   go  build  -o ygx_usercenter_user_grpc-server.exe  main.go
   ygx_usercenter_user_grpc-server.exe
)

echo success