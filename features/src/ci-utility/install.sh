. ./functions.sh

"./installer_$(detect_arch)" \
 -yqVersion="${VERSION:-"system-default"}" \
 -gettextbaseVersion="${VERSION:-"system-default"}" \
 -yamllintVersion="${VERSION:-"system-default"}" \
 -sshpassVersion="${VERSION:-"system-default"}"