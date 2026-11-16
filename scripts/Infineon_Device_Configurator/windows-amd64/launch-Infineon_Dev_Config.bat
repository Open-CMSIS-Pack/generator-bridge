@echo off
:: Copyright 2025 - Infineon Technologies
setlocal EnableExtensions EnableDelayedExpansion
if "%~1"=="" goto usage
if /i "%~1"=="--help" goto usage
set "INDEX=%~f1"
set "BASE=%~dp1"
set "MODE=%~2"
set "DQ=""
set "ERROR="
set "CONTEXT_COUNT=0"
set "PACK_COUNT=0"
set "DESIGN_COUNT=0"
set "LOCAL_COUNT=0"
set "BOARD_COUNT=0"
set "DEVICE_COUNT=0"
set "ACTIVE_GENERATOR=0"
set "UVISION_METADATA="
if not exist "!INDEX!" (
    set "ERROR=Generator index does not exist: !INDEX!"
    goto failed
)
set "INPUT_FILE=!INDEX!"
call :check_input
if defined ERROR goto failed
for /f "usebackq delims=" %%L in ("!INDEX!") do (
    set "LINE=%%L"
    if not defined ERROR call :read_index
)
if defined ERROR goto failed
if not "!CONTEXT_COUNT!"=="1" (
    set "ERROR=Select exactly one InfineonDeviceConfigurator context before launching."
    goto failed
)
if not defined METADATA set "ERROR=Missing cbuild-gen path in generator index."
if not defined CGEN set "ERROR=Missing generator-import name in generator index."
if defined ERROR goto failed
for %%P in ("!METADATA!") do set "BASE=%%~dpP"
set "INPUT_FILE=!METADATA!"
call :check_input
if defined ERROR goto failed
set "SECTION="
set "OWNER="
set "PENDING_PATH="
for /f "usebackq delims=" %%L in ("!METADATA!") do (
    set "LINE=%%L"
    if not defined ERROR call :read_metadata
)
call :flush_design
if defined ERROR goto failed
if "!UVISION_METADATA!"=="1" (
    if not defined PROJECT_DIR set "PROJECT_DIR=!BASE!"
    if not defined SOLUTION_DIR set "SOLUTION_DIR=!BASE!"
)
if not defined DEVICE set "ERROR=Selected device is missing from generator metadata."
if not defined PROJECT_DIR set "ERROR=Project path is missing from generator metadata."
if not defined SOLUTION_DIR set "ERROR=Solution path is missing from generator metadata."
if defined ERROR goto failed
set "PROVIDER_COUNT=0"
set "PREFERRED_COUNT=0"
for /l %%N in (1,1,!PACK_COUNT!) do (
    set "ROOT=!PACK_ROOT_%%N!"
    set "PACK_ID=!PACK_ID_%%N!"
    call :find_libraries
)
if !PREFERRED_COUNT! gtr 0 (
    set "PROVIDER_COUNT=!PREFERRED_COUNT!"
    set "LIBRARY_ROOT=!PREFERRED_ROOT!"
    set "LIBRARY_PACK=!PREFERRED_PACK!"
    set "LIBRARIES=!PREFERRED_LIBRARIES!"
)
if not "!PROVIDER_COUNT!"=="1" (
    set "ERROR=Expected exactly one selected Configurator library provider; found !PROVIDER_COUNT!."
    goto failed
)
for %%P in ("!CGEN!") do set "CONTEXT_DIR=%%~dpP"
set "DESTINATION=!CONTEXT_DIR!design.modus"
if exist "!DESTINATION!" (
    set "CANDIDATE=!DESTINATION!"
    call :matches_device
    if not "!MATCH!"=="1" (
        set "ERROR=Existing context design does not target !DEVICE!; it will not be overwritten: !DESTINATION!"
        goto failed
    )
    set "SOURCE=!DESTINATION!"
    set "SOURCE_KIND=existing project design"
    goto selected
)
for %%P in ("!PROJECT_DIR!design.modus" "!SOLUTION_DIR!design.modus") do (
    if exist "%%~fP" (
        set "CANDIDATE=%%~fP"
        set "CANDIDATE_KIND=LOCAL"
        call :add_candidate
    )
)
for /l %%N in (1,1,!DESIGN_COUNT!) do (
    set "VALUE=!DESIGN_RAW_%%N!"
    set "FILE_OWNER=!DESIGN_OWNER_%%N!"
    call :consider_design
)
if defined ERROR goto failed
if !LOCAL_COUNT! gtr 0 (
    if not "!LOCAL_COUNT!"=="1" (
        set "ERROR=Ambiguous project designs: !LOCAL_COUNT! candidates."
        goto failed
    )
    set "SOURCE=!LOCAL_PATH_1!"
    set "SOURCE_KIND=project design"
    set "CANDIDATE=!SOURCE!"
    call :matches_device
    if not "!MATCH!"=="1" set "ERROR=Project design does not target !DEVICE!; it will not be overwritten."
) else if !BOARD_COUNT! gtr 0 (
    if not "!BOARD_COUNT!"=="1" (
        set "ERROR=Ambiguous selected BSP templates: !BOARD_COUNT! candidates."
        goto failed
    )
    set "SOURCE=!BOARD_PATH_1!"
    set "SOURCE_KIND=selected BSP template"
) else (
    if not "!DEVICE_COUNT!"=="1" (
        set "ERROR=No unique device template for !DEVICE!: found !DEVICE_COUNT! candidates."
        goto failed
    )
    set "SOURCE=!DEVICE_PATH_1!"
    set "SOURCE_KIND=device template"
)
if defined ERROR goto failed
:selected
echo Selected !SOURCE_KIND!: !SOURCE!
echo Configurator libraries: !LIBRARY_PACK!
if /i "!MODE!"=="--plan" (
    echo PLAN_DEVICE=!DEVICE!
    echo PLAN_SOURCE=!SOURCE!
    echo PLAN_DESTINATION=!DESTINATION!
    echo PLAN_SOURCE_KIND=!SOURCE_KIND!
    echo PLAN_CGEN=!CGEN!
    echo PLAN_LIBRARIES=!LIBRARIES!
    echo PLAN_LIBRARY_PACK=!LIBRARY_PACK!
    echo PLAN_LIBRARY_ROOT=!LIBRARY_ROOT!
    exit /b 0
)
if /i "!MODE!"=="--check-paths" (
    call :find_tool
    if defined ERROR goto failed
    echo Configurator executable: !TOOL!
    exit /b 0
)
if /i "!MODE!"=="--update-import" goto import_only
if defined MODE if /i not "!MODE!"=="--prepare-only" (
    set "ERROR=Unknown launcher option: !MODE!"
    goto failed
)
if /i not "!MODE!"=="--prepare-only" (
    call :find_tool
    if defined ERROR goto failed
    call :gui_running
    if "!GUI_RUNNING!"=="1" (
        set "ERROR=Close the existing Device Configurator session before launching another context."
        goto failed
    )
)
if /i not "!SOURCE!"=="!DESTINATION!" (
    if exist "!DESTINATION!" (
        set "ERROR=Destination appeared after selection; it will not be overwritten."
        goto failed
    )
    if not exist "!CONTEXT_DIR!" mkdir "!CONTEXT_DIR!"
    if not exist "!CONTEXT_DIR!" (
        set "ERROR=Cannot create Configurator output directory."
        goto failed
    )
    copy /-y "!SOURCE!" "!DESTINATION!" <nul >nul
    if errorlevel 1 (
        set "ERROR=Cannot copy the selected design without overwriting a file."
        goto failed
    )
)
if /i "!MODE!"=="--prepare-only" exit /b 0
call :update_import
if defined ERROR goto failed
start "" /b "!TOOL!" --library "!LIBRARIES!" --design "!DESTINATION!"
set "START_ATTEMPTS=0"
:wait_start
call :gui_running
if "!GUI_RUNNING!"=="1" goto wait_close
set /a START_ATTEMPTS+=1
if !START_ATTEMPTS! geq 30 (
    set "ERROR=Device Configurator did not start; inspect launch-tool.log."
    goto failed
)
timeout /t 1 /nobreak >nul 2>nul
goto wait_start
:wait_close
echo Device Configurator started. Save/generate and close it to import sources.
:watch_gui
call :update_import
if defined ERROR goto failed
timeout /t 1 /nobreak >nul 2>nul
call :gui_running
if "!GUI_RUNNING!"=="1" goto watch_gui
:import_only
call :update_import
if defined ERROR goto failed
exit /b 0

