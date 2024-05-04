package core

import "time"

type CliResponse struct {
	Total            int           `json:"total"`
	NearEarthObjects []CliNEObject `json:"near_earth_objects"`
}

type CliNEObject struct {
	Date                           time.Time `json:"date"`
	Id                             string    `json:"id"`
	Name                           string    `json:"name"`
	IsPotentiallyHazardousAsteroid bool      `json:"is_potentially_hazardous_asteroid"`
}
