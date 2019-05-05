package go_utils

import (
	"archive/zip"
	"fmt"
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
	FullName  string
	Name      string
	Extension string
	Path      string
	FullPath  string
	Size      int64
}

func NewParams() (params *ScanParams) {
	return &ScanParams{}
}

func (params *ScanParams) ScanFolder() {

	_ = filepath.Walk(params.Request.Path, func(path string, f os.FileInfo, err error) (error) {
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
		if skip == false && hasExtension && ageConstraint {
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
				newFile.Path = strings.Replace(newFile.FullPath, newFile.FullName, "", 1)
				tempSplit := strings.Split(newFile.FullName, ".")
				newFile.Extension = "." + tempSplit[len(tempSplit)-1]
				newFile.Name = strings.Replace(newFile.FullName, newFile.Extension, "", 1)
				newFile.Size = f.Size()
				params.Result.Files = append(params.Result.Files, newFile)
				//fmt.Printf("data : %+v\n", newFile)
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

/**
Method that unzips a file from the given location to the given destination
*/
func Unzip(src string, dest string) (error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	defer r.Close()

	if err != nil {
		return err
	}

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			_ = os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		_ = outFile.Close()
		_ = rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}