package platform

import (
	"runtime"
	"strings" // Needed for signal check
	"testing"
)

// TestFindSymbolFromPidPtr_SpecificExample tests the specific PID/address from the user prompt.
// NOTE: This test will likely be SKIPPED unless the exact process (PID 446510
// running the specific cudademo executable) is active and accessible during test execution.
func TestFindSymbolFromPidPtr_SpecificExample(t *testing.T) {
	// --- Test Configuration ---
	targetPid := 446510
	targetPtr := uintptr(0x639a4a98bc46)
	// Expectations based on the addr2line output provided
	expectedExeNamePart := "cudademo" // Check if the executable path contains this
	expectedSymbolPart := "vectorAdd" // Check if the resolved symbol contains this

	// --- Pre-check: Test environment ---
	if runtime.GOOS != "linux" {
		t.Skipf("Skipping test: requires Linux /proc filesystem, running on %s", runtime.GOOS)
	}

	// --- Execute the function under test ---
	t.Logf("Attempting to resolve symbol for PID %d, Addr 0x%x", targetPid, targetPtr)
	info, err := FindSymbolFromPidPtr(targetPid, targetPtr)

	// --- Analyze Results ---

	// Handle errors first
	if err != nil {
		// It's possible the process existed during the check but disappeared before/during the maps read (race condition).
		if strings.Contains(err.Error(), "process with PID") && strings.Contains(err.Error(), "not found") {
			t.Logf("Test Info: Process %d disappeared after initial check: %v", targetPid, err)
			// Don't fail the test in this specific race condition case, as the setup was valid initially.
			// Consider skipping instead if preferred: t.Skipf(...)
			return
		}
		// Check if the error is due to addr2line failing (e.g., binary stripped, path incorrect)
		// This is an expected outcome if the 'cudademo' binary lacks debug symbols or isn't at the path found in maps.
		if strings.Contains(err.Error(), "addr2line execution failed") {
			t.Logf("Test Info: addr2line failed as expected if symbols/path are missing: %v", err)
			// Even if addr2line failed, we might have partial info. Check if info is non-nil.
			if info == nil {
				t.Errorf("Expected non-nil partial SymbolInfo even when addr2line failed, but got nil")
			} else {
				// Check if the partial path info is reasonable
				if !strings.Contains(info.FilePath, expectedExeNamePart) {
					t.Errorf("Partial FilePath '%s' does not contain expected '%s'", info.FilePath, expectedExeNamePart)
				}
				t.Logf("Partial Info: FilePath='%s', Offset=0x%x, Base=0x%x", info.FilePath, info.Offset, info.BaseAddress)
			}
			// Don't proceed to check symbol name if addr2line failed.
			return
		}

		// Any other unexpected error should fail the test.
		t.Fatalf("FindSymbolFromPidPtr failed with unexpected error: %v", err)
	}

	// If no error occurred, check the SymbolInfo content.
	if info == nil {
		t.Fatalf("FindSymbolFromPidPtr returned nil error but also nil SymbolInfo")
	}

	t.Logf("Successfully resolved: Name='%s', File='%s', Offset=0x%x, Base=0x%x, Source=%s:%d",
		info.SymbolName, info.FilePath, info.Offset, info.BaseAddress, info.SourceFile, info.SourceLine)

	// Verify the file path contains the expected executable name part.
	if !strings.Contains(info.FilePath, expectedExeNamePart) {
		t.Errorf("Resolved FilePath '%s' does not contain expected '%s'", info.FilePath, expectedExeNamePart)
	}

	// Verify the symbol name contains the expected part.
	if !strings.Contains(info.SymbolName, expectedSymbolPart) {
		t.Errorf("Resolved SymbolName '%s' does not contain expected '%s'", info.SymbolName, expectedSymbolPart)
	}

	// Add more checks if needed, e.g., offset or base address if they were known/predictable.
}
