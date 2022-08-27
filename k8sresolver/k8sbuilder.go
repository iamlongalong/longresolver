package k8sresolver

import (
	"fmt"
	"time"

	"github.com/iamlongalong/longresolver"
	"google.golang.org/grpc/resolver"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	resyncInterval = 5 * time.Minute
	nameSelector   = "metadata.name="

	KubernetesScheme = "k8s"

	subsetSize = 32
)

func init() {
	resolver.Register(&kubeBuilder{})
}

type kubeBuilder struct{}

func (b *kubeBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {

	svc, err := ParseTarget(target)
	if err != nil {
		return nil, err
	}
	longresolver.GetDfLogger().Infof("use kube builder : %+v : %+v", target, svc)

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	handler := NewEventHandler(func(endpoints []string) {
		var addrs []resolver.Address
		for _, ep := range endpoints {
			addrs = append(addrs, resolver.Address{
				Addr: fmt.Sprintf("%s:%d", ep, svc.Port),
			})
		}

		if err := cc.UpdateState(resolver.State{
			Addresses: addrs,
		}); err != nil {
			longresolver.GetDfLogger().Errorf(err.Error())
		}
	})
	inf := informers.NewSharedInformerFactoryWithOptions(cs, resyncInterval,
		informers.WithNamespace(svc.Namespace),
		informers.WithTweakListOptions(func(options *v1.ListOptions) {
			options.FieldSelector = nameSelector + svc.Name
		}))
	in := inf.Core().V1().Endpoints()
	in.Informer().AddEventHandler(handler)

	inf.Start(startRun())

	endpoints, err := cs.CoreV1().Endpoints(svc.Namespace).Get(svc.Name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	handler.Update(endpoints)

	return &nopResolver{cc: cc}, nil
}

func (b *kubeBuilder) Scheme() string {
	return KubernetesScheme
}
