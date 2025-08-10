function check_version {
    if ! echo "$1" | grep -q "$2"; then
        echo "Expected '$2' but was '$1'"
        exit 1
    fi
}

function check_file_exists {
    if [ ! -f "$1" ]; then
        echo "File '$1' does not exist"
        exit 1
    fi
}

function check_file_not_exists {
    if [ -f "$1" ]; then
        echo "File '$1' exists"
        exit 1
    fi
}

function check_dir_exists {
    if [ ! -d "$1" ]; then
        echo "Directory '$1' does not exist"
        exit 1
    fi
}

function check_dir_not_exists {
    if [ -d "$1" ]; then
        echo "Directory '$1' exists"
        exit 1
    fi
}

function check_command_exists {
    if ! command -v $1 2>&1 >/dev/null; then
        echo "Command '$1' does not exist"
        exit 1
    fi
}

function check_command_not_exists {
    if command -v $1 2>&1 >/dev/null; then
        echo "Command '$1' exists"
        exit 1
    fi
}

function check_package_installed {
    if ! dpkg -s "$1" &>/dev/null; then
        echo "Package '$1' is not installed"
        exit 1
    fi
}

function check_env_var_exists {
    if [ -z "${!1}" ]; then
        echo "Environment variable '$1' does not exist"
        exit 1
    fi
}
