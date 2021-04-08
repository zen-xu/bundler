package bundle

import (
	"io/ioutil"
	"path/filepath"

	"github.com/gobwas/glob"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Command      string
	ArchivePaths []string `yaml:"archive_paths"`
	IgnorePaths  []string `yaml:"ignore_paths"`
}

func NewConfig(configPath string) (*Config, error) {
	yamlString, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	return newConfig(yamlString)
}

func newConfig(yamlString []byte) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal(yamlString, config)

	var matches []string

	var archivePaths []string
	for _, bundlePath := range config.ArchivePaths {
		matches, _ = filepath.Glob(bundlePath)
		archivePaths = append(archivePaths, matches...)
	}
	config.ArchivePaths = archivePaths

	return config, err
}

// GetArchivePaths will return archiver paths which exclude ignore paths
func (c *Config) GetArchivePaths() (archivePaths []string, ignorePaths []string) {
	for _, archivePath := range c.ArchivePaths {
		needArchive := true
		for _, ignorePath := range c.IgnorePaths {
			if glob.MustCompile(ignorePath).Match(archivePath) {
				needArchive = false
				break
			}
		}
		if needArchive {
			archivePaths = append(archivePaths, archivePath)
		} else {
			ignorePaths = append(ignorePaths, archivePath)
		}
	}

	return archivePaths, ignorePaths
}
