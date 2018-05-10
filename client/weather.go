package main

import (
	"encoding/json"
	"flag"
	owm "github.com/briandowns/openweathermap"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

//Finner IP/hvis man skal ha her
const URL = "http://ip-api.com/json"


const weatherTemplate = `Current weather for {{.Name}}:
    Conditions: 	{{range .Weather}} {{.Description}} {{end}}
    Now:         	{{.Main.Temp}} {{.Unit}}
	Humidity:    	{{.Main.Humidity}} {{.Unit}}
	Wind:			{{.Wind.Speed}} {{.Unit}}
	Wind deg:		{{.Wind.Deg}} {{.Unit}}

`

var (
	whereFlag = flag.String("w", "", "Location to get weather.  If location has a space, wrap the location in double quotes.")
	unitFlag  = flag.String("u", "", "Unit of measure to display temps in")
	langFlag  = flag.String("l", "", "Language to display temps in")
)

type Data struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	City        string  `json:"city"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`

}

// getLocation will get the location details for where this
// application has been run from.
func getLocation() *Data {
	response, err := http.Get(URL)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	r := &Data{}
	err = json.Unmarshal(result, &r)
	if err != nil {
		log.Fatalln(err)
	}
	return r
}


// getCurrent får gjeldende vær for gitt sted i de målene som er gitt
func getCurrent(l, u, lang string) *owm.CurrentWeatherData {
	w, err := owm.NewCurrent(u, lang, os.Getenv("OWM_API_KEY"))
	if err != nil {
		log.Fatalln(err)
	}
	w.CurrentByName(l)
	return w
}

func main() {

	os.Setenv("OWM_API_KEY", "81e8da958c34767cf9621033d5b47ab7")

	flag.Parse()

	// error håndtering
	if len(*whereFlag) <= 1 || len(*unitFlag) != 1 || len(*langFlag) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	// Process request for location of "here"
	if strings.ToLower(*whereFlag) == "here" {
		w := getCurrent(getLocation().City, *unitFlag, *langFlag)
		tmpl, err := template.New("weather").Parse(weatherTemplate)
		if err != nil {
			log.Fatalln(err)
		}

		// Laster template
		err = tmpl.Execute(os.Stdout, w)
		if err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}

	// Process request for the given location
	w := getCurrent(*whereFlag, *unitFlag, *langFlag)
	tmpl, err := template.New("weather").Parse(weatherTemplate)
	if err != nil {
		log.Fatalln(err)
	}

	// Laster template
	err = tmpl.Execute(os.Stdout, w)
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(0)
}
