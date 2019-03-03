package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Api struct {
	apiUrl      string
	apiKey      string
	timestamp   time.Time
	timeSetFlag bool
}

type Member struct {
	Id           float64 `json:"id"`
	Key          float64 `json:"key"`
	Hash         string  `json:"hash"`
	Irc_name     string  `json:"irc_name"`
	Spoken_name  string  `json:"spoken_name"`
	Added_by     float64 `json:"added_by"`
	Date_created string  `json:"date_created"`
	Last_login   string  `json:"last_login"`
	Admin        float64 `json:"admin"`
	Active       float64 `json:active`
	User_id      float64 `json:user_id`
	Created_at   string  `json:"created_at"`
	Updated_at   string  `json:"updated_at"`
}

type Data struct {
	Timestamp struct {
		Date          string  `json:"date"`
		Timezone_type float64 `json:"timezone_type"`
		Timezone      string  `json:"timezone"`
	} `json:"timestamp"`
	Members []Member `json:"members"`
}

type LoginAttempt struct {
	Key       float64
	Timestamp string
	reason    string //success || bad key || bad password
	result    string //success || failure
}

func New(apiUrl string, apiKey string) *Api {
	return &Api{
		apiUrl:      apiUrl,
		apiKey:      apiKey,
		timeSetFlag: false,
	}
}

//Updates Sqlite database if the timestamp has changed in the JSON message. It always
//updates on first run.
func (a *Api) CheckForUpdates() (updateRequired bool, getDat Data) {
	getDat = a.GetUsers()
	t, _ := time.Parse("2006-01-02 15:04:05.000000", getDat.Timestamp.Date)

	if t.Equal(a.timestamp) {
		//fmt.Printf("Database update not required\n")
		updateRequired = false
	} else {
		a.timestamp = t
		updateRequired = true
	}

	if a.timeSetFlag == false {
		a.timestamp = t
		a.timeSetFlag = true
	}

	return

}

//Retreves and fills the Data structure of member data.
func (a *Api) GetUsers() (memberdata Data) {

	httpclient := &http.Client{}
	req, err := http.NewRequest("GET", a.apiUrl, nil)
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	resp, err := httpclient.Do(req)
	if err != nil {
		log.Println("Server did not respond")
		log.Println(err)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Issue digesting server response")
		log.Println(err)
		return
	}

	err = json.Unmarshal(b, &memberdata)
	if err != nil {
		log.Println("Unmasrshaling error")
		log.Println(err)
		return
	}

	for l := range memberdata.Members {
		fmt.Println(memberdata.Members[l].Spoken_name)
	}

	return
}
