@echo off
setlocal enabledelayedexpansion

set "SCRIPT_DIR=%~dp0"
set "REPO_ROOT=%SCRIPT_DIR%..\.."
set "BASH_EXE="

if defined MSYSTEM_PREFIX (
  if exist "%MSYSTEM_PREFIX%\usr\bin\bash.exe" (
    set "BASH_EXE=%MSYSTEM_PREFIX%\usr\bin\bash.exe"
  )
)

if not defined BASH_EXE (
  if exist "C:\msys64\usr\bin\bash.exe" (
    set "BASH_EXE=C:\msys64\usr\bin\bash.exe"
  )
)

if not defined BASH_EXE (
  if exist "C:\msys64\mingw64\bin\bash.exe" (
    set "BASH_EXE=C:\msys64\mingw64\bin\bash.exe"
  )
)

if not defined BASH_EXE (
  for %%I in (bash.exe) do set "BASH_EXE=%%~$PATH:I"
)

if not defined BASH_EXE (
  echo ERROR: Unable to locate bash.exe. Install MSYS2 or add Bash to PATH.
  exit /b 1
)

if not exist "%REPO_ROOT%\scripts\smoke\ikemen-workflow-matrix.sh" (
  echo ERROR: Unable to locate the workflow matrix script at "%REPO_ROOT%\scripts\smoke\ikemen-workflow-matrix.sh"
  exit /b 1
)

"%BASH_EXE%" "%REPO_ROOT%\scripts\smoke\ikemen-workflow-matrix.sh" %*
exit /b %ERRORLEVEL%
