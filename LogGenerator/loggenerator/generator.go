package loggenerator

import (
	"LogGenerator/logger"
	_ "LogGenerator/models"
	"LogGenerator/utils"
	"context"
	"fmt"
	_ "log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)
type Generator struct{}

const maxBatchSizeBytes = 10 * 1024 * 1024

// GenerateLog generates a random log entry string simulating an HTTP request log.
// It simulates various fields like IP address, method, status, and more.
//
// Returns:
//   - A string representing a randomly generated log entry formatted for HTTP access logs.
//
// Example usage:
//   logEntry := GenerateLog()
//   log.Printf("Generated log entry: %s", logEntry)
func GenerateLog() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	ip := utils.Ips[rnd.Intn(len(utils.Ips))]
	method := utils.Methods[rnd.Intn(len(utils.Methods))]
	url := utils.Urls[rnd.Intn(len(utils.Urls))]
	status := utils.Statuses[rnd.Intn(len(utils.Statuses))]
	bodyBytesSent := rnd.Intn(1000) + 500
	referrer := utils.Referrers[rnd.Intn(len(utils.Referrers))]
	userAgent := utils.UserAgents[rnd.Intn(len(utils.UserAgents))]
	xForwardedFor := fmt.Sprintf("%d.%d.%d.%d", rnd.Intn(256), rnd.Intn(256), rnd.Intn(256), rnd.Intn(256))

	request := fmt.Sprintf("%s %s HTTP/1.1", method, url)
	//timeLocal := time.Now()//.Format("02/Jan/2006:15:04:05 -0700")
	timeLocal := time.Now().UTC().Format(time.RFC3339)
	return fmt.Sprintf("%s - - [%s] \"%s\" %d %d \"%s\" \"%s\" \"%s\"",
	ip, timeLocal, request, status, bodyBytesSent, referrer, userAgent, xForwardedFor)

}

// GenerateLogsConcurrently generates logs concurrently across multiple goroutines. The number of logs is 
// distributed among workers based on the optimal number of workers derived from the number of CPU cores and 
// the total number of logs requested. This method also ensures efficient memory usage by batching the logs 
// and sending them to a processor when a batch reaches a certain size.
//
// Parameters:
//   - ctx: The context used to manage cancellation or timeouts during log generation.
//   - numLogs: The total number of logs to be generated.
//   - duration: The duration over which the logs should be generated (e.g., for spreading out log generation).
//   - counter: A WaitGroup used to ensure all goroutines finish before the function returns.
//
// This function generates logs concurrently using multiple workers. The log generation process is 
// controlled by a ticker that spreads out log creation over the specified `duration`. The function 
// also ensures that logs are batched to avoid exceeding memory limits, and batches are sent 
// to the processor when necessary.
//
// Example usage:
//   var wg sync.WaitGroup
//   ctx := context.Background()
//   logGen := Generator{}
//   logGen.GenerateLogsConcurrently(ctx, 10000, 1*time.Minute, &wg)
func (l *Generator) GenerateLogsConcurrently(ctx context.Context, numLogs int, duration time.Duration,counter *sync.WaitGroup) {
	
	logs := make([]string, numLogs)

	numCPU := runtime.NumCPU()
	optimalWorkers := numCPU
	if numLogs > 1000 {
		optimalWorkers = int(float64(numLogs) / float64(1000))
		if optimalWorkers > numCPU*2 {
			optimalWorkers = numCPU * 2
		}
	}

	if optimalWorkers < 1 {
		optimalWorkers = 1
	}

	logsPerWorker := numLogs / optimalWorkers

	var mu sync.Mutex
	var generatedLogs int
	logTicker := time.NewTicker(duration/time.Duration(numLogs))
	defer logTicker.Stop()


	for worker_i := 0; worker_i < optimalWorkers; worker_i++ {
		counter.Add(1)
		go func(workerID int) {
			defer counter.Done()

			startIndex := workerID * logsPerWorker
			endIndex := startIndex + logsPerWorker
			if workerID == optimalWorkers-1 {
				endIndex = numLogs
			}

			batch := []string{}
			totalBatchSize := 0

			for logIndex := startIndex; logIndex < endIndex; logIndex++ {
				select{
				case <-ctx.Done():
					return
				case <-logTicker.C:
						mu.Lock()
						if generatedLogs >= numLogs {
							logger.LogDebug(fmt.Sprintf("\n\n\n Given is size for the given time %v: size", generatedLogs))
							mu.Unlock()
							return
						}
						generatedLogs++
						mu.Unlock()

						logs[logIndex] = GenerateLog()
						logger.LogDebug(fmt.Sprintf("Generated Log: %s\n", logs[logIndex]))

						logSize := len(logs[logIndex])

					if totalBatchSize+logSize > maxBatchSizeBytes {
						logger.LogDebug(fmt.Sprintf("Batch byte size is more:%v", totalBatchSize+logSize))
						go SendLogToProcessor(batch)

						batch = []string{}
						totalBatchSize = 0
					}

					batch = append(batch, logs[logIndex])
					totalBatchSize += logSize

					if len(batch) >= 100 {
						logger.LogDebug(fmt.Sprintf("Batch size is more:%v", len(batch)))
						go SendLogToProcessor(batch)
						batch = []string{}
						totalBatchSize = 0
					}
				}
			}
			if len(batch) > 0 {
				go SendLogToProcessor(batch)
			}
		}(worker_i)
	}
	counter.Wait()
}
