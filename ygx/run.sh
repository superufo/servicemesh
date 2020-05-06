@echo off
echo %1

set os=%1
set ChassisConfDir=D:\\gopro\\src\\github.com\\go-chassis\\go-chassis\\examples\\login\\conf
echo %ChassisConfDir%
set CHASSIS_HOME=D:\\gopro\\src\\github.com\\go-chassis\\go-chassis\\examples\\login\\server
echo %CHASSIS_HOME%

@echo off
if /i "%os%"=="linux" (
   rm -Rf  login
   SET CGO_ENABLED=0
   SET GOOS=linux
   SET GOARCH=amd64
   go  build  -o login  main.go
   login
)  else (
   rm -Rf  login.exe
   SET CGO_ENABLED=0
   SET GOOS=windows
   SET GOARCH=amd64
   go  build  -o login.exe  main.go
   login.exe
)

echo success


