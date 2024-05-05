package core

type NasaAPIResponse struct {
	Total            int                       `json:"element_count"`
	NearEarthObjects map[string][]NasaNEObject `json:"near_earth_objects"`
}

// NasaNEObject is Near Earth Object
type NasaNEObject struct {
	Id                             string `json:"id"`
	Name                           string `json:"name"`
	IsPotentiallyHazardousAsteroid bool   `json:"is_potentially_hazardous_asteroid"`
}
