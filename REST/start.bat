@echo off

:: Start Journal Server
cd /D "C:\Users\User\Desktop\Bain\Go\journal"
start "Journal Server" go run main.go

:: Start Auth Server
cd /D "C:\Users\User\Desktop\Bain\Go\Auth\server"
start "Auth Server" go run server.go
