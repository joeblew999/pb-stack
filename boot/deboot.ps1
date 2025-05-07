# Set Execution Policy to Bypass for the current session
Set-ExecutionPolicy Bypass -Scope Process -Force

# Check Privs

if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Start-Process powershell -Verb runAs -ArgumentList ("-NoProfile -ExecutionPolicy Bypass -File `"" + $PSCommandPath + "`"")
    exit
}


# Uninstalls Bun using winget
try {
    Write-Host "Attempting to uninstall Bun using winget..."
    winget uninstall --id=Oven-sh.Bun -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Bun uninstallation process finished."
} catch {
    Write-Warning "Failed to uninstall Bun using winget. It might not have been installed or an error occurred."
}

# Uninstalls Git using winget
try {
    Write-Host "Attempting to uninstall Git using winget..."
    winget uninstall --id=Git.Git -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Git uninstallation process finished."
} catch {
    Write-Warning "Failed to uninstall Git using winget. It might not have been installed or an error occurred."
}


# Uninstalls Go using winget
try {
    Write-Host "Attempting to uninstall Go using winget..."
    winget uninstall --id=GoLang.Go -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Go uninstallation process finished."
} catch {
    Write-Warning "Failed to uninstall Go using winget. It might not have been installed or an error occurred."
}


# Uninstalls OpenSSH using winget
try {
    Write-Host "Attempting to uninstall OpenSSH using winget..."
    winget uninstall --id=Microsoft.OpenSSH.Preview -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "OpenSSH uninstallation process finished."
} catch {
    Write-Warning "Failed to uninstall OpenSSH using winget. It might not have been installed or an error occurred."
}


# Uninstalls Which using winget
try {
    Write-Host "Attempting to uninstall Which using winget..."
    winget uninstall --id=GnuWin32.Which -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Which uninstallation process finished."
} catch {
    Write-Warning "Failed to uninstall Which using winget. It might not have been installed or an error occurred."
}
