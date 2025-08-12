package main

import (
	"builder/installer"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/roemer/gotaskr"
	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/gttools"
	"github.com/roemer/gotaskr/log"
)

////////////////////////////////////////////////////////////
// Variables
////////////////////////////////////////////////////////////

var featureList = []string{
	"go",
	"zig",
}

////////////////////////////////////////////////////////////
// Main
////////////////////////////////////////////////////////////

func main() {
	os.Exit(gotaskr.Execute())
}

////////////////////////////////////////////////////////////
// Initialize Tasks
////////////////////////////////////////////////////////////

func init() {
	gotaskr.Task("Update-Readme-Files", func() error {
		for _, feature := range featureList {
			if err := BuildReadmeForFeature(fmt.Sprintf("features/src/%s", feature)); err != nil {
				return err
			}
		}
		return nil
	})

	////////// Go
	gotaskr.Task("Feature:go:Package", func() error {
		return packageFeature("go")
	})
	gotaskr.Task("Feature:go:Test", func() error {
		return testFeature("go")
	})
	gotaskr.Task("Feature:go:Publish", func() error {
		return publishFeature("go")
	})

	////////// Zig
	gotaskr.Task("Feature:zig:Package", func() error {
		return packageFeature("zig")
	})
	gotaskr.Task("Feature:zig:Test", func() error {
		return testFeature("zig")
	})
	gotaskr.Task("Feature:zig:Publish", func() error {
		return publishFeature("zig")
	})
}

////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////

func buildGo(workingDirectory string, binaryName string) ([]string, error) {
	// Check if a go installer exists and only compile it then
	if _, err := os.Stat(filepath.Join(workingDirectory, "installer.go")); err != nil {
		if os.IsNotExist(err) {
			log.Information("No go installer found, skip compiling")
			return nil, nil
		} else {
			return nil, err
		}
	}

	// Compile for x86_64 and arm
	totalSize := int64(0)
	buildBinaries := []string{}
	for _, arch := range []string{"amd64", "arm64"} {
		archBinaryName := fmt.Sprintf("%s_%s", binaryName, arch)
		log.Informationf("Compiling for architecture %s", arch)
		// Force static linking
		os.Setenv("CGO_ENABLED", "0")
		// Compile the go installer
		// Optimizations:
		//   ldflags -w (disable DWARF generation), -s (disable symbol table and debug information)
		//   gcflags -l (disable inlining), -B (disable bounds checking)
		cmd := exec.Command("go", "build", "-o", archBinaryName, "-ldflags", "-w -s", "-gcflags", "all=-l -B", ".")
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GOOS=linux")
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", arch))
		cmd.Dir = workingDirectory
		if err := execr.RunCommand(true, cmd); err != nil {
			return buildBinaries, err
		}
		buildBinaries = append(buildBinaries, archBinaryName)
		fullPath := filepath.Join(workingDirectory, archBinaryName)
		fi, err := os.Stat(fullPath)
		if err != nil {
			return buildBinaries, err
		}
		totalSize += fi.Size()
		log.Informationf("Built %s with a size of %s", fullPath, installer.HumanizeBytes(fi.Size(), true))
	}
	log.Informationf("Total size of all binaries: %s", installer.HumanizeBytes(totalSize, true))
	return buildBinaries, nil
}

func packageFeature(featureName string) error {
	featurePath := path.Join("features/src", featureName)

	// Build the installer
	buildBinaries, err := buildGo(featurePath, "installer")
	if err != nil {
		return err
	}
	defer func() {
		for _, binary := range buildBinaries {
			os.Remove(filepath.Join(featurePath, binary))
		}
	}()

	// Package the feature
	settings := &gttools.DevContainerCliFeaturesPackageSettings{
		Target:                 featurePath,
		ForceCleanOutputFolder: gttools.True,
	}
	settings.OutputToConsole = true
	return gotaskr.Tools.DevContainerCli.FeaturesPackage(settings)
}

