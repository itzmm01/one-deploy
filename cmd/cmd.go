package cmd

import (
	"os/exec"
	"runtime"

	"github.com/wonderivan/logger"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func Run(command string) error {
	var result []byte
	var err error

	sysType := runtime.GOOS
	if sysType == "windows" {
		result, err = exec.Command("cmd", "/c", command).CombinedOutput()
		// logger.Error("no support system: ", sysType)
	} else if sysType == "linux" {
		result, err = exec.Command("/bin/bash", "-c", command).CombinedOutput()
	} else {
		logger.Error("no support system: ", sysType)
	}

	if err != nil {
		logger.Error("run cmd failed: ", command)
		logger.Error("run cmd failed: ", err, ConvertByte2String(result, "GB18030"))
	}
	return err
}

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
