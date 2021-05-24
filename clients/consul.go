package clients

import (
	"github.com/hashicorp/consul/api"
	"log"
)

type Consul struct {
	ptr *api.Client
}

func NewConsul() *Consul {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal("Init consul client failed")
	}
	return &Consul{
		ptr: client,
	}
}

func (c *Consul) Register(name string) {
	if err := c.ptr.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name: name,
	}); err != nil {
		log.Fatalf("Register service: %s failed", name)
	}
}

func (c *Consul) Unregister(name string) {
	_ = c.ptr.Agent().ServiceDeregister(name)
}
