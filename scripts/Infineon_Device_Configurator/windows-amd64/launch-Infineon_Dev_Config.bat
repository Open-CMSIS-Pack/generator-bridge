@echo off
setlocal
:: Step 1: Validate input argument
if "%1"=="" (
    echo Error: No *.cbuild-gen-idx.yml file path provided as argument
    exit /b 1
)
set "CBUILD_GEN_IDX=%~1"
set "CBUILD_GEN_IDX=%CBUILD_GEN_IDX:/=\%"
if not exist "%CBUILD_GEN_IDX%" (
    echo Error: %CBUILD_GEN_IDX% file not found
    exit /b 1
)
:: Step 2: Parse Device_Configurator path, project name, and RTE path from *.cbuild-gen-idx.yml
setlocal EnableDelayedExpansion
set "DEVICE_CONFIG_PATH="
set "PROJECT_NAME="
set "RTE_PATH="
for /f "delims=" %%a in ('findstr /C:"output:" "%CBUILD_GEN_IDX%"') do (
    set "DEVICE_CONFIG_PATH=%%a"
    set "DEVICE_CONFIG_PATH=!DEVICE_CONFIG_PATH:output:=!"
    set "DEVICE_CONFIG_PATH=!DEVICE_CONFIG_PATH: =!"
    set "DEVICE_CONFIG_PATH=!DEVICE_CONFIG_PATH:/=\!"
)
for /f "delims=" %%a in ('findstr /C:"project:" "%CBUILD_GEN_IDX%"') do (
    set "PROJECT_NAME=%%a"
    set "PROJECT_NAME=!PROJECT_NAME:project:=!"
    set "PROJECT_NAME=!PROJECT_NAME: =!"
)
for /f "delims=" %%a in ('findstr /C:"name:" "%CBUILD_GEN_IDX%"') do (
    set "RTE_PATH=%%a"
    set "RTE_PATH=!RTE_PATH:name:=!"
    set "RTE_PATH=!RTE_PATH: =!"
    set "RTE_PATH=!RTE_PATH:/=\!"
    :: Derive RTE folder by going one level back from the path in name:
    for %%i in ("!RTE_PATH!\..") do set "RTE_PATH=%%~dpi"
    set "RTE_PATH=!RTE_PATH:~0,-1!"
)
echo DEBUG: DEVICE_CONFIG_PATH: !DEVICE_CONFIG_PATH! >> debug.log
echo DEBUG: PROJECT_NAME: !PROJECT_NAME! >> debug.log
echo DEBUG: RTE_PATH: !RTE_PATH! >> debug.log
if not defined DEVICE_CONFIG_PATH (
    echo Error: Could not parse Device_Configurator path from %CBUILD_GEN_IDX%
    exit /b 1
)
if not defined PROJECT_NAME (
    echo Error: Could not parse project name from %CBUILD_GEN_IDX%
    exit /b 1
)
if not defined RTE_PATH (
    echo Error: Could not parse RTE path from name field in %CBUILD_GEN_IDX%
    exit /b 1
)
endlocal & set "DEVICE_CONFIG_PATH=%DEVICE_CONFIG_PATH%" & set "PROJECT_NAME=%PROJECT_NAME%" & set "RTE_PATH=%RTE_PATH%"

:: Step 3: Define paths based on Device_Configurator folder
set "DESIGN_PATH=%DEVICE_CONFIG_PATH%\design.modus"
set "GEN_DIR=%DEVICE_CONFIG_PATH%\GeneratedSource"
set "YAML_FILE=%DEVICE_CONFIG_PATH%\%PROJECT_NAME%.cgen.yml"
set "TMP_FILE=%YAML_FILE%.tmp"

:: Step 4: Fetch exePath from IDC JSON files, prioritizing latest version of InfineonDeviceConfigurator
setlocal EnableDelayedExpansion
set "TOOL_PATH="
set "LATEST_VERSION=0.0.0"
set "LATEST_PATH="

