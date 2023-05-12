package gsrender

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type TextStandard int

const (
	JSON TextStandard = 0
	XML TextStandard = 1
)

func Write(rw http.ResponseWriter, status uint, ts TextStandard, v interface{}) {
	switch ts {
	case JSON:
		WriteJSON(rw, status, v)
	case XML:
		WriteXML(rw, status, v)
	}
}

// WriteJSON encodes v to json and writes it to the given http.ResponseWriter 
// with the provided status and Content-Type header with application/json value.
func WriteJSON(rw http.ResponseWriter, status uint, v interface{}) {
	json, _ := json.Marshal(v)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(int(status))
	rw.Write(json)
}

// WriteXML encodes v to xml and writes it to the given http.ResponseWriter
// with status code and Content-Type header with text/xml value.
func WriteXML(rw http.ResponseWriter, status uint, v interface{}) {
	xml, _ := xml.Marshal(v)
	rw.Header().Set("Content-Type", "text/xml")
	rw.WriteHeader(int(status))
	rw.Write(xml)
}

// writes an http status to the ResponseWriter
func Status(rw http.ResponseWriter, status int) {
	rw.WriteHeader(status)
}