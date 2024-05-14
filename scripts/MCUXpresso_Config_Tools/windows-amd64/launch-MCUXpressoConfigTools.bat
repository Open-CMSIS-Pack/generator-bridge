@echo off
REM Script launch MCUXpresso Config Tools in OpenCMSIS generator mode. It is supported since MCUXpresso Config Tools v16

REM path to .cbuild-gen.idx.yml file
set idxFile=%1

REM  Get Config Tools Location
set cmd=REG QUERY "HKEY_CLASSES_ROOT\NXP Semiconductors.MCUXpresso Config Tools.mex\shell\open\command" /ve
FOR /F usebackq^ tokens^=2^ delims^=^"  %%F IN (`%cmd%`) DO (
  SET tools_path=%%F
)

REM Exit script config tools was not found
if not defined tools_path (
    echo MCUXpresso config tools was not found!
    exit /b 1
)

REM Get tools folder from tools exe path
SETLOCAL ENABLEDELAYEDEXPANSION
FOR /F "delims=" %%i IN ("%tools_path%") DO (
    SET tools_folder=!%%~dpi!
)
ENDLOCAL & SET "tools_folder=%tools_folder%"

REM Launch config tools from its folder
pushd %tools_folder%
%tools_path% -CreateFromProject %idxFile% -OpencmsisGeneratorCgen
popd
