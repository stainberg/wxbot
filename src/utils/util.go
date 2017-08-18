package utils

import (
	"unsafe"
	"net/url"
	"reflect"
	"regexp"
	"math/rand"
	"time"
	"strconv"
	"encoding/json"
	"strings"
	"fmt"
	"encoding/hex"
	"encoding/base64"
	"crypto/md5"
)

const (
	RAND_KIND_NUM   = 0  // 纯数字
	RAND_KIND_LOWER = 1  // 小写字母
	RAND_KIND_UPPER = 2  // 大写字母
	RAND_KIND_ALL   = 3  // 数字、大小写字母
)

func StringBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func BytesString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func HasParams(vars url.Values, keys ...string) bool {
	if vars != nil {
		for _, key := range keys {
			_, ok := vars[key]
			if !ok {
				return false
			}
		}
	}
	return true
}

func Contain(target interface{}, obj interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

func CheckMobileNumber(mobile string) bool {
	regular := `^1[34578]\d{9}$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobile)
}

func Random(size int, kind int) string {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i :=0; i < size; i++ {
		if is_all {
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base+rand.Intn(scope))
	}
	return string(result)
}

func ReverseString(s string) string {
	str := []rune(s)
	for i, j := 0, len(str) - 1; i < j; i, j = i + 1, j - 1 {
		str[i], str[j] = str[j], str[i]
	}
	return string(str)
}

func DelSlinceItem(s []string, item string) []string {
	var ss []string;
	for i := 0; i < len(s); i++ {
		if item != s[i] {
			ss = append(ss, s[i])
		}
	}
	return ss
}

func GenerateId() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func JsonEncode(nodes interface{}) string {
	body, err := json.Marshal(nodes)
	if err != nil {
		panic(err.Error())
		return `{}`
	}
	return string(body)
}

func JsonDecode(jsonStr string) interface{} {
	jsonStr = strings.Replace(jsonStr, "\n", "", -1)
	var f interface{}
	err := json.Unmarshal([]byte(jsonStr), &f)
	if err != nil {
		fmt.Println(jsonStr)
		return false
	}
	return float2Int(f)
}

func float2Int(input interface{}) interface{} {
	if m, ok := input.([]interface{}); ok {
		for k, v := range m {
			switch v.(type) {
			case float64:
				m[k] = int(v.(float64))
			case []interface{}:
				m[k] = float2Int(m[k])
			case map[string]interface{}:
				m[k] = float2Int(m[k])
			}
		}
	} else if m, ok := input.(map[string]interface{}); ok {
		for k, v := range m {
			switch v.(type) {
			case float64:
				m[k] = int(v.(float64))
			case []interface{}:
				m[k] = float2Int(m[k])
			case map[string]interface{}:
				m[k] = float2Int(m[k])
			}
		}
	} else {
		return false
	}
	return input
}

func SecurityMD5(src string) string {
	hs := md5.New()
	hs.Reset()
	hs.Write(StringBytes(src))
	return strings.ToLower(hex.EncodeToString(hs.Sum(nil)))
}

func Base64Encode(src string) string {
	bs := base64.URLEncoding.EncodeToString(StringBytes(src))
	dst := strings.Replace(string(bs), "/", "_", -1)
	dst = strings.Replace(dst, "+", "-", -1)
	dst = strings.Replace(dst, "=", "", -1)
	return dst
}

func Base64Decode(src string) (string) {
	var missing = (4 - len(src)%4) % 4
	src += strings.Repeat("=", missing)
	db, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		println(err.Error())
	}
	return BytesString(db)
}