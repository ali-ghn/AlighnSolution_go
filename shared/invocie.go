package shared

const (
	InvoiceStatusSettled            = "Settled"
	InvoiceStatusNew                = "New"
	InvoiceStatusNewPaidPartial     = "New (paidPartial)"
	InvoiceStatusExpired            = "Expired"
	InvoiceStatusExpiredPaidPartial = "Expired (paidPartial)"
	InvoiceStatusExpiredPaidLate    = "Expired (paidLate)"
	InvoiceStatusSettledPaidOver    = "Settled (paidOver)"
	InvoiceStatusProcessing         = "Processing"
	InvoiceStatusProcessingPaidOver = "Processing (paidOver)"
	InvoiceStatusSettledMarked      = "Settled (marked)"
	InvoiceStatusInvalid            = "Invalid"
	InvoiceStatusInvalidMarked      = "Invalid (marked)"
	InvoiceStatusInvalidPaidOver    = "Invalid (paidOver)"
)
