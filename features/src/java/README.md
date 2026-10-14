# Java (java)

Installs Java (Temurin/OpenJDK), Maven, Gradle, and Ant.

## Example Usage

```json
"features": {
    "ghcr.io/postfinance/devcontainer-features/java:1.0.0": {
        "version": "latest",
        "mavenVersion": "none",
        "gradleVersion": "none",
        "antVersion": "none",
        "downloadUrl": "",
        "versionsUrl": "",
        "latestUrl": "",
        "mavenDownloadUrl": "",
        "mavenVersionsUrl": "",
        "gradleDownloadUrl": "",
        "gradleVersionsUrl": "",
        "antDownloadUrl": "",
        "antVersionsUrl": ""
    }
}
```

## Options

| Option | Description | Type | Default Value | Proposals |
|-----|-----|-----|-----|-----|
| version | The version of Java (Temurin/OpenJDK) to install. Use a major version (e.g. '21') to get the latest patch release. | string | latest | latest, 21, 17, 11, 21.0.5 |
| mavenVersion | The version of Maven to install. Use 'none' to skip. | string | none | none, latest, 3.9, 3.9.9 |
| gradleVersion | The version of Gradle to install. Use 'none' to skip. | string | none | none, latest, 8, 8.14 |
| antVersion | The version of Ant to install. Use 'none' to skip. | string | none | none, latest, 1.10, 1.10.15 |
| downloadUrl | The download URL to use for Java (Temurin) binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/adoptium-remote/v3/binary |
| versionsUrl | The URL to fetch the available Java (Temurin) versions from. | string | &lt;empty&gt; |  |
| latestUrl | The URL to fetch the latest Java (Temurin) release information from. | string | &lt;empty&gt; |  |
| mavenDownloadUrl | The download URL to use for Maven binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/apache-maven-remote |
| mavenVersionsUrl | The URL to fetch the available Maven versions from. | string | &lt;empty&gt; |  |
| gradleDownloadUrl | The download URL to use for Gradle binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/gradle-distributions-remote |
| gradleVersionsUrl | The URL to fetch the available Gradle versions from. | string | &lt;empty&gt; |  |
| antDownloadUrl | The download URL to use for Ant binaries. | string | &lt;empty&gt; | https://mycompany.com/artifactory/apache-ant-remote |
| antVersionsUrl | The URL to fetch the available Ant versions from. | string | &lt;empty&gt; |  |

## Customizations

### VS Code Extensions

- `vscjava.vscode-java-pack`

## Notes

### System Compatibility

Debian, Ubuntu

### Accessed Urls

Needs access to the following URLs for downloading and resolving Java (Temurin/OpenJDK):
* https://api.adoptium.net

Needs access to the following URL for downloading and resolving Maven:
* https://downloads.apache.org/maven

Needs access to the following URLs for downloading and resolving Gradle:
* https://services.gradle.org

Needs access to the following URL for downloading and resolving Ant:
* https://downloads.apache.org/ant
