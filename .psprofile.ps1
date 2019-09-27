Set-Variable HOME "$env:USERPROFILE" -Scope Global -Force
( Get-PSProvider 'FileSystem' ).Home = $HOME   # replace "~"

$global:ROOT = "$HOME\Projects\terraform-provider-hosts"
$env:PATH = "$ROOT\scripts;$env:PATH"

if ( -not ( Get-Location ).Path.StartsWith("$ROOT") ) {
    Set-Location "$ROOT"
}

Apply-PSConsoleSettings "TERRAFORM-PROVIDER-HOSTS"

#
# for packer
$env:HOME = "$HOME"
$env:ROOT = "$ROOT"

#
# for terraform
$env:TF_ROOT = "$ROOT/examples".Replace("\", "/")
$env:TF_INPUT = "false"
$env:TF_LOG_PATH = "$env:TF_ROOT/_terraform.log"
$env:TF_LOG = "TRACE"

$env:TF_VAR_root = "$env:TF_ROOT"
$env:TF_VAR_terraform = "$env:TF_ROOT"
