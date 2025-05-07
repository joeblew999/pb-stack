# Installs Git using winget
Write-Host "Attempting to install Git using winget..."
winget install --id=Git.Git -e --accept-source-agreements --accept-package-agreements
Write-Host "Git installation process finished."