:: Check per-user IDC JSON files
echo Checking per-user IDC JSON files in %USERPROFILE%\AppData\Local\Infineon_Technologies_AG\Infineon-Toolbox...
if exist "%USERPROFILE%\AppData\Local\Infineon_Technologies_AG\Infineon-Toolbox" (
    for /f "delims=" %%f in ('dir /b "%USERPROFILE%\AppData\Local\Infineon_Technologies_AG\Infineon-Toolbox\*.json" 2^>nul') do (
        set "JSON_FILE=%USERPROFILE%\AppData\Local\Infineon_Technologies_AG\Infineon-Toolbox\%%f"
        echo Processing file: !JSON_FILE! >> debug.log
        set "TEMP_JSON=%TEMP%\temp_json_%RANDOM%.txt"
        type "!JSON_FILE!" > "!TEMP_JSON!"
        echo Searching for featureId... >> debug.log
        findstr /I /C:"com.ifx.tb.tool.deviceconfigurator" "!TEMP_JSON!"
        set "FINDSTR_ERROR=!ERRORLEVEL!"
        if !FINDSTR_ERROR! equ 0 (
            echo featureId found in: !JSON_FILE! >> debug.log
            echo Searching for version... >> debug.log
            findstr /I /C:"\"version\":" "!TEMP_JSON!"
            set "VERSION="
            for /f "tokens=1,* delims=:" %%v in ('findstr /I /C:"\"version\":" "!TEMP_JSON!"') do (
                set "VERSION=%%w"
                set "VERSION=!VERSION:"=!"
                set "VERSION=!VERSION:,=!"
                set "VERSION=!VERSION: =!"
                for /f "tokens=1-3 delims=." %%x in ("!VERSION!") do (
                    set "VERSION=%%x.%%y.%%z"
                )
            )
            if not defined VERSION set "VERSION=0.0.0"
            echo Version found: !VERSION! >> debug.log
            echo Searching for exePath... >> debug.log
            findstr /I /C:"\"exePath\":" "!TEMP_JSON!"
            for /f "tokens=1,* delims=:" %%a in ('findstr /I /C:"\"exePath\":" "!TEMP_JSON!"') do (
                set "EXE_PATH=%%b"
                set "EXE_PATH=!EXE_PATH:"=!"
                set "EXE_PATH=!EXE_PATH:,=!"
                set "EXE_PATH=!EXE_PATH: =!"
                set "EXE_PATH=!EXE_PATH:/=\!"
                set "EXE_PATH=!EXE_PATH:\\=\!"
                echo Extracted exePath: !EXE_PATH! >> debug.log
                if exist "!EXE_PATH!" (
                    echo exePath exists: !EXE_PATH! >> debug.log
                    call :CompareVersions "!VERSION!" "!LATEST_VERSION!"
                    if !ERRORLEVEL! equ 1 (
                        set "LATEST_VERSION=!VERSION!"
                        set "LATEST_PATH=!EXE_PATH!"
                    )
                ) else (
                    echo exePath not found or invalid: !EXE_PATH!
                )
            )
        ) else (
            echo featureId not found in: !JSON_FILE! >> debug.log
        )
        del "!TEMP_JSON!" 2>nul
    )
) else (
    echo Per-user directory not found: %USERPROFILE%\AppData\Local\Infineon_Technologies_AG\Infineon-Toolbox
)
echo Finished processing per-user JSON files.

