package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/roemer/gover"
	"github.com/samber/lo"
)

//////////
// Configuration
//////////

// Full Regex versioning, like v3.15.0-rc.2
var semVerRegex *regexp.Regexp = regexp.MustCompile(`^v?(?P<raw>(\d+)\.(\d+)(?:\.(\d+))?(?:-([a-z]+)(?:\.?(\d+))?)?)$`)

// Regex with 2-3 digits like v1.0 or v2.3.4
var threeDigitRegex *regexp.Regexp = regexp.MustCompile(`^v?(?P<raw>(\d+)\.(\d+)(?:\.(\d+))?)$`)

//////////
// Main
//////////

func main() {
	if err := runMain(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runMain() error {
	// Handle the flags
	version := flag.String("version", "latest", "")
	kubectxVersion := flag.String("kubectxVersion", "latest", "")
	kubensVersion := flag.String("kubensVersion", "latest", "")
	k9sVersion := flag.String("k9sVersion", "none", "")
	helmVersion := flag.String("helmVersion", "none", "")
	kustomizeVersion := flag.String("kustomizeVersion", "none", "")
	kubeconformVersion := flag.String("kubeconformVersion", "none", "")
	kubescoreVersion := flag.String("kubescoreVersion", "none", "")
	downloadUrl := flag.String("downloadUrl", "", "")
	kubectxDownloadUrl := flag.String("kubectxDownloadUrl", "", "")
	kubensDownloadUrl := flag.String("kubensDownloadUrl", "", "")
	k9sDownloadUrl := flag.String("k9sDownloadUrl", "", "")
	helmDownloadUrl := flag.String("helmDownloadUrl", "", "")
	kustomizeDownloadUrl := flag.String("kustomizeDownloadUrl", "", "")
	kubeconformDownloadUrl := flag.String("kubeconformDownloadUrl", "", "")
	kubescoreDownloadUrl := flag.String("kubescoreDownloadUrl", "", "")
	flag.Parse()

	// Load settings from an external file
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	// Apply override logic for URLs
	installer.HandleOverride(downloadUrl, "https://dl.k8s.io/release", "kubectl-download-url")
	installer.HandleGitHubOverride(kubectxDownloadUrl, "ahmetb/kubectx", "kubectl-kubectx-download-url")
	installer.HandleGitHubOverride(kubensDownloadUrl, "ahmetb/kubectx", "kubectl-kubens-download-url") // Yes, this is the kubectx repo
	installer.HandleGitHubOverride(k9sDownloadUrl, "derailed/k9s", "kubectl-k9s-download-url")
	installer.HandleOverride(helmDownloadUrl, "https://get.helm.sh", "kubectl-helm-download-url")
	installer.HandleGitHubOverride(kustomizeDownloadUrl, "kubernetes-sigs/kustomize", "kubectl-kustomize-download-url")
	installer.HandleGitHubOverride(kubeconformDownloadUrl, "yannh/kubeconform", "kubectl-kubeconform-download-url")
	installer.HandleGitHubOverride(kubescoreDownloadUrl, "zegl/kube-score", "kubectl-kubescore-download-url")

	// Create and process the feature
	feature := installer.NewFeature("kubectl", true,
		&kubectlComponent{
			ComponentBase: installer.NewComponentBase("kubectl", *version),
			DownloadUrl:   *downloadUrl,
		},
		&kubectxComponent{
			ComponentBase: installer.NewComponentBase("kubectx", *kubectxVersion),
			DownloadUrl:   *kubectxDownloadUrl,
		},
		&kubensComponent{
			ComponentBase: installer.NewComponentBase("kubens", *kubensVersion),
			DownloadUrl:   *kubensDownloadUrl,
		},
		&k9sComponent{
			ComponentBase: installer.NewComponentBase("k9s", *k9sVersion),
			DownloadUrl:   *k9sDownloadUrl,
		},
		&helmComponent{
			ComponentBase: installer.NewComponentBase("helm", *helmVersion),
			DownloadUrl:   *helmDownloadUrl,
		},
		&kustomizeComponent{
			ComponentBase: installer.NewComponentBase("kustomize", *kustomizeVersion),
			DownloadUrl:   *kustomizeDownloadUrl,
		},
		&kubeconformComponent{
			ComponentBase: installer.NewComponentBase("kubeconform", *kubeconformVersion),
			DownloadUrl:   *kubeconformDownloadUrl,
		},
		&kubescoreComponent{
			ComponentBase: installer.NewComponentBase("kubescore", *kubescoreVersion),
			DownloadUrl:   *kubescoreDownloadUrl,
		},
		&fzfComponent{
			ComponentBase: installer.NewComponentBase("fzf", installer.VERSION_SYSTEM_DEFAULT),
		},
	)
	return feature.Process()
}

//////////
// kubectl
//////////

type kubectlComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *kubectlComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("kubernetes", "kubernetes")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, semVerRegex, true)
}

func (c *kubectlComponent) InstallVersion(version *gover.Version) error {
	// Download
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	downloadUrl := fmt.Sprintf("%s/v%s/bin/linux/%s/kubectl", c.DownloadUrl, version.Raw, archPart)
	fileName := "kubectl"
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "kubectl"); err != nil {
		return err
	}
	// Install
	if err := installer.Tools.System.InstallBinaryToUsrLocalBin(fileName, "kubectl"); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
}

