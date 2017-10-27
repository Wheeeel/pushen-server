package model

import "time"

type Timestamp struct {
	UpdateTimestamp time.Time `json:"update_timestamp"`
	CreateTimestamp time.Time `json:"create_timestamp"`
}

func (ts *Timestamp) BeforeCreate() {
	ts.CreateTimestamp = time.Now()
	ts.UpdateTimestamp = time.Now()
}

func (ts *Timestamp) BeforeUpdate() {
	ts.UpdateTimestamp = time.Now()
}
