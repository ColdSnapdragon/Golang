package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func work(path string) string {
	direntrys, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, de := range direntrys {
		childpath := filepath.Join(path, de.Name())
		if de.IsDir() {
			if res := work(childpath); res != "" {
				return res
			}
		} else if de.Name() == "源.cpp" {
			return childpath
		}
	}
	return ""
}

func main() {
	source := os.Args[1]
	target := os.Args[2]

	des, _ := os.ReadDir(source)
	for _, de := range des {
		tar := filepath.Join(target, de.Name())
		cpp := work(filepath.Join(source, de.Name()))
		if cpp == "" {
			continue
		}
		src, _ := os.Open(cpp)
		os.MkdirAll(tar, 0755)
		dst, _ := os.OpenFile(filepath.Join(tar, "main.cpp"), os.O_RDWR|os.O_CREATE, 0755)
		io.Copy(dst, src)
		fmt.Println("copy", cpp, "to", tar)
		_ = src.Close()
		_ = dst.Close()
	}

}
