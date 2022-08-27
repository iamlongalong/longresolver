# k8sresolver

这是用于 grpc 在 k8s 下的 resolver， 使用时，可参考：
## resolver 初始化
```go
package k8sresolver

import (
	"fmt"

	"github.com/iamlongalong/longresolver/k8sresolver"
)

func Startup() error {
	k8sresolver.SetFinishCallback(func() {
		log.Info("k8s resolver finished callback")
	})

	k8sresolver.SetLogger(&K8sLogger{})

	log.Info("k8s resolver registered")

	return nil
}

type K8sLogger struct {}
// …… 自己实现 logger

```

## 建立连接时
注意，为了达到能够均匀负载到不同pod上，需要使用合适的负载均衡策略，建议使用 roundrobin
```go
import "google.golang.org/grpc/balancer/roundrobin"

var err error
conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithChainUnaryInterceptor(withTracing), grpc.WithBalancerName(roundrobin.Name))
if err != nil {
    log.Fatal(fmt.Sprintf("connect rbac service api failed: %s", err))
	return nil
}
```

## 配置中的 addr
为了能够真正使用 k8s 中的监听 endpoints 变化机制，addr 有特定的格式：

k8s://namespace/service:port

例如：default 下的 rbac 的 7500 端口可以写成：
`k8s://default/rbac:7500` 或者省略 default ，可以写成 `k8s:///rbac:7500`

特别注意 ！！ 一定`不能`写成 `rbac.default:7500` !!! 这是 coredns 的解析格式，但并不是 k8s resolver 的解析格式。

- [ ] 可以考虑将默认的 default 改成 当前 namespace


ps: 该实现大部分来自 go-kratos 中，做了一些裁剪
