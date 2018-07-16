# Test directory
This is the directory containing tests.

Execute the tests.bat file to verify restshell builds and returns success.

```
C:\> tests.bat
Building restshell...
BUILD SUCCEEDED
Running tests...
REM ####################################################################
REM ##  ALL TESTS
REM ####################################################################
REM ## JSON tests
REM ##
REM ## Running some basic tests using some well known REST api's
REM ## All assertions should pass
REM ##
REM ####################################################################
REM ##
REM ## Test getting local ip address; assert finds 3 periods in ip address
REM ##
get -s /?format=json
Assertions Passed (3)
REM ####################################################################
REM ##
REM ## Test getting a post from an online mock rest api
REM ##
get  /posts/1
Assertions Passed (6)
REM ####################################################################
REM ##
REM ## Test an invalid REST api call that will return 404
REM ##
get -s /junk/asdf
GET: HTTP Status: 404 Not Found
Assertions Passed (2)
REM ####################################################################
ALL ASSERTIONS PASSED (11)
Ran 20 commands in 1084.0ms. Exited with Success
REM ####################################################################
REM ## Basic XMLtests
REM ##
REM ## Call a simple XML api that returns IP lookup data and validate
REM ## results
REM ##
REM ####################################################################
REM ##
REM ## Test getting local ip data and assert known fields in XML response
REM ##
get -s /
Assertions Passed (4)
REM ####################################################################
ALL ASSERTIONS PASSED (15)
Ran 9 commands in 139.5ms. Exited with Success
Ran 30 commands in 1226.4ms. Exited with Success
TESTS SUCCEEDED
```

There should be two lines at the end of the run showing the build succeeded and the tests succeeded.
If the tests succeeded, the %ERRORLEVEL% should be 0 otherwise non-zero.