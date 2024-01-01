package config

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

type Decoder interface {
	Decode(b []byte) ([]byte, error)
}

var _ Decoder = DecoderFunc(nil)

type DecoderFunc func([]byte) ([]byte, error)

func (fn DecoderFunc) Decode(b []byte) ([]byte, error) {
	return fn(b)
}

func decodeBase64(b []byte) ([]byte, error) {
	return decodeWithName("base64", func(b []byte) ([]byte, error) {
		if len(b) == 0 {
			return nil, nil
		}

		b, err := base64.StdEncoding.DecodeString(string(b))
		if err != nil {
			return nil, fmt.Errorf("decodeString: %w", err)
		}

		return b, nil
	})(b)
}

func decodeGZIP(b []byte) ([]byte, error) {
	return decodeWithName("gzip", func(b []byte) ([]byte, error) {
		if len(b) == 0 {
			return nil, nil
		}

		r, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, fmt.Errorf("newReader: %w", err)
		}
		defer r.Close()

		b, err = io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("readAll: %w", err)
		}

		return b, nil
	})(b)
}

func decodeWithName(name string, fn DecoderFunc) DecoderFunc {
	return func(b []byte) ([]byte, error) {
		b, err := fn.Decode(b)
		if err != nil {
			return nil, fmt.Errorf("%q decode: %w", name, err)
		}

		return b, nil
	}
}