:usage
echo Usage: %~nx0 index.cbuild-gen-idx.yml [--plan^|--check-paths^|--prepare-only^|--update-import]
exit /b 0
:failed
echo Device Configurator was not launched: !ERROR! 1>&2
exit /b 1

:check_input
if not exist "!INPUT_FILE!" (
    set "ERROR=Generator metadata does not exist: !INPUT_FILE!"
    exit /b
)
setlocal DisableDelayedExpansion
findstr /r /c:"^[ ]*path:" /c:"^[ ]*- file:" /c:"^[ ]*project:" /c:"^[ ]*solution:" /c:"^[ ]*device:" /c:"^[ ]*from-pack:" /c:"^[ ]*board-pack:" /c:"^[ ]*- cbuild-gen:" /c:"^          name:" "%INPUT_FILE%" | findstr /l /c:"!" >nul 2>nul
if not errorlevel 1 (
    endlocal
    set "ERROR=Unsupported exclamation mark in generator metadata."
    exit /b
)
endlocal
exit /b

:split_line
set "TEXT=!LINE!"
set "INDENT=0"
:indent_loop
if "!TEXT:~0,1!"==" " (
    set "TEXT=!TEXT:~1!"
    set /a INDENT+=1
    goto indent_loop
)
set "IS_LIST=0"
if "!TEXT:~0,2!"=="- " (
    set "IS_LIST=1"
    set "TEXT=!TEXT:~2!"
)
set "KEY="
set "VALUE="
for /f "tokens=1,* delims=:" %%K in ("!TEXT!") do (
    set "KEY=%%K"
    set "VALUE=%%L"
)
for /f "tokens=* delims= " %%V in ("!VALUE!") do set "VALUE=%%V"
exit /b

