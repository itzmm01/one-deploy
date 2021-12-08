package tool

import (
	"bufio"
	"fmt"
	"os"
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

func WriteFileA(filePath, contentA string) {
	if !CheckFileIsExist(filePath) {
		os.Create(filePath)
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
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
