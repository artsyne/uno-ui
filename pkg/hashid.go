package pkg

import (
	"fmt"

	"github.com/speps/go-hashids/v2"
)

// HashID 是一个封装了 HashIDData 的结构体
type HashID struct {
	hd *hashids.HashID
}

// GlobalHashID 是一个全局的 HashID 实例
var GlobalHashID *HashID

// NewHashID 创建一个新的 HashID 实例
func NewHashID(salt string, minLength int) *HashID {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, _ := hashids.NewWithData(hd)
	return &HashID{hd: h}
}

// Encode 对整数进行编码
func (h *HashID) Encode(id int) (string, error) {
	e, err := h.hd.Encode([]int{id})
	if err != nil {
		return "", err
	}
	return e, nil
}

// Decode 对编码字符串进行解码
func (h *HashID) Decode(encodedID string) (int, error) {
	d, err := h.hd.DecodeWithError(encodedID)
	if err != nil {
		return 0, err
	}
	if len(d) == 0 {
		return 0, fmt.Errorf("invalid ID")
	}
	return d[0], nil
}
