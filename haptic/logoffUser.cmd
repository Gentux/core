@echo off

rem temp file
set sessionfile=C:\Temp\SAMUSER.txt

rem logoff user
quser | find /I "SAMUSER" > %sessionfile%
FOR /F "tokens=2 delims= " %%i in (%sessionfile%) DO logoff %%i 2> NUL 

rem remove temp file
del %sessionfile%

rem remove connection_user.properties
del C:\Users\SAMUSER\talend\configuration\connection_user.properties 
