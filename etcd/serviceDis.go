package etcd

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type ServiceDis struct {
	endpoints []string
	client    *clientv3.Client

	exit chan bool
	sync.WaitGroup
}

var disInstance *ServiceDis

func InitDiscovery(endpoints []string) error {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}

	disInstance = &ServiceDis{
		endpoints: endpoints,
		client:    c,
		exit:      make(chan bool),
	}

	return nil
}

func GetDisInstance() *ServiceDis {
	return disInstance
}

func (s *ServiceDis) Close() {
	close(s.exit)
	s.Wait()
	s.client.Close()
}

func DefaultWatchHandle(s *ServiceDis, ch clientv3.WatchChan) {
	for {
		select {
		case wresp := <-ch:
			for _, ev := range wresp.Events {
				switch ev.Type {
				case clientv3.EventTypePut:
					fmt.Printf("PUT:%+v \n", ev.Kv)
				case clientv3.EventTypeDelete:
					fmt.Printf("DEL:%+v \n", ev.Kv)
				}
			}
		case <-s.exit:
			return
		}
	}
}

func (s *ServiceDis) Watch(prefix string, handle func(s *ServiceDis, ch clientv3.WatchChan)) {
	ch := s.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	go func() {
		s.Add(1)
		defer s.Done()
		handle(s, ch)
	}()
}

func (s *ServiceDis) GetService(prefix string) (map[string]string, error) {
	resp, err := s.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	nodes := make(map[string]string, 0)
	if resp == nil || resp.Kvs == nil {
		return nil, errors.New("service is null")
	}
	for _, item := range resp.Kvs {
		if v := item.Value; v != nil {
			nodes[string(item.Key)] = string(item.Value)
		}
	}

	return nodes, nil
}
