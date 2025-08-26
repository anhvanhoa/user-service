package discovery

import (
	loggerI "auth-service/domain/service/logger"
	"fmt"
	"net"

	"github.com/hashicorp/consul/api"
)

type Discovery interface {
	Register(serviceName string, servicePort int) error
	GetService(serviceName string) (string, error)
	Close(id string) error
}

type DiscoveryConfig struct {
	ServiceName    string
	ServiceAddress string
	ServicePort    int
	ServiceHost    string
	ServiceID      string
	ServiceTags    []string
	ServiceMeta    map[string]string
}

type DiscoveryImpl struct {
	client *api.Client
	log    loggerI.Log
}

func NewDiscovery(log loggerI.Log) (*DiscoveryImpl, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &DiscoveryImpl{
		client: client,
		log:    log,
	}, nil
}

func (d *DiscoveryImpl) Register(serviceName string, servicePort int) error {
	registration := new(api.AgentServiceRegistration)
	addrs, _ := net.InterfaceAddrs()
	serviceAddress := "127.0.0.1"
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			serviceAddress = ipnet.IP.String()
		}
	}
	registration.Name = serviceName
	registration.Port = servicePort
	registration.Address = serviceAddress

	registration.Check = &api.AgentServiceCheck{
		GRPC:     fmt.Sprintf("%s:%d", serviceAddress, servicePort),
		Interval: "10s",
		Timeout:  "3s",
	}

	if err := d.client.Agent().ServiceRegister(registration); err != nil {
		d.log.Error("Failed to register service", err)
		return err
	}
	d.log.Info("Service registered", "serviceName", serviceName, "serviceAddress", serviceAddress, "servicePort", servicePort)
	return nil
}

func (d *DiscoveryImpl) GetService(serviceName string) (string, error) {
	services, err := d.client.Agent().Services()
	if err != nil {
		d.log.Error("Failed to get service", err)
		return "", err
	}
	for _, service := range services {
		if service.Service == serviceName {
			return service.Address, nil
		}
	}
	return "", fmt.Errorf("service %s not found", serviceName)
}

func (d *DiscoveryImpl) Close(id string) error {
	return d.client.Agent().ServiceDeregister(id)
}