:scalar
if not defined VALUE (
    set "ERROR=Empty scalar for !KEY!."
    exit /b
)
if "!VALUE:~0,1!"=="!DQ!" (
    if not "!VALUE:~-1!"=="!DQ!" (
        set "ERROR=Unterminated quoted scalar for !KEY!."
        exit /b
    )
    set "VALUE=!VALUE:~1,-1!"
    set "VALUE=!VALUE:\\=\!"
) else if "!VALUE:~0,1!"=="'" (
    if not "!VALUE:~-1!"=="'" (
        set "ERROR=Unterminated quoted scalar for !KEY!."
        exit /b
    )
    set "VALUE=!VALUE:~1,-1!"
    set "VALUE=!VALUE:''='!"
) else (
    if not "!VALUE: #=!"=="!VALUE!" set "ERROR=Unsupported inline YAML comment for !KEY!."
)
set "QUOTE_CHECK=!VALUE:"=!"
if not "!QUOTE_CHECK!"=="!VALUE!" set "ERROR=Unsupported embedded quote for !KEY!."
if "!VALUE:~0,1!"=="[" set "ERROR=Flow YAML is not supported for !KEY!."
if "!VALUE:~0,1!"=="{" set "ERROR=Flow YAML is not supported for !KEY!."
if "!VALUE:~0,1!"=="*" set "ERROR=YAML aliases are not supported for !KEY!."
if "!VALUE:~0,1!"=="&" set "ERROR=YAML anchors are not supported for !KEY!."
exit /b

