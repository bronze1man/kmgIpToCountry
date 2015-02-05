package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	//"io"
	"encoding/hex"
	"io/ioutil"
	"os"
	"strconv"
	//"bytes"
)

//current directory should be ./src/github.com/bronze1man/kmgIpToCountry
func main() {
	var inputGzPath string
	var outputGoPath string
	flag.StringVar(&inputGzPath, "inputGzPath", "build/GeoLite2-Country.mmdb.gz", "input gz file path")
	flag.StringVar(&outputGoPath, "outputGoPath", "buildData.go", "output go file path")
	flag.Parse()
	content, err := ioutil.ReadFile(inputGzPath)
	if err != nil {
		fmt.Println("read input file", err)
		return
	}
	err = ioutil.WriteFile(outputGoPath, getGoFileContentB64(content), os.FileMode(0644))
	if err != nil {
		fmt.Println("write output file", err)
		return
	}
	return
}

//生成编译最快,go文件体积最小,ide打开go文件速度最快.运行时读取数据稍慢
func getGoFileContentB64(gzContent []byte) []byte {
	outString := base64.StdEncoding.EncodeToString(gzContent)
	return []byte(`package kmgIpToCountry
//this file is generate by script, do not modify it by hand.
import "encoding/base64"
func getData()[]byte{
	gzContent, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			panic(err)
		}
		return gzContent
}
var data = ` + "`" + outString + "`")
}

//生成编译非常慢,go文件体积大,ide打开go文件速度非常慢.运行时读取数据最快
func getGoFileContentByte(gzContent []byte) []byte {
	outString := make([]byte, 0, len(gzContent)*4+1024)
	outString = append(outString, []byte(`package kmgIpToCountry
//this file is generate by script, do not modify it by hand.
func getData()[]byte{
	return data
}
var data = []byte{`)...)
	for _, b := range gzContent {
		outString = append(outString, []byte(strconv.Itoa(int(b))+",")...)
	}
	if len(gzContent) > 0 {
		outString = outString[:len(outString)-1]
	}
	outString = append(outString, []byte(`}`)...)
	return outString
}

//生成编译稍快,go文件体积稍大,ide打开go文件速度比较慢
func getGoFileContentByteString(gzContent []byte) []byte {
	outString := make([]byte, 0, len(gzContent)*4+1024)
	outString = append(outString, []byte(`package kmgIpToCountry
//this file is generate by script, do not modify it by hand.
func getData()[]byte{
	return data
}
var data = []byte("`)...)
	for i := range gzContent {
		outString = append(outString, []byte(`\x`)...)
		outString = append(outString, []byte(hex.EncodeToString(gzContent[i:i+1]))...)
	}
	outString = append(outString, []byte(`")`)...)
	return outString
}
