package public_tool

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
)

type Des struct {
	block cipher.Block
}

func NewDes() *Des {
	d := new(Des)
	d.Init([]byte{0x1, 0x2, 0x3})
	return d
}
func (self *Des) Init(key []byte) error {
	if len(key) != 8 {
		return errors.New("key size must 8")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return err
	}
	self.block = block
	return nil
}

func (self *Des) Encrypt(data []byte) ([]byte, error) {
	bs := self.block.BlockSize()
	data = pKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		self.block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

func (self *Des) Decrypt(data []byte) ([]byte, error) {
	bs := self.block.BlockSize()
	if len(data)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		self.block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = pKCS5UnPadding(out)
	return out, nil
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
