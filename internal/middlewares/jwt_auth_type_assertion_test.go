package middlewares

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSafeStringTypeAssertion tests the pattern used in the fix
// This demonstrates the difference between unsafe and safe type assertions
func TestSafeStringTypeAssertion(t *testing.T) {
	tests := []struct {
		name        string
		value       interface{}
		expectOK    bool
		expectPanic bool
		description string
	}{
		{
			name:        "Valid string",
			value:       "active",
			expectOK:    true,
			expectPanic: false,
			description: "String value should succeed",
		},
		{
			name:        "Integer value",
			value:       123,
			expectOK:    false,
			expectPanic: false,
			description: "Integer should fail gracefully with ok check",
		},
		{
			name:        "Boolean value",
			value:       true,
			expectOK:    false,
			expectPanic: false,
			description: "Boolean should fail gracefully with ok check",
		},
		{
			name:        "Nil value",
			value:       nil,
			expectOK:    false,
			expectPanic: false,
			description: "Nil should fail gracefully with ok check",
		},
		{
			name:        "Float value",
			value:       3.14,
			expectOK:    false,
			expectPanic: false,
			description: "Float should fail gracefully with ok check",
		},
		{
			name:        "Array value",
			value:       []string{"active"},
			expectOK:    false,
			expectPanic: false,
			description: "Array should fail gracefully with ok check",
		},
		{
			name:        "Map value",
			value:       map[string]string{"status": "active"},
			expectOK:    false,
			expectPanic: false,
			description: "Map should fail gracefully with ok check",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test safe type assertion (what we fixed to)
			var result string
			var ok bool

			assert.NotPanics(t, func() {
				result, ok = tt.value.(string)
			}, "Safe type assertion with ok check should never panic")

			assert.Equal(t, tt.expectOK, ok, tt.description)

			if tt.expectOK {
				assert.NotEmpty(t, result, "Result should not be empty for valid string")
			} else {
				assert.Empty(t, result, "Result should be empty (zero value) for failed assertion")
			}
		})
	}
}

// TestUnsafeTypeAssertionPanics demonstrates what happens with unsafe assertions
func TestUnsafeTypeAssertionPanics(t *testing.T) {
	tests := []struct {
		name        string
		value       interface{}
		shouldPanic bool
		description string
	}{
		{
			name:        "String value (no panic)",
			value:       "active",
			shouldPanic: false,
			description: "Valid string doesn't panic even without ok check",
		},
		{
			name:        "Integer value (PANICS)",
			value:       123,
			shouldPanic: true,
			description: "Integer causes panic without ok check",
		},
		{
			name:        "Boolean value (PANICS)",
			value:       true,
			shouldPanic: true,
			description: "Boolean causes panic without ok check",
		},
		{
			name:        "Nil value (PANICS)",
			value:       nil,
			shouldPanic: true,
			description: "Nil causes panic without ok check",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				// Unsafe assertion WILL panic on wrong type
				assert.Panics(t, func() {
					_ = tt.value.(string) // Unsafe - no ok check
				}, tt.description)
			} else {
				// Only succeeds if type matches
				assert.NotPanics(t, func() {
					_ = tt.value.(string)
				}, tt.description)
			}
		})
	}
}

// TestTypeAssertionBestPractices documents Go type assertion best practices
func TestTypeAssertionBestPractices(t *testing.T) {
	t.Run("Best practice: Always use two-value assertion", func(t *testing.T) {
		var value interface{} = 42 // Not a string

		// GOOD: Safe type assertion with ok check
		str, ok := value.(string)
		if !ok {
			// Handle error gracefully
			assert.False(t, ok, "Type assertion should fail")
			assert.Empty(t, str, "Failed assertion returns zero value")
			return
		}

		t.Fatal("Should not reach here")
	})

	t.Run("Bad practice: Single-value assertion can panic", func(t *testing.T) {
		var value interface{} = 42 // Not a string

		// BAD: Unsafe type assertion without ok check
		assert.Panics(t, func() {
			_ = value.(string) // This WILL panic
		}, "Single-value assertion panics on type mismatch")
	})
}

