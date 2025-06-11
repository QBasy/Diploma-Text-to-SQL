@echo off
set GOROOT=C:\Program Files\Go
set GOPATH=C:\Users\QoidynBasy\go
set PATH=%GOROOT%\bin;%PATH%

:: Билдим и запускаем database_service
echo Building database_service...
go build -o database_service.exe .
start cmd /k database_service.exe

:: Билдим и запускаем API
echo Building API...
go build -o API.exe API
start cmd /k API.exe

:: Билдим и запускаем auth_service
echo Building auth_service...
go build -o auth_service.exe auth-service
start cmd /k auth_service.exe

:: Билдим и запускаем history_service
echo Building history_service...
go build -o history_service.exe history-service
start cmd /k history_service.exe

:: Билдим и запускаем visualisation_service
echo Building visualisation_service...
go build -o visualisation_service.exe visualisation-service/cmd/server
start cmd /k visualisation_service.exe

:: Билдим и запускаем text_to_sql_service
echo Building text_to_sql_service...
go build -o text_to_sql_service.exe text-to-sql/cmd/server
start cmd /k text_to_sql_service.exe

:: Запускаем фронтенд
echo Запуск фронтенда...
cd /d diploma-frontend
start cmd /k npm run dev

echo Все сервисы запущены.
pause
