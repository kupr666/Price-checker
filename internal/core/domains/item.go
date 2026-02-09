package domains

import (
	"time"	
)

type Item struct {
	ID				int64 		`json:"id"`
	URL 			string 		`json:"url"`
	CurrentPrice	float64 	`json:"current_price"`
	TargetPrice 	float64 	`json:"target_price"`
	LastChecked 	time.Time 	`json:"last_checked"`
}

