detect_arch() {
  case $(uname -m | tr '[:upper:]' '[:lower:]') in
    x86_64*) echo "amd64" ;;
    arm*|aarch64*) echo "arm64" ;;
    *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
  esac
}
