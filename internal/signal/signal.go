package signal

import (
	"os"
	"syscall"
)

// CatchableSignals returns the list of signals that can be caught
// SIGKILL and SIGSTOP are excluded as they cannot be caught on Unix-like systems
func CatchableSignals() []os.Signal {
	return []os.Signal{
		syscall.SIGINT,  // Interrupt from keyboard (Ctrl+C)
		syscall.SIGTERM, // Termination signal
		syscall.SIGQUIT, // Quit from keyboard (Ctrl+\)
	}
}

// GracefulShutdownSignal returns the signal that should be used for graceful shutdown
// SIGTERM is preferred as it allows the application to clean up resources
func GracefulShutdownSignal() syscall.Signal {
	return syscall.SIGTERM
}

// IsSignalCatchable returns whether a signal can be caught by the application
func IsSignalCatchable(sig os.Signal) bool {
	// SIGKILL (9) and SIGSTOP (17/19/23 depending on platform) cannot be caught
	sysSig, ok := sig.(syscall.Signal)
	if !ok {
		return false
	}

	// SIGKILL and SIGSTOP are the only uncatchable signals
	return sysSig != syscall.SIGKILL && sysSig != syscall.SIGSTOP
}
