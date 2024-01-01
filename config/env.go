package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

var (
	ErrNotFound = errors.New("not found")
	ErrEmptyKey = errors.New("empty key")
)

type Getter interface {
	Get(key string, opts ...GetOption) (string, error)
	MustGet(key string, opts ...GetOption) string
	GetFloat(key string, opts ...GetOption) (float64, error)
	MustGetFloat(key string, opts ...GetOption) float64
	Have(keys ...string) error
	MustHave(keys ...string)
}

type env struct {
	vars map[string]string
}

func NewFromEnv(opts ...Option) Getter {
	options := &options{}
	for _, opt := range opts {
		opt.apply(options)
	}

	_ = godotenv.Load(options.filenames...)

	return &env{
		vars: make(map[string]string),
	}
}

func (c *env) Have(keys ...string) error {
	missing := make([]string, 0)

	for _, key := range keys {
		if key == "" {
			continue
		}

		if _, ok := os.LookupEnv(key); !ok {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("variables %s %w", strings.Join(missing, ", "), ErrNotFound)
	}

	return nil
}

func (c *env) MustHave(keys ...string) {
	if err := c.Have(keys...); err != nil {
		panic(err)
	}
}

func (c *env) Get(key string, opts ...GetOption) (string, error) {
	value, err := c.getRaw(key)
	if err != nil {
		return "", err
	}

	options := getOptions{}

	for _, opt := range opts {
		opt.apply(&options)
	}

	if value, err = c.decode(value, options.decoders); err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}

	c.vars[key] = value

	return c.vars[key], nil
}

func (c *env) MustGet(key string, opts ...GetOption) string {
	s, err := c.Get(key, opts...)
	if err != nil {
		panic(fmt.Errorf("env: %w", err))
	}

	return s
}

func (c *env) GetFloat(key string, opts ...GetOption) (float64, error) {
	s, err := c.Get(key, opts...)
	if err != nil {
		return 0, err
	}

	if s == "" {
		return 0, nil
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("atoi: %w", err)
	}

	return v, nil
}

func (c *env) MustGetFloat(key string, opts ...GetOption) float64 {
	v, err := c.GetFloat(key, opts...)
	if err != nil {
		panic(fmt.Errorf("env: %w", err))
	}

	return v
}

func (c *env) decode(s string, decoders []Decoder) (string, error) {
	b := []byte(s)

	for _, decoder := range decoders {
		var err error

		if b, err = decoder.Decode(b); err != nil {
			return "", err
		}
	}

	return string(b), nil
}

func (c *env) getRaw(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}

	if v, ok := c.vars[key]; ok {
		return v, nil
	}

	s, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("%q: %w", key, ErrNotFound)
	}

	return s, nil
}
