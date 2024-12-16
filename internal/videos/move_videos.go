package videos

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

func MoveJav(srcFilePath string, destDir string) error {
	// create destiny directory if not exits
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			return err
		}
	}

	_, fileName := path.Split(srcFilePath)
	ext := path.Ext(fileName)
	fileName = strings.ToUpper(fileName[0 : len(fileName)-len(ext)])
	destFilePath := path.Join(destDir, fileName+ext)

	err := os.Rename(srcFilePath, destFilePath)
	if err != nil && strings.Contains(err.Error(), "invalid cross device link") {
		return MoveJavCrossDisk(srcFilePath, destFilePath)
	} else if err != nil {
		return err
	}

	return nil
}

func MoveJavCrossDisk(srcFilePath string, destFilePath string) error {
	src, err := os.Open(srcFilePath)
	if err != nil {
		return errors.Wrap(err, "Open(srcFilePath)")
	}
	dst, err := os.Create(destFilePath)
	if err != nil {
		src.Close()
		return errors.Wrap(err, "Create(destFilePath)")
	}
	_, err = io.Copy(dst, src)
	src.Close()
	dst.Close()

	if err != nil {
		return errors.Wrap(err, "Copy")
	}
	fi, err := os.Stat(srcFilePath)
	if err != nil {
		os.Remove(destFilePath)
		return errors.Wrap(err, "Stat")
	}
	err = os.Chmod(destFilePath, fi.Mode())
	if err != nil {
		os.Remove(destFilePath)
		return errors.Wrap(err, "Stat")
	}
	os.Remove(srcFilePath)
	return nil
}
