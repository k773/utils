package utils

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/SilverCory/golang_discord_rpc"
	"github.com/SpencerSharkey/gomc/query"
	"github.com/parnurzeal/gorequest"
	"github.com/syndtr/goleveldb/leveldb"
	//"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding/charmap"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	//"syscall"
	"time"
	//"unsafe"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36"
const serverAddr = "http://127.0.0.1:8973"

type ScUtils struct {
}

type K773Utils struct {
}

//type Dialog struct {
//	DllFilePath string
//	DllObject   *syscall.LazyDLL
//}

type ServerStruct struct {
	ServerName  string
	LargeText   string
	LargeTextId string
}

type requestStruct struct {
	What string `json:"what"`
	Key  string `json:"key"`
	Data string `json:"data"`
}

type registerJson struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
	Email    string `json:"email"`
}

type registerResponseJsonStruct struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type reputationJson struct {
	User       int    `json:"user"`
	Reputation int    `json:"reputation"`
	Token      string `json:"_token"`
}

type captchaResponseStruct struct {
	Status  int    `json:"status"`
	Request string `json:"request"`
}

type RSA struct {
}

func H2b(encoded string) []byte {
	decoded, _ := hex.DecodeString(encoded)
	return decoded
}

func B2h(text []byte) string {
	return hex.EncodeToString(text)
}

func (RSA) ExportKey(key rsa.PublicKey) []byte {
	bytes1 := x509.MarshalPKCS1PublicKey(&key)
	var pemKey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bytes1,
	}
	return pem.EncodeToMemory(pemKey)
}

func (RSA) ImportKey(key string) rsa.PublicKey {
	data, _ := pem.Decode(H2b(key))
	serverPubKey, err := x509.ParsePKCS1PublicKey(data.Bytes)
	H(err)
	return *serverPubKey
}

func (RSA) EncryptRsa(key rsa.PublicKey, message []byte) []byte {
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key, message, []byte(""))
	H(err)
	return encrypted
}

func (RSA) DecryptRsa(key rsa.PrivateKey, message []byte) []byte {
	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, &key, message, []byte(""))
	H(err)
	return decrypted
}

func (RSA) SignRsa(key rsa.PrivateKey, data []byte) []byte {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	res, err := rsa.SignPSS(rand.Reader, &key, crypto.SHA256, Sha256BtB(data), &opts)
	H(err)
	return res
}

func (RSA) VerifySign(pubKey rsa.PublicKey, data, sign []byte) {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	err := rsa.VerifyPSS(&pubKey, crypto.SHA256, Sha256BtB(data), sign, &opts)
	if err != nil {
		panic(err)
	}
}

func H(err error) {
	if err != nil {
		panic(err)
	}
}

func EncryptBtB(strkey string, text []byte) []byte {
	//fmt.Println(string(text))
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

	//fmt.Println(len(rr))
	return []byte(base64.StdEncoding.EncodeToString(ciphertext))
}

