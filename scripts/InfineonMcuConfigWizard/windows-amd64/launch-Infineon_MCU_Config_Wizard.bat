@echo off
:: Copyright 2025 - Infineon Technologies
setlocal EnableDelayedExpansion

:: Step 1: Validate input argument
if "%1"=="" (
    echo Error: No *.cbuild-gen-idx.yml file path provided as an argument
    exit /b 1
)
set "CBUILD_GEN_IDX=%~1"

:: Step 2: Locate MCU Config Wizard installation from registry
::   The installer writes the install path as a subkey name under
::   HKCU\Software\Infineon Technologies. All matching entries are
::   compared by version (major.minor.patch numerically; build string
::   lexicographically) to select the newest installation.
set "TOOL_BASE_PATH="
set "_best_maj=0"
set "_best_min=0"
set "_best_patch=0"
set "_best_build=0"
for /f "usebackq delims=" %%F in (`reg query "HKCU\Software\Infineon Technologies" 2^>nul ^| findstr /r /i "[/]MCU-Config-Wizard"`) do (
    set "_candidate=%%F"
    set "_candidate=!_candidate:HKEY_CURRENT_USER\Software\Infineon Technologies\=!"
    set "_pathbs=!_candidate:/=\!"
    set "_vmaj=0"
    set "_vmin=0"
    set "_vpatch=0"
    set "_vbuild=0"
    for %%G in ("!_pathbs!") do for /f "tokens=1-4 delims=." %%A in ("%%~nxG") do (
        set "_vmaj=%%A"
        set "_vmin=%%B"
        set "_vpatch=%%C"
        set "_vbuild=%%D"
    )
    set "_update=0"
    set /a "_d=!_vmaj!-!_best_maj!"
    if !_d! gtr 0 set "_update=1"
    if !_d! equ 0 (
        set /a "_d=!_vmin!-!_best_min!"
        if !_d! gtr 0 set "_update=1"
        if !_d! equ 0 (
            set /a "_d=!_vpatch!-!_best_patch!"
            if !_d! gtr 0 set "_update=1"
            if !_d! equ 0 if "!_vbuild!" gtr "!_best_build!" set "_update=1"
        )
    )
    if "!_update!"=="1" (
        set "TOOL_BASE_PATH=!_candidate!"
        set "_best_maj=!_vmaj!"
        set "_best_min=!_vmin!"
        set "_best_patch=!_vpatch!"
        set "_best_build=!_vbuild!"
    )
)
if defined TOOL_BASE_PATH echo Found: !TOOL_BASE_PATH!

if not defined TOOL_BASE_PATH (
    echo Error: MCU Config Wizard not found, please download and install from https://softwaretools.infineon.com/assets/com.ifx.tb.tool.ifxconfigwizardforembeddedpowerics
    exit /b 2
)

:: Step 3: Verify tool existence
if not exist "!TOOL_BASE_PATH!\resources\" (
    echo Error: Resources folder not found at "!TOOL_BASE_PATH!\resources\"
    exit /b 3
)
if not exist "!TOOL_BASE_PATH!\resources\launch-tool.exe" (
    echo Error: launch-tool.exe not found in "!TOOL_BASE_PATH!\resources"
    exit /b 4
)

:: Step 4: Launch MCU Config Wizard
CD /D "!TOOL_BASE_PATH!\resources"
START "" "launch-tool.exe" "--fromCmsis=!CBUILD_GEN_IDX!"
