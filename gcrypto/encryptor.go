package gcrypto

import (
	"gitee.com/monobytes/gcore/glog"
)

type Encryptor interface {
	// Name 名称
	Name() string
	// Encrypt 加密
	Encrypt(data []byte) ([]byte, error)
	// Decrypt 解密
	Decrypt(data []byte) ([]byte, error)
}

var encryptors = make(map[string]Encryptor)

// RegisterEncryptor 注册加密器
func RegisterEncryptor(encryptor Encryptor) {
	if encryptor == nil {
		glog.Fatal("can't register a invalid encryptor")
	}

	name := encryptor.Name()

	if name == "" {
		glog.Fatal("can't register a encryptor without name")
	}

	if _, ok := encryptors[name]; ok {
		glog.Warnf("the old %s encryptor will be overwritten", name)
	}

	encryptors[name] = encryptor
}

// InvokeEncryptor 调用加密器
func InvokeEncryptor(name string) Encryptor {
	encryptor, ok := encryptors[name]
	if !ok {
		glog.Fatalf("%s encryptor is not registered", name)
	}

	return encryptor
}
