package accumulator

import (
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"time"
)

type accumulator struct {
	execute AccumulatorFunc
	name    string
	expr    string
	metrics map[string]map[string]uint64
	mutex   sync.Mutex
}

var accumulators []*accumulator

type AccumulatorFunc interface {
	Execute(map[string]map[string]uint64)
}

func Run() {
	location, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(location))
	for _, a := range accumulators {
		_, err := c.AddFunc(a.expr, a.exec)
		if err != nil {
			log.Fatal(err)
		}
	}
	go c.Run()
}

func SetAccumulator(name string, expr string, f AccumulatorFunc) {
	a := &accumulator{
		name:    name,
		expr:    expr,
		metrics: make(map[string]map[string]uint64),
		mutex:   sync.Mutex{},
		execute: f,
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	accumulators = append(accumulators, a)
}

func NewAccumulator(name string, expr string, f AccumulatorFunc) *accumulator {
	a := &accumulator{
		name:    name,
		expr:    expr,
		metrics: make(map[string]map[string]uint64),
		mutex:   sync.Mutex{},
		execute: f,
	}
	accumulators = append(accumulators, a)
	return a
}

func (a *accumulator) exec() {
	log.Println("start execute job")
	a.execute.Execute(a.metrics)
}

func (a *accumulator) add(property, metric string, val uint64) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	sub, ok := a.metrics[property]
	if !ok {
		sub = make(map[string]uint64)
		a.metrics[property] = sub
	}
	sub[metric] += val
}

func Add(property, metric string, val uint64) {
	for _, a := range accumulators {
		a.add(property, metric, val)
	}
}
