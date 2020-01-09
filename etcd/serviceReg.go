package etcd

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type ServiceReg struct {
	endpoints []string
	client    *clientv3.Client
	lease     clientv3.Lease
	leaseId   clientv3.LeaseID
}

var regInstance *ServiceReg

func InitRegister(endpoints []string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}

	regInstance = &ServiceReg{
		endpoints: endpoints,
		client:    cli,
	}
	regInstance.setLease(5)
	regInstance.keepAlive()

	return nil
}

func GetRegInstance() *ServiceReg {
	return regInstance
}

//设置租约
func (s *ServiceReg) setLease(timeNum int64) error {
	lease := clientv3.NewLease(s.client)

	//设置租约时间
	leaseResp, err := lease.Grant(context.Background(), timeNum)
	if err != nil {
		return err
	}
	s.lease = lease
	s.leaseId = leaseResp.ID

	return nil
}

// 租约续期
func (s *ServiceReg) keepAlive() error {
	leaseRespChan, err := s.lease.KeepAlive(context.Background(), s.leaseId)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case resp := <-leaseRespChan:
				if resp == nil {
					fmt.Printf("已经关闭续租功能\n")
					return
				} else {
					fmt.Printf("续租成功\n")
				}
			}
		}
	}()

	return nil
}

//通过租约 注册服务
func (s *ServiceReg) Register(key string, info string) error {
	kv := clientv3.NewKV(s.client)
	_, err := kv.Put(context.Background(), key, info, clientv3.WithLease(s.leaseId))
	if err != nil {
		return err
	}
	return nil
}

// 通过租约 撤销服务
func (s *ServiceReg) Unregister() {
	println("unregister")
	s.lease.Revoke(context.Background(), s.leaseId)
}

func (s *ServiceReg) Close() {
	s.Unregister()
	s.lease.Close()
	s.client.Close()
}
