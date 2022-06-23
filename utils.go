package utils

import (
	"bufio"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty"
	//"github.com/SilverCory/golang_discord_rpc"
	"github.com/syndtr/goleveldb/leveldb"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	//"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding/charmap"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	//"syscall"
	"time"
	//"unsafe"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36"
const serverAddr = "http://127.0.0.1:8973"

type Ints interface {
	int | int8 | int16 | int32 | int64
}

type Floats interface {
	float32 | float64
}

type Uints interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type complexes interface {
	complex64 | complex128
}

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

type SliceTools struct {
}

func Marshal(a interface{}) (b []byte) {
	b, _ = json.Marshal(a)
	return
}

func WaitChanAny[T comparable, T2 comparable](ch1 <-chan T, ch2 <-chan T2) {
	select {
	case <-ch1:
		break
	case <-ch2:
		break
	}
}

func WaitChanAnyAlwaysTrue[T comparable, T2 comparable](ch1 <-chan T, ch2 <-chan T2) bool {
	select {
	case <-ch1:
		break
	case <-ch2:
		break
	}
	return true
}

type ProxyData struct {
	ProxyType     string `json:"proxyType"`
	ProxyAddress  string `json:"proxyAddress"`
	ProxyPort     int    `json:"proxyPort"`
	ProxyLogin    string `json:"proxyLogin"`
	ProxyPassword string `json:"proxyPassword"`
	UserAgent     string `json:"userAgent,omitempty"`
	Cookies       string `json:"cookies,omitempty"`
}

func (p ProxyData) String() string {
	return p.ProxyType + "://" + p.StringNoType()
}

func (p *ProxyData) StringNoType() string {
	var a = p.ProxyAddress + ":" + strconv.Itoa(p.ProxyPort)
	if p.ProxyLogin != "" {
		if p.ProxyPassword != "" {
			a = p.ProxyLogin + ":" + p.ProxyPassword + "@" + a
		} else {
			a = p.ProxyLogin + "@" + a
		}
	}
	return a
}

func RepeatStringToSlice(s string, n int) []string {
	a := make([]string, n)
	for i := range a {
		a[i] = s
	}
	return a
}

func BuildQuestionMarks(n int) string {
	if n == 0 {
		return ""
	}
	return strings.Join(RepeatStringToSlice("?", n), ",")
}

func PressEnterToExit(msg ...interface{}) {
	fmt.Println(msg...)
	fmt.Println("\n\rPress Enter to continue...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(0)
}

func PrintAndExit(msg ...interface{}) {
	fmt.Println(msg...)
	os.Exit(0)
}

func RemoveOverflowFiles(path string, overflowCount int) error {
	var filesDates []int64
	var filesMap = map[int64][]string{}
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filesDates = append(filesDates, info.ModTime().Unix())
			filesMap[info.ModTime().Unix()] = append(filesMap[info.ModTime().Unix()], path)
		}
		return err
	}); err != nil {
		return err
	}

	if len(filesDates) > overflowCount {
		sort.Slice(filesDates, func(i, j int) bool {
			return filesDates[i] > filesDates[j]
		})

		for i := overflowCount; i < len(filesDates); i++ {
			for _, path := range filesMap[filesDates[i]] {
				if err := os.Remove(path); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func JoinErrors(e1 ...error) (e2 error) {
	var s []string
	for _, e := range e1 {
		if e != nil {
			s = append(s, e.Error())
		}
	}
	if s != nil {
		e2 = errors.New(strings.Join(s, "; "))
	}
	return
}

func Abs[T Ints | Floats](a T) T {
	if a >= 0 {
		return a
	}
	return -a
}

func B64StringFix(s string) string {
	if i := len(s) % 4; i != 0 {
		s += strings.Repeat("=", 4-i)
	}
	return s
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

func RemoveElements(elementToRemove interface{}, elements ...interface{}) []interface{} {
	var ret []interface{}
	for _, element := range elements {
		if element != elementToRemove {
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

func Reverse[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func Sha512B2B(data []byte) []byte {
	h := sha512.New()
	h.Write(data)
	return h.Sum(nil)
}

func Sha512S2B(text string) []byte {
	return Sha512B2B([]byte(text))
}

func Sha512B2H(data []byte) string {
	return hex.EncodeToString(Sha512B2B(data))
}

func Sha512S2H(text string) string {
	return hex.EncodeToString(Sha512B2B([]byte(text)))
}

func Sha256B2B(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

func Sha256S2B(text string) []byte {
	return Sha256B2B([]byte(text))
}

func Sha256B2H(data []byte) string {
	return hex.EncodeToString(Sha256B2B(data))
}

func Sha256S2H(text string) string {
	return hex.EncodeToString(Sha256B2B([]byte(text)))
}

func Sha256File(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
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

func ReadFileByLines(path string) (error, []string) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil { //If error occupied while reading file
		return err, nil
	}

	var ret []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	_ = file.Close() //Closing file finally

	return nil, ret
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
	return FindListOfGroups(text, reg)
}

func FindListOfGroups(text string, reg string) [][]string {
	//Returns data: [[field1.1, field1.2...], [field2.1, field2.2...], ...]
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

func FindRegexText(text, regex string) []string {
	var res []string
	var re = regexp.MustCompile(regex)

	for _, match := range re.FindAllString(text, -1) {
		res = append(res, match)
	}
	return res
}

func FindRegexNamedGroups(data, regex string) []string {
	var namedGroups []string

	r, err := regexp.Compile(regex)
	if err != nil {
		return namedGroups
	}

	allStringSubmatch := r.FindAllStringSubmatch(data, -1)
	for _, val := range allStringSubmatch {
		if len(val) != 2 {
			continue
		}
		namedGroups = append(namedGroups, val[1])
	}
	return namedGroups
}

//func GetQueryServerPlayers(ip string) (bool, []string) {
//	req := query.NewRequest()
//	_ = req.Connect(ip)
//	response, _ := req.Full()
//
//	if response == nil { //Cant connect to server
//		return false, nil
//	}
//
//	return len(response.Players) > 0, response.Players
//}

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

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsF[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsSubSlice[T comparable](parent, child []T) bool {
	var c = true
p:
	for _, a := range child {
		for _, b := range parent {
			if a == b {
				continue p
			}
		}
		c = false
		break
	}
	return c
}

func ContainsAnyOfSubSlice[T comparable](parent, child []T) bool {
	for _, el := range parent {
		for _, ch := range child {
			if ch == el {
				return true
			}
		}
	}
	return false
}

func Md5S2S(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func Md5B2H(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func Md5B2B(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
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

//DEPRECATED
//func (ScUtils) RegisterAccount(ses *gorequest.SuperAgent, ruCaptchaKey string) (string, string, string) {
//	// Returns login string, password string, csrf string
//	const siteKey = "6LeedroUAAAAAK2RUkaNLVBYraeQXNVHX45O227A"
//	const pageUrl = "https://streamcraft.net/api/auth/register"
//	const regexToken = `<meta name="csrf-token" content="(?P<token>.*)">`
//	const action = "register"
//
//register:
//	resp, page, _ := ses.Get(pageUrl).End()
//
//	csrf := FindGroup_(page, regexToken)[1]
//
//	xsrf := (*http.Response)(resp).Cookies()[1].Value
//	//os.Exit(111)
//	email := randomdata.Email()
//	name := randomdata.FirstName(randomdata.Number(1, 2))
//	length := 8
//	if len(name) < 8 {
//		length = len(name)
//	}
//	login := name[:length] + randomdata.RandStringRunes(4) + strconv.Itoa(randomdata.Number(1980, 2017))
//	password := login + login
//
//	//RuCaptcha
//	//CapSolveV3(Ses, pageUrl, action, siteKey, ruCaptchaKey)
//
//	var Json registerJson
//	var capid string
//	Json.Login = login
//	Json.Password = password
//	Json.Captcha, capid = CapSolveV3(ses, pageUrl, action, siteKey, ruCaptchaKey)
//	Json.Email = email
//	Jsonb, _ := json.Marshal(Json)
//
//	_, data, _ := ses.Post(pageUrl).Set("x-csrf-token", csrf).Set("x-xsrf-token", xsrf).Send(string(Jsonb)).EndBytes()
//	var registerResponseJson registerResponseJsonStruct
//	_ = json.Unmarshal(data, &registerResponseJson)
//	if registerResponseJson.Success == true {
//		CapReport(ses, true, ruCaptchaKey, capid)
//	} else {
//		CapReport(ses, true, ruCaptchaKey, capid)
//		goto register
//	}
//
//	fmt.Println(string(data))
//	return login, password, csrf
//}

////DEPRECATED
//func (ScUtils) SetReputation(ses *gorequest.SuperAgent, csrf string, userId int, count int) {
//	//Set user reputation
//	const pageUrl = "https://streamcraft.net/forum/user/reputation"
//
//	var Json reputationJson
//	Json.User = userId
//	Json.Reputation = count
//	Json.Token = csrf
//	JsonB, _ := json.Marshal(Json)
//
//	ses.Post(pageUrl).Send(string(JsonB)).End()
//}
//
////DEPRECATED
//func (ScUtils) GetUserId(ses *gorequest.SuperAgent, nickname string) int {
//	//Get user id
//	const pageUrl = "https://streamcraft.net/user/"
//	const regexUserId = `<i class="fa fa-thumbs-down cursor-pointer" onclick="App\.sendRequest\('/forum/user/reputation', {user: (?P<id>.*), reputation: -1}\);"></i>`
//
//	_, page, _ := ses.Get(pageUrl + nickname).End()
//	id, _ := strconv.Atoi(FindGroup_(page, regexUserId)[1])
//	return id
//}
//
////DEPRECATED
//func (ScUtils) ThreadsIdsParse(ses *gorequest.SuperAgent) []string {
//	const regexThreads = `<a href="/forum/category/(?P<id>.*)"><i class="fa fa-level-down">`
//	const regexThreadsIds = `<a class="btn btn-primary btn-shadow float-right" href="/forum/discussion/create/(?P<id>.*)" role="button">`
//	const ForumUrl = "https://streamcraft.net/forum/"
//	const CategoryUrl = "https://streamcraft.net/forum/category/"
//
//	_, text, _ := ses.Get(ForumUrl).End()
//	temp := FindRegexText(text, regexThreads)
//	var threadsIds []string
//	for _, thread := range temp {
//		_, temp2, _ := ses.Get(CategoryUrl + thread).End()
//		temp3 := FindGroup_(temp2, regexThreadsIds)
//		if len(temp3) < 2 {
//			continue
//		}
//		threadsIds = append(threadsIds, temp3[1])
//	}
//	return threadsIds
//}

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

//func (K773Utils) DecryptAes(data, key string) string {
//	var request requestStruct
//	request.What, request.Data, request.Key = "decrypt", data, key
//	marshalled, _ := json.Marshal(request)
//	_, response, _ := gorequest.New().Post(serverAddr).
//		Send(string(marshalled)).
//		End()
//	return response
//}
//
//func (K773Utils) EncryptAes(data, key string) string {
//	var request requestStruct
//	request.What, request.Data, request.Key = "encrypt", data, key
//	marshalled, _ := json.Marshal(request)
//	_, response, _ := gorequest.New().Post(serverAddr).
//		Send(string(marshalled)).
//		End()
//	return response
//}

//func SetDiscordStatus(server ServerStruct, nickname string) {
//start:
//	win := discordrpc.NewRPCConnection("496419141201297413")
//	err := win.Open()
//	if err != nil {
//		//fmt.Println(err)
//		time.Sleep(5 * time.Second)
//		goto start
//	}
//
//	_, _ = win.Read()
//	//fmt.Println(err)
//	//fmt.Println(str)
//
//	//time.Sleep(time.Second * 3)
//	stamp := time.Now().Unix()
//
//	for {
//		//fmt.Println(os.Getpid())
//		presence := &discordrpc.CommandRichPresenceMessage{
//			CommandMessage: discordrpc.CommandMessage{Command: "SET_ACTIVITY"},
//			Args: &discordrpc.RichPresenceMessageArgs{
//				Pid: os.Getpid(),
//				Activity: &discordrpc.Activity{
//					Details:    "Играет на " + server.ServerName,
//					State:      "С ником " + nickname,
//					Instance:   false,
//					TimeStamps: &discordrpc.TimeStamps{StartTimestamp: stamp},
//					Assets: &discordrpc.Assets{
//						LargeText:    server.LargeText,
//						LargeImageID: server.LargeTextId,
//						SmallText:    "StreamCraft.Net",
//						SmallImageID: "discord",
//					},
//				},
//			},
//		}
//
//		presence.SetNonce()
//		data, err := json.Marshal(presence)
//
//		if err != nil {
//			fmt.Println(err)
//			continue
//		}
//
//		err = win.Write(string(data))
//		if err != nil {
//			fmt.Println(err)
//			continue
//		}
//
//		//str, err := win.Read()
//		//fmt.Println(err)
//		//fmt.Println(str)
//		//
//		//fmt.Println("---\nDone?")
//		time.Sleep(time.Second * 5)
//	}
//}

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

func AreArraysEqual[T comparable](a, b []T, orderIsImportant bool) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	if orderIsImportant {
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
	} else {
		for _, item := range a {
			if !Contains(b, item) {
				return false
			}
		}
	}

	return true
}

func AreArraysEqualF[T, T2 comparable](a []T, b []T2, orderIsImportant bool, equal func(i, j int) bool) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	if orderIsImportant {
		for i := range a {
			if !equal(i, i) {
				return false
			}
		}
	} else {
		for i := range a {
			var anyF bool
			for j := 0; j < len(b); j++ {
				if anyF = equal(i, j); anyF {
					break
				}
			}
			if !anyF {
				return false
			}
		}
	}

	return true
}

func Array1ContainsOnlyOfArray2[T comparable](a, b []T) bool {
	var m = make(map[T]struct{}, len(b))
	for _, v := range b {
		m[v] = struct{}{}
	}
	for _, v := range a {
		if _, h := m[v]; !h {
			return false
		}
	}
	return true
}

func Array1ContainsOnlyKeysOfMap2[T comparable, V comparable](a []T, m map[T]V) bool {
	for _, v := range a {
		if _, h := m[v]; !h {
			return false
		}
	}
	return true
}

func ArrayAny[T comparable](a []T, f func(T) bool) bool {
	for _, b := range a {
		if f(b) {
			return true
		}
	}
	return false
}

// SlicesSubtract returns array of elements that exist in the array 1 but do not exist in the array 2
func SlicesSubtract[T comparable](arr1, arr2 []T) []T {
	var a []T
	var has = make(map[T]struct{})
	for _, el := range arr2 {
		has[el] = struct{}{}
	}
	for _, el := range arr1 {
		if _, h := has[el]; !h {
			a = append(a, el)
		}
	}
	return a
}

func LastElement[T any](arr []T) T {
	return arr[len(arr)-1]
}

func MapAnyKey[T comparable, T2 comparable](a map[T]T2, f func(T) bool) bool {
	for k := range a {
		if f(k) {
			return true
		}
	}
	return false
}

func MapAnyValue[T comparable, T2 comparable](a map[T]T2, f func(T2) bool) bool {
	for _, v := range a {
		if f(v) {
			return true
		}
	}
	return false
}

func MapAnyKeyValue[T comparable, T2 comparable](a map[T]T2, f func(T, T2) bool) bool {
	for k, v := range a {
		if f(k, v) {
			return true
		}
	}
	return false
}

// MapSlice will call f() for each index of the provided slice and combine the results into a single array.
// Return false if no value needs to be inserted.
func MapSlice[T, resT any](in []T, f func(i int) (resT, bool)) []resT {
	var res = make([]resT, 0, len(in))
	for i := range in {
		v, do := f(i)
		if do {
			res = append(res, v)
		}
	}
	return res
}

// MapMap will call f() for each keypair of the provided map and combine the results into a single map.
// Return false if no value needs to be inserted.
func MapMap[K1 comparable, V1 any, K2 comparable, V2 any](m map[K1]V1, f func(K1, V1) (K2, V2, bool)) map[K2]V2 {
	var res = make(map[K2]V2, len(m))
	for k, v := range m {
		k1, v1, do := f(k, v)
		if do {
			res[k1] = v1
		}
	}
	return res
}

// MapSlice2Map will call f() for each index of the provided slice and combine the results into a single map.
// Return false if no value needs to be inserted.
func MapSlice2Map[SliceT any, MapK comparable, MapV any](s []SliceT, f func(i int) (MapK, MapV, bool)) map[MapK]MapV {
	var res = make(map[MapK]MapV, len(s))
	for i := range s {
		k1, v1, do := f(i)
		if do {
			res[k1] = v1
		}
	}
	return res
}

func Slice2String[T any](s []T) string {
	return string(Marshal(&s))
}

func If[T any](a bool, ifTrue T, ifFalse T) T {
	if a {
		return ifTrue
	}
	return ifFalse
}

func CallIf(a bool, ifTrue func(), ifFalse func()) {
	if a {
		if ifTrue != nil {
			ifTrue()
		}
	} else if ifFalse != nil {
		ifFalse()
	}
}

func Avg[T Ints | Uints | Floats | time.Duration](val ...T) float64 {
	if len(val) == 0 {
		return 0
	}

	var a T
	for _, v := range val {
		a += v
	}
	return float64(a) / float64(len(val))
}

func CopyMapValuesToSlice[K, V comparable](m map[K]V) []V {
	var ret = make([]V, len(m))
	i := 0
	for _, v := range m {
		ret[i] = v
	}
	return ret
}

func MemsetRepeat[T comparable](a []T, n int, v T) {
	if len(a) == 0 {
		return
	}
	a[0] = v
	for bp := 1; bp < n; bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}

func Copy[T comparable](a []T) []T {
	var c = make([]T, len(a))
	copy(c, a)
	return c
}

func SliceL22AnySliceL2[T any](a [][]T) [][]any {
	var ret = make([][]any, len(a))
	for i := range a {
		ret[i] = Slice2AnySlice(a[i])
	}
	return ret
}

func Slice2AnySlice[T any](a []T) []any {
	var ret = make([]any, len(a))
	for i := range a {
		ret[i] = a[i]
	}
	return ret
}

// Append is slightly more effective than go's append(): 10k+10k el.: {"go": 12000ns, "utils": 8000ns}
func Append[T any](a1, a2 []T) []T {
	var a3 = make([]T, len(a1)+len(a2))
	copy(a3, a1)
	copy(a3[len(a1):], a2)
	return a3
}

func MapOr[T comparable, T2 any](m1, m2 map[T]T2) map[T]T2 {
	var m3 = make(map[T]T2, If(len(m2) > len(m1), len(m2), len(m1)))

	for k, el := range m1 {
		m3[k] = el
	}
	for k, el := range m2 {
		if _, h := m3[k]; !h {
			m3[k] = el
		}
	}
	return m3
}

func Slice2HasMap[T comparable](a []T) map[T]struct{} {
	var ret = make(map[T]struct{}, len(a))
	for _, v := range a {
		ret[v] = struct{}{}
	}
	return ret
}

func Slice2HasMapExcludeEmpty[T comparable](a []T) map[T]struct{} {
	var ret = make(map[T]struct{}, len(a))
	var def T
	for _, v := range a {
		if v != def {
			ret[v] = struct{}{}
		}
	}
	return ret
}

func SliceEvery[T comparable](a []T, f func(i int) bool) bool {
	for i := range a {
		if !f(i) {
			return false
		}
	}

	return len(a) != 0
}

func PrintAsBinaryAsBigEndian(bytes []byte) {
	for i := len(bytes) - 1; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			zeroOrOne := bytes[i] >> (7 - j) & 1
			fmt.Printf("%c", '0'+zeroOrOne)
		}
		fmt.Print(" ")
	}
	fmt.Println()
}

func PrintAsBinaryAsLittleEndian(bytes []byte) {
	for i := 0; i < len(bytes); i++ {
		for j := 7; j >= 0; j-- {
			zeroOrOne := bytes[i] >> (7 - j) & 1
			fmt.Printf("%c", '0'+zeroOrOne)
		}
		fmt.Print(" ")
	}
	fmt.Println()
}

func Bool2Int(b bool) uint8 {
	var n uint8
	if b {
		n = 1
	}
	return n
}

func SplitNumberIntoPartsByOverflow[T Ints | Uints](num, overflow T) [][2]T {
	var whole = int(num / overflow)
	var mod = num%overflow != 0

	var l = make([][2]T, whole+int(Bool2Int(mod)))
	for i := 0; i < whole; i++ {
		l[i] = [2]T{overflow * T(i), overflow * T(i+1)}
	}
	if mod {
		l[len(l)-1] = [2]T{overflow * T(whole), num}
	}
	return l
}

// SplitNumberIntoParts splits number into n equal parts. If mod:=num%n != 0, the last element will receive all modulo.
// Will panic if num < n.
// Sample: [0 2] [2 4] [4 6]
func SplitNumberIntoParts[T Ints | Uints](num T, n int) [][2]T {
	var l = make([][2]T, n)
	part := num / T(n)
	for i := range l {
		l[i] = [2]T{part * T(i), part * T(i+1)}
	}
	if mod := num % T(n); mod != 0 {
		le := len(l)
		if le != 0 {
			l[le-1][1] += mod
		}
	}
	return l
}

func SumArrFunc[T Ints | Uints | Floats](arrayLength int, sum func(i int) T) T {
	var b T
	for i := 0; i < arrayLength; i++ {
		b += sum(i)
	}
	return b
}

func ProgressCalculator[T Ints | Uints | Floats](total T) func(add T) (total, progress T, abs float64) {
	var progress T
	var s sync.Mutex
	return func(add T) (_, _ T, _ float64) {
		s.Lock()
		defer s.Unlock()
		progress += add
		return total, progress, If(total == 0, 100.0, float64(progress)/float64(total))
	}
}

func Clamp[T Ints | Uints | Floats | time.Duration](num, min, max T) T {
	if num < min {
		return min
	} else if num > max {
		return max
	} else {
		return num
	}
}

func ParseIpFromIpPort(src string) string {
	var i = strings.IndexByte(src, ':')
	if i != -1 {
		return src[:i]
	}
	return src
}

//func VerifyProxyConnection(sesNoProxy, sesProxy *gorequest.SuperAgent) (proxyDelay int64, noProxyIp, proxyIp string, e error) {
//	var r gorequest.Response
//
//	t := UnixMs()
//	r, noProxyIp, _ = sesNoProxy.Get("https://api64.ipify.org?format=plaintext").End()
//	myPing := UnixMs() - t
//	if r == nil {
//		e = errors.New("verifyProxyConnection: sesNoProxy timeout reached")
//	} else {
//		_ = r.Body.Close()
//
//		t2 := UnixMs()
//		r, proxyIp, _ = sesProxy.Get("https://api64.ipify.org?format=plaintext").End()
//		proxyDelay = UnixMs() - t2 - myPing
//		if r == nil {
//			e = errors.New("verifyProxyConnection: sesNoProxy timeout reached")
//		} else {
//			_ = r.Body.Close()
//
//			if proxyIp == noProxyIp {
//				e = errors.New("verifyProxyConnection: proxied ip response is equal to non-proxy ip response (" + noProxyIp + ")")
//			}
//		}
//	}
//
//	return proxyDelay, noProxyIp, proxyIp, e
//}

func GetGoogleDriveDocumentContent(s *resty.Client, docID string) (string, error) {
	r, e := s.R().Get(fmt.Sprintf("https://drive.google.com/uc?id=%v&export=download", docID))
	if e == nil {
		if r.IsError() {
			e = errors.New(r.Status())
		}
	}
	return r.String(), e
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

func CommitTransaction(tx *sql.Tx, log func(msg string), n int) (e error) {
	for i := 0; i == -1 || i < n; i++ {
		if e = tx.Commit(); e == nil {
			break
		} else {
			log("tx commit error: " + e.Error())
			time.Sleep(time.Second)
		}
	}
	return e
}

func RestartTransaction(db *sql.DB, tx *sql.Tx, log func(msg string), n int) (txn *sql.Tx, e error) {
	for i := 0; n == -1 || i < n; i++ {
		if e = tx.Commit(); e == nil {
			for i := 0; n == -1 || i < n; i++ {
				if txn, e = db.Begin(); e == nil {
					break
				} else {
					log("tx begin error: " + e.Error())
					time.Sleep(time.Second)
				}
			}
			break
		} else {
			log("tx commit error: " + e.Error())
			time.Sleep(time.Second)
		}
	}
	return txn, e
}

func NewSimpleLog(prefix string, logTime bool) func(string2 string) {
	return func(str string) {
		var a string
		if logTime {
			a += "[" + time.Now().Format("15:04:05") + "]"
		}
		fmt.Println(a, prefix, ":", str)
	}
}

func SleepWithContext(ctx context.Context, duration time.Duration) (e error) {
	if duration == 0 {
		return nil
	}

	t := time.NewTimer(duration)
	select {
	case <-ctx.Done():
		break
	case <-t.C:
		break
	}
	t.Stop()
	return ctx.Err()
}

func ReadLine(a io.Reader) string {
	reader := bufio.NewReader(a)
	text, _ := reader.ReadString('\n')
	return text
}