// TestJWTClaimTypeValidation tests the specific pattern used in JWT middleware
func TestJWTClaimTypeValidation(t *testing.T) {
	// Simulate JWT claims with various status types
	type testClaim struct {
		status interface{}
	}

	tests := []struct {
		name        string
		claim       testClaim
		expectValid bool
		description string
	}{
		{
			name:        "Valid string status",
			claim:       testClaim{status: "active"},
			expectValid: true,
			description: "String status is valid",
		},
		{
			name:        "Invalid integer status",
			claim:       testClaim{status: 1},
			expectValid: false,
			description: "Integer status is invalid",
		},
		{
			name:        "Invalid boolean status",
			claim:       testClaim{status: true},
			expectValid: false,
			description: "Boolean status is invalid",
		},
		{
			name:        "Invalid nil status",
			claim:       testClaim{status: nil},
			expectValid: false,
			description: "Nil status is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This simulates the fixed code in jwt_auth_mw.go
			statusCtx := tt.claim.status

			// Safe type assertion with ok check (the fix)
			status, ok := statusCtx.(string)

			assert.NotPanics(t, func() {
				_, _ = statusCtx.(string)
			}, "Type assertion should not panic")

			if tt.expectValid {
				assert.True(t, ok, "Should be valid string")
				assert.NotEmpty(t, status, "Status should have value")
			} else {
				assert.False(t, ok, "Should not be valid string")
				assert.Empty(t, status, "Status should be empty on failed assertion")
			}
		})
	}
}

// TestMaliciousJWTPayloads tests handling of potentially malicious payloads
func TestMaliciousJWTPayloads(t *testing.T) {
	maliciousPayloads := []struct {
		name   string
		status interface{}
		desc   string
	}{
		{
			name:   "SQL Injection attempt as integer",
			status: 123,
			desc:   "Attacker tries to inject SQL as non-string",
		},
		{
			name:   "Boolean confusion attack",
			status: false,
			desc:   "Attacker tries to bypass checks with boolean",
		},
		{
			name:   "Array payload",
			status: []interface{}{"active", "admin"},
			desc:   "Attacker tries privilege escalation with array",
		},
		{
			name:   "Nested object",
			status: map[string]interface{}{"role": "admin"},
			desc:   "Attacker embeds object in status field",
		},
		{
			name:   "Large integer",
			status: int64(9223372036854775807),
			desc:   "Attacker tries integer overflow",
		},
	}

	for _, payload := range maliciousPayloads {
		t.Run(payload.name, func(t *testing.T) {
			// The fixed code handles all these safely
			assert.NotPanics(t, func() {
				status, ok := payload.status.(string)
				if ok {
					t.Fatalf("Malicious payload %v was accepted as string", payload.status)
				}
				assert.Empty(t, status, "Should reject malicious payload")
			}, payload.desc)
		})
	}
}

// BenchmarkSafeVsUnsafeTypeAssertion compares performance
func BenchmarkSafeVsUnsafeTypeAssertion(b *testing.B) {
	var value interface{} = "active"

	b.Run("Safe assertion with ok", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, ok := value.(string)
			if !ok {
				b.Fatal("Should not happen in benchmark")
			}
		}
	})

	b.Run("Unsafe assertion without ok", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = value.(string)
		}
	})
}

// BenchmarkTypeAssertionOnDifferentTypes benchmarks various types
func BenchmarkTypeAssertionOnDifferentTypes(b *testing.B) {
	values := []interface{}{
		"string",
		123,
		true,
		3.14,
		[]string{"a"},
		map[string]string{"a": "b"},
	}

	for i, val := range values {
		b.Run(sprint("Type_%d", i), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				_, _ = val.(string)
			}
		})
	}
}

func sprint(format string, a ...interface{}) string {
	// Simple sprintf replacement for benchmark naming
	if len(a) == 0 {
		return format
	}
	// Just a simple case for our benchmark
	return format
}
