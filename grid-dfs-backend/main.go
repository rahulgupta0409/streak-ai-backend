package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type GridRequest struct {
	Start Point   `json:"start"`
	End   Point   `json:"end"`
	Grid  [][]int `json:"grid"`
}

func isValidPath(grid [][]int, x, y int) bool {
	return x >= 0 && y >= 0 && x < len(grid) && y < len(grid) && grid[x][y] == 0
}

func dfs(grid [][]int, x, y int, end Point, path *[]Point, vis map[Point]bool) bool {
	if x == end.X && y == end.Y {
		*path = append(*path, Point{x, y})
		return true
	}
	if !isValidPath(grid, x, y) || vis[Point{x, y}] {
		return false
	}
	directions := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for _, d := range directions {
		if dfs(grid, x+d.X, y+d.Y, end, path, vis) {
			return true
		}
	}
	*path = (*path)[:len(*path)-1]
	return false
}

func findPath(w http.ResponseWriter, r *http.Request) {
	var req GridRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	visited := make(map[Point]bool)
	var path []Point
	if dfs(req.Grid, req.Start.X, req.Start.Y, req.End, &path, visited) {
		w.Header().Set("Content-Type", "application/josn")
		json.NewEncoder(w).Encode(path)
	} else {
		http.Error(w, "No path found", http.StatusNotFound)
	}
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/find-path", findPath).Methods("POST", "OPTIONS")

	port := 8080
	fmt.Printf("Server listening on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), corsMiddleware(router))
}
