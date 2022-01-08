package lb

type ILoadBalanceCondition interface {
	GetCondition() interface{}
}

type ISubject interface {
	Attach(o Observer)
	GetConfig() interface{}
	NotifyAll()
}

type Observer interface {
	GetSubject() ISubject
	Update()
}

type ILoadBalance interface {
	// Get condition为负载均衡的条件，空接口。 返回负载均衡的目标地址
	Get(condition ILoadBalanceCondition) (string, error)
}
