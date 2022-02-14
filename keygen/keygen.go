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

func generateKey(model string) []byte {
	// 根据常量rootKey 生成秘钥
	if model != "file" {
		key := []byte(rootKey)
		return key
	} else {
		key := []byte("bx2mfilesGUxZhA6")
		return key
	}

}

// =================== CBC ======================
func AesEncryptCBC(passStr, model string) (str string) {
	/*
		passStr: 需要加密的内容
		NewCipher该函数限制了输入k的长度必须为16, 24或者32
	*/
	key := generateKey(model)
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                             // 获取秘钥块的长度
	origData := pkcs5Padding([]byte(passStr), blockSize)       // 补全码
	blockMode := cipher.NewCBCEncrypter(block, []byte(rootIV)) // 加密模式
	encrypted := make([]byte, len(origData))                   // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                 // 加密
	strBase64 := base64.StdEncoding.EncodeToString(encrypted)
	return strBase64
}

// AesDecryptCBC
func AesDecryptCBC(encrypted, model string) (decryptedStr string) {
	/*
		encrypted: 需要解密的内容
		key: 秘钥键
		NewCipher该函数限制了输入k的长度必须为16, 24或者32
	*/

	key := generateKey(model)
	block, _ := aes.NewCipher(key)                             // 分组秘钥
	blockMode := cipher.NewCBCDecrypter(block, []byte(rootIV)) // 加密模式
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "decrypted error"
	}

	decrypted := make([]byte, len(encryptedBytes)) // 创建数组
	defer func() {
		if err := recover(); err != nil {
			decryptedStr = "decrypted error"
		}
	}()
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

// ReadBlock
func ReadBlock(filePth, dstPath, op string) {

	if tool.CheckFileIsExist(filePth) {
		content, _ := ioutil.ReadFile(filePth)
		if op == "Encrypt" {
			tool.WriteFileA(dstPath, AesEncryptCBC(string(content), "file"))
		} else if op == "Decrypt" {
			tool.WriteFileA(dstPath, AesDecryptCBC(string(content), "file"))
		}
	} else {
		logger.Error("no such file", filePth)
	}

}

// AesEncryptCBCFile
func AesEncryptCBCFile(srcPath, destPath string) {
	ReadBlock(srcPath, destPath, "Encrypt")
}

// AesDecryptCBCFile
func AesDecryptCBCFile(srcPath, destPath string) {
	ReadBlock(srcPath, destPath, "Decrypt")
}
