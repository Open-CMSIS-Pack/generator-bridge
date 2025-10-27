@echo off
:: Copyright 2025 - Infineon Technologies
setlocal
:: Step 1: Validate input argument
if "%1"=="" (
    echo Error: No *.cbuild-gen-idx.yml file path provided as an argument
    exit /b 1
)
set "CBUILD_GEN_IDX=%~1"
set "CBUILD_GEN_IDX=%CBUILD_GEN_IDX:/=\%"
if not exist "%CBUILD_GEN_IDX%" (
    echo Error: %CBUILD_GEN_IDX% file not found
    exit /b 1
)

:: Step 2: Parse project-name, cbuild-gen path from *.cbuild-gen-idx.yml and device-pack with path from *.cbuild-gen.yml
setlocal EnableDelayedExpansion
set "PROJECT_NAME="
set "CBUILD_GEN_YML_PATH="
set "DESIGN_MODUS_SRC_PATH="
set "CGEN_YML_PATH="
set "DEVICE_PACK="
set "PACK_PATH="

:: Parse project name from *.cbuild-gen-idx.yml
for /f "delims=" %%a in ('findstr /C:"project:" "%CBUILD_GEN_IDX%"') do (
    set "PROJECT_NAME=%%a"
    set "PROJECT_NAME=!PROJECT_NAME:project-name:=!"
    set "PROJECT_NAME=!PROJECT_NAME: =!"
)
:: Parse cbuild-gen path from *.cbuild-gen-idx.yml
for /f "delims=" %%a in ('findstr /C:"cbuild-gen:" "%CBUILD_GEN_IDX%"') do (
    set "CBUILD_GEN_YML_PATH=%%a"
    set "CBUILD_GEN_YML_PATH=!CBUILD_GEN_YML_PATH:cbuild-gen:=!"
    set "CBUILD_GEN_YML_PATH=!CBUILD_GEN_YML_PATH:- =!"
    set "CBUILD_GEN_YML_PATH=!CBUILD_GEN_YML_PATH: =!"
    set "CBUILD_GEN_YML_PATH=!CBUILD_GEN_YML_PATH:/=\!"
)
echo DEBUG: PROJECT_NAME: !PROJECT_NAME! >> debug.log
echo DEBUG: CBUILD_GEN_YML_PATH: !CBUILD_GEN_YML_PATH! >> debug.log

if not defined PROJECT_NAME (
    echo Error: Could not parse project-name from %CBUILD_GEN_IDX%
    exit /b 1
)
if not defined CBUILD_GEN_YML_PATH (
    echo Error: Could not parse cbuild-gen path from %CBUILD_GEN_IDX%
    exit /b 1
)
if not exist "!CBUILD_GEN_YML_PATH!" (
    echo Error: *.cbuild-gen.yml file not found at !CBUILD_GEN_YML_PATH!
    exit /b 1
)

:: Parse device-pack from *.cbuild-gen.yml
for /f "delims=" %%a in ('findstr /C:"device-pack:" "!CBUILD_GEN_YML_PATH!"') do (
    set "DEVICE_PACK=%%a"
    set "DEVICE_PACK=!DEVICE_PACK:device-pack:=!"
    set "DEVICE_PACK=!DEVICE_PACK: =!"
)
echo DEBUG: DEVICE_PACK: !DEVICE_PACK! >> debug.log

if not defined DEVICE_PACK (
    echo Error: Could not parse device-pack from !CBUILD_GEN_YML_PATH!
    exit /b 1
)

