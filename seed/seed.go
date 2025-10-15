package main

import (
	"fmt"
	"os"
	"encoding/json"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/musorok/server/internal/domain"
)

func main() {
	_ = godotenv.Load()
	dsn := getenv("DB_DSN", "postgres://postgres:postgres@localhost:5432/musorok?sslmode=disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil { panic(err) }

	admin := domain.User{ Name:"Admin", Phone:"+77070000000", Role:domain.RoleAdmin }
	db.Where("phone = ?", admin.Phone).FirstOrCreate(&admin)

	test := domain.User{ Name:"Test User", Phone:"+77070000001", Role:domain.RoleUser }
	db.Where("phone = ?", test.Phone).FirstOrCreate(&test)

	// Polygon 4YOU (примерной формы, внутри верх Алматы)
	poly := domain.Polygon{ Name:"ЖК 4YOU", City:"Алматы", IsActive:true }
	gj := map[string]interface{}{
		"type":"Polygon",
		"coordinates":[]interface{}{
			[]interface{}{
				[]float64{76.9100,43.2185},
				[]float64{76.9180,43.2185},
				[]float64{76.9180,43.2230},
				[]float64{76.9100,43.2230},
				[]float64{76.9100,43.2185},
			},
		},
	}
	b, _ := json.Marshal(gj)
	poly.GeoJSON = string(b)
	db.Where("name = ?", poly.Name).FirstOrCreate(&poly)

	addr := domain.Address{ UserID: test.ID, City:"Алматы", Street:"Каскеленская", House:"1", Entrance:"1", Floor:"1", Apartment:"1", Lat:43.2200, Lng:76.9140, IsDefault:true }
	addr.PolygonID = &poly.ID; name := poly.Name; addr.PolygonName = &name
	db.Where("user_id = ? AND apartment = ?", test.ID, addr.Apartment).FirstOrCreate(&addr)

	fmt.Println("Seed completed at", time.Now())
}

func getenv(k, d string) string { if v:=os.Getenv(k); v!="" { return v }; return d }
