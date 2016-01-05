package encrypt

import (
	"errors"
	"public_tool"
)

type Encrypt interface {
	Init(key ...interface{}) error
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

func NewEncrypt(key ...[]byte) (Encrypt, error) {
	if len(key) == 2 {
		pub := key[0]
		pri := key[1]
		R := public_tool.NewRsa()
		if err := R.Init(pub, pri); err != nil {
			return nil, err
		}
		return R, nil
	}
	if len(key) == 1 {
		D := public_tool.NewDes()
		if err := D.Init(key[0]); err != nil {
			return nil, err
		}
		return D, nil
	}
	return nil, errors.New("Args length must 1 or 2")
}
