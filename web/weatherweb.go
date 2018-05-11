

// go get github.com/briandowns/openweathermap
// kjør den i terminalen jokke og jostein





// Kjøres igjennom "http://localhost:8001/here"


package main

import (
	"encoding/json"
	"html/template"
	owm "github.com/briandowns/openweathermap"
	// "io/ioutil"
	"log"
	"net/http"
	"os"
)

// URL for å finne brukerens IP
const URL = "http://ip-api.com/json"

// Data will hold the result of the query to get the IP
// address of the caller.

type Data struct {
	Status      string  `json:"status"`
	CountryCode string  `json:"countryCode"`
	City        string  `json:"city"`
}


// getlocation skaffer detaljer på hvor applikasjonen har blitt kjørt ifra
func getLocation() (*Data, error) {
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	r := &Data{}
	if err = json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}

// getCurrent gets the current weather for the provided location in
// the units provided.
func getCurrent(l, u, lang string) *owm.CurrentWeatherData {
	w, err := owm.NewCurrent(u, lang, os.Getenv("OWM_API_KEY")) // Create the instance with the given unit
	if err != nil {
		log.Fatal(err)
	}
	w.CurrentByName("Bergen, NO") // Setter plasseringen på bynavn
	return w

}

// hereHandler will take are of requests coming in for the "/here" route.
func Handler(w http.ResponseWriter, r *http.Request) {
	location, err := getLocation()
	if err != nil {
		log.Fatal(err)
	}
	wd := getCurrent(location.City, "c", "en")

	// Process our template
	t, err := template.ParseFiles("templates/here.html")
	if err != nil {
		log.Fatal(err)
	}
	// We're doin' naughty things below... Ignoring icon file size and possible errors.
	_, _ = owm.RetrieveIcon("static/img", wd.Weather[0].Icon+".png")

	// Write out the template with the given data
	t.Execute(w, wd)
	t.Execute(w, r)
}




// Run the app
func main() {

	//api Key
	os.Setenv("OWM_API_KEY", "81e8da958c34767cf9621033d5b47ab7")

	//handler
	http.HandleFunc("/here", Handler)

	// Handler til ikonene
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.ListenAndServe(":8001", nil)
}
