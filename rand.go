package utils

import (
	cryptoRand "crypto/rand"
	"errors"
	"github.com/go-resty/resty/v2"
	"io"
	"math/rand"
	"regexp"
	"strconv"
	"unicode"
	"unsafe"
)

// MathRand is the implementation of "math/rand".Reader.
// Someone in the golang team has decided that they just know better.
var MathRand = mathRandImpl{}

type mathRandImpl struct {
}

func (m mathRandImpl) Read(dst []byte) (int, error) {
	return rand.Read(dst)
}

func RandomPassword(reader io.Reader, length int, customCharsets ...[]rune) (string, error) {
	var sources = [][]rune{Numbers, LettersLowercase, LettersUppercase, SpecialChars}
	if customCharsets != nil {
		sources = customCharsets
	}
	var res = make([]rune, length)
	for i := 0; i < length; i++ {
		source, e := RandomChoice(reader, sources)
		if e != nil {
			return "", e
		}
		res[i], e = RandomChoice(reader, source)
		if e != nil {
			return "", e
		}
	}
	return string(res), nil
}

func RandStringMust(reader io.Reader, length int) (res string) {
	res, e := RandString(reader, length)
	if e != nil {
		panic(e)
	}
	return res
}

func RandomChoiceMust[T any](reader io.Reader, data []T) T {
	v, e := RandomChoice(reader, data)
	if e != nil {
		panic(e)
	}
	return v
}

// RandString generates random string of characters [0-9a-Z]
func RandString(reader io.Reader, length int) (res string, e error) {
	// 10 + 26 + 26 = 62
	var getN = func(r uint8) (v uint8) {
		if r >= 10 {
			r += 7
			if r >= 43 {
				r += 6
			}
		}
		return r + '0'
	}

	var buf = make([]byte, length)
	if _, e = reader.Read(buf); e == nil {
		for i, v := range buf {
			buf[i] = getN(v % 62)
		}
	}
	return string(buf), e
}

func RandomChoice[T any](reader io.Reader, src []T) (T, error) {
	var v T
	var r, e = RandValue[uint](reader)
	if e == nil {
		v = src[r%uint(len(src))]
	}
	return v, e
}

func RandomChoiceMulti[T any](reader io.Reader, src []T, resLength int) ([]T, error) {
	var v = make([]T, resLength)
	var r, e = RandArray[uint](reader, make([]uint, resLength))
	if e == nil {
		for i, rv := range r {
			v[i] = src[rv%uint(len(src))]
		}
	}
	return v, e
}

func RandArray[T any](reader io.Reader, arr []T) ([]T, error) {
	var v T
	var tmp2 = unsafe.Slice(&arr[0], len(arr)*int(unsafe.Sizeof(v)))
	var tmp = *(*[]byte)(unsafe.Pointer(&tmp2))

	_, e := reader.Read(tmp)

	return arr, e
}

func RandValue[T any](reader io.Reader) (T, error) {
	var v T
	var tmp2 = unsafe.Slice(&v, int(unsafe.Sizeof(v)))
	var tmp = *(*[]byte)(unsafe.Pointer(&tmp2))

	_, e := reader.Read(tmp)

	return v, e
}

func RandIntN(reader io.Reader, n int64) (value int64, e error) {
	value, e = RandValue[int64](reader)
	if value < 0 {
		value = -value
	}
	value = value % n
	return
}

/*
	Data gens
*/

var steamActualPersonaNameRe = regexp.MustCompile(`(?s)<span class="actual_persona_name">(.*?)</span>`)
var randomUsernameSes = resty.New()

// RandomUsername retrieves a random username from steam persona name (not using steam's username for that because too little people have it).
// According to my tests, steam doesn't mind us sending tons of requests per second, so no need to limit the rps here.
// Tested on maxAttempts = 10, with 100_000 runs none has failed. Average requests made per username = 1.52.
func RandomUsername(maxAttempts int) (username string, attemptsUsed int, e error) {
	for attemptsUsed = 1; attemptsUsed <= maxAttempts; attemptsUsed++ {
		// Generating a random steam id.
		// Since we're generating from a pretty big range of possible ids, it's conceivable that an invalid ID may be produced on rare occasions.
		var randPart int64
		randPart, e = RandIntN(cryptoRand.Reader, 400_000_000)
		if e != nil {
			continue
		}
		var id = 76561197960265728 + 50_000_000 + randPart

		var r *resty.Response
		r, e = randomUsernameSes.R().Get("https://steamcommunity.com/profiles/" + strconv.Itoa(int(id)))

		//var randPart int64
		//randPart, e = RandIntN(cryptoRand.Reader, 1_062_000_000)
		//if e != nil {
		//	continue
		//}
		//
		//var r *resty.Response
		//r, e = randomUsernameSes.R().Get(fmt.Sprintf("https://steamcommunity.com/profiles/[U:1:%v]", randPart))
		if e == nil {
			// Not, that steam returns 200 OK even if the profile doesn't exist
			if !r.IsSuccess() {
				e = errors.New(r.Status())
				continue
			}
			found := steamActualPersonaNameRe.FindStringSubmatch(r.String())
			if len(found) == 0 {
				e = errors.New("regexp entries not found")
				continue
			}

			// Replacing bad chars
			var usernameRunes = []rune(found[1])
			var onlyNumbers = true
			for i, ch := range usernameRunes {
				if !(unicode.Is(unicode.Latin, ch) || unicode.IsDigit(ch) || ch == ' ' || ch == '.' || ch == '_') {
					e = errors.New("symbol not allowed")
					break
				}
				onlyNumbers = onlyNumbers && unicode.IsDigit(ch)

				if ch == ' ' {
					ch, e = RandomChoice(cryptoRand.Reader, []rune{'_', '.'})
					if e != nil {
						break
					}
					usernameRunes[i] = ch
				}
			}
			if e != nil {
				continue
			}
			// Skipping only-numeric usernames
			if onlyNumbers {
				e = errors.New("username contains only numbers")
				continue
			}
			// Skipping too short usernames
			if len(usernameRunes) < 6 {
				e = errors.New("too short")
				continue
			}
			// Skipping with bad prefixes
			if usernameRunes[0] == '.' || usernameRunes[0] == '_' {
				e = errors.New("first symbol is not allowed")
				continue
			}

			// Username satisfies all requirements
			username = string(usernameRunes)
			break
		}
	}
	return
}