:resolve
set "PATH_VALUE=!VALUE:/=\!"
if "!PATH_VALUE:~0,18!"=="${CMSIS_PACK_ROOT}" (
    if not defined CMSIS_PACK_ROOT (
        set "ERROR=CMSIS_PACK_ROOT is not set for the metadata path."
        exit /b
    )
    set "PATH_VALUE=!CMSIS_PACK_ROOT!!PATH_VALUE:~18!"
)
if not "!PATH_VALUE:${=!"=="!PATH_VALUE!" (
    set "ERROR=Unsupported unresolved variable in metadata path."
    exit /b
)
if "!PATH_VALUE:~1,1!"==":" (
    for %%P in ("!PATH_VALUE!") do set "RESOLVED=%%~fP"
) else if "!PATH_VALUE:~0,2!"=="\\" (
    for %%P in ("!PATH_VALUE!") do set "RESOLVED=%%~fP"
) else (
    for %%P in ("!BASE!\!PATH_VALUE!") do set "RESOLVED=%%~fP"
)
exit /b

:read_index
call :split_line
if !INDENT! equ 4 if "!IS_LIST!"=="1" set "ACTIVE_GENERATOR=0"
if !INDENT! equ 4 if "!IS_LIST!"=="1" if "!KEY!"=="id" (
    call :scalar
    if "!VALUE!"=="InfineonDeviceConfigurator" set "ACTIVE_GENERATOR=1"
)
if !INDENT! equ 6 if "!IS_LIST!"=="0" if "!KEY!"=="id" (
    call :scalar
    set "ACTIVE_GENERATOR=0"
    if "!VALUE!"=="InfineonDeviceConfigurator" set "ACTIVE_GENERATOR=1"
)
if not "!ACTIVE_GENERATOR!"=="1" exit /b
if !INDENT! equ 8 if "!IS_LIST!"=="1" if "!KEY!"=="cbuild-gen" (
    set /a CONTEXT_COUNT+=1
    call :scalar
    call :resolve
    set "METADATA=!RESOLVED!"
)
if !INDENT! equ 10 if "!KEY!"=="name" (
    call :scalar
    call :resolve
    set "CGEN=!RESOLVED!"
    for %%P in ("!CGEN!") do if /i not "%%~xP"==".yml" set "ERROR=Expected a generator-import YAML path."
)
exit /b

:read_metadata
call :split_line
if !INDENT! equ 2 (
    call :flush_design
    set "SECTION=!KEY!"
    set "OWNER="
    for %%S in (packs components groups current-generator) do if "!KEY!"=="%%S" if defined VALUE set "ERROR=Expected block YAML for !KEY!."
    if "!KEY!"=="device" (
        call :scalar
        set "DEVICE=!VALUE:*::=!"
        for /f "tokens=1 delims=:" %%D in ("!DEVICE!") do set "DEVICE=%%D"
    )
    if "!KEY!"=="generated-by" (
        call :scalar
        if "!VALUE:~0,16!"=="uVision version " set "UVISION_METADATA=1"
    )
    if "!KEY!"=="board-pack" if defined VALUE (
        call :scalar
        set "BOARD_PACK=!VALUE!"
    )
    if "!KEY!"=="project" (
        call :scalar
        call :resolve
        for %%P in ("!RESOLVED!") do set "PROJECT_DIR=%%~dpP"
    )
    if "!KEY!"=="solution" (
        call :scalar
        call :resolve
        for %%P in ("!RESOLVED!") do set "SOLUTION_DIR=%%~dpP"
    )
    exit /b
)
if "!SECTION!"=="current-generator" if "!KEY!"=="from-pack" (
    call :scalar
    set "GENERATOR_OWNER=!VALUE!"
)
if "!SECTION!"=="packs" (
    if !INDENT! equ 4 if "!KEY!"=="pack" if "!IS_LIST!"=="1" (
        call :scalar
        set /a PACK_COUNT+=1
        set "PACK_ID_!PACK_COUNT!=!VALUE!"
    )
    if !INDENT! equ 6 if "!KEY!"=="path" (
        call :scalar
        call :resolve
        set "PACK_ROOT_!PACK_COUNT!=!RESOLVED!"
    )
    exit /b
)
if not "!SECTION!"=="components" if not "!SECTION!"=="groups" exit /b
if "!SECTION!"=="components" if !INDENT! equ 8 if "!KEY!"=="id" (
    call :scalar
    if "!VALUE!"=="InfineonDeviceConfigurator" if not defined GENERATOR_OWNER set "GENERATOR_OWNER=!OWNER!"
)
if !INDENT! equ 4 if "!IS_LIST!"=="1" (
    call :flush_design
    set "OWNER="
)
if "!KEY!"=="from-pack" (
    call :scalar
    if defined PENDING_PATH if !INDENT! gtr !PENDING_INDENT! (
        set "PENDING_OWNER=!VALUE!"
        exit /b
    )
    set "OWNER=!VALUE!"
)
if "!KEY!"=="file" if "!IS_LIST!"=="1" (
    if !INDENT! lss 8 (
        set "ERROR=Unsupported YAML file-list indentation; use CMSIS block YAML."
        exit /b
    )
    call :flush_design
    set "PENDING_PATH=!VALUE!"
    set "PENDING_OWNER=!OWNER!"
    set "PENDING_INDENT=!INDENT!"
)
exit /b

