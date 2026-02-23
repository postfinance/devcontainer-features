. ./functions.sh

"./installer_$(detect_arch)" \
  -version="${VERSION:-"latest"}" \
  -installLibraries="${INSTALLLIBRARIES:-"true"}" \
  -installDevLibraries="${INSTALLDEVLIBRARIES:-"true"}" \
  -installCompiler="${INSTALLCOMPILER:-"true"}" \
  -installTools="${INSTALLTOOLS:-"true"}" \
  -additionalCudaPackages="${ADDITIONALCUDAPACKAGES:-""}" \
  -downloadUrl="${DOWNLOADURL:-""}"
