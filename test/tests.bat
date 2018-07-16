@echo off
set BUILDRESULT=NO BUILD EXECUTED
set TESTRESULT=NO TESTS EXECUTED
cd ..
echo Building restshell...
go build
if /I "%ERRORLEVEL%" NEQ "0" set BUILDRESULT=BUILD FAILED
if /I "%ERRORLEVEL%" EQU "0" set BUILDRESULT=BUILD SUCCEEDED
echo %BUILDRESULT%

cd test
echo Running tests...
..\restshell run all
if /I "%ERRORLEVEL%" NEQ "0" set TESTRESULT=TESTS FAILED
if /I "%ERRORLEVEL%" EQU "0" set TESTRESULT=TESTS SUCCEEDED
echo %TESTRESULT%