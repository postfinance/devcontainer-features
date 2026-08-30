. ./functions.sh

"./installer_$(detect_arch)" \
    -version="${VERSION:-"latest"}" \
    -mavenVersion="${MAVENVERSION:-"none"}" \
    -gradleVersion="${GRADLEVERSION:-"none"}" \
    -antVersion="${ANTVERSION:-"none"}" \
    -downloadUrl="${DOWNLOADURL:-""}" \
    -versionsUrl="${VERSIONSURL:-""}" \
    -latestUrl="${LATESTURL:-""}" \
    -mavenDownloadUrl="${MAVENDOWNLOADURL:-""}" \
    -mavenVersionsUrl="${MAVENVERSIONSURL:-""}" \
    -gradleDownloadUrl="${GRADLEDOWNLOADURL:-""}" \
    -gradleVersionsUrl="${GRADLEVERSIONSURL:-""}" \
    -antDownloadUrl="${ANTDOWNLOADURL:-""}" \
    -antVersionsUrl="${ANTVERSIONSURL:-""}"
