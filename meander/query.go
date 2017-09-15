package meander

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// APIKey holds our Google Places API key
var APIKey string

// Place will hold the data from the Google Places API response
type Place struct {
	*googleGeometry `json:"geometry"`
	Name            string         `json:"name"`
	Icon            string         `json:"icon"`
	Photos          []*googlePhoto `json:"photos"`
	Vicinity        string         `json:"vicinity"`
}

// Public returns the public view of a Place
func (p *Place) Public() interface{} {
	return map[string]interface{}{
		"name":     p.Name,
		"icon":     p.Icon,
		"photos":   p.Photos,
		"vicinity": p.Vicinity,
		"lat":      p.Lat,
		"lng":      p.Lng,
	}
}

type googleResponse struct {
	Results []*Place `json:"results"`
}

type googleGeometry struct {
	*googleLocation `json:"location"`
}

type googleLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type googlePhoto struct {
	PhotoRef string `json:"photo_reference"`
	URL      string `json:"url"`
}

// Query contains all the information needed to consruct
// a query to send to the Google Places API
type Query struct {
	Lat          float64
	Lng          float64
	Journey      []string
	Radius       int
	CostRangeStr string
}

func (q *Query) find(t string) (*googleResponse, error) {
	u := "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	vals := make(url.Values)
	vals.Set("key", APIKey)
	vals.Set("location", fmt.Sprintf("%g,%g", q.Lat, q.Lng))
	vals.Set("radius", fmt.Sprintf("%d", q.Radius))
	vals.Set("type", t)
	if len(q.CostRangeStr) > 0 {
		r := ParseCostRange(q.CostRangeStr)
		vals.Set("minprice", fmt.Sprintf("%d", int(r.From)-1))
		vals.Set("maxprice", fmt.Sprintf("%d", int(r.To)-1))
	}

	res, err := http.Get(u + "?" + vals.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response googleResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Run makes concurrent requests to the Google Places API and
// randomly selects a result for each place type in a journey.
func (q *Query) Run() []interface{} {
	rand.Seed(time.Now().UnixNano())

	var w sync.WaitGroup
	var l sync.Mutex
	places := make([]interface{}, len(q.Journey))

	for i, t := range q.Journey {
		w.Add(1)
		go func(t string, i int) {
			defer w.Done()
			response, err := q.find(t)
			if err != nil {
				log.Println("Failed to find places:", err)
				return
			}
			if len(response.Results) == 0 {
				log.Println("No places found for", t)
				return
			}
			for _, result := range response.Results {
				for _, photo := range result.Photos {
					photo.URL = "https://maps.googleapis.com/maps/api/place/photo?maxwidth=1000&photoreference=" + photo.PhotoRef + "&key=" + APIKey
				}
			}
			randI := rand.Intn(len(response.Results))
			l.Lock()
			places[i] = response.Results[randI]
			l.Unlock()
		}(t, i)
	}
	w.Wait()
	return places
}
