# Set Execution Policy to Bypass for the current session
Set-ExecutionPolicy Bypass -Scope Process -Force

# Check Privs
Write-Host "Checking if the current user has Administrator privileges..."
if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Start-Process powershell -Verb runAs -ArgumentList ("-NoProfile -ExecutionPolicy Bypass -File `"" + $PSCommandPath + "`"")
    exit
}

# Check for winget and attempt to register if missing
Write-Host ""
Write-Host "Checking for winget package manager..." -ForegroundColor Cyan
if (-not (Get-Command winget -ErrorAction SilentlyContinue)) {
    Write-Host "winget not found. Attempting to ensure it's registered..." -ForegroundColor Yellow
    try {
        Write-Host "Attempting to register Microsoft.DesktopAppInstaller package family..."
        Add-AppxPackage -RegisterByFamilyName -MainPackage Microsoft.DesktopAppInstaller_8wekyb3d8bbwe
        Write-Host "Registration attempt complete. Please re-run the script if winget is now available, or check the Microsoft Store for 'App Installer'."
    } catch {
        Write-Warning "An error occurred during winget registration attempt: $($_.Exception.Message)"
    }

    # Check again if winget is available after registration attempt
    if (-not (Get-Command winget -ErrorAction SilentlyContinue)) {
        Write-Error "winget is still not available after attempting registration."
        Write-Error "Please ensure 'App Installer' is installed from the Microsoft Store and winget is functioning."
        Write-Error "You can try running 'Add-AppxPackage -RegisterByFamilyName -MainPackage Microsoft.DesktopAppInstaller_8wekyb3d8bbwe' manually in an administrative PowerShell."
        Write-Error "Aborting script as winget is required for subsequent operations."
        exit 1 # Exit the script as winget is crucial
    } else {
        Write-Host "winget became available after registration attempt." -ForegroundColor Green
    }
} else {
    Write-Host "winget is already available." -ForegroundColor Green
}

Write-Host ""
Write-Host "---------------------------------------------------------------------" -ForegroundColor Green
Write-Host "Setup script finished basic environment checks (e.g., winget)." -ForegroundColor Green
Write-Host "Package installations and VS Code extensions are now primarily" -ForegroundColor Green
Write-Host "managed by the Go application based on config.json." -ForegroundColor Green
Write-Host "This script can be used for other pre-setup tasks or system" -ForegroundColor Green
Write-Host "configurations not covered by the JSON (e.g., Chocolatey installs if needed)." -ForegroundColor Green
Write-Host "---------------------------------------------------------------------" -ForegroundColor Green
# Add any other non-package, non-vscode-extension setup tasks here.
# For example, specific system configurations, Chocolatey package installations, etc.
