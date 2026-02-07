package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestTimeoutMW_RequestCompletesBeforeTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TimeoutMW())

	handlerExecuted := false
	router.GET("/test", func(c *gin.Context) {
		handlerExecuted = true
		time.Sleep(100 * time.Millisecond) // Short delay, well under timeout
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, handlerExecuted)
	assert.Contains(t, w.Body.String(), "success")
}

func TestTimeoutMW_RequestTimesOut(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TimeoutMW())

	handlerStarted := false
	router.GET("/test", func(c *gin.Context) {
		handlerStarted = true
		// Simulate long-running operation that should timeout
		time.Sleep(15 * time.Second)
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	start := time.Now()
	router.ServeHTTP(w, req)
	duration := time.Since(start)

	assert.Equal(t, http.StatusGatewayTimeout, w.Code)
	assert.True(t, handlerStarted)
	// Should timeout around 10 seconds, not wait for full 15 seconds
	assert.Less(t, duration, 12*time.Second, "Should timeout before handler completes")
	assert.Greater(t, duration, 9*time.Second, "Should take at least the timeout duration")
}

func TestTimeoutMW_ContextCancellationPropagates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TimeoutMW())

	var capturedCtx context.Context
	contextCancelled := false

	router.GET("/test", func(c *gin.Context) {
		capturedCtx = c.Request.Context()

		// Check if context gets cancelled after timeout
		select {
		case <-capturedCtx.Done():
			contextCancelled = true
		case <-time.After(15 * time.Second):
			// Should not reach here if context cancellation works
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Wait a bit for goroutine to finish
	time.Sleep(100 * time.Millisecond)

	assert.NotNil(t, capturedCtx)
	assert.True(t, contextCancelled, "Context should be cancelled on timeout")
	assert.Equal(t, context.DeadlineExceeded, capturedCtx.Err())
}

func TestTimeoutMW_HandlerPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TimeoutMW())

	router.GET("/test", func(c *gin.Context) {
		panic("intentional panic for testing")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTimeoutMW_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TimeoutMW())

	requestCount := 0
	router.GET("/test", func(c *gin.Context) {
		requestCount++
		time.Sleep(50 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"count": requestCount})
	})

	// Make multiple concurrent requests
	done := make(chan bool)
	for i := 0; i < 5; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	assert.Equal(t, 5, requestCount)
}

func TestTimeoutMW_ContextDeadlineSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TimeoutMW())

	var deadline time.Time
	var hasDeadline bool

	router.GET("/test", func(c *gin.Context) {
		deadline, hasDeadline = c.Request.Context().Deadline()
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	before := time.Now()
	router.ServeHTTP(w, req)
	after := time.Now()

	assert.True(t, hasDeadline, "Context should have a deadline")
	assert.True(t, deadline.After(before), "Deadline should be in the future")
	assert.True(t, deadline.Before(after.Add(11*time.Second)), "Deadline should be ~10 seconds from request start")
}

func TestTimeoutMW_RespectsDatabaseContextCancellation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TimeoutMW())

	queryStarted := false
	queryCancelled := false

	router.GET("/test", func(c *gin.Context) {
		queryStarted = true
		ctx := c.Request.Context()

		// Simulate a database query that respects context cancellation
		go func() {
			select {
			case <-ctx.Done():
				queryCancelled = true
			case <-time.After(15 * time.Second):
				// Should not reach if context cancellation works
			}
		}()

		// Wait for timeout to occur
		time.Sleep(15 * time.Second)
		c.JSON(http.StatusOK, gin.H{"message": "should not reach"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Give goroutine time to detect cancellation
	time.Sleep(100 * time.Millisecond)

	assert.True(t, queryStarted, "Query should have started")
	assert.True(t, queryCancelled, "Query should be cancelled when request times out")
	assert.Equal(t, http.StatusGatewayTimeout, w.Code)
}

func TestTimeoutMW_NoGoroutineLeakOnTimeout(t *testing.T) {
	// This test verifies that goroutines properly terminate after timeout
	// by checking that context cancellation propagates

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TimeoutMW())

	goroutineCompleted := make(chan bool, 1)

	router.GET("/test", func(c *gin.Context) {
		ctx := c.Request.Context()

		// Simulate long operation that respects context
		go func() {
			defer func() {
				goroutineCompleted <- true
			}()

			select {
			case <-ctx.Done():
				// Context cancelled, clean exit
				return
			case <-time.After(20 * time.Second):
				// Should not reach if context works properly
				return
			}
		}()

		// Block until timeout
		time.Sleep(15 * time.Second)
		c.JSON(http.StatusOK, gin.H{"message": "done"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Check that goroutine completed due to context cancellation
	select {
	case <-goroutineCompleted:
		// Good - goroutine exited due to context cancellation
	case <-time.After(2 * time.Second):
		t.Fatal("Goroutine did not terminate after context cancellation")
	}

	assert.Equal(t, http.StatusGatewayTimeout, w.Code)
}

func TestTimeoutMW_FastRequestNoDelay(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TimeoutMW())

	router.GET("/test", func(c *gin.Context) {
		// Very fast handler
		c.JSON(http.StatusOK, gin.H{"message": "fast"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	start := time.Now()
	router.ServeHTTP(w, req)
	duration := time.Since(start)

	assert.Equal(t, http.StatusOK, w.Code)
	// Should complete quickly, not wait for timeout
	assert.Less(t, duration, 1*time.Second)
}

func BenchmarkTimeoutMW_NoTimeout(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TimeoutMW())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func TestTimeoutMW_ContextValuesPropagated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	// Middleware that sets a value in context
	router.Use(func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "testKey", "testValue")
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	router.Use(TimeoutMW())

	var capturedValue interface{}
	router.GET("/test", func(c *gin.Context) {
		capturedValue = c.Request.Context().Value("testKey")
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "testValue", capturedValue, "Context values should be propagated through timeout middleware")
}
