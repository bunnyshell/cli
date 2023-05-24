package rdev

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/interactive"
)

const (
	directoryPermissions = 0o755
	filePermissions      = 0o644
)

var (
	errDirectoryNotProvided = errors.New("directory not provided")
	errFileAlreadyExists    = errors.New("file already exists")
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

	if file == nil {
		fmt.Fprintf(os.Stdout, `Skipping "%s" ...`+"\n", filename)

		return nil
	}

	fmt.Fprintf(os.Stdout, `Generating "%s" ...`+"\n", filename)

	return tpl.Execute(file, data)
}

func ensureDirectory(directory string) error {
	if directory == "" {
		return errDirectoryNotProvided
	}

	return os.MkdirAll(directory, directoryPermissions)
}

func openFile(directory string, name string) (*os.File, error) {
	fileName := filepath.Join(directory, name)

	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		if config.GetSettings().NonInteractive {
			return nil, errFileAlreadyExists
		}

		clobber, _ := interactive.Confirm(
			fmt.Sprintf(`The file "%s" already exists, do you want to overwrite it?`, fileName),
		)

		if !clobber {
			return nil, nil
		}
	}

	return os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePermissions)
}
