# Test directory
This is the directory containing tests.

Execute the tests.bat file to verify restshell builds and returns success.

```
C:\> tests.bat
Building restshell...
Running general tests...
REM ## Running the config file for testing RestShell.
REM ####################################################################
REM ## General tests
REM ##
REM ## Running some basic tests using some well known REST api's
REM ## All assertions should pass
REM ##
REM ## This test requires the .rsconfig file to be loaded due to
REM ## dependencies on alias command; used to verify alias. An output
REM ## should have been displayed from the config file above.
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
get -s /posts/1
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
Ran 20 commands in 1177.3ms. Exited with Success
BUILD SUCCEEDED
GENERAL TESTS SUCCEEDED
```

There should be two lines at the end of the run showing the build succeeded and the tests succeeded.
If the tests succeeded, the %ERRORLEVEL% should be 0 otherwise non-zero.