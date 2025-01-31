package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/urfave/cli/v2"
)

func action(c *cli.Context) error {
	if c.Args().Len() != 2 {
		return errors.New("input and output arguments are mandatory")
	}
	input := c.Args().Get(0)
	output := c.Args().Get(1)
	pack := c.String("package")
	name := c.String("name")
	ncols := c.Int("ncols")

	data, err := loadInputFile(input)
	if err != nil {
		return err
	}
	content, err := produce(input, pack, name, ncols, data)
	if err != nil {
		return err
	}

	err = writeOutputFile(output, content)
	if err != nil {
		return err
	}

	return nil
}

func loadInputFile(input string) ([]byte, error) {
	data, err := os.ReadFile(input)
	if err != nil {
		return nil, fmt.Errorf("unable to read file '%s': %w", input, err)
	}
	return data, nil
}

func produce(input, pack, name string, ncols int, data []byte) ([]byte, error) {
	tab := "\t"

	var content strings.Builder
	for i, d := range data {
		if (i % ncols) == 0 {
			content.WriteString(tab + tab)
		}
		content.WriteString(`0x`)
		content.WriteString(hex.EncodeToString([]byte{d}))
		content.WriteString(",")
		if i < (len(data) - 1) {
			if (i % ncols) == (ncols - 1) {
				content.WriteString("\n")
			} else {
				content.WriteString(" ")
			}
		}
	}

	temp, err := template.New("filetemplate").Parse(filetemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to produce: %w", err)
	}

	var buff bytes.Buffer
	err = temp.Execute(&buff, struct {
		Year     string
		Version  string
		DateTime string
		Package  string
		FileName string
		Name     string
		Content  string
	}{
		Year:     strconv.FormatInt(int64(time.Now().Year()), 10),
		Version:  version,
		DateTime: time.Now().Format("Mon, 02 Jan 2006 15:04:05"),
		Package:  pack,
		FileName: filepath.Base(input),
		Name:     name,
		Content:  content.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to produce: %w", err)
	}

	return buff.Bytes(), nil
}

func writeOutputFile(output string, content []byte) error {
	file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("unable to open file '%s': %w", output, err)
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return fmt.Errorf("failed to save file '%s': %w", output, err)
	}
	return nil
}
