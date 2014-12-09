// mkdir -p ~/.ssh
// cat <<"EOF" >> ~/.ssh/config
// Host *
//    StrictHostKeyChecking no
//    UserKnownHostsFile /dev/null
//    LogLevel ERROR
// EOF

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	//"net/http"
	//"github.com/astaxie/beego/httplib"
)

func exec_cmd(cmd string, wg *sync.WaitGroup) {
	fmt.Println("command is ", cmd)

	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("%s", out)
	wg.Done() // Need to signal to waitgroup that this goroutine is done
}

func query_pass(host string, user string) string {
	pass := make(map[string]string)
	pass["root@dev"] = "QyBRKkxywna@$#BB"
	pass["root@192.168.119.164"] = pass["dev"]

	pass["root@192.168.12.48"] = "DBsQl*12.47.sErver"
	pass["root@10.5.9.28"] = "XhXkN.c^-lGv"
	pass["root@10.5.17.68"] = "L9llg*(oH0luUoja"
	pass["root@192.168.16.4"] = "DhJ_DNs^58v5#16.4Svr"
	pass["root@10.58.119.200"] = "4spd9ZhzzUhE9$4+"

	retval := ""

	host_and_user_segs := []string{user, "@", host}
	host_and_user := strings.Join(host_and_user_segs, "")

	p, ok := pass[host_and_user]

	if ok != false {
		retval = p
	} else {
		fmt.Println("not found!")
	}
	return retval
}

//func query_pass_from_ams(host string) string {
//    url_segs := []string{""}
//    url := strings.Join(url_segs, "")
//
//    cookie1 := &http.Cookie{}
//    cookie1.Name = "Sso_Username"
//    cookie1.Value="xiahoufeng"
//
//    cookie2 := &http.Cookie{}
//    cookie2.Name = "Sso_UserID"
//    cookie2.Value = "10232"
//
//    str, err := httplib.Get(url).SetCookie(cookie1).SetCookie(cookie2).SetProtocolVersion("HTTP/1.1").String()
//
//    if err != nil {
//        fmt.Println(str)
//        return str
//    }
//
//    fmt.Println(str)
//    return str
//}

// 获取密码，打印出来
func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	//pass_from_ams := query_pass_from_ams("192.168.119.164")
	//fmt.Println(pass_from_ams)

	// 第二个参数是默认的
	host := flag.String("host", "dev", "a string")
	cmd := flag.String("run", "w; hostname", "a string")
	user := flag.String("user", "work", "a string")

	flag.Parse()

	fmt.Printf("remote host is %s\n", *host)

	pass := query_pass(*host, *user)

	os.Setenv("SSHPASS", pass)
	command := []string{"sshpass -e ", "ssh ", *user, "@", *host, " ", *cmd}
	x := strings.Join(command, "")

	go exec_cmd(x, wg)

	wg.Wait()
}
