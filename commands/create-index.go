package commands

import (
	"fmt"

	"github.com/astroband/astrologer/es"
)

// CreateIndex calls create-indexes command
func CreateIndex() {
	es.CreateIndicies()
	fmt.Println("Indicies created successfully!")
}
