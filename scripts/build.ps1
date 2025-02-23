# build.ps1

# build for current platform
param(
    [string]$Target = "current"
)

# get version info
$Version = $(git describe --tags)
$BuildTime = $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
$GitCommit = $(git rev-parse HEAD)

# create output directory
New-Item -ItemType Directory -Force -Path "build/$Version" | OUT-Null

function Build-Single {
    param (
        [string]$OS,
        [string]$ARCH,
        [string]$Extension = ""
    )
    
    $env:GOOS = $OS
    $env:GOARCH = $ARCH
    $env:CGO_ENABLED = 0
    
    $OutputName = "build/$Version/adc-$OS-$ARCH$Extension"
    
    Write-Host "Building for $OS/$ARCH..."
    
    go build `
        -ldflags="-w -s -X 'adc/cmd.Version=$Version' -X 'adc/cmd.BuildTime=$BuildTime' -X 'adc/cmd.GitCommit=$GitCommit'" `
        -trimpath `
        -o $OutputName
}

# build for current platform
function Build-Current {
    # copy example config.toml
    Copy-Item -Path "config.toml" -Destination "build/$Version"

    $env:CGO_ENABLED = 0
    go build `
        -ldflags="-w -s -X 'adc/cmd.Version=$Version' -X 'adc/cmd.BuildTime=$BuildTime' -X 'adc/cmd.GitCommit=$GitCommit'" `
        -trimpath `
        -o "build/$Version/adc.exe"
}

# build for all platforms
function Build-All {
    # copy example config.toml
    Copy-Item -Path "config.toml" -Destination "build/$Version"

    Build-Single -OS "linux" -ARCH "amd64"
    Build-Single -OS "darwin" -ARCH "amd64"
    Build-Single -OS "windows" -ARCH "amd64" -Extension ".exe"
}



switch ($Target) {
    "current" { Build-Current }
    "all" { Build-All }
    default { Write-Host "Unknown target: $Target" }
}