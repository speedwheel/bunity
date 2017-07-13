package general

import (
	//"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	customsearch "google.golang.org/api/customsearch/v1"
	"io/ioutil"
	"log"
	"app/config"
	//"github.com/kr/pretty"
)


func Run(startV int, id string, q string) ([]*customsearch.Result, int64) {
	data, err := ioutil.ReadFile(config.GetAppPath()+"config/search-key.json")
	if err != nil {
		log.Fatal(err)
	}
	//Get the config from the json key file with the correct scope
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/cse")
	if err != nil {
		log.Fatal(err)
	}

	// Initiate an http.Client. The following GET request will be
	// authorized and authenticated on the behalf of
	// your service account.
	client := conf.Client(oauth2.NoContext)

	cseService, err := customsearch.New(client)
	search := cseService.Cse.List(q)
	search.Cx(id)

	//Thinking about searching a particular place?
	//search.Gl("Chattanooga, TN")
	
	start := int64(startV)
	search.Start(start)
	call, err := search.Do()
	if err != nil {
		log.Fatal(err)
	}
	
	return call.Items, call.SearchInformation.TotalResults
}