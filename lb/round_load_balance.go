package lb

type RoundRobinLoadBalance struct {
	// 显式标记实现了哪些接口
	ILoadBalance
	Observer

	curIndex        int
	loadBalanceConf ISubject
}

func NewRoundRobinLoadBalance(subject ISubject) ILoadBalance {
	balance := &RoundRobinLoadBalance{
		curIndex:        0,
		loadBalanceConf: subject,
	}
	subject.Attach(balance)
	return balance
}

func (r *RoundRobinLoadBalance) GetSubject() ISubject {
	return r.loadBalanceConf
}

func (r *RoundRobinLoadBalance) Update() {
	return
}

func (r *RoundRobinLoadBalance) Get(condition ILoadBalanceCondition) (string, error) {
	return r.next(), nil
}

func (r *RoundRobinLoadBalance) next() string {
	addrs := r.GetSubject().GetConfig().([]string)
	defer func() {
		r.curIndex = (r.curIndex + 1) % len(addrs)
	}()
	return addrs[r.curIndex]
}
