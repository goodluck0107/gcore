package ecc

import (
	"github.com/goodluck0107/gcore/getc"
	"github.com/goodluck0107/gcore/gwrap/hash"
	"strings"
)

const (
	defaultSignerHashKey       = "etc.crypto.rsa.signer.hash"
	defaultSignerDelimiterKey  = "etc.crypto.rsa.signer.delimiter"
	defaultSignerPublicKeyKey  = "etc.crypto.rsa.signer.publicKey"
	defaultSignerPrivateKeyKey = "etc.crypto.rsa.signer.privateKey"
)

type SignerOption func(o *signerOptions)

type signerOptions struct {
	// hash算法。支持sha1、sha224、sha256、sha384、sha512
	// 默认为sha256
	hash hash.Hash

	// 签名分隔符。
	delimiter string

	// 公钥。可设置文件路径或公钥串
	publicKey string

	// 私钥。可设置文件路径或私钥串
	privateKey string
}

func defaultSignerOptions() *signerOptions {
	return &signerOptions{
		hash:       hash.Hash(strings.ToLower(getc.Get(defaultSignerHashKey).String())),
		delimiter:  getc.Get(defaultSignerDelimiterKey, " ").String(),
		publicKey:  getc.Get(defaultSignerPublicKeyKey).String(),
		privateKey: getc.Get(defaultSignerPrivateKeyKey).String(),
	}
}

// WithSignerHash 设置加密hash算法
func WithSignerHash(hash hash.Hash) SignerOption {
	return func(o *signerOptions) { o.hash = hash }
}

// WithSignerDelimiter 设置签名分割符
func WithSignerDelimiter(delimiter string) SignerOption {
	return func(o *signerOptions) { o.delimiter = delimiter }
}

// WithSignerPublicKey 设置验签公钥
func WithSignerPublicKey(publicKey string) SignerOption {
	return func(o *signerOptions) { o.publicKey = publicKey }
}

// WithSignerPrivateKey 设置解密私钥
func WithSignerPrivateKey(privateKey string) SignerOption {
	return func(o *signerOptions) { o.privateKey = privateKey }
}