:flush_design
if not defined PENDING_PATH exit /b
set /a DESIGN_COUNT+=1
set "DESIGN_RAW_!DESIGN_COUNT!=!PENDING_PATH!"
set "DESIGN_OWNER_!DESIGN_COUNT!=!PENDING_OWNER!"
set "PENDING_PATH="
exit /b

:pack_matches
set "PACK_MATCH=0"
if not defined PACK_LEFT exit /b
if not defined PACK_RIGHT exit /b
if /i "!PACK_LEFT!"=="!PACK_RIGHT!" (
    set "PACK_MATCH=1"
    exit /b
)
if not "!PACK_LEFT:@=!"=="!PACK_LEFT!" if not "!PACK_RIGHT:@=!"=="!PACK_RIGHT!" exit /b
for /f "tokens=1 delims=@" %%P in ("!PACK_LEFT!") do set "PACK_LEFT_BASE=%%P"
for /f "tokens=1 delims=@" %%P in ("!PACK_RIGHT!") do set "PACK_RIGHT_BASE=%%P"
if /i "!PACK_LEFT_BASE!"=="!PACK_RIGHT_BASE!" set "PACK_MATCH=1"
exit /b

:find_libraries
if not exist "!ROOT!\device-info\device-db\props.json" exit /b
set "FOUND_LIBRARIES="
for /d %%P in ("!ROOT!\Libraries\mtb-pdl-cat*") do if exist "%%P\props.json" set "FOUND_LIBRARIES=!FOUND_LIBRARIES!%%P\props.json,"
if not defined FOUND_LIBRARIES exit /b
set "FOUND_LIBRARIES=!FOUND_LIBRARIES!!ROOT!\device-info\device-db\props.json"
set /a PROVIDER_COUNT+=1
set "LIBRARY_ROOT=!ROOT!"
set "LIBRARY_PACK=!PACK_ID!"
set "LIBRARIES=!FOUND_LIBRARIES!"
set "PACK_LEFT=!PACK_ID!"
set "PACK_RIGHT=!GENERATOR_OWNER!"
call :pack_matches
if "!PACK_MATCH!"=="1" (
    set /a PREFERRED_COUNT+=1
    set "PREFERRED_ROOT=!ROOT!"
    set "PREFERRED_PACK=!PACK_ID!"
    set "PREFERRED_LIBRARIES=!FOUND_LIBRARIES!"
)
exit /b

:matches_device
set "MATCH=0"
set "INVALID_MPN="
if not exist "!CANDIDATE!" exit /b
for /f "usebackq delims=" %%L in (`findstr /l /c:"mpn=" "!CANDIDATE!" 2^>nul`) do (
    set "MPN_SCAN=%%L"
    call :scan_mpns
)
if defined INVALID_MPN set "MATCH=0"
exit /b

