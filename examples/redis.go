package main

import (
	"fmt"
	"log"
	"time"
	"github.com/garyburd/redigo/redis"
	"github.com/youtube/vitess/go/pools"
)


type ResourceConn struct {
	redis.Conn
}

func (r ResourceConn) Close() {
	r.Conn.Close()
}


func main() {
	p := pools.NewResourcePool(func() (pools.Resource, error) {
		c, err := redis.Dial("tcp", ":6379")
		return ResourceConn{c}, err
	}, 1, 2, time.Minute)

	defer p.Close()

	r, err := p.Get()

	if err != nil {
		fmt.Printf("FATAL: %s\n", err)
		return
	}

	defer p.Put(r)

	c := r.(ResourceConn)
	n, err := c.Do("INFO")

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("INFO: info=%s\n", n)
}
