# Set Execution Policy to Bypass for the current session
Set-ExecutionPolicy Bypass -Scope Process -Force

# Check Privs
Write-Host "Checking if the current user has Administrator privileges..."
if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Start-Process powershell -Verb runAs -ArgumentList ("-NoProfile -ExecutionPolicy Bypass -File `"" + $PSCommandPath + "`"")
    exit
}


# Installs Bun using winget
try {
    Write-Host "Attempting to install Bun using winget..."
    winget install --id=Oven-sh.Bun -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Bun installation process finished."
} catch {
    Write-Warning "Failed to install Bun using winget. Please install it manually."
}

# Installs Git using winget
try {
    Write-Host "Attempting to install Git using winget..."
    winget install --id=Git.Git -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Git installation process finished."
} catch {
    Write-Warning "Failed to install Git using winget. Please install it manually."
}


# Installs Go using winget
try {
    Write-Host "Attempting to install Go using winget..."
    winget install --id=GoLang.Go -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Go installation process finished."
} catch {
    Write-Warning "Failed to install Go using winget. Please install it manually."
}


# Installs OpenSSH using winget
try {
    Write-Host "Attempting to install OpenSSH using winget..."
    winget install --id=Microsoft.OpenSSH.Preview -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "OpenSSH installation process finished."
} catch {
    Write-Warning "Failed to install OpenSSH using winget. Please install it manually."
}



# Installs Which using winget
try {
    Write-Host "Attempting to install Which using winget..."
    winget install --id=GnuWin32.Which -e --accept-source-agreements --accept-package-agreements --silent
    Write-Host "Which installation process finished."
} catch {
    Write-Warning "Failed to install Which using winget. Please install it manually."
}

# Install VS Code Extensions
$extensionsFile = Join-Path (Split-Path -Path $PSCommandPath) "extensions.txt"

if (Test-Path $extensionsFile) {
    try {
        Write-Host "Found extensions.txt. Attempting to install VS Code extensions..."
        
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
