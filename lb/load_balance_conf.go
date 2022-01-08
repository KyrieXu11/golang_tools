package lb

import "sync"

type LoadBalanceConf struct {
	mutex       sync.Mutex
	dstAddrs    []string
	weightAddrs map[int]string
	observers   []Observer
}

func NewLoadBalanceConf(addrs []string) *LoadBalanceConf {
	return &LoadBalanceConf{
		mutex:     sync.Mutex{},
		observers: make([]Observer, 0),
		dstAddrs:  addrs,
	}
}

func NewLoadBalanceConfWeight(weightAddrs map[int]string) ISubject {
	lbconf := &LoadBalanceConf{
		mutex:       sync.Mutex{},
		observers:   make([]Observer, 0),
		weightAddrs: weightAddrs,
	}
	for _, addr := range weightAddrs {
		lbconf.dstAddrs = append(lbconf.dstAddrs, addr)
	}
	return lbconf
}

func (k *LoadBalanceConf) NotifyAll() {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	for _, observer := range k.observers {
		observer.Update()
	}
}

func (k *LoadBalanceConf) Attach(o Observer) {
	k.mutex.Lock()
	defer k.mutex.Unlock()
	k.observers = append(k.observers, o)
}

func (k *LoadBalanceConf) GetConfig() interface{} {
	return k.dstAddrs
}

func (k *LoadBalanceConf) AddAddr(addrs ...string) {
	k.mutex.Lock()
	defer func() {
		k.mutex.Unlock()
		k.NotifyAll()
	}()
	k.dstAddrs = append(k.dstAddrs, addrs...)
}
