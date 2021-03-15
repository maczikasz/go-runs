package server

import (
	"encoding/json"
	"github.com/maczikasz/go-runs/internal/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Writes a value to the passed ResponseWriter
// will return an error in case of any exceptions
// adds Content-Type header and 200 status code
// in case of error adds 500 status code
func WriteJsonResponse(r http.ResponseWriter, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		http.Error(r, "failed to write json", 500)
		return errors.Wrap(err, "failed to write json")
	}

	r.Header().Add("Content-type", "application/json")

	_, err = r.Write(data)

	if err != nil {
		log.Warn("failed to write json to output")
	}

	return nil
}

// Handles an error where we fetch data from a database
// Will write status 404 if the DataNotFoundError is passed
// Will write status 500 if any other error is returned
// the return value indicates whether any error was passed to it
// usage:
// data, err := FindSomeDataById(stepId)
//
//	err = HandleDataError(writer, request, err)
//	if err != nil {
//		return
//	}
// // use the data
func HandleDataError(writer http.ResponseWriter, request *http.Request, err error) error {
	if err != nil {
		if dnfError, ok := err.(*model.DataNotFoundError); ok {
			log.Debug(dnfError)
			http.NotFound(writer, request)
		} else {
			log.Errorf("failed to find data: %s", err)
			http.Error(writer, err.Error(), 500)
		}
		return err
	}
	return nil
}

