package utils

import (
	"unsafe"
	"time"
	"strconv"
	"encoding/json"
	"strings"
	"fmt"
	"encoding/hex"
	"encoding/base64"
	"crypto/md5"
)

func StringBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func BytesString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func GenerateId() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func JsonEncode(nodes interface{}) string {
	body, err := json.Marshal(nodes)
	if err != nil {
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