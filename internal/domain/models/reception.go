package models

import (
	"time"

	"github.com/google/uuid"
)

type ReceptionStatus string

const (
	ReceptionStatusInProgress ReceptionStatus = "in_progress"
	ReceptionStatusClose      ReceptionStatus = "close"
)

type Reception struct {
	ID        uuid.UUID       `json:"id"`
	DateTime  time.Time       `json:"date_time"`
	PVZID     uuid.UUID       `json:"pvz_id"`
	Status    ReceptionStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
}

func NewReception(pvzID uuid.UUID) *Reception {
	now := time.Now()
	return &Reception{
		ID:        uuid.New(),
		DateTime:  now,
		PVZID:     pvzID,
		Status:    ReceptionStatusInProgress,
		CreatedAt: now,
	}
}

func (r *Reception) Close() {
	r.Status = ReceptionStatusClose
}

func (r *Reception) IsInProgress() bool {
	return r.Status == ReceptionStatusInProgress
}
