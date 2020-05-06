@echo off
echo %1

set os=%1
set ChassisConfDir=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\usercenter\\user\\grpc-client\\conf
echo %ChassisConfDir%
set CHASSIS_HOME=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\usercenter\\user\\grpc-client
echo %CHASSIS_HOME%

@echo off
if /i "%os%"=="linux" (
   rm -Rf ygx_usercenter_user_grpc-client
   SET CGO_ENABLED=0
   SET GOOS=linux
   SET GOARCH=amd64
   go  build  -o ygx_usercenter_user_grpc-client  main.go
   ygx_usercenter_demo_rest_grpc-client
)  else (
   rm -Rf ygx_usercenter_user_grpc-client.exe
   SET CGO_ENABLED=0
   SET GOOS=windows
   SET GOARCH=amd64
   go  build  -o ygx_usercenter_user_grpc-client.exe  main.go
   ygx_usercenter_user_grpc-client.exe
)

echo success