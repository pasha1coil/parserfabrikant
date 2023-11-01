package main

import (
	"log"
	"sync"
	"testingParser/internal/readfile"
	"testingParser/internal/worker"
)

func main() {
	var wg sync.WaitGroup
	// канал для передачи значений с тхт в работу
	ch := make(chan string)

	go func() {
		err := readfile.ReadFile("data.txt", ch)
		if err != nil {
			log.Fatalf("Error reading file: %s", err)
		}
		close(ch)
	}()

	for val := range ch {
		wg.Add(1)

		go func(value string) {
			defer wg.Done()
			worker.GetUrl(value)
		}(val)
	}

	wg.Wait()

}
