package main

import (
	"discovery/etcd"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := etcd.InitDiscovery([]string{
		"127.0.0.1:2379",
		"127.0.0.1:22379",
		"127.0.0.1:32379",
	})
	if err != nil {
		println(err.Error())
	}

	list, err := etcd.GetDisInstance().GetService("/test-server/node")
	if err != nil {
		println(err.Error())
	}
	fmt.Println(list)
	etcd.GetDisInstance().Watch("/test-server/node", etcd.DefaultWatchHandle)
	//etcd.GetDisInstance().Watch("/foo/test", etcd.DefaultWatchHandle)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-exit:
		etcd.GetDisInstance().Close()
	}
}
