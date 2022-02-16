@ECHO OFF

SET src=.\src
SET out=build\myrtle.exe

go test %src% -test.v

IF %ERRORLEVEL% EQU 0 (
    go build -trimpath -ldflags "-s -w" -o %out% %src%
)
