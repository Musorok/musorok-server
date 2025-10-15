package geospatial

import "encoding/json"

type Position = [2]float64
type Polygon = [][]Position
type MultiPolygon = []Polygon

type geo struct {
	Type string `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

func pointInPolygon(pt Position, poly Polygon) bool {
	x := pt[0]; y := pt[1]
	inside := false
	for _, ring := range poly {
		for i, j := 0, len(ring)-1; i < len(ring); j, i = i, i+1 {
			xi, yi := ring[i][0], ring[i][1]
			xj, yj := ring[j][0], ring[j][1]
			intersect := ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi)
			if intersect { inside = !inside }
		}
	}
	return inside
}

func ContainsPoint(geojson string, lng, lat float64) bool {
	var g geo
	if err := json.Unmarshal([]byte(geojson), &g); err != nil { return false }
	switch g.Type {
	case "Polygon":
		var p Polygon
		if json.Unmarshal(g.Coordinates, &p) == nil {
			return pointInPolygon([2]float64{lng, lat}, p)
		}
	case "MultiPolygon":
		var mp MultiPolygon
		if json.Unmarshal(g.Coordinates, &mp) == nil {
			for _, p := range mp {
				if pointInPolygon([2]float64{lng, lat}, p) { return true }
			}
		}
	}
	return false
}
