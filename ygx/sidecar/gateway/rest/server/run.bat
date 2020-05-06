@echo off
echo %1

set os=%1
set ChassisConfDir=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\gateway\\rest\\conf
echo %ChassisConfDir%
set CHASSIS_HOME=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\gateway\\rest\\server
echo %CHASSIS_HOME%

@echo off
if /i "%os%"=="linux" (
   rm -Rf  ygx_rest_gateway
   SET CGO_ENABLED=0
   SET GOOS=linux
   SET GOARCH=amd64
   go  build  -o ygx_rest_gateway  main.go
   ygx_rest_gateway
)  else (
   rm -Rf  ygx_rest_gateway.exe
   SET CGO_ENABLED=0
   SET GOOS=windows
   SET GOARCH=amd64
   go  build  -o ygx_rest_gateway.exe  main.go
   ygx_rest_gateway.exe
)

echo success