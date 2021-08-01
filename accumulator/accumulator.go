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

var (
	accumulators   []*accumulator
	accumulatorMap = make(map[string]*accumulator)
)

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
	accumulatorMap[a.name] = a
}

func NewAccumulator(name string, expr string, f AccumulatorFunc) *accumulator {
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
	accumulatorMap[a.name] = a
	return a
}

func (a *accumulator) exec() {
	a.execute.Execute(a.metrics)
}

func (a *accumulator) get(property, metric string) uint64 {
	return a.metrics[property][metric]
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

func AddByName(accumulatorName, property, metric string, val uint64) {
	accumulatorMap[accumulatorName].add(property, metric, val)
}

func Get(accumulatorName, property, metric string) uint64 {
	// for _, a := range accumulators {
	// 	if a.name == accumulatorName {
	// 		return a.get(property, metric)
	// 	}
	// }
	return accumulatorMap[accumulatorName].get(property, metric)
}
