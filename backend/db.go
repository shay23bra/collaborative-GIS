package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

type Area struct {
	ID       int64       `json:"id"`
	Username string      `json:"username"`
	Name     string      `json:"name"`
	Points   [][]float64 `json:"points"`
}

func connectDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	var err error

	for retries := 5; retries > 0; retries-- {
		db, err = pgxpool.Connect(context.Background(), dbURL)
		if err == nil {
			return nil
		}
		log.Printf("Database connection failed: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("error connecting to database: %v", err)
}
func InitDB() error {
	err := connectDB()
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	err = createTables()
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}
	return nil
}

func createTables() error {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS areas (
        id SERIAL PRIMARY KEY,
		username VARCHAR NOT NULL,
        name TEXT NOT NULL,
        geom GEOMETRY(POLYGON, 4326) NOT NULL
    );
    `
	_, err := db.Exec(context.Background(), createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}
	return nil
}

func SaveArea(area Area) error {
	var polygon string
	for i, point := range area.Points {
		if i > 0 {
			polygon += ", "
		}
		polygon += fmt.Sprintf("%f %f", point[1], point[0]) // (lng, lat) format
	}
	polygonWKT := fmt.Sprintf("POLYGON((%s))", polygon)

	query := `
    INSERT INTO areas (name, username, geom)
    VALUES ($1, $2, ST_GeomFromText($3, 4326));
    `
	_, err := db.Exec(context.Background(), query, area.Name, area.Username, polygonWKT)
	if err != nil {
		return fmt.Errorf("error inserting area: %v", err)
	}
	return nil
}

func UpdateArea(area Area) error {
	var polygon string
	for i, point := range area.Points {
		if i > 0 {
			polygon += ", "
		}
		polygon += fmt.Sprintf("%f %f", point[1], point[0]) // (lng, lat) format
	}
	polygonWKT := fmt.Sprintf("POLYGON((%s))", polygon)
	query := `UPDATE areas SET geom = ST_GeomFromText($1, 4326) WHERE id = $2 AND username = $3`
	_, err := db.Exec(context.Background(), query, polygonWKT, area.ID, area.Username)
	if err != nil {
		log.Println("Error updating polygon:", err)
		return err
	}
	return nil
}

func GetAreasWithinBounds(minLat, minLng, maxLat, maxLng float64) ([]Area, error) {
	query := `
    SELECT id, username, name, ST_AsText(geom) AS geom_wkt
    FROM areas
    WHERE geom && ST_MakeEnvelope($1, $2, $3, $4, 4326);
    `

	rows, err := db.Query(context.Background(), query, minLng, minLat, maxLng, maxLat)
	if err != nil {
		return nil, fmt.Errorf("error querying areas: %v", err)
	}
	defer rows.Close()

	var areas []Area
	for rows.Next() {
		var area Area
		var geomWKT string
		if err := rows.Scan(&area.ID, &area.Username, &area.Name, &geomWKT); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		area.Points = parseWKTToPoints(geomWKT)
		areas = append(areas, area)
	}
	return areas, nil
}

func parseWKTToPoints(wkt string) [][]float64 {
	var points [][]float64
	wkt = wkt[len("POLYGON((") : len(wkt)-2]
	pairs := strings.Split(wkt, ", ")
	for _, pair := range pairs {
		var lng, lat float64
		fmt.Sscanf(pair, "%f %f", &lng, &lat)
		points = append(points, []float64{lat, lng})
	}
	return points
}
