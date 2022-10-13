package chirano

import "crypto/rand"

const randomStr = "910abcd4efgh8ijkl#6mno3pqr5st27"

//common Helper functions reusable in other apps
//any var of this type will have access to all the methods with receiver *Helper
type Helpers struct {
	MaxFileSize      int
	AllowedFileTypes []string
}

//gen random string
func (h *Helpers) RandomString(n int) string {

	s, r := make([]rune, n), []rune(randomStr)

	for i := range s {

		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))

		s[i] = r[x%y]
	}

	return string(s)
}
