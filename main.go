package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Execute struct {
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

func main() {

	// execute := &Execute{}
	// accumulator.SetAccumulator("123", "@every 1s", execute)
	// accumulator.Run()
	// go func() {
	// 	for i := 0; i < 100; i++ {
	// 		time.Sleep(time.Second * time.Duration(1))
	// 		accumulator.Add("total", "time", 1)
	// 	}
	// }()
	// go func() {
	// 	for i := 0; i < 100; i++ {
	// 		time.Sleep(time.Second * time.Duration(1))
	// 		accumulator.Add("total", "time", 1)
	// 	}
	// }()

	ticker := time.NewTicker(time.Second * time.Duration(1))
	go func() {
		for{
			select {
			case <-ticker.C:
				fmt.Println("213333")
			}
		}
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	<-done
	fmt.Println("exiting")

}
