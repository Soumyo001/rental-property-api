package models

type BreadCrumb struct {
	LocationID string   `json:"LocationID"`
	Name       string   `json:"Name"`
	Type       string   `json:"Type"`
	Slug       string   `json:"Slug"`
	Display    []string `json:"Display"`
}

type GeoInfo struct {
	Breadcrumbs []BreadCrumb `json:"Breadcrumbs"`
	City        string       `json:"City"`
	Country     string       `json:"Country"`
	CountryCode string       `json:"CountryCode"`
	Name        string       `json:"Name"`
	LocationID  string       `json:"LocationID"`
	Lat         float64      `json:"Lat"`
	Lon         float64      `json:"Lon"`
	State       string       `json:"State"`
	StateAbbr   *string      `json:"StateAbbr"`
}

type Property struct {
	Amenities    []string      `json:"Amenities"`
	Name         string        `json:"Name"`
	Slug         string        `json:"Slug"`
	PropertyType string        `json:"PropertyType"`
	Price        float64       `json:"Price"`
	ReviewScore  float64       `json:"ReviewScore"`
	StarRating   int           `json:"StarRating"`
	Counts       PropertyCount `json:"Counts"`
	Image        PropertyImage `json:"Image"`
}

type PropertyCount struct {
	Bathroom  int `json:"Bathroom"`
	Bedroom   int `json:"Bedroom"`
	Reviews   int `json:"Reviews"`
	Occupancy int `json:"Occupancy"`
}

type PropertyImage struct {
	Count  int      `json:"Count"`
	Images []string `json:"Images"`
}

type PropertyResponse struct {
	ID        string   `json:"ID"`
	Feed      int      `json:"Feed"`
	GeoInfo   GeoInfo  `json:"GeoInfo"`
	Property  Property `json:"Property"`
	Published bool     `json:"Published"`
}

type ListResult struct {
	Count int                `json:"Count"`
	Items []PropertyResponse `json:"Items"`
}

type ListResponse struct {
	Result ListResult `json:"Result"`
}

type ErrorResponse struct {
	Error string `json:"Error"`
}
