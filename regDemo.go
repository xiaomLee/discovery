package main

import (
	"discovery/etcd"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	err := etcd.InitRegister([]string{
		"127.0.0.1:2379",
		"127.0.0.1:22379",
		"127.0.0.1:32379",
	})
	if err != nil {
		println(err.Error())
	}
	rand.Seed(time.Now().UnixNano())
	hostName := strconv.Itoa(rand.Intn(10))
	etcd.GetRegInstance().Register("/test-server/node/"+hostName, hostName)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-exit:
		etcd.GetRegInstance().Close()
	}
}
