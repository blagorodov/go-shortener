package repository

import (
	"encoding/json"
	"os"
)

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func newProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	p := &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}
	return p, nil
}

func (p *producer) writeItem(item *ShortenURL) error {
	err := p.encoder.Encode(item)
	if err != nil {
		return err
	}
	return nil
}

func (p *producer) close() error {
	return p.file.Close()
}
