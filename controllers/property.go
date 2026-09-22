package controllers

import (
	"fmt"
	"net/http"
	"rental-property-api/models"
	"rental-property-api/services"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

type PropertyController struct {
	BaseController
}

func parseFloat(raw string, name string) (*float64, error) {
	if raw == "" {
		return nil, nil
	}

	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return nil, fmt.Errorf("%s must be a number", name)
	}
	return &val, nil
}

func parseInt(raw string, name string) (*int, error) {
	if raw == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(raw)
	if err != nil {
		return nil, fmt.Errorf("%s must be a whole number", name)
	}

	return &val, nil
}

func parseBool(raw string, name string) (*bool, error) {
	switch raw {
	case "":
		return nil, nil
	case "true":
		val := true
		return &val, nil
	case "false":
		val := false
		return &val, nil
	default:
		return nil, fmt.Errorf("%s must be either true or false", name)
	}
}

func parseString(raw string) *string {
	if raw == "" {
		return nil
	}
	return &raw
}

func splitAmenities(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	amenities := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		amenities = append(amenities, trimmed)
	}
	return amenities
}

// make filter structure
func (p *PropertyController) readSearchOptions() (models.FilterOptions, error) {
	var options models.FilterOptions
	var err error

	input := p.Ctx.Input

	if options.MinPrice, err = parseFloat(input.Query("min_price"), "min_price"); err != nil {
		return options, err
	}
	if options.MaxPrice, err = parseFloat(input.Query("max_price"), "max_price"); err != nil {
		return options, err
	}
	if options.MinStarRating, err = parseInt(input.Query("min_star_rating"), "min_star_rating"); err != nil {
		return options, err
	}
	if options.MinReviewScore, err = parseFloat(input.Query("min_review_score"), "min_review_score"); err != nil {
		return options, err
	}
	if options.MinReviews, err = parseInt(input.Query("min_reviews"), "min_reviews"); err != nil {
		return options, err
	}
	if options.Published, err = parseBool(input.Query("published"), "published"); err != nil {
		return options, err
	}
	if options.Feed, err = parseInt(input.Query("feed"), "feed"); err != nil {
		return options, err
	}
	if options.MinBedroom, err = parseInt(input.Query("min_bedroom"), "min_bedroom"); err != nil {
		return options, err
	}
	if options.Limit, err = parseInt(input.Query("limit"), "limit"); err != nil {
		return options, err
	}

	options.PropertyType = parseString(input.Query("property_type"))
	options.Amenities = splitAmenities(input.Query("amenities"))

	if err := services.ValidateSearchParams(options); err != nil {
		return options, err
	}
	return options, nil
}

func (p *PropertyController) GetProperties() {
	searchOptions, err := p.readSearchOptions()
	if err != nil {
		logs.Error("Failed to read search params: %v", err)
		p.writeError(http.StatusBadRequest, err.Error())
		return
	}

	store := services.GetStore()
	if store == nil {
		logs.Critical("Failed to read source data. The source property variable was never initialized")
		p.writeError(http.StatusInternalServerError, "Failed to read source data")
		return
	}

	properties := store.List(searchOptions)
	p.writeJSON(http.StatusOK, models.ListResponse{Result: properties})
}

func (p *PropertyController) GetPropertyByID() {
	propertyID := strings.TrimSpace(p.Ctx.Input.Param(":id"))
	if propertyID == "" {
		logs.Error("Property ID is not defined")
		p.writeError(http.StatusBadRequest, "Property ID is not defined")
		return
	}

	store := services.GetStore()
	if store == nil {
		logs.Critical("Failed to read store data. The source property variable was never initialized")
		p.writeError(http.StatusInternalServerError, "Failed to read source data")
		return
	}

	property, found := store.FindByID(propertyID)
	if !found {
		p.writeError(http.StatusNotFound, "Property not found")
		return
	}
	p.writeJSON(http.StatusOK, property)
}
