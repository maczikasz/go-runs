package model

import "fmt"

type Error struct {
	Name    string
	Message string
	Tags    []string
}

type DataNotFoundError struct {
	dataType string
	id       string
}

func CreateDataNotFoundError(dataType string, id string) error {
	return DataNotFoundError{
		dataType: dataType,
		id: id,
	}
}

func (d DataNotFoundError) Error() string {
	return fmt.Sprintf("could not found %s with Id %s", d.dataType, d.id)
}
