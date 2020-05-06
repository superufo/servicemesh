@echo off
echo %1

set os=%1
set ChassisConfDir=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\example\\demo\\server\\conf
echo %ChassisConfDir%
set CHASSIS_HOME=D:\\gopro\\src\\github.com\\go-chassis\\ygx\\sidecar\\example\\demo\\rest-server
echo %CHASSIS_HOME%

@echo off
if /i "%os%"=="linux" (
   rm -Rf ygx_example_demo_rest_server
   SET CGO_ENABLED=0
   SET GOOS=linux
   SET GOARCH=amd64
   go  build  -o ygx_example_demo_rest_server  main.go
   ygx_example_demo_rest_server
)  else (
   rm -Rf ygx_example_demo_rest_server.exe
   SET CGO_ENABLED=0
   SET GOOS=windows
   SET GOARCH=amd64
   go  build  -o ygx_example_demo_rest_server.exe  main.go
   ygx_example_demo_rest_server.exe
)

echo success