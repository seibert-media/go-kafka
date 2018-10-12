package schema

import (
	"fmt"
	"io"
)

const headerLength = 5

func RemoveMagicHeader(reader io.Reader) error {
	length, err := reader.Read(make([]byte, headerLength))
	if err != nil {
		return err
	}
	if length != headerLength {
		return fmt.Errorf("read %d bytes from reader failed", headerLength)
	}
	return nil
}
