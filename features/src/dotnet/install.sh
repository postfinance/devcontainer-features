. ./functions.sh

"./installer_$(detect_arch)" \
-version="${VERSION:-"latest"}" \
-additionalVersions="${ADDITIONALVERSIONS:-""}" \
-dotnetRuntimeVersions="${DOTNETRUNTIMEVERSIONS:-""}" \
-aspNetCoreRuntimeVersions="${ASPNETCORERUNTIMEVERSIONS:-""}" \
-workloads="${WORKLOADS:-""}" \
-downloadUrl="${DOWNLOADURL:-""}" \
-versionsUrl="${VERSIONSURL:-""}" \
-nugetConfigPath="${NUGETCONFIGPATH:-""}"
