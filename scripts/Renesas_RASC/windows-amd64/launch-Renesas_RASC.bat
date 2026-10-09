@echo off
rem ===========================================================================
rem Opens Renesas RA Smart Configurator (RASC) UI for the CMSIS-Toolbox generator
rem "RenesasRASC".
rem
rem   %1            = <solution>.cbuild-gen-idx.yml (passed by csolution)
rem   RASC_EXE_PATH = path to rasc.exe (same var as "executes:" in cproject.yml)
rem
rem ===========================================================================

setlocal EnableExtensions EnableDelayedExpansion

rem --- Step 1: Validate input argument -------------------------------------
if "%~1"=="" (
    echo [RASC] Error: No *.cbuild-gen-idx.yml file path provided as an argument
    exit /b 1
)
set "IDX_FILE=%~1"
rem findstr treats "/" as a switch, so normalize to "\" first
set "IDX_FILE=%IDX_FILE:/=\%"
if not exist "%IDX_FILE%" (
    echo [RASC] Error: "%IDX_FILE%" not found
    exit /b 1
)

rem --- Step 2: Resolve RASC executable from RASC_EXE_PATH ------------------
if not defined RASC_EXE_PATH (
    echo [RASC] Error: Environment variable RASC_EXE_PATH is not set
    echo [RASC]        Set it to the RA Smart Configurator executable, e.g.
    echo [RASC]        set RASC_EXE_PATH=\path\to\rasc.exe
    exit /b 1
)
set "RASC_EXE=%RASC_EXE_PATH:"=%"
if not exist "%RASC_EXE%" (
    echo [RASC] Error: RASC_EXE_PATH does not point to an existing file: "%RASC_EXE%"
    exit /b 1
)
echo [RASC] Using RASC_EXE_PATH defined in environment: "%RASC_EXE%"

rem --- Step 3: Parse *.cbuild-gen-idx.yml ----------------------------------
set "CBUILD_GEN="
for /f "usebackq tokens=1,* delims=:" %%a in (`findstr /C:"cbuild-gen: " "%IDX_FILE%"`) do (
    if not defined CBUILD_GEN for /f "tokens=*" %%v in ("%%b") do set "CBUILD_GEN=%%v"
)
if not defined CBUILD_GEN (
    echo [RASC] Error: Could not parse "cbuild-gen:" from "%IDX_FILE%"
    exit /b 1
)
set "CBUILD_GEN=%CBUILD_GEN:/=\%"
if not exist "%CBUILD_GEN%" (
    echo [RASC] Error: "%CBUILD_GEN%" not found
    exit /b 1
)

rem --- Step 4: Parse *.cbuild-gen.yml --------------------------------------
set "CPROJECT_FILE="
for /f "usebackq tokens=1,* delims=:" %%a in (`findstr /B /C:"  project: " "%CBUILD_GEN%"`) do (
    if not defined CPROJECT_FILE for /f "tokens=*" %%v in ("%%b") do set "CPROJECT_FILE=%%v"
)
if not defined CPROJECT_FILE (
    echo [RASC] Error: Could not parse "project:" from "%CBUILD_GEN%"
    exit /b 1
)
set "CPROJECT_FILE=%CPROJECT_FILE:/=\%"
for %%p in ("%CPROJECT_FILE%") do set "PROJECT_DIR=%%~dpp"
set "PROJECT_DIR=%PROJECT_DIR:~0,-1%"

rem missing configuration.xml is OK, RASC creates one
set "CONFIG_XML=%PROJECT_DIR%\configuration.xml"

rem --- Step 5: Open RASC UI -------------------------------------------------
echo [RASC] Project directory : "%PROJECT_DIR%"
echo [RASC] Configuration     : "%CONFIG_XML%"

pushd "%PROJECT_DIR%"
rem "call" waits for RASC to close (cmd waits for GUI apps too) and also
rem works if RASC_EXE_PATH points to a wrapper .bat
echo [RASC] Opening RA Smart Configurator, close it when configuration is done...
call "%RASC_EXE%" "%CONFIG_XML%"
set "RASC_RESULT=%ERRORLEVEL%"
popd

if not "%RASC_RESULT%"=="0" (
    echo [RASC] Error: RASC execution failed with exit code %RASC_RESULT%
    exit /b %RASC_RESULT%
)

echo [RASC] RA Smart Configurator closed
exit /b 0
