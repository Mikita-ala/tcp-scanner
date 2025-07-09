package request

type ScanStatus string

const (
	ScanStatusSuccess ScanStatus = "success"
	ScanStatusPending ScanStatus = "pending"
)

type ScanResult struct {
	Status ScanStatus `json:"status"`
	Result []int      `json:"result,omitempty"`
	Error  string     `json:"error,omitempty"`
}

type ScanRequest struct {
	Host        string `json:"host"`
	StartPort   int    `json:"startPort"`
	EndPort     int    `json:"endPort"`
	Concurrency int    `json:"concurrency"`
	Timeout     int    `json:"timeout"` // в миллисекундах
}
