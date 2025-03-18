package config

import (
	l "github.com/faelmori/logz/logger"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"os"
	"path/filepath"
	"strings"
)

func (v *ConfigManager) saveConfig() error {
	if v.globals == nil {
		return nil
	}
	if v.path == "" {
		v.path, _ = v.GetConfigDir()
	}
	if v.name == "" {
		v.name = v.GetConfigName()
	}
	return v.globals.WriteConfigAs(filepath.Join(v.path, v.name))
}
func (v *ConfigManager) getLogz() *l.LogzCore { return v.logz }
func (v *ConfigManager) GetGitDir() string {
	wd, _ := os.Getwd()
	path := strings.Split(wd, "/")
	for i := len(path); i >= 0; i-- {
		subPath := strings.Join(path[:i], "/")
		fileInfo, err := os.Stat(subPath + "/.git")
		if err == nil && fileInfo.IsDir() {
			return subPath
		}
	}
	return ""
}
func (v *ConfigManager) GetGitRemote() string {
	gitdir := v.GetGitDir() + "/.git"
	fs := osfs.New(gitdir)
	s := filesystem.NewStorageWithOptions(fs, cache.NewObjectLRUDefault(), filesystem.Options{KeepDescriptors: true})
	r, err := git.Open(s, fs)
	if err == nil {
		remotes, _ := r.Remotes()
		return remotes[0].Config().URLs[0]
	}
	return ""
}
