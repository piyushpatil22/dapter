package dap

import "fmt"

var (
	ErrDapTableTagNotFound = fmt.Errorf("dap table tag not found")
	ErrMultiplePKFieldFound = fmt.Errorf("multiple primary key found inside same struct")
	ErrDapFieldAttrsTagNotFound = fmt.Errorf("dap field attrs tag not found")
    ErrNoRowsFound = fmt.Errorf("no rows found")
)