func testFeature(featureName string) error {
	// Prepare the temporary folder for the devcontainer spec
	testPath := ".scenario-test"
	os.RemoveAll(testPath)
	os.MkdirAll(testPath, os.ModePerm)

	// Read the images that should be used for testing
	testImagesFile := path.Join("features/test", featureName, "test-images.json")
	testImagesContent, err := os.ReadFile(testImagesFile)
	if err != nil {
		return err
	}
	var testImages []string
	if err := json.Unmarshal(testImagesContent, &testImages); err != nil {
		return err
	}

	// Read and parse the scenario file
	scenariosFile := path.Join("features/test", featureName, "scenarios.json")
	fileContent, err := os.ReadFile(scenariosFile)
	if err != nil {
		return err
	}
	var jsonData map[string]json.RawMessage
	if err := json.Unmarshal(fileContent, &jsonData); err != nil {
		return err
	}
	// Loop thru the scenarios
	for scenarioName, scenarioContent := range jsonData {
		log.Informationf("Processing scenario '%s'", scenarioName)

		// Loop thru the base images
		for _, testImage := range testImages {
			log.Informationf("Testing with image '%s'", testImage)

			// Clear and prepare the devcontainer path
			devcontainerPath := path.Join(testPath, ".devcontainer")
			os.RemoveAll(devcontainerPath)
			os.MkdirAll(devcontainerPath, os.ModePerm)

			// Write the devcontainer spec file
			devcontainerSpecPath := path.Join(devcontainerPath, "devcontainer.json")
			if err := os.WriteFile(devcontainerSpecPath, scenarioContent, os.ModePerm); err != nil {
				return err
			}

			// Copy the verify-script
			data, err := os.ReadFile(path.Join("features/test", featureName, scenarioName+".sh"))
			if err != nil {
				return err
			}
			if err := os.WriteFile(path.Join(devcontainerPath, "check.sh"), data, os.ModePerm); err != nil {
				return err
			}

			// Copy the functions.sh file
			data, err = os.ReadFile(path.Join("features/test/functions.sh"))
			if err != nil {
				return err
			}
			if err := os.WriteFile(path.Join(devcontainerPath, "functions.sh"), data, os.ModePerm); err != nil {
				return err
			}

			// Write the Dockerfile
			if err := os.WriteFile(path.Join(devcontainerPath, "Dockerfile"), []byte(fmt.Sprintf(`
				FROM %s
				ADD check.sh /tmp/check.sh
				ADD functions.sh /tmp/functions.sh
			`, testImage)), os.ModePerm); err != nil {
				return err
			}

			// Copy the required feature
			originalFeaturePath := path.Join("features/src", featureName)
			copiedFeaturePath := path.Join(devcontainerPath, featureName)
			if err := os.CopyFS(copiedFeaturePath, os.DirFS(originalFeaturePath)); err != nil {
				return err
			}

			// Build the go installer inside the feature
			if _, err := buildGo(copiedFeaturePath, "installer"); err != nil {
				return err
			}

			// Build the devcontainer
			imageName := fmt.Sprintf("dev-container-feature-%s-scenario-%s-test", featureName, scenarioName)
			if err := gotaskr.Tools.DevContainerCli.Build(&gttools.DevContainerCliBuildSettings{
				ToolSettingsBase: gttools.ToolSettingsBase{OutputToConsole: true},
				WorkspaceFolder:  testPath,
				ImageNames:       []string{imageName},
			}); err != nil {
				return err
			}
			defer execr.Run(false, "docker", "image", "rm", imageName)

			// Run the check in the container
			checkError := execr.Run(true, "docker", "run", "-t", "--rm", "-v", "/var/run/docker.sock:/var/run/docker.sock", imageName, "sh", "-c", "/tmp/check.sh")
			if checkError != nil {
				return fmt.Errorf("check failed: %w", checkError)
			}
			fmt.Println("Check was successfull")
		}
	}

	return nil

	// TODO: Somewhen in the future this can be done with the devcontainer cli
	/*featurePath := path.Join("features/src", featureName)
	// Build the installer
	if err := buildGo(featurePath, "installer"); err != nil {
		return err
	}
	defer os.Remove(filepath.Join(featurePath, "installer"))

	if err := gotaskr.Tools.DevContainerCli.FeaturesTest(&gttools.DevContainerCliFeaturesTestSettings{
		ToolSettingsBase:  gttools.ToolSettingsBase{OutputToConsole: true},
		ProjectFolder:     "./features",
		Features:          []string{featureName},
		LogLevel:          gttools.DEV_CONTAINER_CLI_LOG_LEVEL_DEBUG,
		SkipAutogenerated: gttools.True,
		SkipDuplicated:    gttools.True,
	}); err != nil {
		return err
	}*/
}

func publishFeature(featureName string) error {
	registry := "ghcr.io"
	namespace := "postfinance/devcontainer-features"

	featurePath := path.Join("features/src", featureName)

	// Build the installer
	buildBinaries, err := buildGo(featurePath, "installer")
	if err != nil {
		return err
	}
	defer func() {
		for _, binary := range buildBinaries {
			os.Remove(filepath.Join(featurePath, binary))
		}
	}()

	// No authentication needed - DevContainerCLI supports GITHUB_TOKEN
	// os.Setenv("DEVCONTAINERS_OCI_AUTH", "ghcr.io|USERNAME|"+os.Getenv("GITHUB_TOKEN"))

	// Build and publish the feature
	settings := &gttools.DevContainerCliFeaturesPublishSettings{
		Target:    featurePath,
		Registry:  registry,
		Namespace: namespace,
	}
	settings.OutputToConsole = true
	return gotaskr.Tools.DevContainerCli.FeaturesPublish(settings)
}
