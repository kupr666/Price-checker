package domains

import (
	"time"	
)

type Item struct {
	ID				int64
	URL 			string
	CurrentPrice	float64
	TargetPrice 	float64
	LastChecked 	time.Time
}

