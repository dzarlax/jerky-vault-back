package constants

import "testing"

func TestIsValidOrderStatus(t *testing.T) {
	validStatuses := []string{
		OrderStatusNew,
		OrderStatusInProgress,
		OrderStatusReady,
		OrderStatusFinished,
		OrderStatusCanceled,
	}

	for _, status := range validStatuses {
		if !IsValidOrderStatus(status) {
			t.Fatalf("expected %q to be valid", status)
		}
	}

	invalidStatuses := []string{"", "done", "cancelled", "new "}
	for _, status := range invalidStatuses {
		if IsValidOrderStatus(status) {
			t.Fatalf("expected %q to be invalid", status)
		}
	}
}
