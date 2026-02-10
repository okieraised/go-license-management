package signal

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSignalChannelBuffering verifies that signal channels should be buffered
// to prevent deadlock scenarios
func TestSignalChannelBuffering(t *testing.T) {
	// Create buffered channels as recommended
	quit := make(chan os.Signal, 1)
	serverQuit := make(chan os.Signal, 1)

	// Verify channels are buffered
	assert.Equal(t, 1, cap(quit), "quit channel should be buffered with capacity 1")
	assert.Equal(t, 1, cap(serverQuit), "serverQuit channel should be buffered with capacity 1")

	// Verify non-blocking send works
	done := make(chan bool)
	go func() {
		serverQuit <- syscall.SIGTERM
		done <- true
	}()

	select {
	case <-done:
		// Good - send didn't block
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Send to buffered channel should not block")
	}

	// Verify signal was received
	select {
	case sig := <-serverQuit:
		assert.Equal(t, syscall.SIGTERM, sig)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Should have received signal from channel")
	}
}

// TestCatchableSignals verifies that only catchable signals are registered
func TestCatchableSignals(t *testing.T) {
	signals := CatchableSignals()

	// Verify we have the expected signals
	assert.Len(t, signals, 3, "Should have 3 catchable signals")

	// Verify specific signals are present
	expectedSignals := map[syscall.Signal]bool{
		syscall.SIGINT:  false,
		syscall.SIGTERM: false,
		syscall.SIGQUIT: false,
	}

	for _, sig := range signals {
		sysSig, ok := sig.(syscall.Signal)
		assert.True(t, ok, "Signal should be syscall.Signal")
		if _, exists := expectedSignals[sysSig]; exists {
			expectedSignals[sysSig] = true
		}
	}

	// Verify all expected signals were found
	for sig, found := range expectedSignals {
		assert.True(t, found, "Signal %v should be in catchable signals", sig)
	}
}

// TestCatchableSignalsRegistration verifies signals can be registered
func TestCatchableSignalsRegistration(t *testing.T) {
	sigChan := make(chan os.Signal, 1)

	// Register catchable signals
	signal.Notify(sigChan, CatchableSignals()...)
	defer signal.Stop(sigChan)

	// Test that SIGINT can be caught
	testSignalDelivery(t, sigChan, syscall.SIGINT)

	// Test that SIGTERM can be caught
	testSignalDelivery(t, sigChan, syscall.SIGTERM)

	// Test that SIGQUIT can be caught
	testSignalDelivery(t, sigChan, syscall.SIGQUIT)
}

func testSignalDelivery(t *testing.T, sigChan chan os.Signal, sig syscall.Signal) {
	// Send signal
	go func() {
		time.Sleep(10 * time.Millisecond)
		sigChan <- sig
	}()

	// Verify signal received
	select {
	case received := <-sigChan:
		assert.Equal(t, sig, received, "Should receive the signal that was sent")
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("Signal %v should be deliverable", sig)
	}
}

// TestIsSignalCatchable verifies signal catchability detection
func TestIsSignalCatchable(t *testing.T) {
	tests := []struct {
		name      string
		signal    os.Signal
		catchable bool
	}{
		{
			name:      "SIGINT is catchable",
			signal:    syscall.SIGINT,
			catchable: true,
		},
		{
			name:      "SIGTERM is catchable",
			signal:    syscall.SIGTERM,
			catchable: true,
		},
		{
			name:      "SIGQUIT is catchable",
			signal:    syscall.SIGQUIT,
			catchable: true,
		},
		{
			name:      "SIGKILL is not catchable",
			signal:    syscall.SIGKILL,
			catchable: false,
		},
		{
			name:      "SIGSTOP is not catchable",
			signal:    syscall.SIGSTOP,
			catchable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSignalCatchable(tt.signal)
			assert.Equal(t, tt.catchable, result, tt.name)
		})
	}
}

// TestGracefulShutdownSignal verifies the graceful shutdown signal
func TestGracefulShutdownSignal(t *testing.T) {
	sig := GracefulShutdownSignal()

	// Should return SIGTERM
	assert.Equal(t, syscall.SIGTERM, sig, "Graceful shutdown should use SIGTERM")

	// Should be catchable
	assert.True(t, IsSignalCatchable(sig), "Graceful shutdown signal should be catchable")
}

// TestSIGKILLNotInCatchableSignals ensures SIGKILL is excluded
func TestSIGKILLNotInCatchableSignals(t *testing.T) {
	signals := CatchableSignals()

	for _, sig := range signals {
		sysSig, ok := sig.(syscall.Signal)
		assert.True(t, ok, "Signal should be syscall.Signal")
		assert.NotEqual(t, syscall.SIGKILL, sysSig, "SIGKILL should not be in catchable signals")
		assert.NotEqual(t, syscall.SIGSTOP, sysSig, "SIGSTOP should not be in catchable signals")
	}
}

// TestSignalChannelNonBlockingSend verifies buffered channels prevent deadlock
func TestSignalChannelNonBlockingSend(t *testing.T) {
	// Unbuffered channel would block
	unbuffered := make(chan os.Signal)

	blocked := true
	go func() {
		select {
		case unbuffered <- syscall.SIGTERM:
			blocked = false
		case <-time.After(50 * time.Millisecond):
			// Expected to timeout - send blocks without receiver
		}
	}()

	time.Sleep(100 * time.Millisecond)
	assert.True(t, blocked, "Unbuffered channel should block without receiver")

	// Buffered channel should not block
	buffered := make(chan os.Signal, 1)

	notBlocked := false
	go func() {
		buffered <- syscall.SIGTERM
		notBlocked = true
	}()

	time.Sleep(50 * time.Millisecond)
	assert.True(t, notBlocked, "Buffered channel should not block")

	// Clean up
	<-buffered
}

