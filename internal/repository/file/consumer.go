package file

import (
	"bufio"
	"encoding/json"
	"os"
)

type consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func newConsumer(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *consumer) readItem() (*ShortenURL, error) {
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	data := c.scanner.Bytes()
	item := ShortenURL{}

	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (c *consumer) close() error {
	return c.file.Close()
}
