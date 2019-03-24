package go_utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ScanParams struct {
	Request Request
	Result  Result
}

type Request struct {
	Path       string
	Extensions []string
	Ignore     []string
	MinAge     time.Duration
}

type Result struct {
	Folders []string
	Files   []*Files
	Errors  []error
}

type Files struct {
	Name      string
	FullName  string
	FullPath  string
	Path      string
	Extension string
}

func NewParams() (params *ScanParams) {
	return &ScanParams{}
}

func (params *ScanParams) ScanFolder() {

	filepath.Walk(params.Request.Path, func(path string, f os.FileInfo, err error) (error) {
		skip := false
		hasExtension := true
		ageConstraint := true
		for _, i := range params.Request.Ignore {

			if strings.Index(path, i) != -1 {
				skip = true
			}
		}

		if len(params.Request.Extensions) > 0 {
			hasExtension = false
			for _, i := range params.Request.Extensions {
				if strings.HasSuffix(path, i) {
					hasExtension = true
				}
			}
		}
		if params.Request.MinAge > 0 {
			if f.ModTime().After(time.Now().Local().Add(- params.Request.MinAge)) {
				ageConstraint = false
			}
		}
		if skip == false && hasExtension && ageConstraint{
			f, err = os.Stat(path)
			if err != nil {
				params.Result.Errors = append(params.Result.Errors, err)
			}
			f_mode := f.Mode()
			if f_mode.IsDir() {

				params.Result.Folders = append(params.Result.Folders, path)
			} else if f_mode.IsRegular() {
				newFile := new(Files)
				newFile.FullName = f.Name()
				newFile.FullPath = path
				newFile.Path = strings.Replace(newFile.FullPath, newFile.Name, "", 1)
				tempSplit := strings.Split(newFile.FullName, ".")
				newFile.Extension = "." + tempSplit[len(tempSplit)-1]
				newFile.Name = strings.Replace(newFile.FullName, newFile.Extension, "", 1)
				params.Result.Files = append(params.Result.Files, newFile)
			}
		}
		return nil
	})
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