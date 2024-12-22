// Package plasmactlpackage implements a package launchr plugin
package plasmactlpackage

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/launchrctl/launchr"
	"github.com/launchrctl/launchr/pkg/action"
)

//go:embed action.yaml
var actionYaml []byte

func init() {
	launchr.RegisterPlugin(&Plugin{})
}

// Plugin is [launchr.Plugin] providing package action.
type Plugin struct{}

// PluginInfo implements [launchr.Plugin] interface.
func (p *Plugin) PluginInfo() launchr.PluginInfo {
	return launchr.PluginInfo{}
}

// DiscoverActions implements [launchr.ActionDiscoveryPlugin] interface.
func (p *Plugin) DiscoverActions(_ context.Context) ([]*action.Action, error) {
	a := action.NewFromYAML("package", actionYaml)
	a.SetRuntime(action.NewFnRuntime(func(_ context.Context, _ *action.Action) error {
		return createArtifact()
	}))
	return []*action.Action{a}, nil
}

func createArtifact() error {
	repoName, lastCommitMessage, lastCommitShortSHA, err := getRepoInfo()
	if err != nil {
		launchr.Log().Error("error", "error", err)
		return errors.New("error getting repository information")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	archiveDir := filepath.Join(homeDir, "artifact")
	artifactDir := ".compose/artifacts"
	archiveFile := fmt.Sprintf("%s-%s-plasma-src.tar.gz", repoName, lastCommitShortSHA)

	launchr.Log().Info("initialize artifact",
		"REPO_NAME", repoName,
		"LAST_COMMIT_MESSAGE", lastCommitMessage,
		"LAST_COMMIT_SHORT_SHA", lastCommitShortSHA,
		"ARCHIVE_FILE", archiveFile,
		"HOME", homeDir,
	)

	buildDir := ".compose/build"
	// ensure ./compose/build exists before archiving it.
	_, err = os.Stat(buildDir)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("build dir doesn't exist, please check your composition")
		}
	}

	launchr.Term().Printfln("Creating artifact %s/%s...", artifactDir, archiveFile)
	if err = createArchive(buildDir, archiveDir, artifactDir, archiveFile); err != nil {
		launchr.Log().Error("error", "error", err)
		return errors.New("error creating archive")
	}

	launchr.Term().Success().Println("Done")
	return nil
}
