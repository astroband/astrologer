package es

// AccountThresholds represents account thresholds from XDR
type AccountThresholds struct {
	Low    *byte `json:"low,omitempty"`
	Medium *byte `json:"medium,omitempty"`
	High   *byte `json:"high,omitempty"`
	Master *byte `json:"master,omitempty"`
}
