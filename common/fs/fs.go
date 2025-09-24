package fs

import (
	"fmt"
	"os"
)

func MustDir(path string) error {

	if s, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if !s.IsDir() {
		return fmt.Errorf("%s 不是目录", path)
	}

	return nil
}

func MustFile(path string) error {
	if s, err := os.Stat(path); err != nil {
		return err
	} else if s.IsDir() {
		return fmt.Errorf("%s 不是文件", path)
	}
	return nil
}
