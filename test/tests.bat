@echo off
set BUILDRESULT=NO BUILD EXECUTED
set TESTRESULT=NO TESTS EXECUTED
cd ..
echo Building restshell...
go build
if /I "%ERRORLEVEL%" NEQ "0" set BUILDRESULT=BUILD FAILED
if /I "%ERRORLEVEL%" EQU "0" set BUILDRESULT=BUILD SUCCEEDED
cd test
echo Running general tests...
..\restshell run general
if /I "%ERRORLEVEL%" NEQ "0" set TESTRESULT=GENERAL TESTS FAILED
if /I "%ERRORLEVEL%" EQU "0" set TESTRESULT=GENERAL TESTS SUCCEEDED

echo %BUILDRESULT%
echo %TESTRESULT%