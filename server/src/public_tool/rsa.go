package public_tool

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type Rsa struct {
	pub *rsa.PublicKey
	pri *rsa.PrivateKey
}

func NewRsa() *Rsa {
	pub := `-----BEGIN 公钥-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4XXxm9EAsOKVX/YgMHSM
AM44gyPx8PdkXPeh0de94qTS2kyvFG5RJYSHI3xC+cyltnMzfFuDbjKB0/sfnfYP
yNANFvzbVGnhEPvulsX+vbJdnuCVcVI/H/cjOJGJ+MKTGQz55n+Idp/BsQhBc+9u
b2PvdnFxng6zki3HFJ4Kq9onWjV4smFe0JlqIqDg+xD8U+zLvlQdSDD7oNhBYAML
Zn1uYbbh9Adoh7nDkqLdLLRsJqi5lAftRme06eSAy1p5fs8suHUUKCc8Sh1RbgHI
qdTZKbKBvcEmyklQJrGdXpPOLBQkTh+NkIzs0zk5h08qOLv/t2y4QUvBJ69/L1HW
2wIDAQAB
-----END 公钥-----`
	pri := `-----BEGIN 私钥-----
MIIEpAIBAAKCAQEA4XXxm9EAsOKVX/YgMHSMAM44gyPx8PdkXPeh0de94qTS2kyv
FG5RJYSHI3xC+cyltnMzfFuDbjKB0/sfnfYPyNANFvzbVGnhEPvulsX+vbJdnuCV
cVI/H/cjOJGJ+MKTGQz55n+Idp/BsQhBc+9ub2PvdnFxng6zki3HFJ4Kq9onWjV4
smFe0JlqIqDg+xD8U+zLvlQdSDD7oNhBYAMLZn1uYbbh9Adoh7nDkqLdLLRsJqi5
lAftRme06eSAy1p5fs8suHUUKCc8Sh1RbgHIqdTZKbKBvcEmyklQJrGdXpPOLBQk
Th+NkIzs0zk5h08qOLv/t2y4QUvBJ69/L1HW2wIDAQABAoIBADm/Td0NEVI9Ful3
TxNaJqnmKA9e249OrkQpoSbwTVCJyv/i+E0RXHNxmHN8VmNJCYDLzPojLmzBPwOe
pKB+79gNgezLYxoh7GW/QYMYv4Cy+MvC1kQqPsTfNgQ9MsumeqrF6hQhwkAv9xpy
9ODPCMg7HpCFygYj2datJvqe6DjPsk3zGMPyFMRyHt/XGXb6SEP0i7CeUsDMnygC
2ZrjK86lTEB5D6F8rGT7WlYyyLn5J11uMML8y7aB9+ccwVE74rW1oPWQjbVXIyQn
eStSG3JYgzGUsEc68at3jWONZEQDvUHomRyxEjsW0spce3BHq7X2giVJVZYSnsJP
iIOmqUECgYEA/xqVb2P4ThhgxYB8to+jPuNHfrtKNfQfW9Z0VLO2k1bKbvtCr0X7
w2qgLKBALmONAAz/KPpOjNz/DZ00AqSfXT/b/ovtUDi5bjWlHlTt1AWFXIQuTgkC
tvgr1Enzn3pxa9rE4ryuuvl1tRYTMtYX24ugNzvb8j9dnJMfohyJxqsCgYEA4kCz
qyHGDciOfCiIIUxpkqSEPhQuLfvGtoEvtiuxLLmjJ1SPV8LvOCL9mXw8VABdHX0a
uQ7jNtF0Y+9/OeOMSo2zLYPk9Xi0KjSjrLnsiS2B4fDjP0Ry5q3kILTlKJo6lTrs
ciM0gcIsYz8ggiaw+neDjEaPyLsxoVkj77EY8JECgYEAuqdoz6gF9p4/sELi/XD7
sPf7R+8hzXhhuYCgfZlA7W2DkNCnajd6jvFlYUGftFGCyZa42/LJpqfMttlfRM/P
Cxu+i/E2IoxeoRT/S8I4gfnIKnlMqCxPoDDVYO77IqUkeBYKGRyVfJkqVuVgBsI/
kpQHFmLl+8oBZJ8BdkwLQyECgYEAnhc5s60wt3bY4LZtkF7VMesUoE/3iJfx3Jpe
HTtgXHEGKLg0RM1n2+DPNM0TVlq+tZkx7/cQGsC2RBIX4vo1j+59MaOEe2Uw9oC7
kTiEp8GNjLOGBjIs2zTMP3JG4V0K7DU0+/fPe4+S9nIoo+inJwwVdhHj7A2o+yXP
L2+ejpECgYBjqpzvBTOKOBlAsYI55JRd45lu7iZv/9napBiK8LN3pdZZPVceUQNU
JU75yOPwUGqUVW+8JDF1GUC6pjQCFZzZ7d0qdPD+vMHqM0S3hpX1sdqp2C9nql4F
IEpw7qCwdfxiqWu4Be/TihDoKXSgz+i95tgiJxDDApv3VS/wTW2cYA==
-----END 私钥-----`

	r := new(Rsa)
	r.Init([]byte(pub), []byte(pri))
	return r
}

//args 1 public_key args 2 private_key
func (self *Rsa) Init(key ...interface{}) error {
	if len(key) != 2 {
		return errors.New("Args error.")
	}
	pub, ok := key[0].([]byte)
	if !ok {
		return errors.New("Type error.")
	}
	pri, ok := key[1].([]byte)
	if !ok {
		return errors.New("Type error.")
	}
	block, _ := pem.Decode(pub)
	if block == nil {
		return errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	p := pubInterface.(*rsa.PublicKey)
	self.pub = p
	block, _ = pem.Decode(pri)
	if block == nil {
		return errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	self.pri = priv
	return nil
}

func (self *Rsa) Encrypt(Data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, self.pub, Data)
}

func (self *Rsa) Decrypt(data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, self.pri, data)
}
