package discovery

import (
	"gitee.com/monobytes/gcore/gerrors"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gregistry"
	"gitee.com/monobytes/gcore/gwrap/endpoint"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

type Resolver struct {
	builder     *Builder
	cc          resolver.ClientConn
	serviceName string
}

func newResolver(builder *Builder, serviceName string, cc resolver.ClientConn) *Resolver {
	return &Resolver{
		builder:     builder,
		cc:          cc,
		serviceName: serviceName,
	}
}

func (r *Resolver) updateInstances(instances []*gregistry.ServiceInstance) {
	state := resolver.State{Addresses: make([]resolver.Address, 0, len(instances))}
	for _, instance := range instances {
		exists := false

		for _, service := range instance.Services {
			if service == r.serviceName {
				exists = true
				break
			}
		}

		if !exists {
			continue
		}

		ep, err := endpoint.ParseEndpoint(instance.Endpoint)
		if err != nil {
			glog.Errorf("parse discovery endpoint failed: %v", err)
			continue
		}

		state.Addresses = append(state.Addresses, resolver.Address{
			Addr:       ep.Address(),
			ServerName: r.serviceName,
		})
	}

	if err := r.cc.UpdateState(state); err != nil {
		r.cc.ReportError(err)

		if !(len(state.Addresses) == 0 && gerrors.Is(err, balancer.ErrBadResolverState)) {
			glog.Errorf("update client conn state failed: %v", err)
		}
	}
}

func (r *Resolver) ResolveNow(_ resolver.ResolveNowOptions) {

}

// Close closes the resolver.
func (r *Resolver) Close() {
	r.builder.removeResolver(r.serviceName)
}
