// Package plasmactlpackage implements a package launchr plugin
package plasmactlpackage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/launchrctl/launchr"
)

func init() {
	launchr.RegisterPlugin(&Plugin{})
}

// Plugin is [launchr.Plugin] providing package action.
type Plugin struct{}

// PluginInfo implements [launchr.Plugin] interface.
func (p *Plugin) PluginInfo() launchr.PluginInfo {
	return launchr.PluginInfo{}
}

// CobraAddCommands implements [launchr.CobraPlugin] interface to provide package functionality.
func (p *Plugin) CobraAddCommands(rootCmd *launchr.Command) error {
	var pkgCmd = &launchr.Command{
		Use:   "package",
		Short: "Creates an archive to contain composed-compiled-propagated artifact",
		RunE: func(cmd *launchr.Command, _ []string) error {
			// Don't show usage help on a runtime error.
			cmd.SilenceUsage = true

			return createArtifact()
		},
	}

	rootCmd.AddCommand(pkgCmd)
	return nil
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
