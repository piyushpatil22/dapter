package dapter

import "fmt"

var (
	ErrDapTableTagNotFound = fmt.Errorf("dap table tag not found")
	ErrMultiplePKFieldFound = fmt.Errorf("multiple primary key found inside same struct")
	ErrDapFieldAttrsTagNotFound = fmt.Errorf("dap field attrs tag not found")
)
