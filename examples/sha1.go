package main

import (
    "crypto/sha1"
    "crypto/md5"
    "fmt"
)


func Sha1sum(s string) string {
    sum := sha1.Sum([]byte(s))
    return fmt.Sprintf("%x", sum)
}


func Md5sum(s string) string {
    sum := md5.Sum([]byte(s))
    return fmt.Sprintf("%x", sum)
}


func main() {
    str := "12333"

	fmt.Println(Sha1sum(str))
	fmt.Println(Md5sum(str))
}
