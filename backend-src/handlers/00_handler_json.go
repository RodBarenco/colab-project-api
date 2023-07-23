package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/datatypes"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code >= 500 {
		log.Println("Responding with 5XX error:", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(toGORMJSON(payload))
	if err != nil {
		log.Printf("Failed to marshal JSON response %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func toGORMJSON(data interface{}) interface{} {
	switch v := data.(type) {
	case datatypes.JSON:
		var result interface{}
		if err := v.UnmarshalJSON([]byte(v)); err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
		}
		return result
	case []datatypes.JSON:
		var result []interface{}
		for _, d := range v {
			var r interface{}
			if err := d.UnmarshalJSON([]byte(d)); err != nil {
				log.Printf("Failed to unmarshal JSON: %v", err)
				continue
			}
			result = append(result, r)
		}
		return result
	default:
		return data
	}
}
