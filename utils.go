package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/koteezy/ruCaptcha"
	"github.com/parnurzeal/gorequest"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

type scUtils struct {
}

type registerJson struct {
	Login           string `json:"login"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirmation"`
	Email           string `json:"email"`
	Captcha         string `json:"captcha"`
	Token           string `json:"_token"`
}

func EncryptBtB(strkey string, text []byte) []byte {
	key, _ := hex.DecodeString(strkey)
	plaintext := text

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext
}

func DecryptBtB(strkey string, bytes []byte) []byte {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString(strkey)
	ciphertext := bytes

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}

func EncryptStH(strkey string, str string) string {
	return hex.EncodeToString(EncryptBtB(strkey, []byte(str)))
}

func DecryptHtS(strkey string, hexStr string) string {
	ciphertext, _ := hex.DecodeString(hexStr)
	return string(DecryptBtB(strkey, ciphertext))
}

func Sha256StH(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func readFile(path string) string {
	text, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}
	return string(text)
}

func findGroup(text string, reg string) []string {
	regex, _ := regexp.Compile(reg)
	return regex.FindStringSubmatch(text)
}

func findGroups(text string, reg string) []string {
	regex, _ := regexp.Compile(reg)
	temp := regex.FindAllStringSubmatch(text, -1)
	var temp4 []string
	for _, temp2 := range temp {
		temp4 = append(temp4, temp2[1])
	}
	return temp4
}

func (scUtils) registerAccount(ses *gorequest.SuperAgent, ruCaptchaKey string) (string, string, string) {
	// Returns login string, password string, csrf string
	siteKey := "6LcUwBgUAAAAAAyJnKWJvhBNNzItS7DlHoARaQbG"
	pageUrl := "https://streamcraft.net/register"
	regexToken := `<meta name="csrf-token" content="(?P<token>.*)">`
register:
	_, page, _ := ses.Get(pageUrl).End()

	csrf := findGroup(page, regexToken)[1]
	email := randomdata.Email()
	name := randomdata.FirstName(randomdata.Number(1, 2))
	length := 8
	if len(name) < 8 {
		length = len(name)
	}
	login := name[:length] + randomdata.RandStringRunes(4) + strconv.Itoa(randomdata.Number(1980, 2017))
	password := login + login

	//RuCaptcha
	re := rucaptcha.New(ruCaptchaKey)
	captcha, err := re.ReCaptcha(pageUrl, siteKey)
	if err != nil {
		goto register
		//panic(err)
	}

	var Json = registerJson{}
	Json.Login = login
	Json.Password = password
	Json.PasswordConfirm = password
	Json.Email = email
	Json.Captcha = captcha
	Json.Token = csrf
	Jsonb, _ := json.Marshal(Json)

	_, data, _ := ses.Post(pageUrl).Send(string(Jsonb)).End()
	if data != `{"success":true,"redirect":"https:\/\/streamcraft.net\/home"}` {
		fmt.Println("Error while solving captcha! Trying again...", data)
		goto register
	}
	return login, password, csrf
}
