package cli

type Decoder interface {
	Decode(s string) error
}

type Encoder interface {
	Encode() string
}
