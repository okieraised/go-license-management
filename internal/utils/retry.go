package utils

import "time"

// RetryWithBackoff retries the function call with waiting time between each retries
func RetryWithBackoff(fn func() error, maxRetry int, startBackoff, maxBackoff time.Duration) {

	for attempt := 0; ; attempt++ {
		err := fn()
		if err == nil {
			return
		}
		if attempt == maxRetry-1 {
			return
		}

		time.Sleep(startBackoff)
		if startBackoff < maxBackoff {
			startBackoff *= 2
		}
	}
}

// Retry retries the function call without waiting time between each retries
func Retry(fn func() error, maxRetry int) {
	for attempt := 0; ; attempt++ {
		err := fn()
		if err == nil {
			return
		}
		if attempt == maxRetry-1 {
			return
		}
	}
}
