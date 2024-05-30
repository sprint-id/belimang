package geo

import "math"

type pos struct {
	φ float64 // latitude, radians
	ψ float64 // longitude, radians
}

const rEarth = 6372.8 // km

// Reference: https://rosettacode.org/wiki/Haversine_formula#Go
func Haversine(θ float64) float64 {
	return .5 * (1 - math.Cos(θ))
}

func DegPos(lat, lon float64) pos {
	return pos{lat * math.Pi / 180, lon * math.Pi / 180}
}

func HsDist(p1, p2 pos) float64 {
	return 2 * rEarth * math.Asin(math.Sqrt(Haversine(p2.φ-p1.φ)+
		math.Cos(p1.φ)*math.Cos(p2.φ)*Haversine(p2.ψ-p1.ψ)))
}

// CalculateDistance to calculate distance between two points
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	return HsDist(DegPos(lat1, lon1), DegPos(lat2, lon2))
}
