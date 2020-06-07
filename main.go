package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type Event struct {
	Id          graphql.ID `json:"id"`
	Description string     `json:"description"`
	Location    string     `json:"location"`
	Time        string     `json:"time"`
	Cost        float64    `json:"cost"`
	Pictures    []string   `json:"pictures"`
}

type RootResolver struct{}

func getEventsfromJson() []Event {
	file, _ := ioutil.ReadFile("./events.json")
	var events []Event
	_ = json.Unmarshal([]byte(file), &events)
	return events
}

func getSchemaFromFile() string {
	file, _ := ioutil.ReadFile("./schema.gql")
	return string(file)
}

func (r *RootResolver) Events() ([]Event, error) {
	return getEventsfromJson(), nil
}

func (r *RootResolver) Event(args struct{ Id graphql.ID }) (Event, error) {
	for _, event := range getEventsfromJson() {
		if args.Id == event.Id {
			return event, nil
		}
	}
	return Event{}, nil
}

func main() {
	var opts = []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	var schema = graphql.MustParseSchema(getSchemaFromFile(), &RootResolver{}, opts...)
	http.Handle("/", &relay.Handler{Schema: schema})
	log.Println(http.ListenAndServe("127.0.0.1:8080", nil))
}
