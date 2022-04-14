package main

type StateQueryResponse struct {
	end     bool
	incList []*Incop
}

func NewStateQueryResponse(end bool, incList []*Incop) *StateQueryResponse {
	return &StateQueryResponse{end, incList}
}
