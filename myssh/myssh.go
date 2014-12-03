// mkdir -p ~/.ssh
// cat <<"EOF" >> ~/.ssh/config
// Host *
//    StrictHostKeyChecking no
//    UserKnownHostsFile /dev/null
//    LogLevel ERROR
// EOF

package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "sync"
    "flag"
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

func query_pass(host string) string {
    pass := make(map[string]string)
    pass["dev"] = "QyBRKkxywna@$#BB"
    pass["192.168.119.164"] = pass["dev"]

    pass["192.168.12.48"] = "DBsQl*12.47.sErver"

    retval := ""
    p, ok := pass[host]; 

    if ok != false {
        retval = p
    } else {
        fmt.Println("not found!")
    }
    return retval
}

//func query_pass_from_ams(host string) string {
//    url_segs := []string{"https://ams.58corp.com/assets/device/get_data?_dc=1417600086833&idtype=&page=1&start=0&limit=25&sort=last_save_time&dir=DESC&query=", host, "&callback=Ext.data.JsonP.callback3"}
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

    host := flag.String("host", "dev", "a string")
    flag.Parse()

    fmt.Printf("%s\n", *host)
    pass := query_pass(*host)
    os.Setenv("SSHPASS", pass)
    command := []string{"sshpass -e ", "ssh ", *host, " w; hostname"}
    x := strings.Join(command, "")

    go exec_cmd(x, wg)

    wg.Wait()
}
