package etcd

import (
	"time"

	"go.etcd.io/etcd/clientv3"
)

type Service struct {
	endpoints []string
	client    *clientv3.Client
}

var instance *Service

func Init(endpoints []string) error {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		return err
	}

	instance = &Service{
		endpoints: endpoints,
		client:    c,
	}

	return nil
}

func GetInstance() *Service {
	return instance
}

func (s *Service) Register(key string, info interface{}) {

	go s.keepAlive()
}

func (s *Service) Unregister(key string) {

}

func (s *Service) keepAlive() {

}

func (s *Service) Watch(key string) {

}

func (s *Service) GetService(Key string) interface{} {

	return nil
}

func (s *Service) GetServiceList() interface{} {

	return nil
}
