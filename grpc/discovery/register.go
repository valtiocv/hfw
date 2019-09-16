package discovery

import (
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/hsyan2008/go-logger"
	"github.com/hsyan2008/hfw/configs"
	"github.com/hsyan2008/hfw/grpc/discovery/register"
	"github.com/hsyan2008/hfw/grpc/discovery/resolver"
)

func RegisterServer(cc configs.ServerConfig, address string) (r register.Register, err error) {
	if cc.ResolverType == "" || len(cc.ResolverAddresses) == 0 || cc.ServerName == "" {
		logger.Mix("ResolverType or ResolverAddresses or ServerName is empty, so do not Registered")
		return nil, nil
	}
	if cc.UpdateInterval == 0 {
		cc.UpdateInterval = 10
	}
	host, port, err := getHostPort(cc, address)
	if err != nil {
		return nil, err
	}
	switch cc.ResolverType {
	case resolver.ConsulResolver:
		r = register.NewConsulRegister(cc.ResolverAddresses[0], int(cc.UpdateInterval)*2)
	case resolver.EtcdResolver:
		r = register.NewEtcdRegister(cc.ResolverAddresses, int(cc.UpdateInterval)*2)
	default:
		return nil, errors.New("unsupport ResolverType")
	}
	err = r.Register(register.RegisterInfo{
		Host:           host,
		Port:           port,
		ServiceName:    cc.ServerName,
		UpdateInterval: cc.UpdateInterval * time.Second,
	})
	return r, err
}

func getHostPort(cc configs.ServerConfig, address string) (host string, port int, err error) {
	var p string
	host, p, err = net.SplitHostPort(address)
	if err != nil {
		return
	}

	if cc.Interface != "" {
		var iface *net.Interface
		iface, err = net.InterfaceByName(cc.Interface)
		if err != nil {
			return
		}
		var addrs []net.Addr
		addrs, err = iface.Addrs()
		if err != nil {
			return
		}

		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				host = ipnet.IP.String()
				break
			}
		}
	}

	port, err = strconv.Atoi(p)

	return
}