:scan_mpns
set "MPN_REST=!MPN_SCAN:*mpn=!"
if "!MPN_REST!"=="!MPN_SCAN!" exit /b
if not "!MPN_REST:~0,1!"=="=" (
    set "INVALID_MPN=1"
    exit /b
)
set "MPN_SCAN=!MPN_REST:~1!"
if not "!MPN_SCAN:~0,1!"=="!DQ!" if not "!MPN_SCAN:~0,1!"=="'" (
    set "INVALID_MPN=1"
    exit /b
)
set "MPN_ATTRIBUTE=!MPN_SCAN:"=#!"
set "MPN_ATTRIBUTE=!MPN_ATTRIBUTE:'=#!"
for /f "tokens=1 delims=#" %%D in ("!MPN_ATTRIBUTE!") do (
    if "%%D"=="!DEVICE!" (set "MATCH=1") else set "INVALID_MPN=1"
)
goto scan_mpns

:consider_design
call :scalar
if defined ERROR exit /b
for %%P in ("!VALUE!") do if /i not "%%~xP"==".modus" exit /b
call :resolve
if defined ERROR exit /b
set "CANDIDATE=!RESOLVED!"
if not exist "!CANDIDATE!" exit /b
if not defined FILE_OWNER (
    for /l %%P in (1,1,!PACK_COUNT!) do (
        set "PREFIX=!PACK_ROOT_%%P!\"
        set "PREFIX_LENGTH=0"
        set "PREFIX_REST=!PREFIX!"
        call :prefix_length
        for %%S in (!PREFIX_LENGTH!) do if /i "!CANDIDATE:~0,%%S!"=="!PREFIX!" set "FILE_OWNER=!PACK_ID_%%P!"
    )
)
if not defined FILE_OWNER (
    set "CANDIDATE_KIND=LOCAL"
    call :add_candidate
    exit /b
)
call :matches_device
if not "!MATCH!"=="1" exit /b
set "PACK_LEFT=!FILE_OWNER!"
set "PACK_RIGHT=!BOARD_PACK!"
call :pack_matches
if "!PACK_MATCH!"=="1" (
    set "CANDIDATE_KIND=BOARD"
    call :add_candidate
    exit /b
)
set "PACK_RIGHT=!LIBRARY_PACK!"
call :pack_matches
if "!PACK_MATCH!"=="1" (
    set "CANDIDATE_KIND=DEVICE"
    call :add_candidate
)
exit /b

:prefix_length
if not defined PREFIX_REST exit /b
set "PREFIX_REST=!PREFIX_REST:~1!"
set /a PREFIX_LENGTH+=1
goto prefix_length

:add_candidate
for %%K in (!CANDIDATE_KIND!) do (
    for /l %%N in (1,1,!%%K_COUNT!) do if /i "!%%K_PATH_%%N!"=="!CANDIDATE!" exit /b
    set /a %%K_COUNT+=1
    for %%N in (!%%K_COUNT!) do set "%%K_PATH_%%N=!CANDIDATE!"
)
exit /b

:find_tool
set "TOOL="
set "BEST_VERSION=0.0.0.0"
if defined CONFIGURATOR_EXE (
    set "TOOL=!CONFIGURATOR_EXE!"
    if not exist "!TOOL!" set "ERROR=CONFIGURATOR_EXE does not exist: !TOOL!"
    exit /b
)
for %%R in ("!LOCALAPPDATA!\Infineon_Technologies_AG\Infineon-Toolbox" "!ALLUSERSPROFILE!\Infineon_Technologies_AG\Infineon-Toolbox") do if exist "%%~R" (
    set "REGISTRATION_ROOT=%%~R"
    call :scan_registrations
)
if not defined TOOL set "ERROR=Device Configurator is not installed or registered in Infineon Toolbox."
exit /b

:scan_registrations
for %%J in ("!REGISTRATION_ROOT!\*.json") do if exist "%%~fJ" (
    findstr /l /c:"com.ifx.tb.tool.deviceconfigurator" "%%~fJ" >nul 2>nul
    if not errorlevel 1 (
        set "REGISTRATION=%%~fJ"
        call :read_registration
    )
)
exit /b

