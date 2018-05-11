

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

	func getMessage(temp int, ws float64, wt string) (string) {
		var message string
		if wt == "Regn" || wt == "Kraftig regn" || wt == "Kraftige regnbyger" {
			message = "I dag kan du late som du bor i Bergen"
			if temp > 15 {
				message += "PLASKEPARTY!!!!"
			} else if temp > 10 {
				message += "Regnjakke. Alt du trenger"
			} else if temp > 0 {
				message += " Perfekt vær å sove i"
			} else {
				message += " Arnold Schwarzenegger i Batman & Robin. Get it?"
			}
		} else if wt == "Skyet" || wt == "Lettskyet" || wt == "Delvis skyet" {
			message = "Solen er litt shy i dag"
			if temp > 20 {
				message += "Alt annet enn shorts i dag er ulovlig, jfr UD-21"
			} else if temp > 10 {
				message += " Nice temp, do what you want, brother"
			} else if temp > 0 {
				message += "Surt. Bare surt"
			} else {
				message += " You gonna freez boy!"
			}
		} else if wt == "Klarvær" || wt == "Sol" {
			message = "Suns out, guns out!."
			if temp > 20 {
				message += " Varmere enn djevelsens baller i dag."
			} else if temp > 10 {
				message += "Alt er lov i dag"
			} else if temp > 0 {
				message += " Ufyselig kaldt, men du overlever."
			} else {
				message += "Sibirtilstander i dag"
			}
		} else {
			message = ""
			if temp > 20 {
				message += "Spådom: Air Condition er din bestevenn i dag."
			} else if temp > 10 {
				message += "Hvis du er kul tar du på shorts, er du smart tar du på litt mer klær."
			} else if temp > 0 {
				message += "You gonna freez boy!."
			}
			if ws > 10 {
				message += " Det blåser mer enn Stavanger en sommer dag, på med allværsjakke!."
			}
		}
		return message
	}

	
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
