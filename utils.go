package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/SpencerSharkey/gomc/query"
	"github.com/koteezy/ruCaptcha"
	"github.com/parnurzeal/gorequest"
	"github.com/syndtr/goleveldb/leveldb"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36"

type scUtils struct {
}

type k773Utils struct {
}

type registerJson struct {
	Login           string `json:"login"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirmation"`
	Email           string `json:"email"`
	Captcha         string `json:"captcha"`
	Token           string `json:"_token"`
}

type reputationJson struct {
	User       int    `json:"user"`
	Reputation int    `json:"reputation"`
	Token      string `json:"_token"`
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

func FindGroup(text string, reg string) []string {
	regex, _ := regexp.Compile(reg)
	return regex.FindStringSubmatch(text)
}

func FindGroups(text string, reg string) []string {
	regex, _ := regexp.Compile(reg)
	temp := regex.FindAllStringSubmatch(text, -1)
	var temp4 []string
	for _, temp2 := range temp {
		temp4 = append(temp4, temp2[1])
	}
	return temp4
}

func FindAllGroups(text string, reg string) [][]string {
	regex, _ := regexp.Compile(reg)
	temp := regex.FindAllStringSubmatch(text, -1)
	var temp4 [][]string
	for _, temp2 := range temp {
		var temp6 []string
		for _, temp5 := range temp2 {
			if temp5 == temp2[0] {
				continue
			}
			temp6 = append(temp6, temp5)
		}
		temp4 = append(temp4, temp6)
	}
	return temp4
}

func GetServerPlayers(ip string) []string {
	req := query.NewRequest()
	_ = req.Connect(ip)
	response, _ := req.Full()

	if response == nil {
		///fmt.Println("Error", ip, err)
		return []string{}
	}

	var playersArray []string
	for _, player := range response.Players {
		//fmt.Println(i)
		playersArray = append(playersArray, player)
	}
	return playersArray
}

func DbGet(db *leveldb.DB, key string) []byte {
	val, err := db.Get([]byte(key), nil)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			return []byte{}
		}
		panic(err)
	}
	return val
}

func DbSet(db *leveldb.DB, key string, value interface{}) {
	_ = db.Delete([]byte(key), nil)
	_ = db.Put([]byte(key), value.([]byte), nil)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetRequest(ses *gorequest.SuperAgent, url string, args ...string) string {
	if len(args)%2 != 0 {
		return ""
	}

	for i, arg := range args {
		if i%2 == 0 {
			url += arg + "="
		} else {
			url += arg + "&"
		}
	}

	_, response, _ := ses.Get(url).End()
	return response
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (scUtils) RegisterAccount(ses *gorequest.SuperAgent, ruCaptchaKey string) (string, string, string) {
	// Returns login string, password string, csrf string
	const siteKey = "6LcUwBgUAAAAAAyJnKWJvhBNNzItS7DlHoARaQbG"
	const pageUrl = "https://streamcraft.net/register"
	const regexToken = `<meta name="csrf-token" content="(?P<token>.*)">`
register:
	_, page, _ := ses.Get(pageUrl).End()

	csrf := FindGroup(page, regexToken)[1]
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

func (scUtils) SetReputation(ses *gorequest.SuperAgent, csrf string, userId int, count int) {
	//Set user reputation
	const pageUrl = "https://streamcraft.net/forum/user/reputation"

	var Json reputationJson
	Json.User = userId
	Json.Reputation = count
	Json.Token = csrf
	JsonB, _ := json.Marshal(Json)

	ses.Post(pageUrl).Send(string(JsonB)).End()
}

func (scUtils) GetUserId(ses *gorequest.SuperAgent, nickname string) int {
	//Get user id
	const pageUrl = "https://streamcraft.net/user/"
	const regexUserId = `<i class="fa fa-thumbs-down cursor-pointer" onclick="App\.sendRequest\('/forum/user/reputation', {user: (?P<id>.*), reputation: -1}\);"></i>`

	_, page, _ := ses.Get(pageUrl + nickname).End()
	id, _ := strconv.Atoi(FindGroup(page, regexUserId)[1])
	return id
}

func (scUtils) ThreadsIdsParse(ses *gorequest.SuperAgent) []string {
	const regexThreads = `<a href="/forum/category/(?P<id>.*)"><i class="fa fa-level-down">`
	const regexThreadsIds = `<a class="btn btn-primary btn-shadow float-right" href="/forum/discussion/create/(?P<id>.*)" role="button">`
	const ForumUrl = "https://streamcraft.net/forum/"
	const CategoryUrl = "https://streamcraft.net/forum/category/"

	_, text, _ := ses.Get(ForumUrl).End()
	temp := FindGroups(text, regexThreads)
	var threadsIds []string
	for _, thread := range temp {
		_, temp2, _ := ses.Get(CategoryUrl + string(thread)).End()
		temp3 := FindGroup(temp2, regexThreadsIds)
		if len(temp3) < 2 {
			continue
		}
		threadsIds = append(threadsIds, temp3[1])
	}
	return threadsIds
}

func (k773Utils) H2s(hex string) string {
	src := []byte(hex)
	n, _ := decode(src, src)
	return string(src[:n])
}

func (k773Utils) S2h(text string) string {
	src := []byte(text)
	dst := make([]byte, encodedLen(len(src)))
	encode(dst, src)
	return string(dst)
}

func encodedLen(n int) int { return n * 2 }

func encode(dst, src []byte) int {
	for i, v := range src {
		v += 4
		dst[i*2] = "0123456789abcdef"[v>>4]
		dst[i*2+1] = "0123456789abcdef"[v&0x0f]
	}

	return len(src) * 2
}

func decode(dst, src []byte) (int, error) {
	var i int
	for i = 0; i < len(src)/2; i++ {
		a, ok := fromHexChar(src[i*2])

		if !ok {
		}
		b, ok := fromHexChar(src[i*2+1])

		if !ok {
		}
		dst[i] = ((a << 4) | b) - 4
	}
	if len(src)%2 == 1 {
		// Check for invalid char before reporting bad length,
		// since the invalid char (if present) is an earlier problem.
		if _, ok := fromHexChar(src[i*2]); !ok {
		}
	}
	return i, nil
}

func fromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}

func decryptAes128Ecb(data, key string) string {
	_, response, _ := gorequest.New().Post("http://212.237.2.10/gethashes.php?what=decrypt").Send(fmt.Sprintf("key=%s&data=%s", key, data)).End()
	return response
}
