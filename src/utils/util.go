package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"unsafe"
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
		Log("Utils JsonDecode", err.Error())
		Log("Utils JsonDecode", jsonStr)
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

func CheckToken(p string) bool {
	if p == Conf.HttpConf.Token {
		return true
	}
	return false
}