:: Check all-users IDC JSON files if no valid path found
if "!LATEST_PATH!"=="" (
    echo Checking all-users IDC JSON files in %ALLUSERSPROFILE%\Infineon_Technologies_AG\Infineon-Toolbox... >> debug.log
    if exist "%ALLUSERSPROFILE%\Infineon_Technologies_AG\Infineon-Toolbox" (
        for /f "delims=" %%f in ('dir /b "%ALLUSERSPROFILE%\Infineon_Technologies_AG\Infineon-Toolbox\*.json" 2^>nul') do (
            set "JSON_FILE=%ALLUSERSPROFILE%\Infineon_Technologies_AG\Infineon-Toolbox\%%f"
            echo Processing file: !JSON_FILE! >> debug.log
            set "TEMP_JSON=%TEMP%\temp_json_%RANDOM%.txt"
            type "!JSON_FILE!" > "!TEMP_JSON!"
            echo Searching for featureId... >> debug.log
            findstr /I /C:"com.ifx.tb.tool.deviceconfigurator" "!TEMP_JSON!"
            set "FINDSTR_ERROR=!ERRORLEVEL!"
            if !FINDSTR_ERROR! equ 0 (
                echo featureId found in: !JSON_FILE! >> debug.log
                echo Searching for version... >> debug.log
                findstr /I /C:"\"version\":" "!TEMP_JSON!"
                set "VERSION="
                for /f "tokens=1,* delims=:" %%v in ('findstr /I /C:"\"version\":" "!TEMP_JSON!"') do (
                    set "VERSION=%%w"
                    set "VERSION=!VERSION:"=!"
                    set "VERSION=!VERSION:,=!"
                    set "VERSION=!VERSION: =!"
                    for /f "tokens=1-3 delims=." %%x in ("!VERSION!") do (
                        set "VERSION=%%x.%%y.%%z"
                    )
                )
                if not defined VERSION set "VERSION=0.0.0"
                echo Version found: !VERSION! >> debug.log
                echo Searching for exePath... >> debug.log
                findstr /I /C:"\"exePath\":" "!TEMP_JSON!"
                for /f "tokens=1,* delims=:" %%a in ('findstr /I /C:"\"exePath\":" "!TEMP_JSON!"') do (
                    set "EXE_PATH=%%b"
                    set "EXE_PATH=!EXE_PATH:"=!"
                    set "EXE_PATH=!EXE_PATH:,=!"
                    set "EXE_PATH=!EXE_PATH: =!"
                    set "EXE_PATH=!EXE_PATH:/=\!"
                    set "EXE_PATH=!EXE_PATH:\\=\!"
                    echo Extracted exePath: !EXE_PATH! >> debug.log
                    if exist "!EXE_PATH!" (
                        echo exePath exists: !EXE_PATH! >> debug.log
                        call :CompareVersions "!VERSION!" "!LATEST_VERSION!"
                        if !ERRORLEVEL! equ 1 (
                            set "LATEST_VERSION=!VERSION!"
                            set "LATEST_PATH=!EXE_PATH!"
                        )
                    ) else (
                        echo exePath not found or invalid: !EXE_PATH!
                    )
                )
            ) else (
                echo featureId not found in: !JSON_FILE! >> debug.log
            )
            del "!TEMP_JSON!" 2>nul
        )
    ) else (
        echo All-users directory not found: %ALLUSERSPROFILE%\Infineon_Technologies_AG\Infineon-Toolbox >> debug.log
    )
)
echo Finished processing all-users JSON files.

:: Output the latest exePath
if defined LATEST_PATH (
    echo Selected latest version: !LATEST_VERSION!
    echo Selected exePath: !LATEST_PATH!
    set "TOOL_PATH=!LATEST_PATH!"
) else (
    echo Error: Could not locate valid exePath for device-configurator in IDC JSON files
    exit /b 1
)

:: Verify TOOL_PATH was set
if "!TOOL_PATH!"=="" (
    echo Error: Could not locate valid exePath for device-configurator in IDC JSON files
    exit /b 1
)
endlocal & set "TOOL_PATH=%TOOL_PATH%"

:: Step 5: Parse cbuild.yml from the same directory as RTE and extract props.json paths
setlocal EnableDelayedExpansion
set "LIBRARY_PATH="
set "CBUILD_FILE="
set "MTB_PDL_FOUND=0"
set "DEVICE_DB_FOUND=0"
set "MISSING_FILES="
set "MTB_PDL_DIRS="
set "DEVICE_DB_DIRS="

:: Derive directory containing RTE
for %%i in ("%RTE_PATH%") do set "PARENT_DIR=%%~dpi"
set "PARENT_DIR=%PARENT_DIR:~0,-1%"

:: Find *.cbuild.yml in parent directory
echo Searching for *.cbuild.yml in %PARENT_DIR%... >> debug.log
for /f "delims=" %%f in ('dir /b "%PARENT_DIR%\*.cbuild.yml" 2^>nul') do (
    set "CBUILD_FILE=%PARENT_DIR%\%%f"
)

if not defined CBUILD_FILE (
    echo Error: No *.cbuild.yml file found in %PARENT_DIR%
    exit /b 1
)
if not exist "!CBUILD_FILE!" (
    echo Error: cbuild file not found at !CBUILD_FILE!
    exit /b 1
)

:: Log the cbuild file being processed
echo Processing cbuild file: !CBUILD_FILE!

:: Parse device-pack from cbuild.yml
set "DEVICE_PACK="
for /f "delims=" %%a in ('findstr /C:"device-pack:" "!CBUILD_FILE!"') do (
    set "DEVICE_PACK=%%a"
    set "DEVICE_PACK=!DEVICE_PACK:device-pack:=!"
    set "DEVICE_PACK=!DEVICE_PACK: =!"
)
if not defined DEVICE_PACK (
    echo Error: Could not parse device-pack from !CBUILD_FILE!
    exit /b 1
)

