package utils

import (
	"errors"
	"fmt"
	"os"
)

func ReadFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			file, err = os.Create(path)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("open file error: %w", err)
		}
	}
	return file, err
}
