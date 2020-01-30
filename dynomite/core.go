package dynomite

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/majewsky/schwift"
	"github.com/majewsky/schwift/gopherschwift"
)

// Dynomite represents a dynomite instance
type Dynomite struct {
	Host    string
	Port    int16
	Backend *Redis
}

// NewDynomite returns a new instance of Dynomite
func NewDynomite(host string, port int16) *Dynomite {
	return &Dynomite{
		Host: host,
		Port: port,
	}
}

// NewDynomiteRedis creates a new Redis struct with initialized ConnectionPool
func NewDynomiteRedis(host string, port, backendPort int16, password string) *Dynomite {
	redis := NewRedis(host, backendPort, password)

	return &Dynomite{
		Host:    host,
		Port:    port,
		Backend: redis,
	}
}

func newObjectStoreAccount() (*schwift.Account, error) {
	var account *schwift.Account

	ao, err := clientconfig.AuthOptions(nil)
	if err != nil {
		return account, err
	}
	provider, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return account, err
	}
	err = openstack.Authenticate(provider, *ao)
	client, err := openstack.NewObjectStorageV1(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return account, err
	}
	account, err = gopherschwift.Wrap(client, nil)
	if err != nil {
		return account, err
	}

	return account, nil
}
