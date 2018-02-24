package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HostStatus struct {
	Hostname string `json:"hostname"`
	Status   string `json:"status"`
}
type HostStatusTemplate struct {
	Hostname string
	OK       string
	Down     string
}

var defaultTemplate = `
<!DOCTYPE html>

<html>
    <head>
        <meta charset="utf-8">
        <title>Status page</title>
        <!-- UIkit CSS -->
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/uikit/3.0.0-beta.40/css/uikit.min.css" />
        <!-- UIkit JS -->
        <script src="https://cdnjs.cloudflare.com/ajax/libs/uikit/3.0.0-beta.40/js/uikit.min.js"></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/uikit/3.0.0-beta.40/js/uikit-icons.min.js"></script>
    </head>
    <body>
        <h1 class="uk-heading-primary"><span uk-icon="server"></span> Infrastructure status</h1>
        <div class="uk-container uk-container-large uk-background-muted">
                <ul class="uk-grid-small uk-child-width-1-2 uk-child-width-1-4@s uk-text-center" uk-sortable="handle: .uk-card" uk-grid>
                    {{range .}}
                        <li>
                            <div class="uk-card uk-card-default uk-card-body">
                            <p>{{.Hostname}}</p>
                            <p>
                            {{if .OK}}
                            <span class="uk-label uk-label-success">{{.OK}}</span>
                            {{end}}
                            {{if .Down}}
                            <span class="uk-label uk-label-warning">{{.Down}}</span>
                            {{end}}
                            </p>
                        </div>
                        </li> 
                    {{end}}
                </ul>
        </div>
    </body>
</html>
`
var (
	AgentService string
)

func main() {
	flag.StringVar(&AgentService, "agent", "localhost:18080", "The hostname and port for the agent service")
	var listenPort = flag.String("port", "18081", "This is the port the application will listen on")
	flag.Parse()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/status/{token}", getHosts).Methods("GET")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *listenPort), router))
}
func getHosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp, err := http.Get(fmt.Sprintf("http://%s/agent/%s", AgentService, vars["token"]))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var hosts []HostStatus
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&hosts)
	//fmt.Println("%+v", hosts)
	var thosts []HostStatusTemplate
	for _, h := range hosts {
		if h.Status == "DOWN" {
			thosts = append(thosts, HostStatusTemplate{Hostname: h.Hostname, Down: h.Status})
		} else {
			thosts = append(thosts, HostStatusTemplate{Hostname: h.Hostname, OK: h.Status})
		}
	}
	t := template.Must(template.New("").Parse(defaultTemplate))
	t.ExecuteTemplate(w, "", thosts)
}