:: Parse packs section to find the path for the device-pack
set "FOUND_PACK="
set "PACK_PATH="
set "IN_PACKS=0"
for /f "delims=" %%b in ('type "!CBUILD_GEN_YML_PATH!"') do (
    set "LINE=%%b"
    set "TRIMMED_LINE=!LINE: =!"
    if "!TRIMMED_LINE!"=="packs:" (
        set "IN_PACKS=1"
    ) else if !IN_PACKS! equ 1 (
        if "!TRIMMED_LINE!"=="-pack:!DEVICE_PACK!" (
            set "FOUND_PACK=1"
        ) else if defined FOUND_PACK (
            if "!TRIMMED_LINE:~0,5!"=="path:" (
                set "PACK_PATH=!LINE!"
                set "PACK_PATH=!PACK_PATH:path:=!"
                set "PACK_PATH=!PACK_PATH: =!"
                set "PACK_PATH=!PACK_PATH:/=\!"
                goto :end_pack_parse
            )
        )
    )
)
:end_pack_parse
echo DEBUG: PACK_PATH: !PACK_PATH! >> debug.log

if not defined PACK_PATH (
    echo Error: Could not parse path for device-pack !DEVICE_PACK! from !CBUILD_GEN_YML_PATH!
    exit /b 1
)
if not exist "!PACK_PATH!" (
    echo Error: Pack path not found at !PACK_PATH!
    exit /b 1
)

:: Parse cgen.yml, and design.modus paths from *.cbuild-gen.yml
for /f "delims=" %%a in ('findstr /C:"cgen.yml" "!CBUILD_GEN_YML_PATH!"') do (
    set "CGEN_YML_PATH=%%a"
    set "CGEN_YML_PATH=!CGEN_YML_PATH:path:=!"
    set "CGEN_YML_PATH=!CGEN_YML_PATH: =!"
    set "CGEN_YML_PATH=!CGEN_YML_PATH:/=\!"
)
for /f "delims=" %%a in ('findstr /C:"design.modus" "!CBUILD_GEN_YML_PATH!"') do (
    set "DESIGN_MODUS_SRC_PATH=%%a"
    set "DESIGN_MODUS_SRC_PATH=!DESIGN_MODUS_SRC_PATH:- file:=!"
    set "DESIGN_MODUS_SRC_PATH=!DESIGN_MODUS_SRC_PATH: =!"
    set "DESIGN_MODUS_SRC_PATH=!DESIGN_MODUS_SRC_PATH:/=\!"
)
echo DEBUG: CGEN_YML_PATH: !CGEN_YML_PATH! >> debug.log
echo DEBUG: DESIGN_MODUS_SRC_PATH: !DESIGN_MODUS_SRC_PATH! >> debug.log

if not defined CGEN_YML_PATH (
    echo Error: Could not parse cgen.yml path from !CBUILD_GEN_YML_PATH!
    exit /b 1
)
if not defined DESIGN_MODUS_SRC_PATH (
    echo Error: Could not parse design.modus path from !CBUILD_GEN_YML_PATH!
    exit /b 1
)

:: Check if cgen.yml directory exists, create if not
for %%i in ("!CGEN_YML_PATH!") do set "CGEN_DIR=%%~dpi"
set "CGEN_DIR=!CGEN_DIR:~0,-1!"
if not exist "!CGEN_DIR!\" (
    mkdir "!CGEN_DIR!" >nul 2>&1
    if exist "!CGEN_DIR!\" (
        echo Created directory: !CGEN_DIR! >> debug.log
    ) else (
        echo Error: Failed to create directory !CGEN_DIR! >> debug.log
        exit /b 1
    )
) else (
    echo Directory already exists: !CGEN_DIR! >> debug.log
)