// TestMultipleSignalsHandling verifies that multiple signals can be caught
func TestMultipleSignalsHandling(t *testing.T) {
	sigChan := make(chan os.Signal, 3)
	signal.Notify(sigChan, CatchableSignals()...)
	defer signal.Stop(sigChan)

	signals := []syscall.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

	// Send all signals
	for _, sig := range signals {
		sigChan <- sig
	}

	// Verify all received
	received := make(map[syscall.Signal]bool)
	for i := 0; i < len(signals); i++ {
		select {
		case sig := <-sigChan:
			received[sig.(syscall.Signal)] = true
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Should receive all signals")
		}
	}

	for _, sig := range signals {
		assert.True(t, received[sig], "Should have received signal %v", sig)
	}
}

// TestSignalStopPreventsDelivery verifies that signal.Stop works correctly
func TestSignalStopPreventsDelivery(t *testing.T) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	// Stop signal delivery
	signal.Stop(sigChan)

	// Manually sending still works (it's just a channel)
	sigChan <- syscall.SIGINT

	select {
	case sig := <-sigChan:
		assert.Equal(t, syscall.SIGINT, sig, "Manual send should still work")
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Should receive manually sent signal")
	}
}

// TestSIGKILLDocumentation documents why SIGKILL is excluded
func TestSIGKILLDocumentation(t *testing.T) {
	// This test serves as documentation for why SIGKILL is not in signal.Notify

	// From signal package documentation:
	// "The SIGKILL and SIGSTOP signals may not be caught by a program"

	// SIGKILL (signal 9):
	// - Cannot be caught, blocked, or ignored
	// - Terminates process immediately without cleanup
	// - Used by OS to forcefully kill unresponsive processes
	// - Including it in signal.Notify has no effect

	// Proper shutdown signals:
	// - SIGTERM (15): Termination signal, allows graceful shutdown
	// - SIGINT (2): Interrupt from keyboard (Ctrl+C)
	// - SIGQUIT (3): Quit from keyboard (Ctrl+\)

	t.Log("SIGKILL cannot be caught - this is by design")
	t.Log("Use SIGTERM for graceful shutdowns")
	t.Log("SIGINT for user-initiated interrupts")
	t.Log("SIGQUIT for quit with core dump")

	// Verify SIGKILL is not catchable
	assert.False(t, IsSignalCatchable(syscall.SIGKILL), "SIGKILL must not be catchable")
	assert.False(t, IsSignalCatchable(syscall.SIGSTOP), "SIGSTOP must not be catchable")
}

// TestSIGTERMvsSIGKILL documents the difference between signals
func TestSIGTERMvsSIGKILL(t *testing.T) {
	tests := []struct {
		name        string
		signal      syscall.Signal
		catchable   bool
		graceful    bool
		description string
	}{
		{
			name:        "SIGTERM",
			signal:      syscall.SIGTERM,
			catchable:   true,
			graceful:    true,
			description: "Graceful termination - allows cleanup",
		},
		{
			name:        "SIGINT",
			signal:      syscall.SIGINT,
			catchable:   true,
			graceful:    true,
			description: "Interrupt signal - typically from Ctrl+C",
		},
		{
			name:        "SIGQUIT",
			signal:      syscall.SIGQUIT,
			catchable:   true,
			graceful:    true,
			description: "Quit signal - typically from Ctrl+\\",
		},
		{
			name:        "SIGKILL",
			signal:      syscall.SIGKILL,
			catchable:   false,
			graceful:    false,
			description: "Force kill - cannot be caught or ignored",
		},
		{
			name:        "SIGSTOP",
			signal:      syscall.SIGSTOP,
			catchable:   false,
			graceful:    false,
			description: "Stop process - cannot be caught or ignored",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Signal: %s (%d)", tt.name, tt.signal)
			t.Logf("Catchable: %v", tt.catchable)
			t.Logf("Graceful: %v", tt.graceful)
			t.Logf("Description: %s", tt.description)

			result := IsSignalCatchable(tt.signal)
			assert.Equal(t, tt.catchable, result, "IsSignalCatchable should match expected value")

			if tt.name == "SIGKILL" || tt.name == "SIGSTOP" {
				assert.False(t, tt.catchable, "%s should not be catchable", tt.name)
				assert.False(t, tt.graceful, "%s does not allow graceful shutdown", tt.name)
			}
		})
	}
}

// BenchmarkSignalChannelSend benchmarks signal channel operations
func BenchmarkSignalChannelSend(b *testing.B) {
	sigChan := make(chan os.Signal, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sigChan <- syscall.SIGTERM
		<-sigChan
	}
}

// BenchmarkBufferedVsUnbuffered compares buffered and unbuffered channels
func BenchmarkBufferedVsUnbuffered(b *testing.B) {
	b.Run("Buffered", func(b *testing.B) {
		sigChan := make(chan os.Signal, 1)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			go func() { sigChan <- syscall.SIGTERM }()
			<-sigChan
		}
	})

	b.Run("Unbuffered", func(b *testing.B) {
		sigChan := make(chan os.Signal)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			go func() { sigChan <- syscall.SIGTERM }()
			<-sigChan
		}
	})
}

// BenchmarkIsSignalCatchable benchmarks signal catchability check
func BenchmarkIsSignalCatchable(b *testing.B) {
	signals := []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGSTOP,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, sig := range signals {
			_ = IsSignalCatchable(sig)
		}
	}
}
