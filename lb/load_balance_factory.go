package lb

const (
	TypeRoundRobin = 1
	TypeWeight     = 3
)

func GetLoadBalancer(lbType int, subject ISubject) ILoadBalance {
	switch lbType {
	case TypeRoundRobin:
		return NewRoundRobinLoadBalance(subject)
	default:
		return nil
	}
}
