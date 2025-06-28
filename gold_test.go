package gold_test

import (
	"io/ioutil"
	"testing"

	. "github.com/mfcochauxlaberge/gold"

	"github.com/stretchr/testify/assert"
)

func TestNewRunner(t *testing.T) {
	assert := assert.New(t)

	runner := Runner{
		Directory: "some/path",
	}
	assert.Equal(false, runner.Update)
	assert.Equal("some/path", runner.Directory)
	assert.Len(runner.Filters, 0)
}

func TestRunnerPrepare(t *testing.T) {
	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "")
	assert.NoError(err)

	// Put a file in the directory.
	_, err = ioutil.TempFile(dir, "")
	assert.NoError(err)

	// Make sure it exists.
	files, err := ioutil.ReadDir(dir)
	assert.NoError(err)
	assert.Equal(1, len(files))

	// Create runner and prepare it.
	runner := Runner{
		Directory: dir,
	}

	// Nothing happens here because
	// update mode is disabled.
	err = runner.Prepare()
	assert.NoError(err)

	files, err = ioutil.ReadDir(dir)
	assert.NoError(err)
	assert.Equal(1, len(files))

	// The directory gets deleted and
	// recreated in Prepare, so it
	// should be empty after the call.
	runner.Update = true

	err = runner.Prepare()
	assert.NoError(err)

	files, err = ioutil.ReadDir(dir)
	assert.NoError(err)
	assert.Equal(0, len(files))
}

func TestRunnerTest(t *testing.T) {
	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "")
	assert.NoError(err)

	// Create runner and prepare it.
	runner := Runner{
		Directory: dir,
		Update:    true,
	}

	err = runner.Prepare()
	assert.NoError(err)

	// With update mode on, the following calls will
	// create files.
	err = runner.Test("test1", []byte("This is a test."))
	assert.NoError(err)

	err = runner.Test("test2", []byte("This is another test."))
	assert.NoError(err)

	// The update mode is turned off. The following should
	// not return an error since the files have already
	// been created.
	runner.Update = false

	err = runner.Test("test1", []byte("This is a test."))
	assert.NoError(err)

	err = runner.Test("test2", []byte("This is another test."))
	assert.NoError(err)
}

func TestRunnerWithFilter(t *testing.T) {
	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "")
	assert.NoError(err)

	// Create runner and prepare it.
	runner := Runner{
		Directory: dir,
		Update:    true,
	}

	err = runner.Prepare()
	assert.NoError(err)

	runner.Filters = append(runner.Filters, func(src []byte) []byte {
		return []byte("The output was altered.")
	})

	// Create file.
	err = runner.Test("test1", []byte("The output was not altered."))
	assert.NoError(err)

	// Turn off update mode and check file.
	runner.Update = false

	err = runner.Test("test1", []byte("The output was altered."))
	assert.NoError(err)
}

func TestRunnerTestComparisonError(t *testing.T) {
	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "")
	assert.NoError(err)

	// Create runner and prepare it.
	runner := Runner{
		Directory: dir,
		Update:    true,
	}

	err = runner.Prepare()
	assert.NoError(err)

	// Create file.
	err = runner.Test("test1", []byte("This is a test."))
	assert.NoError(err)

	// Turn off update mode and check file.
	runner.Update = false

	err = runner.Test("test1", []byte("The content is different."))
	assert.Error(err)

	cErr, ok := err.(ComparisonError)
	assert.True(ok)
	assert.Equal("output and file are different", cErr.Error())
}
