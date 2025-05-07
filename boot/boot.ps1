# Check Privs

if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Start-Process powershell -Verb runAs -ArgumentList ("-NoProfile -ExecutionPolicy Bypass -File `"" + $PSCommandPath + "`"")
    exit
}


# Installs Git using winget
try {
    Write-Host "Attempting to install Git using winget..."
    winget install --id=Git.Git -e --accept-source-agreements --accept-package-agreements
    Write-Host "Git installation process finished."
} catch {
    Write-Warning "Failed to install Git using winget. Please install it manually."
}
