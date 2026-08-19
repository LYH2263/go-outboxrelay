package codec

import "encoding/json"

// Marshal 编码。
func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal 解码到新缓冲。
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(CloneBytes(data), v)
}

// MustJSON 测试辅助。
func MustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
