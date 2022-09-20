@ECHO OFF

if %1==start start /b main.exe run
if %1==status tasklist | findstr main.exe
if %1==stop taskkill /F /im main.exe
if %1==log more VoiceServerLog.log