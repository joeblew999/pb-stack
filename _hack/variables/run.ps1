& {
    $env:HTML_DEEP_PATH = Join-Path -Path $PWD.Path -ChildPath "over-from-shell"
    $env:HTML_DEEP_NAME = "over-from-shell"
    task.exe
}