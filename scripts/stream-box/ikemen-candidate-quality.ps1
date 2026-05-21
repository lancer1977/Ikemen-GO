param(
  [string]$SourceRoot = 'T:\chars-candidates',
  [string]$FactsOutRoot = 'T:\chars\_facts',
  [string]$NormalizedOutRoot = 'T:\chars\_normalized',
  [string]$ApiIkemenScriptsRoot = 'C:\code\Api.Ikemen\scripts\stack-fact-finding',
  [string]$RunId = '',
  [int]$Limit = 0
)

$ErrorActionPreference = 'Stop'

function Invoke-PythonStep {
  param(
    [Parameter(Mandatory = $true)][string]$ScriptPath,
    [Parameter(Mandatory = $true)][string[]]$Arguments
  )

  & python $ScriptPath @Arguments
  if ($LASTEXITCODE -ne 0) {
    throw "Python step failed: $ScriptPath"
  }
}

function Resolve-ExistingPath {
  param(
    [Parameter(Mandatory = $true)][string[]]$Candidates,
    [Parameter(Mandatory = $true)][string]$Label
  )

  foreach ($candidate in $Candidates) {
    if (-not [string]::IsNullOrWhiteSpace($candidate) -and (Test-Path -LiteralPath $candidate)) {
      return $candidate
    }
  }

  $rendered = $Candidates | Where-Object { -not [string]::IsNullOrWhiteSpace($_) } | ForEach-Object { "  - $_" }
  throw "$Label not found. Tried:`n$($rendered -join [Environment]::NewLine)"
}

function Read-FirstReadyCandidate {
  param([Parameter(Mandatory = $true)][string]$GateRoot)

  $readyPath = "$GateRoot\ready.tsv"
  if (-not (Test-Path $readyPath)) {
    return $null
  }

  $rows = Import-Csv -Path $readyPath -Delimiter "`t"
  if ($null -eq $rows -or $rows.Count -lt 1) {
    return $null
  }

  $first = $rows | Select-Object -First 1
  if ($null -eq $first -or [string]::IsNullOrWhiteSpace($first.source_def)) {
    return $null
  }

  return (Split-Path -Parent $first.source_def)
}

if (-not (Test-Path -LiteralPath $ApiIkemenScriptsRoot)) {
  throw "Api.Ikemen script root not found: $ApiIkemenScriptsRoot"
}

if (-not (Get-Command python -ErrorAction SilentlyContinue)) {
  throw "python was not found on stream-box"
}

if ([string]::IsNullOrWhiteSpace($RunId)) {
  $RunId = [DateTime]::UtcNow.ToString('yyyy-MM-ddTHH-mm-ss')
}

$runtimeRoot = 'C:\Apps\mugen'
$TargetRoot = Resolve-ExistingPath -Label 'Promotion target root' -Candidates @(
  'T:\chars',
  "$runtimeRoot\chars"
)
if ($PSBoundParameters.ContainsKey('SourceRoot')) {
  $SourceRoot = Resolve-ExistingPath -Label 'Candidate source root' -Candidates @($SourceRoot)
} else {
  $SourceRoot = Resolve-ExistingPath -Label 'Candidate source root' -Candidates @(
    $SourceRoot,
    "$runtimeRoot\chars-candidates"
  )
}
if ($PSBoundParameters.ContainsKey('FactsOutRoot')) {
  if (-not (Test-Path -LiteralPath $FactsOutRoot)) {
    throw "Facts output root not found: $FactsOutRoot"
  }
} else {
  $FactsOutRoot = if (Test-Path -LiteralPath $FactsOutRoot) { $FactsOutRoot } else { "$runtimeRoot\_facts" }
}
if ($PSBoundParameters.ContainsKey('NormalizedOutRoot')) {
  if (-not (Test-Path -LiteralPath $NormalizedOutRoot)) {
    throw "Normalized output root not found: $NormalizedOutRoot"
  }
} else {
  $NormalizedOutRoot = if (Test-Path -LiteralPath $NormalizedOutRoot) { $NormalizedOutRoot } else { "$runtimeRoot\_normalized" }
}

$inventoryScript = Join-Path $ApiIkemenScriptsRoot 'analyze-mugen-disk.py'
$normalizeScript = Join-Path $ApiIkemenScriptsRoot 'normalize-mugen-candidates.py'
$triageScript = Join-Path $ApiIkemenScriptsRoot 'summarize-mugen-normalization-run.py'
$dedupeScript = Join-Path $ApiIkemenScriptsRoot 'summarize-mugen-duplicate-characters.py'
$gateScript = Join-Path $ApiIkemenScriptsRoot 'summarize-mugen-promotion-gate.py'
$promoteScript = Join-Path $ApiIkemenScriptsRoot 'promote-mugen-candidate.py'

$normalizedRunRoot = "$NormalizedOutRoot\runs\$RunId"
$gateRoot = "$normalizedRunRoot\promotion-gate"

Write-Host "candidate-quality: inventory"
Invoke-PythonStep -ScriptPath $inventoryScript -Arguments @('--source-root', $SourceRoot, '--out-root', $FactsOutRoot, '--run-id', $RunId)

Write-Host "candidate-quality: normalize"
$normalizeArgs = @('--source-root', $SourceRoot, '--out-root', $NormalizedOutRoot, '--run-id', $RunId)
if ($Limit -gt 0) {
  $normalizeArgs += @('--limit', [string]$Limit)
}
Invoke-PythonStep -ScriptPath $normalizeScript -Arguments $normalizeArgs

Write-Host "candidate-quality: triage"
Invoke-PythonStep -ScriptPath $triageScript -Arguments @('--run-root', $normalizedRunRoot)

Write-Host "candidate-quality: dedupe"
Invoke-PythonStep -ScriptPath $dedupeScript -Arguments @('--run-root', $normalizedRunRoot)

Write-Host "candidate-quality: promotion gate"
Invoke-PythonStep -ScriptPath $gateScript -Arguments @('--run-root', $normalizedRunRoot)

$candidate = Read-FirstReadyCandidate -GateRoot $gateRoot
if ($null -ne $candidate) {
  Write-Host ("candidate-quality: dry-run promote " + $candidate)
  Invoke-PythonStep -ScriptPath $promoteScript -Arguments @(
    '--candidate', $candidate,
    '--source-root', $SourceRoot,
    '--target-root', $TargetRoot,
    '--run-id', $RunId,
    '--dry-run'
  )
} else {
  Write-Host "candidate-quality: no ready candidate found for dry-run promotion"
}

Write-Host "candidate-quality: run-root $normalizedRunRoot"
Write-Host "candidate-quality: gate-root $gateRoot"
Write-Host "candidate-quality: complete"
