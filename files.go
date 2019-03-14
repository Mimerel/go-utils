package go_utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ScanDirectory(path string, extensions []string, ignore []string) (folders []string, files []string, scanError []error) {
	filepath.Walk(path, func(path string, f os.FileInfo, err error) (error) {
		skip := false
		for _, i := range ignore {

			if strings.Index(path, i) != -1 {
				skip = true
			}
		}
		hasExtension := false
		for _, i := range extensions {
			if strings.HasSuffix(path, i) {
				hasExtension = true
			}
		}

		if skip == false && hasExtension {
			f, err = os.Stat(path)
			if err != nil {
				scanError = append(scanError, err)
			}
			f_mode := f.Mode()
			if f_mode.IsDir() {
				folders = append(folders, path)
			} else if f_mode.IsRegular(){
				files = append(files, path)
			}
		}
		return nil
	})
	return folders, files, scanError
}

func CopyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}