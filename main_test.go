package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	expected = `
package myresources

var (
	// test generated resource from file "test.jpg".
	test = []byte{
		0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01,
	}
)
`
)

func TestRunMain(t *testing.T) {
	t.Parallel()

	progname := "goconvres"
	currentdir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	input := filepath.Join(currentdir, "test.jpg")
	output := filepath.Join(currentdir, "test.go")
	pack := "myresources"
	name := "test"
	data := []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01}

	err := os.WriteFile(input, data, 0644)
	assert.NoError(t, err)

	os.Args = []string{
		progname,
		"--help",
	}
	err = parseCLI()
	assert.NoError(t, err)

	os.Args = []string{
		progname,
		"--version",
	}
	err = parseCLI()
	assert.NoError(t, err)

	os.Args = []string{
		progname,
		"--name", name,
		"--package", pack,
		"--ncols", "12",
		input,
		output,
	}
	err = parseCLI()
	assert.NoError(t, err)

	res, err := os.ReadFile(output)
	assert.NoError(t, err)
	fmt.Printf("res = %v\n", string(res))
	assert.Contains(t, string(res), expected)

	os.Args = []string{
		progname,
		"--name", name,
		"--package", pack,
		"--ncols", "12",
		input,
	}
	err = parseCLI()
	assert.Error(t, err)

	os.Args = []string{
		progname,
		"--name", name,
		"--package", pack,
		"--ncols", "12",
		"xxxxxx",
		output,
	}
	err = parseCLI()
	assert.Error(t, err)
}

func TestFileSizeLimit(t *testing.T) {
	t.Parallel()

	progname := "goconvres"
	currentdir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// Create a file larger than the maxInputSize (10 MB)
	largeInput := filepath.Join(currentdir, "largefile.bin")
	largeData := make([]byte, maxInputSize+1) // One byte over the limit

	err := os.WriteFile(largeInput, largeData, 0644)
	assert.NoError(t, err)
	defer os.Remove(largeInput) // Clean up

	output := filepath.Join(currentdir, "largefile.go")
	os.Args = []string{
		progname,
		"--name", "test",
		"--package", "testpkg",
		largeInput,
		output,
	}

	err = parseCLI()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum allowed size")
}