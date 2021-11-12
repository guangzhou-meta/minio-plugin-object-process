package main

import (
	"bytes"
	"fmt"
	"github.com/guangzhou-meta/minio-plugin-object-process/process"
	"io/fs"
	"io/ioutil"
)

func main() {
	buf, err := ioutil.ReadFile("test/img/example.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	reader, _, _ := process.ProcessObject(bytes.NewReader(buf), "image/sharpen,100")
	if reader == nil {
		fmt.Println("process error")
		return
	}
	buf, err = ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile("test/img/result.jpg", buf, fs.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
}