:read_registration
set "REG_FEATURE="
set "REG_VERSION="
set "REG_EXE="
for /f "usebackq tokens=1,* delims=:" %%K in ("!REGISTRATION!") do (
    set "REG_KEY=%%K"
    set "REG_VALUE=%%L"
    set "REG_KEY=!REG_KEY: =!"
    set "REG_KEY=!REG_KEY:"=!"
    for /f "tokens=* delims= " %%V in ("!REG_VALUE!") do set "REG_VALUE=%%V"
    if "!REG_VALUE:~-1!"=="," set "REG_VALUE=!REG_VALUE:~0,-1!"
    set "REG_VALUE=!REG_VALUE:"=!"
    set "REG_VALUE=!REG_VALUE:\\=\!"
    if "!REG_KEY!"=="featureId" set "REG_FEATURE=!REG_VALUE!"
    if "!REG_KEY!"=="version" set "REG_VERSION=!REG_VALUE!"
    if "!REG_KEY!"=="exePath" set "REG_EXE=!REG_VALUE!"
)
if not "!REG_FEATURE!"=="com.ifx.tb.tool.deviceconfigurator" exit /b
if not defined REG_VERSION exit /b
if not exist "!REG_EXE!" exit /b
set "VERSION_NEW=!REG_VERSION!"
set "VERSION_OLD=!BEST_VERSION!"
call :version_newer
if "!IS_NEWER!"=="1" (
    set "BEST_VERSION=!REG_VERSION!"
    set "TOOL=!REG_EXE!"
)
exit /b

:version_newer
set "IS_NEWER=0"
set "NEW_1=0"
set "NEW_2=0"
set "NEW_3=0"
set "NEW_4=0"
set "OLD_1=0"
set "OLD_2=0"
set "OLD_3=0"
set "OLD_4=0"
for /f "tokens=1-4 delims=." %%A in ("!VERSION_NEW!") do (
    set "NEW_1=%%A"
    set "NEW_2=%%B"
    set "NEW_3=%%C"
    set "NEW_4=%%D"
)
for /f "tokens=1-4 delims=." %%A in ("!VERSION_OLD!") do (
    set "OLD_1=%%A"
    set "OLD_2=%%B"
    set "OLD_3=%%C"
    set "OLD_4=%%D"
)
for /l %%N in (1,1,4) do (
    set "CMP_NEW=00000000!NEW_%%N!"
    set "CMP_OLD=00000000!OLD_%%N!"
    set "CMP_NEW=!CMP_NEW:~-8!"
    set "CMP_OLD=!CMP_OLD:~-8!"
    if "!CMP_NEW!" GTR "!CMP_OLD!" (
        set "IS_NEWER=1"
        exit /b
    )
    if "!CMP_NEW!" LSS "!CMP_OLD!" exit /b
)
exit /b

:gui_running
set "GUI_RUNNING=0"
tasklist /fi "IMAGENAME eq device-configurator.exe" /fo csv /nh 2>nul | findstr /i /c:"device-configurator.exe" >nul
if not errorlevel 1 set "GUI_RUNNING=1"
exit /b

:update_import
if not exist "!CONTEXT_DIR!" (
    set "ERROR=Configurator output directory does not exist."
    exit /b
)
set "IMPORT_TEMP=!CGEN!.tmp"
>"!IMPORT_TEMP!" echo generator-import:
>>"!IMPORT_TEMP!" echo   generated-by: 'Infineon Device Configurator'
>>"!IMPORT_TEMP!" echo   groups:
>>"!IMPORT_TEMP!" echo   - group: ConfigTools
>>"!IMPORT_TEMP!" echo     files:
for %%E in (c h) do for %%F in ("!CONTEXT_DIR!GeneratedSource\*.%%E") do if exist "%%~fF" (
    set "GENERATED_NAME=%%~nxF"
    set "GENERATED_NAME=!GENERATED_NAME:'=''!"
    >>"!IMPORT_TEMP!" echo     - file: 'GeneratedSource/!GENERATED_NAME!'
)
move /y "!IMPORT_TEMP!" "!CGEN!" >nul
if errorlevel 1 set "ERROR=Cannot publish generated-source import: !CGEN!"
exit /b