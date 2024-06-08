package geo

func IsValidLatitude(lat float64) bool {
	return lat >= -90 && lat <= 90
}

func IsValidLongitude(long float64) bool {
	return long >= -180 && long <= 180
}
