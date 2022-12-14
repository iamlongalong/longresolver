package k8sresolver

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/resolver"
)

const colon = ":"

var defaultNamespace = "default"

func SetDefaultNameSpace(n string) {
	defaultNamespace = n
}

var emptyService Service

// Service represents a service with namespace, name and port.
type Service struct {
	Namespace string
	Name      string
	Port      int
}

// ParseTarget parses the resolver.Target.
func ParseTarget(target resolver.Target) (Service, error) {
	var service Service
	service.Namespace = target.Authority
	if len(service.Namespace) == 0 {
		service.Namespace = defaultNamespace
	}

	segs := strings.SplitN(target.Endpoint, colon, 2)
	if len(segs) < 2 {
		return emptyService, fmt.Errorf("bad endpoint: %s", target.Endpoint)
	}

	service.Name = segs[0]
	port, err := strconv.Atoi(segs[1])
	if err != nil {
		return emptyService, err
	}

	service.Port = port

	return service, nil
}