:: Check if design.modus exists in the same directory as *.cgen.yml, copy if not
set "DESIGN_MODUS_DEST_PATH=!CGEN_DIR!\design.modus"
if not exist "!DESIGN_MODUS_DEST_PATH!" (
    if exist "!DESIGN_MODUS_SRC_PATH!" (
        copy "!DESIGN_MODUS_SRC_PATH!" "!DESIGN_MODUS_DEST_PATH!" >nul 2>&1
        if exist "!DESIGN_MODUS_DEST_PATH!" (
            echo Copied design.modus to !DESIGN_MODUS_DEST_PATH! >> debug.log
        ) else (
            echo Error: Failed to copy design.modus to !DESIGN_MODUS_DEST_PATH!
            exit /b 1
        )
    ) else (
        echo Error: Source design.modus not found at !DESIGN_MODUS_SRC_PATH!
        exit /b 1
    )
) else (
    echo design.modus already exists at !DESIGN_MODUS_DEST_PATH! >> debug.log
)
endlocal & set "PROJECT_NAME=%PROJECT_NAME%" & set "CGEN_YML_PATH=%CGEN_YML_PATH%" & set "PACK_PATH=%PACK_PATH%"

:: Step 3: Define paths based on *.cgen.yml path
setlocal EnableDelayedExpansion
for %%i in ("!CGEN_YML_PATH!") do set "DESIGN_DIR=%%~dpi"
set "DESIGN_DIR=!DESIGN_DIR:~0,-1!"
set "DESIGN_PATH=!DESIGN_DIR!\design.modus"
set "GEN_DIR=!DESIGN_DIR!\GeneratedSource"
set "YAML_FILE=!CGEN_YML_PATH!"
set "TMP_FILE=!YAML_FILE!.tmp"
endlocal & set "DESIGN_PATH=%DESIGN_PATH%" & set "GEN_DIR=%GEN_DIR%" & set "YAML_FILE=%YAML_FILE%" & set "TMP_FILE=%TMP_FILE%"

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
echo Finished processing per-user JSON files. >> debug.log

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
echo Finished processing all-users JSON files. >> debug.log

:: Output the latest exePath
if defined LATEST_PATH (
    echo Selected latest version: !LATEST_VERSION!
    echo Selected exePath: !LATEST_PATH!
    set "TOOL_PATH=!LATEST_PATH!"
) else (
    echo Error: Infineon Device-Configurator not found, please download and install from https://softwaretools.infineon.com/tools/com.ifx.tb.tool.deviceconfigurator
    exit /b 1
)

:: Verify TOOL_PATH was set
if "!TOOL_PATH!"=="" (
    echo Error: Infineon Device-Configurator not found, please download and install from https://softwaretools.infineon.com/tools/com.ifx.tb.tool.deviceconfigurator
    exit /b 1
)
endlocal & set "TOOL_PATH=%TOOL_PATH%"

:: Step 5: Extract props.json paths using PACK_PATH from Step 2
setlocal EnableDelayedExpansion
set "LIBRARY_PATH="
set "MTB_PDL_FOUND=0"
set "DEVICE_DB_FOUND=0"
set "MISSING_FILES="
set "MTB_PDL_DIRS="
set "DEVICE_DB_DIRS="

:: Use PACK_PATH from Step 2
echo Checking pack path from Step 2: !PACK_PATH! >> debug.log

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
    echo No mtb-pdl-cat* directory found at: !LIBRARIES_PATH!\mtb-pdl-cat* >> debug.log
)

:: Check for device-info\device-db\props.json
set "DEVICE_DB_PATH=!PACK_PATH!\device-info\device-db"
if exist "!DEVICE_DB_PATH!\" (
    set "DEVICE_DB_DIRS=!DEVICE_DB_DIRS!!DEVICE_DB_PATH!;"
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
            set "FULL_PATH=%%f"
            for %%i in ("!FULL_PATH!") do set "REL_PATH=GeneratedSource\%%~nxi"
            set "REL_PATH=!REL_PATH:\=/!"
            echo     - file: !REL_PATH! >> "%TMP_FILE%"
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
    set "FULL_PATH=%%f"
    for %%i in ("!FULL_PATH!") do set "REL_PATH=GeneratedSource\%%~nxi"
    set "REL_PATH=!REL_PATH:\=/!"
    echo     - file: !REL_PATH! >> "%TMP_FILE%"
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