:: Extract pack name and version
set "PACK_NAME="
set "PACK_VERSION="
for /f "tokens=1,2 delims=@" %%a in ("!DEVICE_PACK!") do (
    set "PACK_NAME=%%a"
    set "PACK_VERSION=%%b"
)
if not defined PACK_NAME (
    echo Error: Could not parse pack name from device-pack: !DEVICE_PACK!
    exit /b 1
)
if not defined PACK_VERSION (
    echo Error: Could not parse pack version from device-pack: !DEVICE_PACK!
    exit /b 1
)

:: Remove vendor prefix (Infineon::) from pack name
set "PACK_NAME_PATH=!PACK_NAME:Infineon::=!"
echo Parsed pack name: !PACK_NAME!, version: !PACK_VERSION! >> debug.log

:: Construct pack path
set "PACK_PATH=%USERPROFILE%\AppData\Local\Arm\Packs\Infineon\!PACK_NAME_PATH!\!PACK_VERSION!"
set "PACK_PATH=!PACK_PATH:/=\!"
echo Checking pack path: !PACK_PATH! >> debug.log

:: Verify pack path exists
if not exist "!PACK_PATH!" (
    echo Error: Pack path not found at !PACK_PATH!
    exit /b 1
)

:: Check for Libraries folder
set "LIBRARIES_PATH=!PACK_PATH!\Libraries"
if not exist "!LIBRARIES_PATH!" (
    echo Error: Libraries folder not found at !LIBRARIES_PATH!
    exit /b 1
)
echo Found Libraries folder: !LIBRARIES_PATH! >> debug.log

:: Check for mtb-pdl-cat* directories
set "MTB_PDL_FOUND_LOCAL=0"
for /d %%d in ("!LIBRARIES_PATH!\mtb-pdl-cat*") do (
    set "MTB_PDL_DIRS=!MTB_PDL_DIRS!%%d;"
    echo Found mtb-pdl-cat* directory: %%d >> debug.log
    if exist "%%d\props.json" (
        set "LIBRARY_PATH=!LIBRARY_PATH!%%d\props.json,"
        set "MTB_PDL_FOUND=1"
        set "MTB_PDL_FOUND_LOCAL=1"
        echo Found mtb-pdl props.json: %%d\props.json >> debug.log
    ) else (
        set "MISSING_FILES=!MISSING_FILES!mtb-pdl-cat* props.json at %%d\props.json;"
        echo Missing mtb-pdl-cat* props.json: %%d\props.json >> debug.log
    )
)
if !MTB_PDL_FOUND_LOCAL! equ 0 (
    set "MISSING_FILES=!MISSING_FILES!mtb-pdl-cat* directory at !LIBRARIES_PATH!\mtb-pdl-cat*;"
    echo No mtb-pdl-cat* directory found at: !LIBRARIES_PATH!\mtb-pdl-cat* >> debug.log
)

:: Check for device-info\device-db\props.json
set "DEVICE_DB_PATH=!PACK_PATH!\device-info\device-db"
if exist "!DEVICE_DB_PATH!\" (
    set "DEVICE_DB_DIRS=!DEVICE_DB_PATH!;"
    echo Found device-db directory: !DEVICE_DB_PATH! >> debug.log
    if exist "!DEVICE_DB_PATH!\props.json" (
        set "LIBRARY_PATH=!LIBRARY_PATH!!DEVICE_DB_PATH!\props.json,"
        set "DEVICE_DB_FOUND=1"
        echo Found device-db props.json: !DEVICE_DB_PATH!\props.json >> debug.log
    ) else (
        set "MISSING_FILES=!MISSING_FILES!device-db props.json at !DEVICE_DB_PATH!\props.json;"
        echo Missing device-db props.json: !DEVICE_DB_PATH!\props.json >> debug.log
    )
) else (
    set "MISSING_FILES=!MISSING_FILES!device-db directory at !DEVICE_DB_PATH!;"
    echo No device-db directory found at: !DEVICE_DB_PATH! >> debug.log
)

:: Echo mtb-pdl-cat* and device-db directories
echo Mtb-pdl-cat* directories found: !MTB_PDL_DIRS:;=, ! >> debug.log
echo Device-db directories found: !DEVICE_DB_DIRS:;=, ! >> debug.log

