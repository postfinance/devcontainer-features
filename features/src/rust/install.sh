. ./functions.sh

"./installer_$(detect_arch)" \
    -version="${VERSION:-"latest"}" \
    -rustupVersion="${RUSTUPVERSION:-"latest"}" \
    -profile="${PROFILE:-"minimal"}" \
    -components="${COMPONENTS:-"rustfmt,rust-analyzer,rust-src,clippy"}" \
    -enableWindowsTarget="${ENABLEWINDOWSTARGET:-"false"}"
