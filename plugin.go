// Package plasmactlpackage implements a package launchr plugin
package plasmactlpackage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/launchrctl/launchr"
	"github.com/launchrctl/launchr/pkg/log"
	"github.com/spf13/cobra"
)

func init() {
	launchr.RegisterPlugin(&Plugin{})
}

// Plugin is launchr plugin providing bump action.
type Plugin struct{}

// PluginInfo implements launchr.Plugin interface.
func (p *Plugin) PluginInfo() launchr.PluginInfo {
	return launchr.PluginInfo{}
}

// OnAppInit implements launchr.Plugin interface.
func (p *Plugin) OnAppInit(_ launchr.App) error {
	return nil
}

// CobraAddCommands implements launchr.CobraPlugin interface to provide bump functionality.
func (p *Plugin) CobraAddCommands(rootCmd *cobra.Command) error {
	var pkgCmd = &cobra.Command{
		Use:   "package",
		Short: "Creates an archive to contain composed-compiled-propagated artifact",
		RunE: func(cmd *cobra.Command, args []string) error {
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
		log.Debug("%s", err)
		return errors.New("error getting repository information")
	}

	archiveDir := filepath.Join(os.Getenv("HOME"), "artifact")
	artifactDir := ".compose/artifacts"
	archiveFile := fmt.Sprintf("%s-%s-plasma-src.tar.gz", repoName, lastCommitShortSHA)

	log.Info("REPO_NAME=%s", repoName)
	log.Info("LAST_COMMIT_MESSAGE=%s", lastCommitMessage)
	log.Info("LAST_COMMIT_SHORT_SHA=%s", lastCommitShortSHA)
	log.Info("ARCHIVE_FILE=%s", archiveFile)
	log.Info("HOME=%s", os.Getenv("HOME"))

	buildDir := ".compose/build"
	// ensure ./compose/build exists before archiving it.
	_, err = os.Stat(buildDir)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("build dir doesn't exist, please check your composition")
		}
	}

	fmt.Printf("Creating artifact %s/%s...\n", artifactDir, archiveFile)
	if err = createArchive(buildDir, archiveDir, artifactDir, archiveFile); err != nil {
		log.Debug("%v", err)
		return errors.New("error creating archive")
	}

	fmt.Println("Done")
	return nil
}
