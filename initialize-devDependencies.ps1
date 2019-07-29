

docker stack deploy -c "$PSScriptRoot/examples/stack.yml" sql
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$terraformPath = "$($env:LOCALAPPDATA)/terraform"

if (-not (Test-Path $terraformPath))
{
    Invoke-WebRequest `
        -Uri "https://releases.hashicorp.com/terraform/0.12.5/terraform_0.12.5_windows_amd64.zip" `
        -OutFile "$($env:TEMP)/terraform.zip"
    Expand-Archive `
        -Path "$($env:TEMP)/terraform.zip" `
        -DestinationPath $terraformPath
}

$commandSnippet = '$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "user")'
if (-not (Test-Path $PROFILE) -or -not (Get-Content -Raw $PROFILE).Contains($commandSnippet))
{
    $commandSnippet | Out-File -Append -Force $PROFILE
}

if (-not $env:Path.Contains($terraformPath))
{
    [System.Environment]::SetEnvironmentVariable("Path", "$([System.Environment]::GetEnvironmentVariable("Path", "user"));$($terraformPath)", "user")
    $env:Path = "$($env:Path);$($terraformPath)"
}