:: Check if any props.json files were found
if not defined LIBRARY_PATH (
    echo Error: No props.json files found in mtb-pdl-cat* or device-db directories
    if defined MISSING_FILES (
        echo Missing library files: !MISSING_FILES:;=, !
    )
    exit /b 1
)

:: Report missing files if any
if "!MTB_PDL_FOUND!"=="0" (
    echo Warning: No mtb-pdl-cat* props.json found
)
if "!DEVICE_DB_FOUND!"=="0" (
    echo Warning: No device-db props.json found
)
if defined MISSING_FILES (
    echo Missing library files: !MISSING_FILES:;=, !
)

:: Remove trailing comma from LIBRARY_PATH
set "LIBRARY_PATH=%LIBRARY_PATH:~0,-1%"
echo LIBRARY_PATH set to: %LIBRARY_PATH%
endlocal & set "LIBRARY_PATH=%LIBRARY_PATH%"

:: Step 6: Verify tool and design file existence
if not exist "%TOOL_PATH%" (
    echo Error: Infineon Device Configurator not found at %TOOL_PATH%
    exit /b 1
)
if not exist "%DESIGN_PATH%" (
    echo Error: Design file not found at %DESIGN_PATH%
    exit /b 1
)

:: Step 7: Check for GeneratedSource directory and existing files
echo Checking for existing files in GeneratedSource... >> debug.log
setlocal EnableDelayedExpansion
if not exist "%GEN_DIR%" (
    echo GeneratedSource folder not found at %GEN_DIR%. It will be created when files are generated. >> debug.log
    echo generator-import: > "%TMP_FILE%"
    echo   generated-by: 'Infineon Device configurator' >> "%TMP_FILE%"
    echo   groups: >> "%TMP_FILE%"
    echo   - group: ConfigTools >> "%TMP_FILE%"
    echo     files: >> "%TMP_FILE%"
) else (
    set "FILE_COUNT=0"
    for %%f in ("%GEN_DIR%\*.c" "%GEN_DIR%\*.h") do (
        set /a FILE_COUNT+=1
    )
    echo generator-import: > "%TMP_FILE%"
    echo   generated-by: 'Infineon Device configurator' >> "%TMP_FILE%"
    echo   groups: >> "%TMP_FILE%"
    echo   - group: ConfigTools >> "%TMP_FILE%"
    echo     files: >> "%TMP_FILE%"
    if !FILE_COUNT! gtr 0 (
        echo Existing files found in %GEN_DIR%. Adding to %YAML_FILE%... >> debug.log
        for %%f in ("%GEN_DIR%\*.c" "%GEN_DIR%\*.h") do (
            set "REL_PATH=%%f"
            set "REL_PATH=!REL_PATH:%DEVICE_CONFIG_PATH%\=!"
            set "REL_PATH=!REL_PATH:\=/!"
            echo     - file: !REL_PATH!>> "%TMP_FILE%"
        )
    )
)
if exist "%TMP_FILE%" (
    move /Y "%TMP_FILE%" "%YAML_FILE%"
    if !ERRORLEVEL! equ 0 (
        echo %YAML_FILE% updated with existing files successfully! >> debug.log
    ) else (
        echo Error: Failed to move %TMP_FILE% to %YAML_FILE%
        exit /b 1
    )
) else (
    echo Error: Temporary file %TMP_FILE% was not created!
    exit /b 1
)
endlocal

:: Step 8: Get initial timestamp of design.modus
for %%F in ("%DESIGN_PATH%") do set "INITIAL_TIMESTAMP=%%~tF"
echo Initial timestamp of design.modus: %INITIAL_TIMESTAMP%

:: Step 9: Launch Device Configurator with LIBRARY_PATH as a single comma-separated argument
echo Building Device Configurator command...
setlocal EnableDelayedExpansion
echo DEBUG: LIBRARY_PATH: %LIBRARY_PATH% >> debug.log
echo Launching Infineon Device Configurator with command: "%TOOL_PATH%" --library "%LIBRARY_PATH%" --design "%DESIGN_PATH%"
start /B "" "%TOOL_PATH%" --library "%LIBRARY_PATH%" --design "%DESIGN_PATH%"
endlocal

