. ./functions.sh

"./installer_$(detect_arch)" \
    -chromeVersion="${CHROMEVERSION:-"none"}" \
    -useChromeForTesting="${USECHROMEFORTESTING:-"true"}" \
    -firefoxVersion="${FIREFOXVERSION:-"none"}" \
    -chromeDownloadUrl="${CHROMEDOWNLOADURL:-""}" \
    -chromeVersionsUrl="${CHROMEVERSIONSURL:-""}" \
    -chromeTestingVersionsUrl="${CHROMETESTINGVERSIONSURL:-""}" \
    -firefoxDownloadUrl="${FIREFOXDOWNLOADURL:-""}" \
    -firefoxVersionsUrl="${FIREFOXVERSIONSURL:-""}" 