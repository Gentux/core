@echo off

icacls \\winad.intra.nxbay.com\NETLOGON\SAMUSER.config.bat /inheritance:r /grant:r intra.nxbay.com\SAMUSER:RX /Q

icacls \\winad.intra.nxbay.com\NETLOGON\SAMUSER.config.bat /remove:g "Domain Users" /remove:g "Authenticated Users" /remove:g "Everyone" /Q
