package util

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
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
func HandleDataError(context *gin.Context, err error) error {
	if err != nil {
		if dnfError, ok := err.(*model.DataNotFoundError); ok {
			log.Debug("failed to find data: %s", dnfError)
			_ = context.Error(dnfError)
			context.Status(http.StatusNotFound)
			_ = context.Error(err)
		} else {
			log.Error(err)
			_ = context.Error(err)
			context.Status(http.StatusInternalServerError)
			_ = context.Error(err)
		}
		return err
	}
	return nil
}

func ToSet(origins []string) (result map[string]bool) {
	result = make(map[string]bool)
	for _, v := range origins {
		result[v] = true
	}

	return
}
