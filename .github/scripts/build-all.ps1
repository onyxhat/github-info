[CmdletBinding()]
param (
    [Parameter()]
    [string]
    $WorkDir = $(Get-Location),

    [Parameter()]
    [string]
    $BinDir = $(if (!$env:BINARY_OUTDIR) { "bin" } else { $env:BINARY_OUTDIR }),

    [Parameter()]
    [string]
    $ProjectName = $(if (!$env:CI_REPOSITORY_NAME_SLUG) { $(Get-Location | Get-Item).BaseName } else { $env:CI_REPOSITORY_NAME_SLUG }),

    [Parameter()]
    [string]
    $BuildArch = $(if (!$env:BUILD_ARCH) { ((& go tool dist list) -match '^(darwin|linux|windows)/(arm|arm64|386|amd64)$') -join ',' } else { $env:BUILD_ARCH })
)

if (!(Test-Path $WorkDir/$BinDir)) { New-Item -Path $WorkDir/$BinDir -ItemType Directory | Out-Null }

Push-Location $WorkDir

& go get .

foreach ($g in $BuildArch.Split(',')) {
    $env:GOOS = $g.split('/')[0]
    $env:GOARCH = $g.split('/')[1]

    if ($env:GOOS -eq "windows") { $Ext = ".exe" } else { $Ext = $null }

    try {
        Write-Host "Building: ${env:GOOS} (${env:GOARCH})"
        & go build -ldflags="-s -w" -o "${WorkDir}/${BinDir}/${ProjectName}-${env:GOOS}-${env:GOARCH}${Ext}"
    }

    catch {
        Write-Warning ("Build Failed:`r`n" + $_.Exception.Message)
    }
}

Pop-Location
