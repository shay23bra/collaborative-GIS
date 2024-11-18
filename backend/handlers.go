package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func saveAreaHandler(w http.ResponseWriter, r *http.Request) {
	var area Area
	if err := json.NewDecoder(r.Body).Decode(&area); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := SaveArea(area)
	if err != nil {
		http.Error(w, "Failed to save area", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func updateAreaHandler(w http.ResponseWriter, r *http.Request) {
	var area Area
	if err := json.NewDecoder(r.Body).Decode(&area); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := UpdateArea(area)
	if err != nil {
		http.Error(w, "Failed to update area", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Polygon updated successfully"))
}

func getAreasInBoundsHandler(w http.ResponseWriter, r *http.Request) {
	minLat, _ := strconv.ParseFloat(r.URL.Query().Get("minLat"), 64)
	minLng, _ := strconv.ParseFloat(r.URL.Query().Get("minLng"), 64)
	maxLat, _ := strconv.ParseFloat(r.URL.Query().Get("maxLat"), 64)
	maxLng, _ := strconv.ParseFloat(r.URL.Query().Get("maxLng"), 64)

	areas, err := GetAreasWithinBounds(minLat, minLng, maxLat, maxLng)
	if err != nil {
		http.Error(w, "Failed to retrieve areas", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(areas)
}
