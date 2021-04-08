package bundle

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/zen-xu/bundler/pkg/utils"
)

type Bundler struct {
	config *Config
}

func NewBundler(configPath string) (*Bundler, error) {
	config, err := NewConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &Bundler{
		config: config,
	}, nil
}

func (b *Bundler) Bundle(outputPath string) (ignorePaths []string) {
	archivePath := "bundle.tar.gz"
	archiver := NewArchiver(archivePath)
	archivePaths, ignorePaths := b.config.GetArchivePaths()

	var err error
	for _, archivePath := range archivePaths {
		err = archiver.Archive(archivePath, true)
		utils.CheckError(err, fmt.Sprintf("Unable to add '%s' to bundle", archivePath))
	}
	archiver.Close()

	execute := `#!/bin/bash
set -eu
export TMPDIR=$(mktemp -d /tmp/bundler.XXXXXX)
ARCHIVE=$(awk '/^__BUNDLER_ARCHIVE__/ {print NR + 1; exit 0; }' $0)

tail -n+$ARCHIVE $0 | tar -xz -C $TMPDIR

pushd $TMPDIR > /dev/null
{{.Command}} $*
popd > /dev/null
rm -rf $TMPDIR

exit 0

__BUNDLER_ARCHIVE__
`
	var buff bytes.Buffer
	var values = struct {
		Command string
	}{
		Command: b.config.Command,
	}

	tmpl := template.New("bundle")
	tmpl, err = tmpl.Parse(execute)
	utils.CheckError(err, "Failed to parse Execute template")
	err = tmpl.Execute(&buff, values)
	utils.CheckError(err, "Failed to render Execute template")

	runnerFh, err := os.Create(outputPath)
	utils.CheckError(err, fmt.Sprintf("Unable to create runner executable file: %s", outputPath))
	defer runnerFh.Close()

	_, err = runnerFh.Write(buff.Bytes())
	utils.CheckError(err, fmt.Sprintf("Unable to write bootstrap script to runner executable file: %s", outputPath))

	archiveFh, err := os.Open(archivePath)
	utils.CheckError(err, "Unable to open payload file")
	defer archiveFh.Close()
	defer os.Remove(archivePath)

	_, err = io.Copy(runnerFh, archiveFh)
	utils.CheckError(err, fmt.Sprintf("Unable to write payload to runner executable file: %s", outputPath))

	err = os.Chmod(outputPath, 0755)
	utils.CheckError(err, "Unable to change runner permissions")

	return ignorePaths
}
