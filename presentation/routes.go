package presentation

import "daemon/domain"

func CreateRoutes(service domain.URLService) map[string]domain.RouteDef {
	URLH := NewURLHandler(service)
	return map[string]domain.RouteDef{
		"url/add": {
			Method:      "POST",
			HandlerFunc: URLH.Add,
		},
		"url/retrieve": {
			Method:      "GET",
			HandlerFunc: URLH.Retrieve,
		},
	}
}
