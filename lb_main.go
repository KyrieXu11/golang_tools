package main

import (
	"golang_tools/lb"
	"log"
)

func main() {
	conf := lb.NewLoadBalanceConf([]string{
		"1", "2", "3", "4",
	})
	balance := lb.NewRoundRobinLoadBalance(conf)
	s, _ := balance.Get(nil)
	log.Println(s)
	s, _ = balance.Get(nil)
	log.Println(s)
	s, _ = balance.Get(nil)
	log.Println(s)

	s, _ = balance.Get(nil)
	log.Println(s)
	s, _ = balance.Get(nil)
	log.Println(s)

	conf.AddAddr("5", "6")

	s, _ = balance.Get(nil)
	log.Println(s)

	s, _ = balance.Get(nil)
	log.Println(s)
	s, _ = balance.Get(nil)
	log.Println(s)
	s, _ = balance.Get(nil)
	log.Println(s)
	s, _ = balance.Get(nil)
	log.Println(s)

	balance = lb.GetLoadBalancer(lb.TypeRoundRobin, conf)

	s, _ = balance.Get(nil)
	log.Println(s)
}
