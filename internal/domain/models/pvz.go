package models

import (
	"time"

	"github.com/google/uuid"
)

type City string

const (
	CityMoscow          City = "Москва"
	CitySaintPetersburg City = "Санкт-Петербург"
	CityKazan           City = "Казань"
)

type PVZ struct {
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	City             City      `json:"city"`
	CreatedAt        time.Time `json:"created_at"`
}

func NewPVZ(city City) *PVZ {
	now := time.Now()
	return &PVZ{
		ID:               uuid.New(),
		RegistrationDate: now,
		City:             city,
		CreatedAt:        now,
	}
}

func IsValidCity(city City) bool {
	return city == CityMoscow || city == CitySaintPetersburg || city == CityKazan
}
