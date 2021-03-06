@echo off
echo %1

set os=%1
set ChassisConfDir=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\admin\\login\\conf
echo %ChassisConfDir%
set CHASSIS_HOME=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\admin\\login\\rest-server
echo %CHASSIS_HOME%

@echo off
if /i "%os%"=="linux" (
   rm -Rf  ygx_admin_login_rest_server
   SET CGO_ENABLED=0
   SET GOOS=linux
   SET GOARCH=amd64
   go  build  -o ygx_admin_login_rest_server  main.go
   ygx_admin_login_rest_server
)  else (
   rm -Rf  ygx_admin_login_rest_server.exe
   SET CGO_ENABLED=0
   SET GOOS=windows
   SET GOARCH=amd64
   go  build  -o ygx_admin_login_rest_server.exe  main.go
   ygx_admin_login_rest_server.exe
)

echo success