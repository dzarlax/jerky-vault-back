package constants

// Order statuses
const (
	OrderStatusNew        = "new"
	OrderStatusInProgress = "in_progress"
	OrderStatusReady      = "ready"
	OrderStatusFinished   = "finished"
	OrderStatusCanceled   = "canceled"
)

// IsValidOrderStatus reports whether status is one of the supported order states.
func IsValidOrderStatus(status string) bool {
	switch status {
	case OrderStatusNew, OrderStatusInProgress, OrderStatusReady, OrderStatusFinished, OrderStatusCanceled:
		return true
	default:
		return false
	}
}
