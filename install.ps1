[CmdletBinding()]
param(
  [string]$GOPATH="$($env:USERPROFILE)/go",
  [string]$GITHUB_USER="estenrye",
  [string]$PROVIDER_NAME='terraform-provider-sqlserver',
  [switch]$skipGet = $false
)
$BIN_PATH="${GOPATH}/bin/${PROVIDER_NAME}.exe"
$SRC_PATH=$PSScriptRoot

Push-Location $SRC_PATH

if ($skipGet)
{
  $PROVIDER_REPO=""
}
else
{
  $PROVIDER_REPO="github.com/${GITHUB_USER}/${PROVIDER_NAME}"
  New-Item -ItemType Directory ${GOPATH} -Force | Out-Null
}

go get ${PROVIDER_REPO}
Remove-Item ${BIN_PATH}
go build
go install

$TERRAFORM_PLUGINS_DIR="$($env:APPDATA)\terraform.d\plugins\windows_amd64"
New-Item -ItemType Directory ${TERRAFORM_PLUGINS_DIR} -Force | Out-Null

$PROVIDER_PATH=(Join-Path ${TERRAFORM_PLUGINS_DIR} "${PROVIDER_NAME}.exe")

Write-Host ${BIN_PATH}
Write-Host ${PROVIDER_PATH}
if (Test-Path "${BIN_PATH}")
{
  Copy-Item ${BIN_PATH} ${PROVIDER_PATH} -Force
  Write-Host "Copy Successful.  ${PROVIDER_PATH}"
}
else
{
  Write-Output 'Build Failed, Copy Aborted.'
  exit 1
}
Pop-Location
exit 0