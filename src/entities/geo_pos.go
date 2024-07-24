package entities

type GeoPos struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (loc GeoPos) IsValid() bool {
	return loc.Latitude >= -90 && loc.Latitude <= 90 && loc.Longitude >= -180 && loc.Longitude <= 180

}
