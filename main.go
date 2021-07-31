package main

import (
	"fmt"
	"golang_tools/accumulator"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Execute struct {
	accumulator.AccumulatorFunc
}

func (e *Execute) Execute(m map[string]map[string]uint64) {
	for k, map1 := range m {
		fmt.Print(k)
		for k1, v := range map1 {
			fmt.Printf("\t%s\t%d", k1, v)
		}
		fmt.Println()
	}
}

func testAccumulator() {
	execute := &Execute{}
	accumulator.SetAccumulator("123", "@every 1s", execute)
	accumulator.Run()
	done := make(chan bool)
	go func() {
		for i := 0; i < 100; i++ {
			time.Sleep(time.Second * time.Duration(1))
			accumulator.Add("total", "time", 1)
		}
		done <- true
	}()
	// go func() {
	// 	for i := 0; i < 100; i++ {
	// 		time.Sleep(time.Second * time.Duration(1))
	// 		res := accumulator.Get("123", "total", "time")
	// 		fmt.Println(res)
	// 	}
	// 	done <- true
	// }()

	go func() {
		for true {
			select {
			case <-done:
				break
			default:
				res := accumulator.Get("123", "total", "time")
				fmt.Println(res)
			}
			time.Sleep(time.Second*time.Duration(1))
		}
	}()
}

func testTimeTicker() {
	ticker := time.NewTicker(time.Second * time.Duration(1))
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("213333")
			}
		}
	}()
}

func exit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()
	<-done
	fmt.Println("exiting")
}

func main() {
	testAccumulator()
	exit()
}