:: Step 10: Wait for Device Configurator to start and then close, updating cgen.yml on design.modus saves
echo Waiting for Infineon Device Configurator to start...
setlocal EnableDelayedExpansion
set "ATTEMPTS=10"
set "COUNT=0"
:start_loop
if !COUNT! geq %ATTEMPTS% (
    echo Error: Infineon Device Configurator failed to start within %ATTEMPTS%
    exit /b 1
)
ping 127.0.0.1 -n 2 -w 500 >nul 2>nul
tasklist /FI "IMAGENAME eq device-configurator.exe" 2>nul | find /I "device-configurator.exe" >nul
if %ERRORLEVEL% equ 0 (
    echo Infineon Device Configurator has started.
    goto wait_loop
)
set /a COUNT+=1
goto start_loop

:wait_loop
ping 127.0.0.1 -n 1 -w 200 >nul 2>nul
for %%F in ("%DESIGN_PATH%") do set "NEW_TIMESTAMP=%%~tF"
if "!NEW_TIMESTAMP!" neq "!INITIAL_TIMESTAMP!" (
    echo design.modus was saved with new timestamp: !NEW_TIMESTAMP!
    call :UpdateCgenYml
    set "INITIAL_TIMESTAMP=!NEW_TIMESTAMP!"
    ping 127.0.0.1 -n 1 -w 500 >nul 2>nul
)
tasklist /FI "IMAGENAME eq device-configurator.exe" 2>nul | find /I "device-configurator.exe" >nul
if %ERRORLEVEL% equ 0 (
    goto wait_loop
)
echo Infineon Device Configurator has closed.
call :UpdateCgenYml
endlocal
goto :eof

:: Subroutine to update cgen.yml
:UpdateCgenYml
setlocal EnableDelayedExpansion
echo Checking for generated files... >> debug.log
if not exist "%GEN_DIR%" (
    echo Error: GeneratedSource folder not found at %GEN_DIR%!
    exit /b 1
)
set "FILE_COUNT=0"
for %%f in ("%GEN_DIR%\*.c" "%GEN_DIR%\*.h") do (
    set /a FILE_COUNT+=1
)
if !FILE_COUNT! equ 0 (
    echo Error: No .c or .h files found in %GEN_DIR%!
    exit /b 1
)
echo Updating %YAML_FILE% with generated files... >> debug.log
echo generator-import: > "%TMP_FILE%"
echo   generated-by: 'Infineon Device configurator' >> "%TMP_FILE%"
echo   groups: >> "%TMP_FILE%"
echo   - group: ConfigTools >> "%TMP_FILE%"
echo     files: >> "%TMP_FILE%"
for %%f in ("%GEN_DIR%\*.c" "%GEN_DIR%\*.h") do (
    set "REL_PATH=%%f"
    set "REL_PATH=!REL_PATH:%DEVICE_CONFIG_PATH%\=!"
    set "REL_PATH=!REL_PATH:\=/!"
    echo     - file: !REL_PATH!>> "%TMP_FILE%"
)
if exist "%TMP_FILE%" (
    move /Y "%TMP_FILE%" "%YAML_FILE%"
    if !ERRORLEVEL! equ 0 (
        echo %YAML_FILE% updated with generated files successfully! >> debug.log
    ) else (
        echo Error: Failed to move %TMP_FILE% to %YAML_FILE%
        exit /b 1
    )
) else (
    echo Error: Temporary file %TMP_FILE% was not created!
    exit /b 1
)
endlocal
goto :eof

:: Function to compare versions (returns 1 if version1 > version2, 0 otherwise)
:CompareVersions
set "version1=%~1"
set "version2=%~2"
for /f "tokens=1-3 delims=." %%a in ("!version1!") do (
    set "v1_major=%%a"
    set "v1_minor=%%b"
    set "v1_patch=%%c"
)
for /f "tokens=1-3 delims=." %%a in ("!version2!") do (
    set "v2_major=%%a"
    set "v2_minor=%%b"
    set "v2_patch=%%c"
)
:: Pad with zeros if parts are missing
if not defined v1_major set "v1_major=0"
if not defined v1_minor set "v1_minor=0"
if not defined v1_patch set "v1_patch=0"
if not defined v2_major set "v2_major=0"
if not defined v2_minor set "v2_minor=0"
if not defined v2_patch set "v2_patch=0"
:: Compare major, minor, patch
if !v1_major! gtr !v2_major! exit /b 1
if !v1_major! lss !v2_major! exit /b 0
if !v1_minor! gtr !v2_minor! exit /b 1
if !v1_minor! lss !v2_minor! exit /b 0
if !v1_patch! gtr !v2_patch! exit /b 1
exit /b 0
