package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"rental-property-api/models"

	"github.com/beego/beego/v2/core/logs"
)

type Store struct {
	properties []models.SourceProperty
	index      map[string]int
}

// private methods & variables
var defaultStore *Store
var allowedFeeds = map[int]bool{
	11: true,
	12: true,
	22: true,
	24: true,
}
var allowedPropertyTypes = map[string]bool{
	"Hotel":     true,
	"House":     true,
	"Apartment": true,
	"Villa":     true,
	"Resort":    true,
	"Hostel":    true,
}

// response transformation
func emptyIfNil(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func latOf(point models.LonLat) float64 {
	if len(point.Coordinates) < 2 {
		return 0
	}
	return point.Coordinates[1]
}

func lonOf(point models.LonLat) float64 {
	if len(point.Coordinates) < 2 {
		return 0
	}
	return point.Coordinates[0]
}

func parseBreadcrumbs(category string, propertyID string) []models.BreadCrumb {
	breadCrumbs := make([]models.BreadCrumb, 0)
	if category == "" {
		return breadCrumbs
	}

	var decoded []models.SourceCategory
	if err := json.Unmarshal([]byte(category), &decoded); err != nil {
		logs.Error("failed to decode category of property %s: %v", propertyID, err)
		return breadCrumbs
	}

	for _, categoryItem := range decoded {
		breadCrumbs = append(breadCrumbs, models.BreadCrumb{
			LocationID: categoryItem.LocationID,
			Name:       categoryItem.Name,
			Type:       categoryItem.Type,
			Slug:       categoryItem.Slug,
			Display:    emptyIfNil(categoryItem.Display),
		})
	}
	return breadCrumbs
}

func transform(src models.SourceProperty) models.PropertyResponse {
	return models.PropertyResponse{
		ID:   src.ID,
		Feed: src.Feed,
		GeoInfo: models.GeoInfo{
			Breadcrumbs: parseBreadcrumbs(src.Categories, src.ID),
			City:        src.City,
			Country:     src.Country,
			CountryCode: src.CountryCode,
			Name:        src.Display,
			LocationID:  src.LocationID,
			Lat:         latOf(src.LonLat),
			Lon:         lonOf(src.LonLat),
			State:       src.State,
			StateAbbr:   src.StateAbbr,
		},
		Property: models.Property{
			Amenities:    emptyIfNil(src.AmenityCategories),
			Name:         src.PropertyName,
			Slug:         src.PropertySlug,
			PropertyType: src.PropertyTypeCategory,
			Price:        src.USDPrice,
			ReviewScore:  src.ReviewScoreGeneral,
			StarRating:   src.StarRating,
			Counts: models.PropertyCount{
				Bathroom:  src.BathroomCount,
				Bedroom:   src.BedroomCount,
				Reviews:   src.NumberOfReview,
				Occupancy: src.Occupancy,
			},
			Image: models.PropertyImage{
				Count:  len(src.Images),
				Images: emptyIfNil(src.Images),
			},
		},
		Published: src.Published,
	}
}

// filter logics
func matchesAmenity(propertyAmenities []string, filterAmenities []string) bool {
	hasFilterAmenity := make(map[string]bool)

	for _, propertyAmenity := range propertyAmenities {
		hasFilterAmenity[propertyAmenity] = true
	}

	for _, filterAmenity := range filterAmenities {
		if hasFilterAmenity[filterAmenity] {
			return true
		}
	}
	return false
}

func filter(property models.SourceProperty, options models.FilterOptions) bool {
	if options.MinPrice != nil && *options.MinPrice > property.USDPrice {
		return false
	}
	if options.MaxPrice != nil && *options.MaxPrice < property.USDPrice {
		return false
	}
	if options.MinStarRating != nil && *options.MinStarRating > property.StarRating {
		return false
	}
	if options.MinReviewScore != nil && *options.MinReviewScore > property.ReviewScoreGeneral {
		return false
	}
	if options.MinReviews != nil && *options.MinReviews > property.NumberOfReview {
		return false
	}
	if options.Published != nil && property.Published != *options.Published {
		return false
	}
	if options.PropertyType != nil && *options.PropertyType != property.PropertyTypeCategory {
		return false
	}
	if options.Feed != nil && *options.Feed != property.Feed {
		return false
	}
	if options.MinBedroom != nil && *options.MinBedroom > property.BedroomCount {
		return false
	}
	if len(options.Amenities) > 0 && !matchesAmenity(property.AmenityCategories, options.Amenities) {
		return false
	}
	return true
}

// source json loader
func newStore(records []models.SourceProperty) *Store {
	index := make(map[string]int, len(records))
	for i, record := range records {
		index[record.ID] = i
	}

	return &Store{
		properties: records,
		index:      index,
	}
}

func loadFromFile(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var decoded []models.SourceProperty
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, fmt.Errorf("failed to decode file %s: %w", path, err)
	}

	return newStore(decoded), nil
}

// public methods
func GetStore() *Store {
	return defaultStore
}

func (s *Store) Count() int {
	return len(s.properties)
}

func (s *Store) FindByID(propertyID string) (models.PropertyResponse, bool) {
	position, found := s.index[propertyID]
	if !found {
		return models.PropertyResponse{}, false
	}

	return transform(s.properties[position]), true
}

func (s *Store) List(options models.FilterOptions) models.ListResult {
	filteredProperties := make([]models.PropertyResponse, 0)

	for _, property := range s.properties {
		if !filter(property, options) {
			continue
		}
		filteredProperties = append(filteredProperties, transform(property))
	}

	if options.Limit != nil && *options.Limit >= 0 && len(filteredProperties) > *options.Limit {
		filteredProperties = filteredProperties[:*options.Limit]
	}

	return models.ListResult{
		Count: len(filteredProperties),
		Items: filteredProperties,
	}
}

func Init(path string) error {
	store, err := loadFromFile(path)
	if err != nil {
		return err
	}

	defaultStore = store
	return nil
}

func ValidateSearchParams(options models.FilterOptions) error {
	if options.MinPrice != nil && *options.MinPrice < 0 {
		return errors.New("min_price cannot be negative")
	}
	if options.MaxPrice != nil && *options.MaxPrice < 0 {
		return errors.New("max_price cannot be negative")
	}
	if options.MinPrice != nil && options.MaxPrice != nil && *options.MaxPrice < *options.MinPrice {
		return errors.New("max_price cannot be less than min_price")
	}
	if options.MinStarRating != nil && *options.MinStarRating < 0 {
		return errors.New("min_star_rating cannot be negative")
	}
	if options.MinReviewScore != nil && *options.MinReviewScore < 0 {
		return errors.New("min_review_score cannot be negative")
	}
	if options.MinReviews != nil && *options.MinReviews < 0 {
		return errors.New("min_reviews cannot be negative")
	}
	if options.MinBedroom != nil && *options.MinBedroom < 0 {
		return errors.New("min_bedroom cannot be negative")
	}
	if options.Limit != nil && *options.Limit < 1 {
		return errors.New("limit must be a positive number")
	}
	if options.PropertyType != nil && !allowedPropertyTypes[*options.PropertyType] {
		return errors.New("property_type must be either Hotel, House, Apartment, Villa, Resort or Hostel")
	}
	if options.Feed != nil && !allowedFeeds[*options.Feed] {
		return errors.New("feed must be either 11, 12, 22 or 24")
	}
	return nil
}
