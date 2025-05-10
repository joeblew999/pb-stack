# Set Execution Policy to Bypass for the current session
Set-ExecutionPolicy Bypass -Scope Process -Force

& {
    $env:HTML_DEEP_PATH = Join-Path -Path $PWD.Path -ChildPath "over-from-shell"
    $env:HTML_DEEP_NAME = "over-from-shell"
    task.exe
}