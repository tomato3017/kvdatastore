package io

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/tomato3017/kvdatastore/util/errcodes"
	"os"
)

func TouchFile(fs afero.Fs, filename string) error {
	if _, err := fs.Stat(filename); os.IsNotExist(err) {
		f, err := fs.Create(filename)
		if err != nil {
			return fmt.Errorf("unable to touch file %s. Err: %w", filename, err)
		}

		if err := f.Close(); err != nil {
			return err
		}

		return nil
	} else if err != nil && !os.IsExist(err) {
		return err
	}

	return errcodes.ErrFileExists{Filename: filename}
}
