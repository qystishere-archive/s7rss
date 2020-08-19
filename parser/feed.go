package parser

import (
	"time"
)

type Feed struct {
	URL        string
	Processing bool

	UpdatedAt time.Time
}
