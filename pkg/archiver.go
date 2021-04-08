package bundle

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zen-xu/bundler/pkg/utils"
)

type Archiver interface {
	Archive(srcPath string, preservePath bool) error
	Close()
}

func NewArchiver(dest string) Archiver {
	fw, err := os.Create(dest)
	utils.CheckError(err, "Could not create archiver file")
	gw := gzip.NewWriter(fw)
	tw := tar.NewWriter(gw)
	return &archiver{
		outputFile: fw,
		gzipWriter: gw,
		tarWriter:  tw,
	}
}

type archiver struct {
	outputFile *os.File
	gzipWriter *gzip.Writer
	tarWriter  *tar.Writer
}

func (archiver *archiver) Close() {
	_ = archiver.tarWriter.Close()
	_ = archiver.gzipWriter.Close()
	_ = archiver.outputFile.Close()
}

func isDir(pth string) (bool, error) {
	fi, err := os.Stat(pth)
	if err != nil {
		return false, err
	}

	return fi.Mode().IsDir(), nil
}

func (archiver *archiver) Archive(srcPath string, preservePath bool) error {
	absPath, err := filepath.Abs(srcPath)
	if err != nil {
		return err
	}

	isDirectory, err := isDir(srcPath)
	utils.CheckError(err, "Could not determine if this is a directory.")

	if isDirectory || !preservePath {
		err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			var relative string
			if os.IsPathSeparator(srcPath[len(srcPath)-1]) {
				relative, err = filepath.Rel(absPath, path)
			} else {
				relative, err = filepath.Rel(filepath.Dir(absPath), path)
			}

			relative = filepath.ToSlash(relative)

			if err != nil {
				return err
			}

			return archiver.addTarFile(path, relative)
		})
	} else {
		fields := strings.Split(srcPath, string(os.PathSeparator))
		for idx := range fields {
			path := strings.Join(fields[:idx+1], string(os.PathSeparator))
			err := archiver.addTarFile(path, path)
			utils.CheckError(err, "Unable to archiver file")
		}
	}

	return err
}

func (archiver *archiver) addTarFile(path, name string) error {
	if strings.Contains(path, "..") {
		return errors.New("Path cannot contain a relative marker of '..': " + path)
	}
	fi, err := os.Lstat(path)
	if err != nil {
		return err
	}

	link := ""
	if fi.Mode()&os.ModeSymlink != 0 {
		if link, err = os.Readlink(path); err != nil {
			return err
		}
	}

	hdr, err := tar.FileInfoHeader(fi, link)
	if err != nil {
		return err
	}

	if fi.IsDir() && !os.IsPathSeparator(name[len(name)-1]) {
		name = name + "/"
	}

	if hdr.Typeflag == tar.TypeReg && name == "." {
		// archiving a single file
		hdr.Name = filepath.ToSlash(filepath.Base(path))
	} else {
		hdr.Name = filepath.ToSlash(name)
	}

	if err := archiver.tarWriter.WriteHeader(hdr); err != nil {
		return err
	}

	if hdr.Typeflag == tar.TypeReg {
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		_, err = io.Copy(archiver.tarWriter, file)
		if err != nil {
			return err
		}
	}

	return nil
}
