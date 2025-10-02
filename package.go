package plasmactlpackage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/launchrctl/launchr/pkg/archive"
)

func getRepoInfo() (repoName, lastCommitMessage, lastCommitShortSHA string, err error) {
	// Open repository
	r, err := git.PlainOpen(".")
	if err != nil {
		return "", "", "", err
	}

	// Get repository name
	remote, err := r.Remote("origin")
	if err != nil {
		return "", "", "", err
	}
	repoName = remote.Config().URLs[0]
	repoName = filepath.Base(repoName)
	repoName = repoName[:len(repoName)-4]

	// Get last commit information
	ref, err := r.Head()
	if err != nil {
		return "", "", "", err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return "", "", "", err
	}
	lastCommitMessage = commit.Message
	lastCommitShortSHA = ref.Hash().String()[:7]

	return repoName, lastCommitMessage, lastCommitShortSHA, nil
}

func createArchive(srcDir, archiveTempDir, archiveFinalDir, archiveDestFile string) error {
	// Ensure archive directory exists
	if err := os.MkdirAll(archiveTempDir, 0750); err != nil {
		return err
	}
	if err := os.MkdirAll(archiveFinalDir, 0750); err != nil {
		return err
	}

	// Create tar.gz archive
	archivePath := filepath.Join(archiveTempDir, archiveDestFile)
	artifactPath := filepath.Join(archiveFinalDir, archiveDestFile)
	tarDst, err := os.Create(filepath.Clean(archivePath))
	if err != nil {
		return err
	}
	defer tarDst.Close()

	tarSrc, err := archive.Tar(archive.CopyInfo{
		Path: archivePath,
	}, archive.CopyInfo{}, &archive.TarOptions{
		Compression: archive.Gzip,
		ExcludePatterns: []string{
			".git",
			".compose",
			".plasmactl",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to archive the package: %w", err)
	}
	if _, err = io.Copy(tarDst, tarSrc); err != nil {
		return fmt.Errorf("failed to archive the package: %w", err)
	}

	// Delete temp file from archive temp directory.
	err = os.Rename(
		filepath.Clean(archivePath),
		filepath.Clean(artifactPath),
	)
	if err != nil {
		return fmt.Errorf("error deleting temp file: %v", err)
	}

	return nil
}
