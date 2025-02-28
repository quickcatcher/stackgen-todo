package controller

type CreateTodoItemResponse struct {
	Id int `json:"int"`
}

type GetAllTodoItemsResponse struct {
	Items []*Item `json:"items"`
}
