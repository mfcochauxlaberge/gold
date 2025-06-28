package gold

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// A Runner stores the configuration for the golden file runner.
type Runner struct {
	// Update reports whether or not the runner
	// must run in update mode. If it is the
	// case, files are created instead of being
	// compared.
	Update bool

	// The directory in which the golden files
	// will be stored. If empty, "testdata" is
	// used by default.
	Directory string

	// Filters holds a list of filters to apply
	// on the contents given to Test. They are
	// applied in the order they appear in the
	// slice.
	Filters []Filter
}

// Prepare prepare the runner's directory if in update mode by deleting it
// (including its content) and recreating it.
func (r *Runner) Prepare() error {
	if !r.Update {
		return nil
	}

	err := os.RemoveAll(r.Directory)
	if err != nil {
		return err
	}

	err = os.Mkdir(filepath.Join(r.Directory), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Test takes a path and some contents to check.
//
// An error is returned if the comparison failed or any error occurred during
// the process. Only nil means the compared contents where the same.
//
// In update mode, the file is created instead of being compared.
func (r *Runner) Test(path string, content []byte) error {
	path = filepath.Join(r.Directory, path)

	for _, filter := range r.Filters {
		content = filter(content)
	}

	if r.Update {
		// Make sure the necessary directories exist.
		dir, _ := filepath.Split(path)

		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("gold: could not create directory: %s", err)
		}

		err = ioutil.WriteFile(path, content, os.ModePerm)
		if err != nil {
			return fmt.Errorf("gold: could not write file: %s", err)
		}
	} else {
		// Compare the file with the given content.
		file, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("gold: could not read file: %s", err)
		}

		if !bytes.Equal(file, content) {
			return ComparisonError{}
		}
	}

	return nil
}
