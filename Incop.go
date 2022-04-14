package main

type Incop struct {
	key  string
	data *Data
}

func NewIncop(key string, data *Data) *Incop {
	return &Incop{key, data}
}
