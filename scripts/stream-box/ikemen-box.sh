#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STREAM_BOX_HOST="${STREAM_BOX_HOST:-stream-box}"
REMOTE_ROOT="${IKEMEN_STREAM_BOX_ROOT:-C:\\Apps\\mugen}"
REMOTE_APPS_ROOT="${IKEMEN_STREAM_BOX_APPS_ROOT:-C:\\Apps}"
REMOTE_TOOLS_ROOT="${IKEMEN_STREAM_BOX_TOOLS_ROOT:-$REMOTE_APPS_ROOT\\ikemen-tools}"

usage() {
  cat <<'EOF'
Usage: scripts/stream-box/ikemen-box.sh <install|check|launch|launch-menu|status|tail|stop|repoint|bundle|normalize-chars|normalize-chars-apply|candidate-quality>

Installs and runs a tiny PowerShell toolbox on stream-box so the common Ikemen
operations do not need long command lines.
EOF
}

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

run_remote() {
  ssh "$STREAM_BOX_HOST" "$@"
}

remote_file() {
  local relative="$1"
  printf '%s\\%s' "$REMOTE_TOOLS_ROOT" "$relative"
}

install_toolbox() {
  local workdir=""
  workdir="$(mktemp -d "${TMPDIR:-/tmp}/ikemen-box.XXXXXX")"
  trap 'workdir_path="${workdir:-}"; [[ -n "$workdir_path" ]] && rm -rf "$workdir_path"' RETURN

  cat > "$workdir/ikemen-check.ps1" <<'EOF'
$ErrorActionPreference = 'Stop'
$root = 'C:\Apps\mugen'
$required = @(
  "$root\Ikemen_GO.exe",
  "$root\external\script\main.lua",
  "$root\external\icons\IkemenCylia_256.png",
  "$root\external\icons\IkemenCylia_96.png",
  "$root\external\icons\IkemenCylia_48.png",
  "$root\external\script\start.lua"
)
$missing = @()
foreach ($path in $required) {
  if (-not (Test-Path $path)) {
    $missing += $path
  }
}
if ($missing.Count -gt 0) {
  Write-Error ("Missing stream-box files:`n" + ($missing -join "`n"))
  exit 1
}

$def = Get-ChildItem -Path (Join-Path $root 'chars') -Recurse -Filter '*.def' -File |
  Where-Object {
    $n = $_.Name.ToLower()
    -not ($n -match 'ending|select|fight|gofx|intro')
  } |
  Select-Object -First 1
if ($null -eq $def) {
  Write-Error "No usable character .def found under $root\chars"
  exit 1
}
Write-Host "ikemen-check: ok"
EOF

  cat > "$workdir/ikemen-stop.ps1" <<'EOF'
$ErrorActionPreference = 'SilentlyContinue'
Get-Process Ikemen_GO -ErrorAction SilentlyContinue | Stop-Process -Force
Write-Host "ikemen-stop: done"
EOF

  cat > "$workdir/ikemen-launch-quickvs.ps1" <<'EOF'
$ErrorActionPreference = 'SilentlyContinue'
$root = 'C:\Apps\mugen'
$apps = 'C:\Apps'
$resultDir = Join-Path $apps 'ikemen-results'
$resultPathLog = Join-Path $apps 'ikemen-workflow-smoke.result-path.txt'
$runId = "{0}-{1}" -f (Get-Date -Format 'yyyyMMdd-HHmmss-fff'), ([Guid]::NewGuid().ToString('N').Substring(0,8))
$resultFile = Join-Path $resultDir ("ikemen-workflow-smoke.$runId.json")
$stdout = Join-Path $apps 'ikemen-workflow-smoke.stdout.log'
$stderr = Join-Path $apps 'ikemen-workflow-smoke.stderr.log'
$pidLog = Join-Path $apps 'ikemen-workflow-smoke.pid.txt'
$exitLog = Join-Path $apps 'ikemen-workflow-smoke.exit.txt'
New-Item -ItemType Directory -Force -Path $apps | Out-Null
New-Item -ItemType Directory -Force -Path $resultDir | Out-Null
New-Item -Path 'HKCU:\Software\Microsoft\Windows\Windows Error Reporting' -Force | Out-Null
Set-ItemProperty -Path 'HKCU:\Software\Microsoft\Windows\Windows Error Reporting' -Name DontShowUI -Type DWord -Value 1
$env:IKEMEN_SUPPRESS_ERROR_DIALOG = '1'
Stop-Process -Name Ikemen_GO -Force -ErrorAction SilentlyContinue
foreach ($icon in @(
  (Join-Path $root 'external\icons\IkemenCylia_256.png'),
  (Join-Path $root 'external\icons\IkemenCylia_96.png'),
  (Join-Path $root 'external\icons\IkemenCylia_48.png')
)) {
  if (-not (Test-Path $icon)) {
    Write-Error "Missing required icon asset: $icon"
    exit 1
  }
}
Remove-Item (Join-Path $root 'Ikemen.log') -Force -ErrorAction SilentlyContinue
Remove-Item $stdout -Force -ErrorAction SilentlyContinue
Remove-Item $stderr -Force -ErrorAction SilentlyContinue
Remove-Item $pidLog -Force -ErrorAction SilentlyContinue
Remove-Item $exitLog -Force -ErrorAction SilentlyContinue
$resultFile | Set-Content -Path $resultPathLog -Encoding ASCII
Remove-Item $resultFile -Force -ErrorAction SilentlyContinue
$defaultDef = Get-ChildItem -Path (Join-Path $root 'chars') -Recurse -Filter '*.def' -File |
  Where-Object {
    $n = $_.Name.ToLower()
    -not ($n -match 'ending|select|fight|gofx|intro')
  } |
  Select-Object -First 1
if ($null -eq $defaultDef) {
  Write-Error "No usable character .def found under $root\chars"
  exit 1
}
$relativeDef = $defaultDef.FullName.Substring($root.Length + 1)
$argsList = @(
  '-windowed',
  '-nojoy',
  '-nomusic',
  '-nosound',
  '-noerrordialog',
  '-jsonstdout',
  '-resultfile', $resultFile,
  '-p1', $relativeDef,
  '-p2', $relativeDef,
  '-rounds', '1',
  '-time', '5',
  '-debugstartup'
)
Push-Location $root
try {
  $proc = Start-Process -FilePath (Join-Path $root 'Ikemen_GO.exe') -WorkingDirectory $root -ArgumentList $argsList -RedirectStandardOutput $stdout -RedirectStandardError $stderr -PassThru
  $proc.Id | Set-Content -Path $pidLog -Encoding ASCII
  if ($proc.WaitForExit(15000)) {
    $code = $proc.ExitCode
    if ($null -eq $code) {
      $code = 1
    }
    $code | Set-Content -Path $exitLog -Encoding ASCII
    if ((Test-Path $stderr) -and ((Get-Content $stderr -Raw) -match '(?m)attempt to call a non-function object|runtime error:|panic:|segmentation fault|fatal error|stack traceback')) {
      exit 1
    }
    if ($code -ne 0) {
      exit $code
    }
  } else {
    'still-running' | Set-Content -Path $exitLog -Encoding ASCII
    if ((Test-Path $stderr) -and ((Get-Content $stderr -Raw) -match '(?m)attempt to call a non-function object|runtime error:|panic:|segmentation fault|fatal error|stack traceback')) {
      exit 1
    }
  }
} finally {
  Pop-Location
}
EOF

  cat > "$workdir/ikemen-launch-menu.ps1" <<'EOF'
$ErrorActionPreference = 'SilentlyContinue'
$root = 'C:\Apps\mugen'
$apps = 'C:\Apps'
$stdout = Join-Path $apps 'ikemen-menu.stdout.log'
$stderr = Join-Path $apps 'ikemen-menu.stderr.log'
$pidLog = Join-Path $apps 'ikemen-menu.pid.txt'
$exitLog = Join-Path $apps 'ikemen-menu.exit.txt'
New-Item -ItemType Directory -Force -Path $apps | Out-Null
New-Item -Path 'HKCU:\Software\Microsoft\Windows\Windows Error Reporting' -Force | Out-Null
Set-ItemProperty -Path 'HKCU:\Software\Microsoft\Windows\Windows Error Reporting' -Name DontShowUI -Type DWord -Value 1
$env:IKEMEN_SUPPRESS_ERROR_DIALOG = '1'
Stop-Process -Name Ikemen_GO -Force -ErrorAction SilentlyContinue
foreach ($icon in @(
  (Join-Path $root 'external\icons\IkemenCylia_256.png'),
  (Join-Path $root 'external\icons\IkemenCylia_96.png'),
  (Join-Path $root 'external\icons\IkemenCylia_48.png')
)) {
  if (-not (Test-Path $icon)) {
    Write-Error "Missing required icon asset: $icon"
    exit 1
  }
}
Remove-Item (Join-Path $root 'Ikemen.log') -Force -ErrorAction SilentlyContinue
Remove-Item $stdout -Force -ErrorAction SilentlyContinue
Remove-Item $stderr -Force -ErrorAction SilentlyContinue
Remove-Item $pidLog -Force -ErrorAction SilentlyContinue
Remove-Item $exitLog -Force -ErrorAction SilentlyContinue
Push-Location $root
try {
  $proc = Start-Process -FilePath (Join-Path $root 'Ikemen_GO.exe') -WorkingDirectory $root -ArgumentList @('-windowed','-nojoy','-nomusic','-nosound','-noerrordialog') -RedirectStandardOutput $stdout -RedirectStandardError $stderr -PassThru
  $proc.Id | Set-Content -Path $pidLog -Encoding ASCII
  if ($proc.WaitForExit(15000)) {
    $code = $proc.ExitCode
    if ($null -eq $code) {
      $code = 1
    }
    $code | Set-Content -Path $exitLog -Encoding ASCII
    if ((Test-Path $stderr) -and ((Get-Content $stderr -Raw) -match '(?m)attempt to call a non-function object|runtime error:|panic:|segmentation fault|fatal error|stack traceback')) {
      exit 1
    }
    if ($code -ne 0) {
      exit $code
    }
  } else {
    'still-running' | Set-Content -Path $exitLog -Encoding ASCII
    if ((Test-Path $stderr) -and ((Get-Content $stderr -Raw) -match '(?m)attempt to call a non-function object|runtime error:|panic:|segmentation fault|fatal error|stack traceback')) {
      exit 1
    }
  }
} finally {
  Pop-Location
}
EOF

  cat > "$workdir/ikemen-status.ps1" <<'EOF'
$ErrorActionPreference = 'SilentlyContinue'
$apps = 'C:\Apps'
$root = 'C:\Apps\mugen'
$link = Get-Item $root -Force -ErrorAction SilentlyContinue
if ($link) {
  Write-Host "--- root ---"
  $link | Select-Object Name,FullName,Attributes,LinkType,Target | Format-List
}
Get-Process Ikemen_GO -ErrorAction SilentlyContinue | Select-Object Id,ProcessName,Responding,StartTime | Format-Table -AutoSize
foreach ($path in @(
  "$root\Ikemen.log",
  "$apps\ikemen-workflow-smoke.stdout.log",
  "$apps\ikemen-workflow-smoke.stderr.log",
  "$apps\ikemen-workflow-smoke.exit.txt",
  "$apps\ikemen-menu.stdout.log",
  "$apps\ikemen-menu.stderr.log",
  "$apps\ikemen-menu.exit.txt",
  "$apps\ikemen-workflow-smoke.result-path.txt"
)) {
  if (Test-Path $path) {
    Write-Host "--- $path ---"
    Get-Content $path -Tail 40
  }
}

$resultPathLog = "$apps\ikemen-workflow-smoke.result-path.txt"
if (Test-Path $resultPathLog) {
  $resultPath = (Get-Content $resultPathLog -Raw).Trim()
  if ($resultPath -and (Test-Path $resultPath)) {
    Write-Host "--- $resultPath ---"
    Get-Content $resultPath -Tail 40
  }
}
EOF

  cat > "$workdir/ikemen-tail.ps1" <<'EOF'
$ErrorActionPreference = 'SilentlyContinue'
$apps = 'C:\Apps'
$root = 'C:\Apps\mugen'
foreach ($path in @(
  "$root\Ikemen.log",
  "$apps\ikemen-workflow-smoke.stdout.log",
  "$apps\ikemen-workflow-smoke.stderr.log",
  "$apps\ikemen-menu.stdout.log",
  "$apps\ikemen-menu.stderr.log",
  "$apps\ikemen-workflow-smoke.result-path.txt"
)) {
  if (Test-Path $path) {
    Write-Host "--- $path ---"
    Get-Content $path -Tail 80
  }
}

$resultPathLog = "$apps\ikemen-workflow-smoke.result-path.txt"
if (Test-Path $resultPathLog) {
  $resultPath = (Get-Content $resultPathLog -Raw).Trim()
  if ($resultPath -and (Test-Path $resultPath)) {
    Write-Host "--- $resultPath ---"
    Get-Content $resultPath -Tail 80
  }
}
EOF

  cat > "$workdir/ikemen-repoint.ps1" <<'EOF'
param(
  [Parameter(Mandatory = $true)]
  [string]$TargetRoot,

  [switch]$ForceDirectory
)

$ErrorActionPreference = 'Stop'
$apps = 'C:\Apps'
$link = Join-Path $apps 'mugen'

if (-not (Test-Path $TargetRoot)) {
  throw "TargetRoot not found: $TargetRoot"
}

$targetFull = (Resolve-Path $TargetRoot).Path
$existing = Get-Item $link -Force -ErrorAction SilentlyContinue

if ($existing) {
  if ($existing.LinkType -eq 'Junction' -or $existing.LinkType -eq 'SymbolicLink') {
    if ($existing.Target -contains $targetFull -or $existing.Target -eq $targetFull) {
      Write-Host ("ikemen-repoint: already points to " + $targetFull)
      exit 0
    }
    cmd /c rmdir "$link"
    if ($LASTEXITCODE -ne 0) {
      throw "Failed to remove existing link: $link"
    }
  } elseif ($ForceDirectory) {
    $backup = "$link.backup.$(Get-Date -Format 'yyyyMMdd-HHmmss')"
    Move-Item -Path $link -Destination $backup
    Write-Host ("ikemen-repoint: moved non-link directory to " + $backup)
  } else {
    throw "Refusing to remove non-link directory $link. Re-run with -ForceDirectory to move it aside."
  }
}
New-Item -ItemType Junction -Path $link -Target $targetFull | Out-Null
Write-Host ("ikemen-repoint: " + $link + " -> " + $targetFull)
EOF

  cat > "$workdir/ikemen-bundle.ps1" <<'EOF'
param(
  [string]$Label = 'manual'
)

$ErrorActionPreference = 'SilentlyContinue'
$apps = 'C:\Apps'
$root = 'C:\Apps\mugen'
$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$safeLabel = ($Label -replace '[^a-zA-Z0-9_.-]', '-')
$diagRoot = Join-Path $apps 'ikemen-diagnostics'
$bundleRoot = Join-Path $diagRoot "$stamp-$safeLabel"
$zipPath = "$bundleRoot.zip"

New-Item -ItemType Directory -Force -Path $bundleRoot | Out-Null

function Add-IfExists {
  param(
    [Parameter(Mandatory = $true)][string]$Path,
    [Parameter(Mandatory = $true)][string]$DestinationName
  )
  if (Test-Path $Path) {
    Copy-Item -Path $Path -Destination (Join-Path $bundleRoot $DestinationName) -Force
  }
}

function Write-Section {
  param(
    [Parameter(Mandatory = $true)][string]$Title,
    [Parameter(Mandatory = $true)][string]$Path,
    [scriptblock]$Body
  )
  "--- $Title ---" | Out-File -FilePath $Path -Encoding UTF8 -Append
  & $Body | Out-File -FilePath $Path -Encoding UTF8 -Append
  "" | Out-File -FilePath $Path -Encoding UTF8 -Append
}

$summary = Join-Path $bundleRoot 'summary.txt'
"Ikemen stream-box diagnostic bundle" | Out-File -FilePath $summary -Encoding UTF8
("Created: " + (Get-Date -Format o)) | Out-File -FilePath $summary -Encoding UTF8 -Append
("Root: " + $root) | Out-File -FilePath $summary -Encoding UTF8 -Append
("Label: " + $safeLabel) | Out-File -FilePath $summary -Encoding UTF8 -Append
"" | Out-File -FilePath $summary -Encoding UTF8 -Append

Write-Section 'root-link' $summary {
  Get-Item $root -Force -ErrorAction SilentlyContinue | Select-Object Name,FullName,Attributes,LinkType,Target | Format-List
}
Write-Section 'ikemen-process' $summary {
  Get-Process Ikemen_GO -ErrorAction SilentlyContinue | Select-Object Id,ProcessName,Responding,StartTime,Path | Format-List
}
Write-Section 'desktop-monitor' $summary {
  Get-CimInstance Win32_DesktopMonitor | Select-Object Name,ScreenWidth,ScreenHeight,MonitorType,Status | Format-Table -AutoSize
}
Write-Section 'video-controller' $summary {
  Get-CimInstance Win32_VideoController | Select-Object Name,VideoModeDescription,CurrentHorizontalResolution,CurrentVerticalResolution,DriverVersion,Status | Format-Table -AutoSize
}
Write-Section 'top-root-files' $summary {
  Get-ChildItem $root -Force -ErrorAction SilentlyContinue | Select-Object Mode,Length,LastWriteTime,Name | Format-Table -AutoSize
}

$files = @(
  @{ Path = "$root\Ikemen.log"; Name = 'Ikemen.log' },
  @{ Path = "$apps\ikemen-workflow-smoke.stdout.log"; Name = 'ikemen-workflow-smoke.stdout.log' },
  @{ Path = "$apps\ikemen-workflow-smoke.stderr.log"; Name = 'ikemen-workflow-smoke.stderr.log' },
  @{ Path = "$apps\ikemen-workflow-smoke.exit.txt"; Name = 'ikemen-workflow-smoke.exit.txt' },
  @{ Path = "$apps\ikemen-workflow-smoke.pid.txt"; Name = 'ikemen-workflow-smoke.pid.txt' },
  @{ Path = "$apps\ikemen-menu.stdout.log"; Name = 'ikemen-menu.stdout.log' },
  @{ Path = "$apps\ikemen-menu.stderr.log"; Name = 'ikemen-menu.stderr.log' },
  @{ Path = "$apps\ikemen-menu.exit.txt"; Name = 'ikemen-menu.exit.txt' },
  @{ Path = "$apps\ikemen-menu.pid.txt"; Name = 'ikemen-menu.pid.txt' },
  @{ Path = "$apps\ikemen-workflow-smoke.result-path.txt"; Name = 'ikemen-workflow-smoke.result-path.txt' },
  @{ Path = "$root\save\last-match.json"; Name = 'last-match.json' }
)

foreach ($file in $files) {
  Add-IfExists -Path $file.Path -DestinationName $file.Name
}

$bridgeResults = Join-Path $root 'save\bridge-results'
if (Test-Path $bridgeResults) {
  $resultsOut = Join-Path $bundleRoot 'bridge-results'
  New-Item -ItemType Directory -Force -Path $resultsOut | Out-Null
  Get-ChildItem $bridgeResults -Filter '*.json' -File |
    Sort-Object LastWriteTime -Descending |
    Select-Object -First 10 |
    Copy-Item -Destination $resultsOut -Force
}

$resultPathLog = Join-Path $apps 'ikemen-workflow-smoke.result-path.txt'
if (Test-Path $resultPathLog) {
  $resultPath = (Get-Content $resultPathLog -Raw).Trim()
  if ($resultPath -and (Test-Path $resultPath)) {
    Copy-Item -Path $resultPath -Destination (Join-Path $bundleRoot (Split-Path $resultPath -Leaf)) -Force
  }
}

if (Test-Path $zipPath) {
  Remove-Item $zipPath -Force
}
Compress-Archive -Path (Join-Path $bundleRoot '*') -DestinationPath $zipPath -Force
Write-Host ("ikemen-bundle: " + $zipPath)
EOF

  cat > "$workdir/ikemen-normalize-chars.ps1" <<'EOF'
param(
  [string]$CharsRoot = 'C:\Apps\mugen\chars',
  [switch]$Apply,
  [string]$MapPath = ''
)

$ErrorActionPreference = 'Stop'

if (-not (Test-Path -LiteralPath $CharsRoot)) {
  throw "Chars root not found: $CharsRoot"
}

function Normalize-Token {
  param([string]$Name)
  $value = $Name.ToLowerInvariant()
  $value = [Regex]::Replace($value, '[^a-z0-9]+', '_')
  $value = $value.Trim('_')
  if ([string]::IsNullOrWhiteSpace($value)) {
    return 'item'
  }
  return $value
}

function Get-UniqueLeafName {
  param(
    [string]$Directory,
    [string]$BaseName,
    [string]$Extension,
    [hashtable]$ReservedPaths,
    [string]$CurrentLeaf
  )

  $index = 1
  while ($true) {
    $suffix = if ($index -eq 1) { '' } else { "_$index" }
    $leaf = "$BaseName$suffix$Extension"
    $full = Join-Path $Directory $leaf
    $key = $full.ToLowerInvariant()
    if (-not [string]::IsNullOrWhiteSpace($CurrentLeaf) -and $leaf.ToLowerInvariant() -eq $CurrentLeaf.ToLowerInvariant()) {
      $ReservedPaths[$key] = $true
      return $leaf
    }
    if (-not (Test-Path -LiteralPath $full) -and -not $ReservedPaths.ContainsKey($key)) {
      $ReservedPaths[$key] = $true
      return $leaf
    }
    $index++
  }
}

function Clean-Field {
  param([string]$Value)
  if ($null -eq $Value) {
    return ''
  }
  return (($Value -replace "`t", ' ') -replace "`r?`n", ' ')
}

if ([string]::IsNullOrWhiteSpace($MapPath)) {
  $stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
  $mode = if ($Apply) { 'apply' } else { 'dryrun' }
  $MapPath = "C:\Apps\ikemen-char-normalization.$mode.$stamp.tsv"
}

$records = New-Object System.Collections.Generic.List[object]
$folderReserved = @{}
$folderChanges = 0
$defChanges = 0

# Normalize top-level character folder names first.
$folders = Get-ChildItem -LiteralPath $CharsRoot -Directory | Sort-Object Name
foreach ($folder in $folders) {
  $targetBase = Normalize-Token $folder.Name
  $targetLeaf = Get-UniqueLeafName -Directory $CharsRoot -BaseName $targetBase -Extension '' -ReservedPaths $folderReserved -CurrentLeaf $folder.Name
  $oldPath = $folder.FullName
  $newPath = Join-Path $CharsRoot $targetLeaf

  if ($targetLeaf -eq $folder.Name) {
    continue
  }

  if ($Apply) {
    Rename-Item -LiteralPath $oldPath -NewName $targetLeaf
    $records.Add([PSCustomObject]@{
      Kind = 'folder'
      Status = 'renamed'
      OldPath = $oldPath
      NewPath = $newPath
      Note = ''
    })
  } else {
    $records.Add([PSCustomObject]@{
      Kind = 'folder'
      Status = 'planned'
      OldPath = $oldPath
      NewPath = $newPath
      Note = ''
    })
  }
  $folderChanges++
}

# Normalize all .def file names recursively after folder pass.
$fileReservedByDir = @{}
$defs = Get-ChildItem -LiteralPath $CharsRoot -Recurse -File -Filter '*.def' | Sort-Object FullName
foreach ($def in $defs) {
  $directory = $def.DirectoryName
  $dirKey = $directory.ToLowerInvariant()
  if (-not $fileReservedByDir.ContainsKey($dirKey)) {
    $fileReservedByDir[$dirKey] = @{}
  }
  $reserved = $fileReservedByDir[$dirKey]

  $stem = [IO.Path]::GetFileNameWithoutExtension($def.Name)
  $targetBase = Normalize-Token $stem
  $targetLeaf = Get-UniqueLeafName -Directory $directory -BaseName $targetBase -Extension '.def' -ReservedPaths $reserved -CurrentLeaf $def.Name
  $oldPath = $def.FullName
  $newPath = Join-Path $directory $targetLeaf

  if ($targetLeaf -eq $def.Name) {
    continue
  }

  if ($Apply) {
    Rename-Item -LiteralPath $oldPath -NewName $targetLeaf
    $records.Add([PSCustomObject]@{
      Kind = 'def'
      Status = 'renamed'
      OldPath = $oldPath
      NewPath = $newPath
      Note = ''
    })
  } else {
    $records.Add([PSCustomObject]@{
      Kind = 'def'
      Status = 'planned'
      OldPath = $oldPath
      NewPath = $newPath
      Note = ''
    })
  }
  $defChanges++
}

$lines = New-Object System.Collections.Generic.List[string]
$lines.Add("kind`tstatus`toldPath`tnewPath`tnote")
foreach ($record in $records) {
  $line = "{0}`t{1}`t{2}`t{3}`t{4}" -f `
    (Clean-Field $record.Kind), `
    (Clean-Field $record.Status), `
    (Clean-Field $record.OldPath), `
    (Clean-Field $record.NewPath), `
    (Clean-Field $record.Note)
  $lines.Add($line)
}
Set-Content -Path $MapPath -Value $lines -Encoding UTF8

$modeLabel = if ($Apply) { 'apply' } else { 'dry-run' }
Write-Host ("ikemen-normalize-chars: mode={0} foldersChanged={1} defsChanged={2}" -f $modeLabel, $folderChanges, $defChanges)
Write-Host ("ikemen-normalize-chars map: " + $MapPath)
if (-not $Apply) {
  Write-Host "Re-run with -Apply to persist the rename plan."
}
EOF

  run_remote "powershell -NoProfile -Command \"New-Item -ItemType Directory -Force -Path '$REMOTE_TOOLS_ROOT' | Out-Null\"" >/dev/null
  scp "$workdir/ikemen-check.ps1" "$STREAM_BOX_HOST:$(remote_file 'ikemen-check.ps1')" >/dev/null
  scp "$workdir/ikemen-stop.ps1" "$STREAM_BOX_HOST:$(remote_file 'ikemen-stop.ps1')" >/dev/null
  scp "$workdir/ikemen-launch-quickvs.ps1" "$STREAM_BOX_HOST:$(remote_file 'ikemen-launch-quickvs.ps1')" >/dev/null
  scp "$workdir/ikemen-launch-menu.ps1" "$STREAM_BOX_HOST:$(remote_file 'ikemen-launch-menu.ps1')" >/dev/null
  scp "$workdir/ikemen-status.ps1" "$STREAM_BOX_HOST:$(remote_file 'ikemen-status.ps1')" >/dev/null
  scp "$workdir/ikemen-tail.ps1" "$STREAM_BOX_HOST:$(remote_file 'ikemen-tail.ps1')" >/dev/null
  scp "$workdir/ikemen-repoint.ps1" "$STREAM_BOX_HOST:$(remote_file 'ikemen-repoint.ps1')" >/dev/null
  scp "$workdir/ikemen-bundle.ps1" "$STREAM_BOX_HOST:$(remote_file 'ikemen-bundle.ps1')" >/dev/null
  scp "$workdir/ikemen-normalize-chars.ps1" "$STREAM_BOX_HOST:$(remote_file 'ikemen-normalize-chars.ps1')" >/dev/null
  scp "$SCRIPT_DIR/ikemen-candidate-quality.ps1" "$STREAM_BOX_HOST:$(remote_file 'ikemen-candidate-quality.ps1')" >/dev/null
  printf 'Installed toolbox scripts under %s\n' "$REMOTE_TOOLS_ROOT"
}

run_tool() {
  local script_name="$1"
  run_remote "powershell -NoProfile -ExecutionPolicy Bypass -File $(remote_file "$script_name")"
}

main() {
  [[ $# -ge 1 ]] || { usage; exit 1; }
  case "$1" in
    install)
      install_toolbox
      ;;
    check)
      run_tool 'ikemen-check.ps1'
      ;;
    launch)
      run_tool 'ikemen-launch-quickvs.ps1'
      ;;
    launch-menu)
      run_tool 'ikemen-launch-menu.ps1'
      ;;
    status)
      run_tool 'ikemen-status.ps1'
      ;;
    tail)
      run_tool 'ikemen-tail.ps1'
      ;;
    stop)
      run_tool 'ikemen-stop.ps1'
      ;;
    bundle)
      shift
      label="${1:-manual}"
      run_remote "powershell -NoProfile -ExecutionPolicy Bypass -File $(remote_file 'ikemen-bundle.ps1') -Label \"$label\""
      ;;
    repoint)
      shift
      [[ $# -ge 1 ]] || die "repoint requires a target root path"
      extra_args=""
      if [[ "${2:-}" == "--force-directory" ]]; then
        extra_args=" -ForceDirectory"
      fi
      run_remote "powershell -NoProfile -ExecutionPolicy Bypass -File $(remote_file 'ikemen-repoint.ps1') -TargetRoot \"$1\"$extra_args"
      ;;
    normalize-chars)
      run_remote "powershell -NoProfile -ExecutionPolicy Bypass -File $(remote_file 'ikemen-normalize-chars.ps1')"
      ;;
    normalize-chars-apply)
      run_remote "powershell -NoProfile -ExecutionPolicy Bypass -File $(remote_file 'ikemen-normalize-chars.ps1') -Apply"
      ;;
    candidate-quality)
      shift
      source_root=""
      facts_out_root=""
      normalized_out_root=""
      run_id=""
      limit=""
      while [[ $# -gt 0 ]]; do
        case "$1" in
          --source-root)
            source_root="$2"
            shift 2
            ;;
          --facts-out-root)
            facts_out_root="$2"
            shift 2
            ;;
          --normalized-out-root)
            normalized_out_root="$2"
            shift 2
            ;;
          --run-id)
            run_id="$2"
            shift 2
            ;;
          --limit)
            limit="$2"
            shift 2
            ;;
          *)
            die "Unknown candidate-quality argument: $1"
          ;;
        esac
      done
      install_toolbox
      cmd="powershell -NoProfile -ExecutionPolicy Bypass -File $(remote_file 'ikemen-candidate-quality.ps1')"
      if [[ -n "$source_root" ]]; then
        cmd+=" -SourceRoot \"$source_root\""
      fi
      if [[ -n "$facts_out_root" ]]; then
        cmd+=" -FactsOutRoot \"$facts_out_root\""
      fi
      if [[ -n "$normalized_out_root" ]]; then
        cmd+=" -NormalizedOutRoot \"$normalized_out_root\""
      fi
      if [[ -n "$run_id" ]]; then
        cmd+=" -RunId \"$run_id\""
      fi
      if [[ -n "$limit" ]]; then
        cmd+=" -Limit $limit"
      fi
      run_remote "$cmd"
      ;;
    -h|--help|help)
      usage
      ;;
    *)
      die "Unknown command: $1"
      ;;
  esac
}

main "$@"
