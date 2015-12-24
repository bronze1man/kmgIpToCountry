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
	"bytes"
	"encoding/binary"
	"compress/gzip"
)

//current directory should be ./src/github.com/bronze1man/kmgIpToCountry
func main() {
	var inputGzPath string
	var outputGoPath string
	flag.StringVar(&inputGzPath, "inputGzPath", "build/GeoLite2-Country.mmdb.gz", "input gz file path")
	flag.StringVar(&outputGoPath, "outputGoPath", "buildData.go", "output go file path")
	flag.Parse()
	gzContent, err := ioutil.ReadFile(inputGzPath)
	if err != nil {
		fmt.Println("read input file", err)
		return
	}
	r := bytes.NewReader(gzContent)
	gzReader, err := gzip.NewReader(r)
	if err != nil {
		panic(err)
	}
	uncompressData, err := ioutil.ReadAll(gzReader)
	if err != nil {
		panic(err)
	}
	//writeGoFileWithAsm(uncompressData)
	mustWriteFile(outputGoPath,getGoFileContentByteStringV2(uncompressData))
	/*
	err = ioutil.WriteFile(outputGoPath, getGoFileContentByteStringV2(content), os.FileMode(0644))
	if err != nil {
		fmt.Println("write output file", err)
		return
	}*/
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

//生成编译稍快,go文件体积稍大,ide打开go文件速度比较慢
func getGoFileContentByteStringV2(gzContent []byte) []byte {
	outString := make([]byte, 0, len(gzContent)*4+1024)
	outString = append(outString, []byte(`package kmgIpToCountry
//this file is generate by script, do not modify it by hand.
func getData()[]byte{
	return data
}
var data = []byte(`+fmt.Sprintf("%#v",string(gzContent))+`)`)...)
	return outString
}

func getGoFileContentRice(gzContent []byte) []byte {
	outString := make([]byte, 0, len(gzContent)*4+1024)
	outString = append(outString, []byte(`package kmgIpToCountry
//this file is generate by script, do not modify it by hand.
import "github.com/GeertJohan/go.rice"
func getData()[]byte{
	return rice.FindBox("build/GeoLite2-Country.mmdb.gz")
}`)...)
	return outString
}

func writeGoFileWithAsm(gzContent []byte){
	mustWriteFile("buildData.go",[]byte(`package kmgIpToCountry
//this file is generate by script, do not modify it by hand.
var my_data_string []byte
func getData()[]byte{
	return my_data_string
}`))
	size:=len(gzContent)/8*8+8
	thisBuf:=make([]byte,8)
	buf:=&bytes.Buffer{}
	buf.WriteString(`#include "textflag.h"
#define  g(g_a,g_b) DATA ·my_data_string_content+(g_a)(SB)/8, $(g_b)
`)
	for i:=0;i<len(gzContent);i+=8{
		end:=i+8
		if end>len(gzContent){
			end = len(gzContent)
		}
		copy(thisBuf,gzContent[i:end])
		num:=binary.LittleEndian.Uint64(thisBuf)
		binary.BigEndian.PutUint64(thisBuf,num)
		//buf.WriteString(`DATA ·my_data_string_content+`+strconv.Itoa(i)+`(SB)/8, $0x`+hex.EncodeToString(thisBuf)+"\n")
		buf.WriteString(`g(`+strconv.Itoa(i)+`,0x`+hex.EncodeToString(thisBuf)+")\n")
	}
	buf.WriteString(`GLOBL ·my_data_string_content(SB),RODATA,$`+strconv.Itoa(size)+"\n")
	buf.WriteString(`DATA ·my_data_string+0(SB)/8, $·my_data_string_content(SB)
DATA ·my_data_string+8(SB)/8, $`+strconv.Itoa(len(gzContent))+`
DATA ·my_data_string+16(SB)/8, $`+strconv.Itoa(len(gzContent))+`
GLOBL ·my_data_string(SB),RODATA,$24`+"\n")
	mustWriteFile("buildData.s",buf.Bytes())
}

func mustWriteFile(path string,content []byte){
	err := ioutil.WriteFile(path, content, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
	fmt.Println(path,len(content))
}