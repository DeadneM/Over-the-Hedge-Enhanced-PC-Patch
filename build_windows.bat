@echo off
setlocal
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -trimpath -buildvcs=false -ldflags="-s -w" -o OverTheHedge_Patcher_R1.exe ./patcher
if errorlevel 1 exit /b 1
certutil -hashfile OverTheHedge_Patcher_R1.exe SHA256
