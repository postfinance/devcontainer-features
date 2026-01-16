. ./functions.sh

"./installer_$(detect_arch)" \
  -version="${VERSION:-"latest"}" \
  -kubectxVersion="${KUBECTXVERSION:-"latest"}" \
  -kubensVersion="${KUBENSVERSION:-"latest"}" \
  -k9sVersion="${K9SVERSION:-"none"}" \
  -helmVersion="${HELMVERSION:-"none"}" \
  -kustomizeVersion="${KUSTOMIZEVERSION:-"none"}" \
  -kubeconformVersion="${KUBECONFORMVERSION:-"none"}" \
  -kubescoreVersion="${KUBESCOREVERSION:-"none"}" \
  -downloadUrl="${DOWNLOADURL:-""}" \
  -kubectxDownloadUrl="${KUBECTXDOWNLOADURL:-""}" \
  -kubensDownloadUrl="${KUBENSDOWNLOADURL:-""}" \
  -k9sDownloadUrl="${K9SDOWNLOADURL:-""}" \
  -helmDownloadUrl="${HELMDOWNLOADURL:-""}" \
  -kustomizeDownloadUrl="${KUSTOMIZEDOWNLOADURL:-""}" \
  -kubeconformDownloadUrl="${KUBECONFORMDOWNLOADURL:-""}" \
  -kubescoreDownloadUrl="${KUBESCOREDOWNLOADURL:-""}"
