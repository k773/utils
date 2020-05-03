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
	"database/sql"
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
	"log"
	"strings"

	//"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding/charmap"
	"io"
	"io/ioutil"
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

type MySqlUtils struct {
	Db *sql.DB
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

const LevelError = 0
const LevelInfo = 1
const LevelDebug = 2

type RSA struct {
}

type SliceTools struct {
}

type Logger struct {
	LogLevel   int
	LoggerName string
}

func SplitStringByCount(str string, maxCount int) []string {
	var ret []string
	for i := 0; true; i += maxCount {
		str = getSymbols(str, i, maxCount)
		ret = append(ret, str)
		if len(str) != maxCount {
			break
		}
	}
	return ret
}

func getSymbols(str string, startIndex, count int) string {
	var ret string
	str = str[startIndex:]
	if len(str) < count {
		count = len(str)
	}
	for i := 0; i < count; i++ {
		ret += string(str[i])
	}
	return ret
}

func H2b(encoded string) []byte {
	decoded, _ := hex.DecodeString(encoded)
	return decoded
}

func B2h(text []byte) string {
	return hex.EncodeToString(text)
}

func S2h(text string) string {
	return hex.EncodeToString([]byte(text))
}

func H2s(h string) string {
	if val, err := hex.DecodeString(h); err == nil {
		return string(val)
	}
	return ""
}

func GoLog(args ...interface{}) {
	var string_ string
	for _, element := range args {
		string_ += fmt.Sprintf("%v", element) + " "
	}
	timeString := time.Now().String()[:19]
	fmt.Println(fmt.Sprintf("[%v]", timeString)+":", string_) //aaa
}

func (logger Logger) Debug(data ...interface{}) {
	if logger.LogLevel >= LevelDebug {
		var string_ string
		for _, element := range data {
			string_ += fmt.Sprintf("%v", element) + " "
		}
		GoLog(fmt.Sprintf("[%v] [%v]:", logger.LoggerName, "DEBUG"), string_)
	}
}

func (logger Logger) Error(data ...interface{}) {
	if logger.LogLevel >= LevelError {
		var string_ string
		for _, element := range data {
			string_ += fmt.Sprintf("%v", element) + " "
		}
		GoLog(fmt.Sprintf("[%v] [%v]:", logger.LoggerName, "ERROR"), string_)
	}
}

func (logger Logger) Info(data ...interface{}) {
	if logger.LogLevel >= LevelInfo {
		var string_ string //
		for _, element := range data {
			string_ += fmt.Sprintf("%v", element) + " "
		}
		GoLog(fmt.Sprintf("[%v] [%v]:", logger.LoggerName, "INFO"), string_)
	}
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
	data, _ := pem.Decode([]byte(key))
	serverPubKey, err := x509.ParsePKCS1PublicKey(data.Bytes)
	if err != nil {
		GoLog("Error while importing key:", err)
	}
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
	var opts rsa.PSSOptions = rsa.PSSOptions{SaltLength: 20}
	err := rsa.VerifyPSS(&pubKey, crypto.SHA256, Sha256BtB(data), sign, &opts)
	if err != nil {
		panic(err)
	}
}

func UnixMs() int64 {
	return time.Now().UnixNano() / 1000000
}

func ClearEmptyStrings(elements []string) []string {
	var ret []string
	for _, val := range elements {
		if len(val) > 0 {
			ret = append(ret, val)
		}
	}
	return ret
}

func ClearEmtyElements(badExample interface{}, elements ...interface{}) []interface{} {
	var ret []interface{}
	for _, element := range elements {
		if element != badExample {
			ret = append(ret, element)
		}
	}
	return ret
}

func H(err error) {
	if err != nil {
		panic(err)
	}
}

func CountMapElementsStartsWith(m map[string]interface{}, text string) int {
	count := 0
	for key := range m {
		if strings.HasPrefix(key, text) {
			count++
		}
	}
	return count
}

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Encrypt(key, data []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext
}

func Decrypt(key, data, iv []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(data) < aes.BlockSize {
		log.Println("ciphertext too short: " + strconv.Itoa(len(data)))
		return []byte{}
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return data
}

func DecryptBtB(strkey string, bytes []byte) []byte {
	return Decrypt(H2b(strkey), bytes[aes.BlockSize:], bytes[:aes.BlockSize])
}

func DecryptHtB(strkey, hexData string) []byte {
	data := H2b(hexData)
	//fmt.Println(hexData, data)
	return Decrypt(H2b(strkey), data[aes.BlockSize:], data[:aes.BlockSize])
}

func DecryptHtS(strkey, hexStr string) string {
	ciphertext, _ := hex.DecodeString(hexStr)
	//fmt.Println("IV:", B2h(ciphertext[:aes.BlockSize]))
	//fmt.Println("Data:", B2h(ciphertext[aes.BlockSize:]))
	return string(Decrypt(H2b(strkey), ciphertext[aes.BlockSize:], ciphertext[:aes.BlockSize]))
}

func DecryptB64tB(strkey, b64 string) []byte {
	ciphertext, _ := base64.StdEncoding.DecodeString(b64)
	return Decrypt(H2b(strkey), ciphertext[aes.BlockSize:], ciphertext[:aes.BlockSize])
}

func DecryptB64tB_safe(strkey, b64 string) []byte {
	return DecryptB64tB(strkey, b64)[16:]
}

func EncryptBtB(strkey string, data []byte) []byte {
	return Encrypt(H2b(strkey), data)
}

func EncryptBtH(strkey string, data []byte) string {
	return hex.EncodeToString(Encrypt(H2b(strkey), data))
}

func EncryptBtB64(strkey string, data []byte) string {
	return base64.StdEncoding.EncodeToString(Encrypt(H2b(strkey), data))
}

func EncryptBtB64_safe(strkey string, data []byte) string {
	iv := make([]byte, 16)
	rand.Read(iv)
	return EncryptBtB64(strkey, append(iv, data...))
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
	defer H(f.Close())

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

//DEPRECATED
func FindGroup_(text string, reg string) []string {
	regex, _ := regexp.Compile(reg)
	return regex.FindStringSubmatch(text)
}

//DEPRECATED
func FindGroups_(text string, reg string) []string {
	regex, _ := regexp.Compile(reg)
	temp := regex.FindAllStringSubmatch(text, -1)
	var temp4 []string
	for _, temp2 := range temp {
		temp4 = append(temp4, temp2[1])
	}
	return temp4
}

//DEPRECATED
func FindAllGroups_(text string, reg string) [][]string {
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

func FindGroups(text, regex string) []string {
	var res []string
	var re = regexp.MustCompile(regex)

	for _, match := range re.FindAllString(text, -1) {
		res = append(res, match)
	}
	return res
}

func GetServerPlayers(ip string) (bool, []string) {
	req := query.NewRequest()
	_ = req.Connect(ip)
	response, _ := req.Full()

	if response == nil { //Cant connect to server
		return false, nil
	}

	return len(response.Players) > 0, response.Players
}

func BytesToBool(bytes []byte) bool {
	ret, err := strconv.ParseBool(string(bytes))
	H(err)
	return ret
}

func DbGet_(db *leveldb.DB, key string, defVal []byte) []byte {
	if DbHas(db, key) {
		return DbGet(db, key)
	}
	return defVal
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

func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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

func (SliceTools) GetIntIndex(slice []int, element int) int {
	for i, val := range slice {
		if val == element {
			return i
		}
	}
	return -1
}

func DecodeWindows1251(ba []uint8) []uint8 {
	dec := charmap.Windows1251.NewDecoder()
	out, _ := dec.Bytes(ba)
	return out
}

func CapSolveV3(ses *gorequest.SuperAgent, url, action, key, apikey string) (string, string) {
begin:
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
	if capchaResponse2.Request == "ERROR_CAPTCHA_UNSOLVABLE" {
		goto begin
	}
	return capchaResponse2.Request, capchaResponse1.Request
}

func CapSolveV2(ses *gorequest.SuperAgent, url, key, apiKey string) (string, string) { //returns cap-response, cap-id
startCaptcha:
	gurl := "https://rucaptcha.com/in.php?key=%s&method=userrecaptcha&version=v2&" +
		"googlekey=%s&pageurl=%s&json=1"
	gurl = fmt.Sprintf(gurl, apiKey, key, url)
	_, response, _ := ses.Get(gurl).EndBytes()
	var captchaResponse1 captchaResponseStruct
	_ = json.Unmarshal(response, &captchaResponse1)
	if captchaResponse1.Status != 1 {
		goto startCaptcha
	}

	var capchaResponse2 captchaResponseStruct
waitForCaptcha:
	switch capchaResponse2.Request {
	case "CAPCHA_NOT_READY", "":
		_, capchaResponse2B, _ := ses.Get(fmt.Sprintf("https://rucaptcha.com/res.php?key=%s&action=get&taskinfo=0&json=1&id=%s", apiKey, captchaResponse1.Request)).EndBytes()
		_ = json.Unmarshal(capchaResponse2B, &capchaResponse2)
		time.Sleep(2000 * time.Millisecond)
		goto waitForCaptcha
	case "ERROR_CAPTCHA_UNSOLVABLE":
		goto startCaptcha
	}
	return capchaResponse2.Request, captchaResponse1.Request
}

func CapReport(ses *gorequest.SuperAgent, good bool, apikey, capid string) {
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
	//CapSolveV3(ses, pageUrl, action, siteKey, ruCaptchaKey)

	var Json registerJson
	var capid string
	Json.Login = login
	Json.Password = password
	Json.Captcha, capid = CapSolveV3(ses, pageUrl, action, siteKey, ruCaptchaKey)
	Json.Email = email
	Jsonb, _ := json.Marshal(Json)

	_, data, _ := ses.Post(pageUrl).Set("x-csrf-token", csrf).Set("x-xsrf-token", xsrf).Send(string(Jsonb)).EndBytes()
	var registerResponseJson registerResponseJsonStruct
	_ = json.Unmarshal(data, &registerResponseJson)
	if registerResponseJson.Success == true {
		CapReport(ses, true, ruCaptchaKey, capid)
	} else {
		CapReport(ses, true, ruCaptchaKey, capid)
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
		_, temp2, _ := ses.Get(CategoryUrl + thread).End()
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
