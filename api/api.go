package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
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
	Date    string   `json:"timestamp"`
	Members []Member `json:"members"`
}

type LoginAttempt struct {
	Key       int    `json:"key"`
	Timestamp string `json:"timestamp"`
	Reason    string `json:"reason"`
	Result    string `json:"result"`
}

func New(apiUrl string, apiKey string) *Api {
	return &Api{
		apiUrl:      apiUrl,
		apiKey:      apiKey,
		timeSetFlag: false,
	}
}

//CheckForUpdates updates the Sqlite database if the timestamp has changed in the JSON message. It always
//updates on first run.
func (a *Api) CheckForUpdates() (updateRequired bool, getDat Data) {
	getDat = a.GetUsers()
	t, _ := time.Parse("2006-01-02T15:04:05.000000Z", getDat.Date)

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

//GetUsers Retreves and fills the Data structure of member data.
func (a *Api) GetUsers() (memberdata Data) {

	//For some reason there is a return being tacked onto the apikey
	var bearer = "Bearer " + strings.TrimSpace(a.apiKey)

	req, err := http.NewRequest("GET", a.apiUrl+"/members", nil)

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")

	//Capture debug information
	dump, err := httputil.DumpRequestOut(req, true)

	httpclient := &http.Client{}
	resp, err := httpclient.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
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
		log.Println("SENT:\n" + string(dump) + "RETURNED:\n" + string(b))
		return
	}

	return
}

//SendLoginAttempt sends the result of a login to the web api
func (a *Api) SendLoginAttempt(key int, reason string, result string) {

	unixtime := strconv.FormatInt(time.Now().Unix(), 10)
	la := LoginAttempt{key, unixtime, reason, result}

	json, err := json.Marshal(la)
	if err != nil {
		log.Println("Masrshaling error")
		return
	}

	req, err := http.NewRequest("POST", a.apiUrl+"/login-attempt", bytes.NewBuffer(json))

	//For some reason there is a return being tacked onto the apikey
	var bearer = "Bearer " + strings.TrimSpace(a.apiKey)
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	httpclient := &http.Client{}
	resp, err := httpclient.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if strings.Contains(string(b), "true") == false {
		log.Println("Login attempt was not posted to the API Server")
	}

	return
}
