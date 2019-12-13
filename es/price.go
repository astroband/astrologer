package es

// Price represents Price struct from XDR
type Price struct {
	N int `json:"n"`
	D int `json:"d"`
}

func (p *Price) float() float64 {
	return float64(p.N) / float64(p.D)
}
