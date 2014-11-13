package main

import (
    "bufio"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "log"
    "net"
    "runtime"
    "strconv"
    "strings"
    "syscall"
    "time"
)

import (
    "github.com/Shopify/sarama"
    "github.com/stathat/jconfig"
)

func makeConfig() *sarama.ProducerConfig {
    config := sarama.NewProducerConfig()
    config.RequiredAcks = 0
    config.MaxBufferTime = uint32(CONFIG.GetInt("max_buffer_time"))
    config.MaxBufferedBytes = uint32(CONFIG.GetInt("max_buffered_bytes"))
    return config
}

func makeProducer(addr []string) (producer *sarama.Producer, err error) {

    client, err := sarama.NewClient("client_id", addr, nil)
    if err != nil {
        return
    }

    producer, err = sarama.NewProducer(client, makeConfig())
    if err != nil {
        return
    }

    defer producer.Close()

    return
}

func publish(producer *sarama.Producer, topic string, msg string) error {
    return producer.QueueMessage(topic, nil, sarama.StringEncoder(msg))
}

func selectProducer(msg map[string]interface{}) (p *sarama.Producer, err error) {

    cs := CONFIG.GetArray("clusters")
    lb := CONFIG.GetString("load_balancer")

    for i := 0; i < len(cs); i++ {
        c := cs[i].(map[string]interface{})

        target, ok := c[lb].([]interface{})
        if !ok {
            continue
        }

        for j := 0; j < len(target); j++ {
            t := target[j].(string)
            idc := msg["idc"]
            topic := msg["topic"]
            if t == idc.(string) || t == topic.(string) {
                p = PRODUCERS[c["name"].(string)]
                //      log.Println("idc=", idc, ",topic=", topic, ",lb=", lb, ",cluster=", c["name"], topic, 
",lb=", lb, ",cluster=", c["name"].(string))
                return
            }
        }
    }

    if p == nil {
        for _, v := range PRODUCERS {
            p = v
            break
        }
    }

    return
}

func verify(s []byte) (producer *sarama.Producer, topic string, msg string, err error) {

    var j map[string]interface{}
    err = json.Unmarshal(s, &j)
    if err != nil {
        return
    }

    product, ok := j["product"].(string)
    if !ok {
        err = errors.New("product not find in message")
        return
    }
    service, ok := j["service"].(string)
    if !ok {
        err = errors.New("service not find in message")
        return
    }

    topic = product + "_" + service + "_topic"

    j["topic"] = topic

    syscall.Gettimeofday(&TV)
    j["collector_time"] = (int64(TV.Sec)*1e3 + int64(TV.Usec)/1e3)

    producer, err = selectProducer(j)
    if nil != err {
        return
    }

    msgInByte, err := json.Marshal(j)

    msg = string(msgInByte)

    return
}

func handle(conn net.Conn, idx uint) {
    defer conn.Close()

    if tcpConn, ok := conn.(*net.TCPConn); ok {
        tcpConn.SetKeepAlive(true)
    }

    reader := bufio.NewReader(conn)

    state := 0

    for {
        if tcpConn, ok := conn.(*net.TCPConn); ok {
            tcpConn.SetReadDeadline(time.Now().Add(20 * time.Second))
            tcpConn.SetWriteDeadline(time.Now().Add(20 * time.Second))
        }

        txn, err := reader.ReadString(' ')
        if err != nil {
            log.Println(idx, err)
            log.Println(idx, "closed")
            return
        }
        txn = strings.TrimSpace(txn)

        cmd, err := reader.ReadString(' ')
        if err != nil {
            log.Println(idx, err)
            log.Println(idx, "closed")
            return
        }
        cmd = strings.TrimSpace(cmd)

        // TODO: handle 0 datalen -- loop on bytes until non-digit?
        dataLenString, err := reader.ReadString(' ')
        if err != nil {
            log.Println(idx, err)
            log.Println(idx, "closed")
            return
        }
        dataLen, err := strconv.Atoi(strings.TrimSpace(dataLenString))
        if err != nil {
            log.Println(idx, err)
            log.Println(idx, "closed")
            return
        }

        dataBytes := make([]byte, dataLen)
        _, err = io.ReadFull(reader, dataBytes)
        if err != nil {
            log.Println(idx, err)
            log.Println(idx, "closed")
            return
        }

        switch cmd {
        case "open":
            _, err := conn.Write([]byte(fmt.Sprintf("%s rsp 96 200 
OK\nrelp_version=0\nrelp_software=relp2kafka,1.0.0,http://aqueducts.baidu.com\ncommands=syslog\n", txn)))
            if err != nil {
                log.Println(idx, err)
                log.Println(idx, "closed")
                return
            }
            state = 1
        default:
            producer, topic, msg, err := verify(dataBytes)
            if err != nil {
                state = 0
            } else {
                go publish(producer, topic, msg)
            }

            if state != 1 {
                _, err := conn.Write([]byte(fmt.Sprintf("%s rsp 7 500 ERR\n", txn)))
                if err != nil {
                    log.Println(idx, err)
                    log.Println(idx, "closed")
                }
                return
            } else {
                _, err := conn.Write([]byte(fmt.Sprintf("%s rsp 6 200 OK\n", txn)))
                if err != nil {
                    log.Println(idx, err)
                    log.Println(idx, "closed")
                    return
                }
            }
        }
    }
    log.Println(idx, "closed")
}

func initProducers(config *jconfig.Config) error {

    PRODUCERS = make(map[string]*sarama.Producer)

    cs := config.GetArray("clusters")

    for i := 0; i < len(cs); i++ {
        c := cs[i].(map[string]interface{})

        tmp_broker_list, ok := c["broker_list"].([]interface{})
        if !ok {
            return errors.New("broker_list not find in config")
        }

        broker_list := make([]string, 5)
        for j := 0; j < len(tmp_broker_list); j++ {
            broker_list = append(broker_list, tmp_broker_list[j].(string))
        }

        if len(broker_list) < 1 {
            return errors.New("broker_list is null")
        }

        p, e := makeProducer(broker_list)
        if nil != e {
            log.Fatal(e)
            return e
        } else {
            name, ok := c["name"].(string)
            if ok {
                PRODUCERS[name] = p
            } else {
                return errors.New("cluster name is not provided")
            }
        }
    }
    return nil
}

// gloal variable
var PRODUCERS map[string]*sarama.Producer
var TV syscall.Timeval
var CONFIG *jconfig.Config

func main() {

    const Compiler = "gc"
    const GOARCH string = "amd64"
    const GOOS string = "linux"

    runtime.GOMAXPROCS(runtime.NumCPU())

    CONFIG = jconfig.LoadConfig("default.json")

    port := CONFIG.GetString("port")
    if port == "" {
        port = "8090"
    }

    ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
    if err != nil {
        log.Fatal(err)
    }

    err = initProducers(CONFIG)
    if err != nil {
        log.Println(err)
        return
    }

    var idx uint = 0
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Println(err)
            continue
        }
        log.Println("got a connection:", idx)
        go handle(conn, idx)
        idx++
    }
}
