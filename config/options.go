package config

type options struct {
	filenames []string
}

type Option interface {
	apply(options *options)
}

var _ Option = optionFunc(nil)

type optionFunc func(options *options)

func (fn optionFunc) apply(options *options) {
	fn(options)
}

func WithFilenames(filenames ...string) Option {
	return optionFunc(func(options *options) {
		options.filenames = append(options.filenames, filenames...)
	})
}

type GetOption interface {
	apply(options *getOptions)
}

type getOptions struct {
	decoders []Decoder
}

var _ GetOption = getOptionFunc(nil)

type getOptionFunc func(options *getOptions)

func (fn getOptionFunc) apply(options *getOptions) {
	fn(options)
}

func WithBase64() GetOption {
	return getOptionFunc(func(options *getOptions) {
		options.decoders = append(options.decoders, DecoderFunc(decodeBase64))
	})
}

func WithGZIP() GetOption {
	return getOptionFunc(func(options *getOptions) {
		options.decoders = append(options.decoders, DecoderFunc(decodeBase64), DecoderFunc(decodeGZIP))
	})
}

func WithDecoder(decoder Decoder) GetOption {
	return getOptionFunc(func(options *getOptions) {
		options.decoders = append(options.decoders, decoder)
	})
}

func newOptions(opts ...Option) *options {
	options := &options{}
	for _, opt := range opts {
		opt.apply(options)
	}

	return options
}
