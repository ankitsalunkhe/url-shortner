package shortner

const (
	base        = 62
	base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Shortner interface {
	Generate(num int) string
}

var _ Shortner = (*Base62)(nil)

type Base62 struct {
}

func (b Base62) Generate(num int) string {
	return base10ToBase62(num)
}

func base10ToBase62(in int) string {
	result := ""
	for in > 0 {
		remainder := in % base
		result = string(base62Chars[remainder]) + result
		in = in / base
	}
	return string(result)
}