//////////
// kubectx
//////////

type kubectxComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *kubectxComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("ahmetb", "kubectx")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, threeDigitRegex, true)
}

func (c *kubectxComponent) InstallVersion(version *gover.Version) error {
	// Download
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x86_64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("kubectx_v%s_linux_%s.tar.gz", version.Raw, archPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "v"+version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "kubectx"); err != nil {
		return err
	}
	// Extract
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "kubectx", false); err != nil {
		return err
	}
	// Install
	if err := installer.Tools.System.InstallBinaryToUsrLocalBin("kubectx/kubectx", "kubectx"); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	if err := os.RemoveAll("kubectx"); err != nil {
		return err
	}
	return nil
}

//////////
// kubens
//////////

type kubensComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *kubensComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("ahmetb", "kubectx") // Yes, this is the kubectx repo
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, threeDigitRegex, true)
}

func (c *kubensComponent) InstallVersion(version *gover.Version) error {
	// Download
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x86_64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("kubens_v%s_linux_%s.tar.gz", version.Raw, archPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "v"+version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "kubens"); err != nil {
		return err
	}
	// Extract
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "kubens", false); err != nil {
		return err
	}
	// Install
	if err := installer.Tools.System.InstallBinaryToUsrLocalBin("kubens/kubens", "kubens"); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	if err := os.RemoveAll("kubens"); err != nil {
		return err
	}
	return nil
}

//////////
// k9s
//////////

type k9sComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *k9sComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("derailed", "k9s")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, threeDigitRegex, true)
}

func (c *k9sComponent) InstallVersion(version *gover.Version) error {
	// Download
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("k9s_Linux_%s.tar.gz", archPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "v"+version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "k9s"); err != nil {
		return err
	}
	// Extract
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "k9s", false); err != nil {
		return err
	}
	// Install
	if err := installer.Tools.System.InstallBinaryToUsrLocalBin("k9s/k9s", "k9s"); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	if err := os.RemoveAll("k9s"); err != nil {
		return err
	}
	return nil
}

//////////
// helm
//////////

type helmComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *helmComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("helm", "helm")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, semVerRegex, true)
}

func (c *helmComponent) InstallVersion(version *gover.Version) error {
	// Download
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("helm-v%s-linux-%s.tar.gz", version.Raw, archPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "helm"); err != nil {
		return err
	}
	// Extract
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "helm", false); err != nil {
		return err
	}
	// Install
	if err := installer.Tools.System.InstallBinaryToUsrLocalBin(fmt.Sprintf("helm/linux-%s/helm", archPart), "helm"); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	if err := os.RemoveAll("helm"); err != nil {
		return err
	}
	return nil
}

//////////
// kustomize
//////////

type kustomizeComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *kustomizeComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("kubernetes-sigs", "kustomize")
	if err != nil {
		return nil, err
	}
	filterRegex := regexp.MustCompile("kustomize/(.*)")
	filteredTags := lo.FilterMap(tags, func(tag string, _ int) (string, bool) {
		matches := filterRegex.FindStringSubmatch(tag)
		if len(matches) > 1 {
			return matches[1], true
		}
		return "", false
	})
	return installer.Tools.Versioning.ParseVersionsFromList(filteredTags, threeDigitRegex, true)
}

func (c *kustomizeComponent) InstallVersion(version *gover.Version) error {
	// Download
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("kustomize_v%s_linux_%s.tar.gz", version.Raw, archPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "kustomize/v"+version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "kustomize"); err != nil {
		return err
	}
	// Extract
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "kustomize", false); err != nil {
		return err
	}
	// Install
	if err := installer.Tools.System.InstallBinaryToUsrLocalBin("kustomize/kustomize", "kustomize"); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	if err := os.RemoveAll("kustomize"); err != nil {
		return err
	}
	return nil
}

//////////
// kubeconform
//////////

type kubeconformComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *kubeconformComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("yannh", "kubeconform")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, threeDigitRegex, true)
}

func (c *kubeconformComponent) InstallVersion(version *gover.Version) error {
	// Download
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("kubeconform-linux-%s.tar.gz", archPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "v"+version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "kubeconform"); err != nil {
		return err
	}
	// Extract
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "kubeconform", false); err != nil {
		return err
	}
	// Install
	if err := installer.Tools.System.InstallBinaryToUsrLocalBin("kubeconform/kubeconform", "kubeconform"); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	if err := os.RemoveAll("kubeconform"); err != nil {
		return err
	}
	return nil
}

//////////
// kubescore
//////////

type kubescoreComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *kubescoreComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("zegl", "kube-score")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, threeDigitRegex, true)
}

func (c *kubescoreComponent) InstallVersion(version *gover.Version) error {
	// Download
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("kube-score_%s_linux_%s.tar.gz", version.Raw, archPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "v"+version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "kubescore"); err != nil {
		return err
	}
	// Extract
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "kubescore", false); err != nil {
		return err
	}
	// Install
	if err := installer.Tools.System.InstallBinaryToUsrLocalBin("kubescore/kube-score", "kube-score"); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	if err := os.RemoveAll("kubescore"); err != nil {
		return err
	}
	return nil
}

//////////
// fzf
//////////

type fzfComponent struct {
	*installer.ComponentBase
}

func (c *fzfComponent) InstallVersion(version *gover.Version) error {
	return installer.Tools.System.InstallPackages("fzf")
}
