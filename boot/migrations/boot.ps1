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

# Installs Bun using winget
try {
    Write-Host "" # Empty line
    Write-Host "Attempting to install Bun using winget..." -ForegroundColor Red
    winget install --id=Oven-sh.Bun -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Bun installation process finished."
} catch {
    Write-Warning "Failed to install Bun using winget. Please install it manually."
}

# Installs Git using winget
try {
    Write-Host "" # Empty line
    Write-Host "Attempting to install Git using winget..." -ForegroundColor Red
    winget install --id=Git.Git -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Git installation process finished."
} catch {
    Write-Warning "Failed to install Git using winget. Please install it manually."
}


# Installs Go using winget
try {
    Write-Host "" # Empty line
    Write-Host "Attempting to install Go using winget..." -ForegroundColor Red
    winget install --id=GoLang.Go -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Go installation process finished."
} catch {
    Write-Warning "Failed to install Go using winget. Please install it manually."
}

# Installs Visual Studio Code using winget
try {
    Write-Host "" # Empty line
    Write-Host "Attempting to install Visual Studio Code using winget..." -ForegroundColor Red
    winget install --id=Microsoft.VisualStudioCode -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Visual Studio Code installation process finished."
} catch {
    Write-Warning "Failed to install Visual Studio Code using winget. Please install it manually."
}

# Installs Task (taskfile.dev) using winget
try {
    Write-Host "" # Empty line
    Write-Host "Attempting to install Task (taskfile.dev) using winget..." -ForegroundColor Red
    winget install --id=GoTask.Task -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Task (taskfile.dev) installation process finished."
} catch {
    Write-Warning "Failed to install Task (taskfile.dev) using winget. Please install it manually."
}




# Installs OpenSSH using winget
try {
    Write-Host "" # Empty line
    Write-Host "Attempting to install OpenSSH using winget..." -ForegroundColor Red
    winget install --id=Microsoft.OpenSSH.Preview -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "OpenSSH installation process finished."
} catch {
    Write-Warning "Failed to install OpenSSH using winget. Please install it manually."
}



# Installs Which using winget
try {
    Write-Host "" # Empty line
    Write-Host "Attempting to install Which using winget..." -ForegroundColor Red
    winget install --id=GnuWin32.Which -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Which installation process finished."
} catch {
    Write-Warning "Failed to install Which using winget. Please install it manually."
}

# Install VS Code Extensions
$extensionsFile = Join-Path (Split-Path -Path $PSCommandPath) "extensions.txt"

if (Test-Path $extensionsFile) {
    try {
        Write-Host "" # Empty line
        Write-Host "Found extensions.txt. Attempting to install VS Code extensions..." -ForegroundColor Red
        
        # Read each line from the extensions.txt file
        Get-Content $extensionsFile | ForEach-Object {
            $extensionId = $_.Trim()
            
            # Check if the line is not empty or a comment
            if ($extensionId -and $extensionId -notmatch "^#") {
                Write-Host "Installing extension: $extensionId"
                code --install-extension $extensionId --force
            }
        }
        Write-Host "VS Code extension installation process finished."
    } catch {
        Write-Warning "Failed to install VS Code extensions. Please check the extensions.txt file and ensure VS Code is installed."
    }
} else {
    Write-Host "extensions.txt not found. Skipping VS Code extension installation."
}

# List installed VS Code Extensions
try {
    Write-Host "Listing installed VS Code extensions..."
    code --list-extensions
} catch {
    Write-Warning "Failed to list VS Code extensions. Ensure VS Code is installed and in your PATH."
}
