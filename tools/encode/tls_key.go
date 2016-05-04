package encode

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

//func main() {
//	info := CertInformation{Country: []string{"中国"}, Organization: []string{"游戏蜗牛"},
//		OrganizationalUnit: []string{"blog.csdn.net/fyxichen"}, EmailAddress: []string{"czxichen@163.com"},
//		StreetAddress: []string{"中心大道西171"}, Province: []string{"江苏省·工业园区"}, SubjectKeyId: []byte{1, 2, 3, 4, 5, 6},
//		Certificate: "client.pem", PrivateKey: "client.key", ROOTCertificate: "server.pem", ROOTPrivateKey: "server.key"}
//	err := CreateCerts(info)
//	if err != nil {
//		println(err.Error())
//	}
//}

type CertInformation struct {
	Country            []string
	Organization       []string
	OrganizationalUnit []string
	EmailAddress       []string
	Province           []string
	StreetAddress      []string
	SubjectKeyId       []byte
	Certificate        string
	PrivateKey         string
	ROOTCertificate    string
	ROOTPrivateKey     string
}

func CreateCerts(info CertInformation) error {

	var rootPrivateKey *rsa.PrivateKey
	var rootcertificate *x509.Certificate
	var err error
	ca := newCertificate(info)
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)

	//读取根证书
	if info.Certificate != "" && info.PrivateKey != "" {
		rootcertificate, rootPrivateKey, err = parseCerts(info.ROOTCertificate, info.ROOTPrivateKey)
		if os.IsNotExist(err) {
			rootPrivateKey, _ = rsa.GenerateKey(rand.Reader, 2048)
			rootcertificate = ca
			ca_b, err := x509.CreateCertificate(rand.Reader, rootcertificate, rootcertificate, &rootPrivateKey.PublicKey, rootPrivateKey)
			if err != nil {
				return err
			}
			err = write(info.ROOTCertificate, "CERTIFICATE", ca_b)
			if err != nil {
				return err
			}
			priv_b := x509.MarshalPKCS1PrivateKey(rootPrivateKey)
			err = write(info.ROOTPrivateKey, "PRIVATE KEY", priv_b)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	} else {
		rootcertificate = ca
		rootPrivateKey = priv
	}

	ca_b, err := x509.CreateCertificate(rand.Reader, ca, rootcertificate, &priv.PublicKey, rootPrivateKey)
	if err != nil {
		return err
	}
	err = write(info.Certificate, "CERTIFICATE", ca_b)
	if err != nil {
		return err
	}

	priv_b := x509.MarshalPKCS1PrivateKey(priv)
	err = write(info.PrivateKey, "PRIVATE KEY", priv_b)
	if err != nil {
		return err
	}

	return nil
}

func write(filename, Type string, p []byte) error {
	File, err := os.Create(filename)
	defer File.Close()
	if err != nil {
		return err
	}
	var b *pem.Block = &pem.Block{Bytes: p, Type: Type}
	err = pem.Encode(File, b)
	if err != nil {
		return err
	}
	return nil
}

func parseCerts(ROOTCertificate, ROOTPrivateKey string) (*x509.Certificate, *rsa.PrivateKey, error) {
	var rootPrivateKey *rsa.PrivateKey
	var rootcertificate *x509.Certificate

	buf, err := ioutil.ReadFile(ROOTCertificate)
	if err != nil {
		return nil, nil, err
	}
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	rootcertificate, err = x509.ParseCertificate(p.Bytes)
	if err != nil {
		return nil, nil, err
	}
	buf, err = ioutil.ReadFile(ROOTPrivateKey)
	if err != nil {
		return nil, nil, err
	}
	p, buf = pem.Decode(buf)
	rootPrivateKey, err = x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return nil, nil, err
	}
	return rootcertificate, rootPrivateKey, nil
}

func newCertificate(info CertInformation) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Country:            info.Country,
			Organization:       info.Organization,
			OrganizationalUnit: info.OrganizationalUnit,
			Province:           info.Province,
			StreetAddress:      info.StreetAddress,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          info.SubjectKeyId,
		BasicConstraintsValid: true,
		IsCA:           true,
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:       x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		EmailAddresses: info.EmailAddress,
	}
}

func ParseCert(crt string) ([]byte, error) {
	buf, err := ioutil.ReadFile(crt)
	if err != nil {
		return nil, err
	}
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	return p.Bytes, nil
}
