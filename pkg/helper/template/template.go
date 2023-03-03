package template

import (
	"errors"
	"fmt"
	"os"
)

const (
	directoryPermissions = 0o755
	filePermissions      = 0o644
)

var (
	errDirectoryNotProvided   = errors.New("directory not provided")
	errDirectoryAlreadyExists = errors.New("directory already exists")
)

func NewData(directory string) *Data {
	return &Data{
		Directory: directory,
	}
}

func Generate(directory string) error {
	if err := ensureDirectory(directory); err != nil {
		return fmt.Errorf("%w (%s)", err, directory)
	}

	filenames, err := getAllFilenames()
	if err != nil {
		return err
	}

	templateData := NewData(directory)

	for _, filename := range filenames {
		if err = create(*templateData, filename); err != nil {
			return err
		}
	}

	return nil
}

func create(data Data, filename string) error {
	tpl, err := GetTemplate(filename)
	if err != nil {
		return err
	}

	file, err := openFile(data.Directory, filename)
	if err != nil {
		return err
	}

	fmt.Println("Generating " + filename + "...")
	return tpl.Execute(file, data)
}

func ensureDirectory(directory string) error {
	if directory == "" {
		return errDirectoryNotProvided
	}

	// check if directory exists
	_, err := os.Stat(directory)
	if !os.IsNotExist(err) {
		return errDirectoryAlreadyExists
	}

	return os.MkdirAll(directory, directoryPermissions)
}

func openFile(directory, name string) (*os.File, error) {
	return os.OpenFile(directory+"/"+name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePermissions)
}
