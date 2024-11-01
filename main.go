package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Pokemon struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Stats   Stats    `json:"stats"`
	Types   []string `json:"types"`
	PodInfo PodInfo  `json:"pod_info"`
}

type Stats struct {
	HP        int `json:"hp"`
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	SpAttack  int `json:"special-attack"`
	SpDefense int `json:"special-defense"`
	Speed     int `json:"speed"`
}

type PodInfo struct {
	IP       string
	Hostname string
}

type PokeAPIResponse struct {
	Name    string `json:"name"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func getPodInfo() PodInfo {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return PodInfo{
		IP:       hostname,
		Hostname: hostname,
	}
}

func main() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle the main route
	http.HandleFunc("/", handleHome)

	log.Printf("Server starting on :8080... (Pod: %s)", getPodInfo().Hostname)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	// Log incoming request
	log.Printf("Received request from %s for path %s", r.RemoteAddr, r.URL.Path)

	// Get random Pokemon data
	pokemon, err := getRandomPokemon()
	if err != nil {
		log.Printf("Error getting Pokemon: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add pod info
	pokemon.PodInfo = getPodInfo()
	log.Printf("Serving Pokemon %s from pod %s", pokemon.Name, pokemon.PodInfo.Hostname)

	// Parse and execute template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pokemon)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getRandomPokemon() (*Pokemon, error) {
	rand.Seed(time.Now().UnixNano())
	pokemonID := rand.Intn(898) + 1 // There are 898 Pokemon in the API

	resp, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", pokemonID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp PokeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	// Convert API response to our Pokemon struct
	pokemon := &Pokemon{
		Name:  capitalize(apiResp.Name),
		Image: apiResp.Sprites.FrontDefault,
		Stats: Stats{},
		Types: make([]string, 0),
	}

	// Extract stats
	for _, stat := range apiResp.Stats {
		switch stat.Stat.Name {
		case "hp":
			pokemon.Stats.HP = stat.BaseStat
		case "attack":
			pokemon.Stats.Attack = stat.BaseStat
		case "defense":
			pokemon.Stats.Defense = stat.BaseStat
		case "special-attack":
			pokemon.Stats.SpAttack = stat.BaseStat
		case "special-defense":
			pokemon.Stats.SpDefense = stat.BaseStat
		case "speed":
			pokemon.Stats.Speed = stat.BaseStat
		}
	}

	// Extract types
	for _, t := range apiResp.Types {
		pokemon.Types = append(pokemon.Types, t.Type.Name)
	}

	return pokemon, nil
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
