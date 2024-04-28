package models

import "time"

type FilterOptions struct {
	Page      uint64
	Limit     uint64
	StartDate *time.Time
	EndDate   *time.Time
	Author    *string
}
