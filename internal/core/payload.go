package core

type NASAResponse struct {
	Total            int                   `json:"element_count"`
	NearEarthObjects map[string][]NEObject `json:"near_earth_objects"`
}

// NEObject is Near Earth Object
type NEObject struct {
	// Date                           time.Time `json:"date"`
	Id                             string `json:"id"`
	Name                           string `json:"name"`
	IsPotentiallyHazardousAsteroid bool   `json:"is_potentially_hazardous_asteroid"`
}
