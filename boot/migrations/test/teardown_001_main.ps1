# Set Execution Policy to Bypass for the current session
Set-ExecutionPolicy Bypass -Scope Process -Force

# Check Privs

if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Start-Process powershell -Verb runAs -ArgumentList ("-NoProfile -ExecutionPolicy Bypass -File `"" + $PSCommandPath + "`"")
    exit
}

Write-Host ""
Write-Host "---------------------------------------------------------------------" -ForegroundColor Green
Write-Host "Teardown script finished basic environment checks." -ForegroundColor Green
Write-Host "Package uninstallation and VS Code extensions are now primarily" -ForegroundColor Green
Write-Host "managed by the Go application based on config.json." -ForegroundColor Green
Write-Host "This script can be used for other post-teardown tasks or system" -ForegroundColor Green
Write-Host "cleanup not covered by the JSON." -ForegroundColor Green
Write-Host "---------------------------------------------------------------------" -ForegroundColor Green
# Add any other non-package, non-vscode-extension teardown tasks here.
