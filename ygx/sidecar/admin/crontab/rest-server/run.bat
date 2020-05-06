@echo off
echo %1

set os=%1
set ChassisConfDir=..\\conf
echo %ChassisConfDir%
set CHASSIS_HOME=.
echo %CHASSIS_HOME%

@echo off
if /i "%os%"=="linux" (
   SET CGO_ENABLED=0
   SET GOOS=linux
   SET GOARCH=amd64
   go  build  -o ygx_admin_crontab_rest_server  main.go
   ygx_admin_crontab_rest_server
)  else (
   SET CGO_ENABLED=0
   SET GOOS=windows
   SET GOARCH=amd64

   rem   -H windowsgui -w -s 隐藏程序自身黑窗口
   go  build  -o  ygx_admin_crontab_rest_server.exe  main.go
   ygx_admin_crontab_rest_server.exe
)

echo success