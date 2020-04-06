package internal

import "errors"

type Parser struct {
	MaxNewLine int32
}

func New(maxNewLine int32) (*Parser, error) {
	if maxNewLine < 0 {
		return nil, errors.New("maxNewLine must be a positive integer")
	}
	p := &Parser{MaxNewLine: maxNewLine}
	return p, nil
}
