package tool

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// WriteFileA
func WriteFileA(filePath, contentA string) {
	if !CheckFileIsExist(filePath) {
		os.Create(filePath)
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(contentA)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

/*
 * 生成随机字符串
 */
func RandomString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
func printGreen(message string) {
	conf := 1  // 配置、终端默认设置
	bg := 32   // 背景色、终端默认设置
	text := 40 // 前景色、红色
	fmt.Printf("\n %c[%d;%d;%dm%s%c[0m\n\n", 0x1B, conf, bg, text, message, 0x1B)
}

// 请确认
func PleaseConfirm(message string) string {
	var choice string
	if message != "" {
		printGreen(message)
	} else {
		printGreen("This operation will overwrite the original data. Please operate with care. Confirm whether to continue!!!")
	}
	fmt.Println("Yes|NO")
	fmt.Scanln(&choice)
	return *&choice
}
