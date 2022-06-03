package googleAuthenticator

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"github.com/k773/utils"
	"strconv"
	"strings"
	"time"
)

const timeWindowLength = 30 * time.Second

// GeneratePasscode generates passcode to a given shared secret based on a given time. For verification better use VerifyPasscode.
func GeneratePasscode(sharedSecretBase32 string, at time.Time) (code string, e error) {
	sharedSecret, e := base32.StdEncoding.DecodeString(strings.ToUpper(sharedSecretBase32))
	if e == nil {
		var hash = hmac.New(sha1.New, sharedSecret)
		_ = binary.Write(hash, binary.BigEndian, at.Unix()/30)
		sum := hash.Sum(nil)
		offset := sum[19] & 0x0f // resetting all bits but the first four
		// code = uint32(hash[offset:offset+4]) /* removing most significant bit */ & 0x7fffffff % 1'000'000
		code = utils.PadLeft(strconv.FormatUint((uint64(binary.BigEndian.Uint32(sum[offset:offset+4]))&0x7fffffff)%1000000, 10), 6, '0')
	}
	return
}

// VerifyPasscode verifies passcode for the given shared secret based on the given time && time window.
// Time window is the number of seconds in which the passcode would be tried to match ([time-time window:time+time window]);
// 	time window may be 0 -> passcode will be verified only with the given time.
func VerifyPasscode(passcode, sharedSecretBase32 string, at time.Time, timeWindow time.Duration) (pass bool, e error) {
	var atDuration = time.Duration(at.Truncate(time.Second).UnixNano())

	var windowNBefore = int((atDuration - timeWindow) / timeWindowLength)
	var windowNCurrent = int(atDuration / timeWindowLength)
	var windowNAfter = int((atDuration + timeWindow) / timeWindowLength)

	var code string
	for i := windowNBefore; i <= windowNAfter && e == nil; i++ {
		windowTime := time.Unix(0, int64(atDuration+time.Duration(i-windowNCurrent)*timeWindowLength))
		if code, e = GeneratePasscode(sharedSecretBase32, windowTime); e == nil {
			fmt.Println("time:", windowTime, ", code:", code)
			if pass = pass || code == passcode; pass {
				break
			}
		}
	}
	return
}
