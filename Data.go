package main

type Data struct {
	createdOn int
	value     int
}

func NewData(createdOn, value int) *Data {
	return &Data{
		createdOn: createdOn,
		value:     value,
	}
}
