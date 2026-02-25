. ./functions.sh

"./installer_$(detect_arch)" \
  -version="${VERSION:-"latest"}" \
  -keyringVersion="${KEYRINGVERSION:-""}" \
  -installLibraries="${INSTALLLIBRARIES:-"true"}" \
  -installDevLibraries="${INSTALLDEVLIBRARIES:-"true"}" \
  -installCompiler="${INSTALLCOMPILER:-"true"}" \
  -installTools="${INSTALLTOOLS:-"true"}" \
  -additionalCudaPackages="${ADDITIONALCUDAPACKAGES:-""}" \
  -downloadUrl="${DOWNLOADURL:-""}"
