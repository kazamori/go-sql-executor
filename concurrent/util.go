package concurrent

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/cybozu-go/log"
)

const (
	bufferedSize = 1024
)

func Call(
	ctx context.Context, concurrent int, f Func,
) []Data {
	ch := make(chan Data, bufferedSize)
	wg := &sync.WaitGroup{}
	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := f(ctx, ch); err != nil {
				log.Error("failed to run", map[string]interface{}{
					"i":     i,
					"error": err,
				})
				return
			}
			log.Debug(fmt.Sprintf("completed goroutine-%d", i), nil)
		}(i)
	}

	wg.Wait()
	close(ch)

	results := make([]Data, 0, concurrent)
	for data := range ch {
		results = append(results, data)
	}

	log.Debug("results", map[string]interface{}{
		"concurrent":    concurrent,
		"result_length": len(results),
	})

	return results
}

func ReadLine(r io.ReadCloser) <-chan Line {
	ch := make(chan Line, bufferedSize)
	go func() {
		defer func() {
			r.Close()
			close(ch)
		}()
		reader := bufio.NewReader(r)
		for {
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				return
			}
			lineData := Line{
				Value: strings.TrimSuffix(string(line), "\n"),
				Error: err,
			}
			ch <- lineData
			if err != nil {
				log.Error("failed to read", map[string]interface{}{
					"error": err,
				})
				return
			}
		}
	}()
	return ch
}

func ReadLines(r io.ReadCloser) ([]string, error) {
	lines := make([]string, 0)
	for line := range ReadLine(r) {
		if line.Error != nil {
			return nil, line.Error
		}
		lines = append(lines, line.Value)
	}
	return lines, nil
}
