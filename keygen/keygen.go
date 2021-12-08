package keygen

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io/ioutil"
	"one-backup/tool"

	"github.com/wonderivan/logger"
)

const rootKey = "bx2maw66sGUxZhA6"
const rootIV = "4qnDK5ZITgVuegUR"

func generateKey() []byte {
	// 根据常量rootKey 生成秘钥
	key := []byte(rootKey)
	return key

}

// =================== CBC ======================
func AesEncryptCBC(passStr string) (str string) {
	/*
		passStr: 需要加密的内容
		NewCipher该函数限制了输入k的长度必须为16, 24或者32
	*/
	key := generateKey()
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                             // 获取秘钥块的长度
	origData := pkcs5Padding([]byte(passStr), blockSize)       // 补全码
	blockMode := cipher.NewCBCEncrypter(block, []byte(rootIV)) // 加密模式
	encrypted := make([]byte, len(origData))                   // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                 // 加密
	strBase64 := base64.StdEncoding.EncodeToString(encrypted)
	return strBase64
}

func AesDecryptCBC(encrypted string) (decryptedStr string) {
	/*
		encrypted: 需要解密的内容
		key: 秘钥键
		NewCipher该函数限制了输入k的长度必须为16, 24或者32
	*/
	key := generateKey()
	block, _ := aes.NewCipher(key)                             // 分组秘钥
	blockMode := cipher.NewCBCDecrypter(block, []byte(rootIV)) // 加密模式
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "base64 error"
	}
	decrypted := make([]byte, len(encryptedBytes))   // 创建数组
	blockMode.CryptBlocks(decrypted, encryptedBytes) // 解密
	decrypted = pkcs5UnPadding(decrypted)            // 去除补全码
	decryptedStr = string(decrypted)

	return decryptedStr
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func ReadBlock(filePth, dstPath, op string) {

	if tool.CheckFileIsExist(filePth) {
		content, _ := ioutil.ReadFile(filePth)
		if op == "Encrypt" {
			tool.WriteFileA(dstPath, AesEncryptCBC(string(content)))
		} else if op == "Decrypt" {
			tool.WriteFileA(dstPath, AesDecryptCBC(string(content)))
		}
	} else {
		logger.Error("no such file", filePth)
	}

}

func AesEncryptCBCFile(srcPath, destPath string) {
	ReadBlock(srcPath, destPath, "Encrypt")
}

func AesDecryptCBCFile(srcPath, destPath string) {
	ReadBlock(srcPath, destPath, "Decrypt")
}
