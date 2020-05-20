package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
	"strings"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func main() {
	command := "netsh"
	params := []string{"wlan","show","profiles"}
	cmd := exec.Command(command, params...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd.Start()
	in := bufio.NewScanner(stdout)
	count := 0
	for in.Scan() {
		cmdRe:=ConvertByte2String(in.Bytes(),"GB18030")
		if strings.Contains(cmdRe,"所有用户配置文件") {
			count++
			index := strings.Index(cmdRe,":")
			getKey(cmdRe[index+2:])
		}
	}
	cmd.Wait()
	if count == 0 {
		fmt.Println("此电脑未连过wifi")
	}
	fmt.Scanf("按任意键结束")
}

func getKey(name string) {
	command := "netsh"
	params := []string{"wlan","show","profiles","name=",name,"key=clear"}
	cmd := exec.Command(command, params...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd.Start()
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		cmdRe:=ConvertByte2String(in.Bytes(),"GB18030")
		if strings.Contains(cmdRe,"关键内容") {
			index := strings.Index(cmdRe,":")
			st := cmdRe[index+2:]
			fmt.Println(name,":",st)

		}
	}
	cmd.Wait()
	return
}
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes,_=simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str= string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}