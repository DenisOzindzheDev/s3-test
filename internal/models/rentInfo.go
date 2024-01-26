package models

import "time"

type RentInfo struct {
	StartedAt    time.Time
	CompletedAt  time.Time
	ImagesBefore []string
	ImagesAfter  []string
}
