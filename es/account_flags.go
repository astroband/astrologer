package es

// AccountFlags represents account flags from XDR
type AccountFlags struct {
	AuthRequired  bool `json:"required,omitempty"`
	AuthRevocable bool `json:"revocable,omitempty"`
	AuthImmutable bool `json:"immutable,omitempty"`
}