func Rev(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func DecryptBtB(strkey string, bytes []byte) []byte {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString(strkey)
	ciphertext, err := base64.StdEncoding.DecodeString(string(bytes))
	if err != nil {
		fmt.Println(string(bytes))
		log.Println(err)
		return []byte{}
	}

	block, err := aes.NewCipher(key)
	//fmt.Println(len(bytes))
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		log.Println("ciphertext too short: " + strconv.Itoa(len(bytes)))
		return []byte{}
		//panic("ciphertext too short: " + strconv.Itoa(len(bytes)))
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}

//func EncryptBtB(password string, text []byte) []byte {
//	return []byte(aes256.Encrypt(string(text), password))
//}
//
//func DecryptBtB(password string, text []byte) []byte {
//	return []byte(aes256.Decrypt(string(text), password))
//}

func EncryptStH(strkey string, str string) string {
	return hex.EncodeToString(EncryptBtB(strkey, []byte(str)))
}

func DecryptHtS(strkey string, hexStr string) string {
	ciphertext, _ := hex.DecodeString(hexStr)
	return string(DecryptBtB(strkey, ciphertext))
}

func Sha256StH(text string) string {
	return hex.EncodeToString(Sha256BtB([]byte(text)))
}

func Sha256BtB(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

func Sha256File(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func ReadFile(path string) string {
	text, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}
	return string(text)
}

func ReadFileBytes(path string) []byte {
	text, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}
	return text
}

func ReadFiles(paths []string) []string {
	var data []string
	for _, path := range paths {
		data = append(data, ReadFile(path))
	}
	return data
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

func DbHas(db *leveldb.DB, key string) bool {
	has, _ := db.Has([]byte(key), nil)
	return has
}

func DbDelete(db *leveldb.DB, key string) {
	db.Delete([]byte(key), nil)
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
		fmt.Println(len(args))
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

func Md5(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Md5b(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

func DecodeWindows1251(ba []uint8) []uint8 {
	dec := charmap.Windows1251.NewDecoder()
	out, _ := dec.Bytes(ba)
	return out
}

func capSolve(ses *gorequest.SuperAgent, url, action, key, apikey string) (string, string) {
	gurl := "https://rucaptcha.com/in.php?key=%s&method=userrecaptcha&version=v3&action=%s&min_score=0.9&" +
		"googlekey=%s&pageurl=%s&json=1"
	gurl = fmt.Sprintf(gurl, apikey, action, key, url)
	_, response, _ := ses.Get(gurl).EndBytes()
	var capchaResponse1 captchaResponseStruct
	_ = json.Unmarshal(response, &capchaResponse1)
	if capchaResponse1.Status != 1 {
		os.Exit(-10)
	}
	var capchaResponse2 captchaResponseStruct
	for capchaResponse2.Request == "" || capchaResponse2.Request == "CAPCHA_NOT_READY" {
		_, capchaResponse2B, _ := ses.Get(fmt.Sprintf("https://rucaptcha.com/res.php?key=%s&action=get&taskinfo=0&json=1&id=%s", apikey, capchaResponse1.Request)).EndBytes()
		_ = json.Unmarshal(capchaResponse2B, &capchaResponse2)
		time.Sleep(2000 * time.Millisecond)
	}
	return capchaResponse2.Request, capchaResponse1.Request
}

func capReport(ses *gorequest.SuperAgent, good bool, apikey, capid string) {
	var action string
	if good {
		action = "reportgood"
	} else {
		action = "reportbad"
	}
	_, aga, _ := ses.Get(fmt.Sprintf("https://rucaptcha.com/res.php?key=%s&action=%s&id=%s", apikey, action, capid)).End()
	fmt.Println(aga)
}

//func MachineID() (string, error) {
//	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.QUERY_VALUE|registry.WOW64_64KEY)
//	if err != nil {
//		return "", err
//	}
//	defer k.Close()
//
//	s, _, err := k.GetStringValue("MachineGuid")
//	if err != nil {
//		return "", err
//	}
//	return s, nil
//}

func (ScUtils) RegisterAccount(ses *gorequest.SuperAgent, ruCaptchaKey string) (string, string, string) {
	// Returns login string, password string, csrf string
	const siteKey = "6LeedroUAAAAAK2RUkaNLVBYraeQXNVHX45O227A"
	const pageUrl = "https://streamcraft.net/api/auth/register"
	const regexToken = `<meta name="csrf-token" content="(?P<token>.*)">`
	const action = "register"

register:
	resp, page, _ := ses.Get(pageUrl).End()

	csrf := FindGroup(page, regexToken)[1]

	xsrf := (*http.Response)(resp).Cookies()[1].Value
	//os.Exit(111)
	email := randomdata.Email()
	name := randomdata.FirstName(randomdata.Number(1, 2))
	length := 8
	if len(name) < 8 {
		length = len(name)
	}
	login := name[:length] + randomdata.RandStringRunes(4) + strconv.Itoa(randomdata.Number(1980, 2017))
	password := login + login

	//RuCaptcha
	//capSolve(ses, pageUrl, action, siteKey, ruCaptchaKey)

	var Json registerJson
	var capid string
	Json.Login = login
	Json.Password = password
	Json.Captcha, capid = capSolve(ses, pageUrl, action, siteKey, ruCaptchaKey)
	Json.Email = email
	Jsonb, _ := json.Marshal(Json)

	_, data, _ := ses.Post(pageUrl).Set("x-csrf-token", csrf).Set("x-xsrf-token", xsrf).Send(string(Jsonb)).EndBytes()
	var registerResponseJson registerResponseJsonStruct
	_ = json.Unmarshal(data, &registerResponseJson)
	if registerResponseJson.Success == true {
		capReport(ses, true, ruCaptchaKey, capid)
	} else {
		capReport(ses, true, ruCaptchaKey, capid)
		goto register
	}

	fmt.Println(string(data))
	return login, password, csrf
}

func (ScUtils) SetReputation(ses *gorequest.SuperAgent, csrf string, userId int, count int) {
	//Set user reputation
	const pageUrl = "https://streamcraft.net/forum/user/reputation"

	var Json reputationJson
	Json.User = userId
	Json.Reputation = count
	Json.Token = csrf
	JsonB, _ := json.Marshal(Json)

	ses.Post(pageUrl).Send(string(JsonB)).End()
}

func (ScUtils) GetUserId(ses *gorequest.SuperAgent, nickname string) int {
	//Get user id
	const pageUrl = "https://streamcraft.net/user/"
	const regexUserId = `<i class="fa fa-thumbs-down cursor-pointer" onclick="App\.sendRequest\('/forum/user/reputation', {user: (?P<id>.*), reputation: -1}\);"></i>`

	_, page, _ := ses.Get(pageUrl + nickname).End()
	id, _ := strconv.Atoi(FindGroup(page, regexUserId)[1])
	return id
}

func (ScUtils) ThreadsIdsParse(ses *gorequest.SuperAgent) []string {
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

func (K773Utils) H2s(hex string) string {
	src := []byte(hex)
	n, _ := decode(src, src)
	return string(src[:n])
}

func (K773Utils) S2h(text string) string {
	src := []byte(text)
	dst := make([]byte, encodedLen(len(src)))
	encode(dst, src)
	return string(dst)
}

func (K773Utils) DecryptAes(data, key string) string {
	var request requestStruct
	request.What, request.Data, request.Key = "decrypt", data, key
	marshalled, _ := json.Marshal(request)
	_, response, _ := gorequest.New().Post(serverAddr).
		Send(string(marshalled)).
		End()
	return response
}

func (K773Utils) EncryptAes(data, key string) string {
	var request requestStruct
	request.What, request.Data, request.Key = "encrypt", data, key
	marshalled, _ := json.Marshal(request)
	_, response, _ := gorequest.New().Post(serverAddr).
		Send(string(marshalled)).
		End()
	return response
}

func SetDiscordStatus(server ServerStruct, nickname string) {
start:
	win := discordrpc.NewRPCConnection("496419141201297413")
	err := win.Open()
	if err != nil {
		//fmt.Println(err)
		time.Sleep(5 * time.Second)
		goto start
	}

	_, _ = win.Read()
	//fmt.Println(err)
	//fmt.Println(str)

	//time.Sleep(time.Second * 3)
	stamp := time.Now().Unix()

	for {
		//fmt.Println(os.Getpid())
		presence := &discordrpc.CommandRichPresenceMessage{
			CommandMessage: discordrpc.CommandMessage{Command: "SET_ACTIVITY"},
			Args: &discordrpc.RichPresenceMessageArgs{
				Pid: os.Getpid(),
				Activity: &discordrpc.Activity{
					Details:    "Играет на " + server.ServerName,
					State:      "С ником " + nickname,
					Instance:   false,
					TimeStamps: &discordrpc.TimeStamps{StartTimestamp: stamp},
					Assets: &discordrpc.Assets{
						LargeText:    server.LargeText,
						LargeImageID: server.LargeTextId,
						SmallText:    "StreamCraft.Net",
						SmallImageID: "discord",
					},
				},
			},
		}

		presence.SetNonce()
		data, err := json.Marshal(presence)

		if err != nil {
			fmt.Println(err)
			continue
		}

		err = win.Write(string(data))
		if err != nil {
			fmt.Println(err)
			continue
		}

		//str, err := win.Read()
		//fmt.Println(err)
		//fmt.Println(str)
		//
		//fmt.Println("---\nDone?")
		time.Sleep(time.Second * 5)
	}
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

//
//func U(something string) uintptr {
//	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(something)))
//}

//func (dialog *Dialog) CallDll() {
//	if dialog.DllFilePath != "" && dialog.DllObject == nil {
//		dialog.DllObject = syscall.NewLazyDLL(dialog.DllFilePath)
//	}
//}
//
//func (dialog *Dialog) YesNo(title, label, yesButtonText, noButtonText string) (bool, bool) {
//	proc := dialog.DllObject.NewProc("YesNo")
//	code, _, _ := proc.Call(U(title), U(label), U(yesButtonText), U(noButtonText))
//	switch code {
//	case 0:
//		return false, false
//	case 100:
//		return false, true
//	case 101:
//		return true, true
//	}
//	return false, false
//}
//
//func (dialog *Dialog) TextInput(title, label, buttonText string) (uintptr, uintptr, error) {
//	proc := dialog.DllObject.NewProc("TextInputDialog")
//	return proc.Call(U(title), U(label), U(buttonText))
//}
