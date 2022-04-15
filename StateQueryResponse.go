package main

type StateQueryResponse struct {
	incList []*Incop
}

func NewStateQueryResponse(incList []*Incop) *StateQueryResponse {
	return &StateQueryResponse{incList}
}
