package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	startNow := time.Now()
	cities := []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix"}

	// withoutGoroutine(cities)
	// withGoroutine(cities)
	withGoroutineAndChannel(cities)

	fmt.Println("Without Goroutines:", time.Since(startNow))

}

func withoutGoroutine(cities []string) {
	var totalTemperature float64
	for _, city := range cities {
		temp := getTemperature()
		totalTemperature += temp
		fmt.Printf("The temperature in %s is %.2f degrees.\n", city, temp)
	}
	fmt.Printf("Total temperature: %.2f degrees.\n", totalTemperature)
}

// NOTE:
// This is the anti-pattern version of using goroutines.
// Poor timing could result in overwritten values being rewritten by another go routine.
// Better to use channels or Lock()/UnLock().

func withGoroutine(cities []string) {
	var totalTemperature float64
	var wg sync.WaitGroup
	for _, city := range cities {
		wg.Add(1)
		go func(city string) {
			defer wg.Done()
			temp := getTemperature()
			totalTemperature += temp
			fmt.Printf("The temperature in %s is %.2f degrees.\n", city, temp)
		}(city)
	}
	wg.Wait()
	fmt.Printf("Total temperature: %.2f degrees.\n", totalTemperature)
}

// NOTE:
// Regarding https://pkg.go.dev/sync#pkg-overview
// Other than the Once and WaitGroup types, most are intended for use by low-level library routines.
// Higher-level synchronization is better done via channels and communication.
func withGoroutineAndChannel(cities []string) {
	ch := make(chan float64)
	for _, city := range cities {
		go func(city string) {
			temp := getTemperature()
			fmt.Printf("The temperature in %s is %.2f degrees.\n", city, temp)
			ch <- temp
		}(city)
	}

	var totalTemperature float64
	for range cities {
		temp := <-ch
		totalTemperature += temp
	}

	fmt.Printf("Total temperature: %.2f degrees.\n", totalTemperature)
	close(ch)
}

func getTemperature() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	wait := time.Duration(r.Intn(5)) * time.Second
	time.Sleep(wait)

	temp := 0.0 + r.Float64()*(40.0-0.0)
	return temp
}
