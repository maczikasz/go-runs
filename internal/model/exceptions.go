package model

import "fmt"

func CreateDataNotFoundError(dataType string, id string) error {
	return DataNotFoundError{
		dataType: dataType,
		id:       id,
	}
}

type DataNotFoundError struct {
	dataType string
	id       string
}

func (d DataNotFoundError) Error() string {
	return fmt.Sprintf("could not found %s with Id %s", d.dataType, d.id)
}
