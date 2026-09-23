package services

import (
	"reflect"
	"rental-property-api/models"
	"testing"
)

func intPtr(val int) *int {
	return &val
}

func floatPtr(val float64) *float64 {
	return &val
}

func stringPtr(val string) *string {
	return &val
}

func boolPtr(val bool) *bool {
	return &val
}

func testRecords() []models.SourceProperty {
	kanto := "KT"
	california := "CA"

	return []models.SourceProperty{
		{
			ID:                   "BC-1000001",
			Feed:                 11,
			Country:              "Japan",
			CountryCode:          "JP",
			State:                "Kanto",
			StateAbbr:            &kanto,
			City:                 "Tokyo",
			Display:              "Shinjuku, Tokyo",
			LocationID:           "512",
			PropertyName:         "Shinjuku Garden Hotel",
			PropertySlug:         "shinjuku-garden-hotel",
			PropertyTypeCategory: "Hotel",
			USDPrice:             100.00,
			Occupancy:            4,
			BedroomCount:         2,
			BathroomCount:        1,
			NumberOfReview:       15,
			ReviewScoreGeneral:   7.5,
			StarRating:           3,
			AmenityCategories:    []string{"Internet", "Pool"},
			LonLat:               models.LonLat{Coordinates: []float64{139.69, 35.68}},
			Categories:           `[{"LocationID":"89","Name":"Japan","Type":"country","Slug":"japan","Display":["japan"]},{"LocationID":"512","Name":"Tokyo","Type":"city","Slug":"tokyo","Display":["japan","tokyo"]}]`,
			Published:            false,
			Images:               []string{"one.jpg", "two.jpg"},
		},
		{
			ID:                   "HA-2000002",
			Feed:                 12,
			Country:              "United States",
			CountryCode:          "US",
			State:                "California",
			StateAbbr:            &california,
			City:                 "San Diego",
			Display:              "La Jolla, San Diego",
			LocationID:           "781",
			PropertyName:         "La Jolla Beach House",
			PropertySlug:         "la-jolla-beach-house",
			PropertyTypeCategory: "House",
			USDPrice:             250.50,
			Occupancy:            8,
			BedroomCount:         4,
			BathroomCount:        3,
			NumberOfReview:       200,
			ReviewScoreGeneral:   9.1,
			StarRating:           5,
			AmenityCategories:    []string{"Parking", "Kitchen"},
			LonLat:               models.LonLat{Coordinates: []float64{-117.27, 32.84}},
			Categories:           `[{"LocationID":"1","Name":"United States","Type":"country","Slug":"united-states","Display":["united-states"]}]`,
			Published:            true,
			Images:               []string{"a.jpg", "b.jpg", "c.jpg"},
		},
		{
			ID:                   "HG-3000003",
			Feed:                 11,
			Country:              "Japan",
			CountryCode:          "JP",
			State:                "Kanto",
			StateAbbr:            nil,
			City:                 "Tokyo",
			Display:              "Shibuya, Tokyo",
			LocationID:           "513",
			PropertyName:         "Shibuya Compact Stay",
			PropertySlug:         "shibuya-compact-stay",
			PropertyTypeCategory: "Apartment",
			USDPrice:             75.25,
			Occupancy:            2,
			BedroomCount:         1,
			BathroomCount:        1,
			NumberOfReview:       40,
			ReviewScoreGeneral:   6.0,
			StarRating:           2,
			AmenityCategories:    []string{"Internet"},
			LonLat:               models.LonLat{Coordinates: []float64{139.70, 35.66}},
			Categories:           `[{"LocationID":"89","Name":"Japan","Type":"country","Slug":"japan","Display":["japan"]}]`,
			Published:            true,
			Images:               []string{"only.jpg"},
		},
		{
			ID:                   "EP-4000004",
			Feed:                 24,
			Country:              "Spain",
			CountryCode:          "ES",
			State:                "Andalusia",
			StateAbbr:            nil,
			City:                 "Marbella",
			Display:              "Marbella, Andalusia",
			LocationID:           "904",
			PropertyName:         "Marbella Hillside Villa",
			PropertySlug:         "marbella-hillside-villa",
			PropertyTypeCategory: "Villa",
			USDPrice:             340.00,
			Occupancy:            6,
			BedroomCount:         5,
			BathroomCount:        3,
			NumberOfReview:       60,
			ReviewScoreGeneral:   8.8,
			StarRating:           4,
			AmenityCategories:    []string{"Pool", "Gym"},
			LonLat:               models.LonLat{Coordinates: []float64{-4.88, 36.51}},
			Categories:           `[{"LocationID":"7","Name":"Spain","Type":"country","Slug":"spain","Display":["spain"]}]`,
			Published:            false,
			Images:               []string{"v1.jpg", "v2.jpg", "v3.jpg", "v4.jpg"},
		},
	}
}

func collectIDs(items []models.PropertyResponse) []string {
	ids := make([]string, 0, len(items))

	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func TestTransform(t *testing.T) {
	full := testRecords()[0]

	noCoordinates := testRecords()[1]
	noCoordinates.LonLat = models.LonLat{Coordinates: []float64{}}

	brokenCategories := testRecords()[2]
	brokenCategories.Categories = "no json at all !"

	emptyCategories := testRecords()[3]
	emptyCategories.Categories = ""

	tests := []struct {
		name           string
		record         models.SourceProperty
		wantName       string
		wantPrice      float64
		wantLat        float64
		wantLon        float64
		wantImageCount int
		wantCrumbCount int
		wantFirstCrumb string
	}{
		{
			name:           "maps every field of a complete record",
			record:         full,
			wantName:       "Shinjuku Garden Hotel",
			wantPrice:      100.00,
			wantLat:        35.68,
			wantLon:        139.69,
			wantImageCount: 2,
			wantCrumbCount: 2,
			wantFirstCrumb: "Japan",
		},
		{
			name:           "falls back to zero when coordinates are missing",
			record:         noCoordinates,
			wantName:       "La Jolla Beach House",
			wantPrice:      250.50,
			wantLat:        0,
			wantLon:        0,
			wantImageCount: 3,
			wantCrumbCount: 1,
			wantFirstCrumb: "United States",
		},
		{
			name:           "returns no breadcrumbs when categories cannot be decoded",
			record:         brokenCategories,
			wantName:       "Shibuya Compact Stay",
			wantPrice:      75.25,
			wantLat:        35.66,
			wantLon:        139.70,
			wantImageCount: 1,
			wantCrumbCount: 0,
		},
		{
			name:           "returns no breadcrumbs when categories is empty",
			record:         emptyCategories,
			wantName:       "Marbella Hillside Villa",
			wantPrice:      340.00,
			wantLat:        36.51,
			wantLon:        -4.88,
			wantImageCount: 4,
			wantCrumbCount: 0,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := transform(testCase.record)

			if got.ID != testCase.record.ID {
				t.Errorf("ID = %q, want %q", got.ID, testCase.record.ID)
			}
			if got.Feed != testCase.record.Feed {
				t.Errorf("Feed = %d, want %d", got.Feed, testCase.record.Feed)
			}
			if got.Published != testCase.record.Published {
				t.Errorf("Published = %v, want %v", got.Published, testCase.record.Published)
			}
			if got.GeoInfo.Name != testCase.record.Display {
				t.Errorf("GeoInfo.Name = %q, want %q", got.GeoInfo.Name, testCase.record.Display)
			}
			if got.GeoInfo.Lat != testCase.wantLat {
				t.Errorf("GeoInfo.Lat = %v, want %v", got.GeoInfo.Lat, testCase.wantLat)
			}
			if got.GeoInfo.Lon != testCase.wantLon {
				t.Errorf("GeoInfo.Lon = %v, want %v", got.GeoInfo.Lon, testCase.wantLon)
			}
			if len(got.GeoInfo.Breadcrumbs) != testCase.wantCrumbCount {
				t.Fatalf("len(Breadcrumbs) = %d, want %d", len(got.GeoInfo.Breadcrumbs), testCase.wantCrumbCount)
			}
			if got.GeoInfo.Breadcrumbs == nil {
				t.Error("Breadcrumbs is nil, it must transform to [] not nil")
			}
			if got.Property.Name != testCase.wantName {
				t.Errorf("Property.Name = %q, want %q", got.Property.Name, testCase.wantName)
			}
			if got.Property.Price != testCase.wantPrice {
				t.Errorf("Property.Price = %v, want %v", got.Property.Price, testCase.wantPrice)
			}
			if got.Property.Image.Count != testCase.wantImageCount {
				t.Errorf("Image.Count = %d, want %d", got.Property.Image.Count, testCase.wantImageCount)
			}
			if len(got.Property.Image.Images) != testCase.wantImageCount {
				t.Errorf("len(Image.Images) = %d, want %d", len(got.Property.Image.Images), testCase.wantImageCount)
			}
			if testCase.wantCrumbCount > 0 && got.GeoInfo.Breadcrumbs[0].Name != testCase.wantFirstCrumb {
				t.Errorf("Breadcrumbs[0].Name = %q, want %q", got.GeoInfo.Breadcrumbs[0].Name, testCase.wantFirstCrumb)
			}
		})
	}
}

func TestTransformKeepSliceEmptyThanNil(t *testing.T) {
	record := testRecords()[1]

	record.AmenityCategories = nil
	record.Images = nil

	got := transform(record)

	if got.Property.Amenities == nil {
		t.Error("Amenities is nil. It must be transformed to []")
	}
	if got.Property.Image.Images == nil {
		t.Error("got.Property.Image.Images is nil. It must be transformed to []")
	}
	if got.Property.Image.Count != 0 {
		t.Errorf("got.Property.Image.Count = %d. want 0", got.Property.Image.Count)
	}
}

func TestStoreList(t *testing.T) {
	store := newStore(testRecords())

	tests := []struct {
		name    string
		options models.FilterOptions
		wantIDs []string
	}{
		{
			name:    "no filters returns everything",
			options: models.FilterOptions{},
			wantIDs: []string{"BC-1000001", "HA-2000002", "HG-3000003", "EP-4000004"},
		},
		{
			name: "feed and published together",
			options: models.FilterOptions{
				Feed:      intPtr(11),
				Published: boolPtr(true),
			},
			wantIDs: []string{"HG-3000003"},
		},
		{
			name: "price range is inclusive at both ends",
			options: models.FilterOptions{
				MinPrice: floatPtr(75.25),
				MaxPrice: floatPtr(250.50),
			},
			wantIDs: []string{"BC-1000001", "HA-2000002", "HG-3000003"},
		},
		{
			name: "star rating and review count together",
			options: models.FilterOptions{
				MinStarRating: intPtr(4),
				MinReviews:    intPtr(100),
			},
			wantIDs: []string{"HA-2000002"},
		},
		{
			name: "amenities match with or logic",
			options: models.FilterOptions{
				Amenities: []string{"Internet", "Parking"},
			},
			wantIDs: []string{"BC-1000001", "HA-2000002", "HG-3000003"},
		},
		{
			name: "amenity matching is case sensitive",
			options: models.FilterOptions{
				Amenities: []string{"internet"},
			},
			wantIDs: []string{},
		},
		{
			name: "amenities or combined with an and filter",
			options: models.FilterOptions{
				Feed:      intPtr(11),
				Amenities: []string{"Internet", "Parking"},
			},
			wantIDs: []string{"BC-1000001", "HG-3000003"},
		},
		{
			name: "property type with a minimum bedroom count",
			options: models.FilterOptions{
				PropertyType: stringPtr("Villa"),
				MinBedroom:   intPtr(5),
			},
			wantIDs: []string{"EP-4000004"},
		},
		{
			name: "filters that match nothing",
			options: models.FilterOptions{
				PropertyType: stringPtr("Resort"),
			},
			wantIDs: []string{},
		},
		{
			name: "limit caps the number of results",
			options: models.FilterOptions{
				Limit: intPtr(2),
			},
			wantIDs: []string{"BC-1000001", "HA-2000002"},
		},
		{
			name: "limit larger than the match count changes nothing",
			options: models.FilterOptions{
				Feed:  intPtr(24),
				Limit: intPtr(10),
			},
			wantIDs: []string{"EP-4000004"},
		},
		{
			name: "minimum review score",
			options: models.FilterOptions{
				MinReviewScore: floatPtr(8.0),
			},
			wantIDs: []string{"HA-2000002", "EP-4000004"},
		},
		{
			name: "minimum price on its own",
			options: models.FilterOptions{
				MinPrice: floatPtr(100),
			},
			wantIDs: []string{"BC-1000001", "HA-2000002", "EP-4000004"},
		},
		{
			name: "minimum bedroom count on its own",
			options: models.FilterOptions{
				MinBedroom: intPtr(4),
			},
			wantIDs: []string{"HA-2000002", "EP-4000004"},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := store.List(testCase.options)

			if got.Items == nil {
				t.Fatal("Items is nil, it must be []")
			}
			if got.Count != len(got.Items) {
				t.Errorf("Count = %d but there are %d items", got.Count, len(got.Items))
			}

			IDs := collectIDs(got.Items)
			if !reflect.DeepEqual(IDs, testCase.wantIDs) {
				t.Errorf("ids = %v, want %v", IDs, testCase.wantIDs)
			}
		})
	}
}

func TestStoreFindByID(t *testing.T) {
	store := newStore(testRecords())

	tests := []struct {
		name      string
		id        string
		wantFound bool
	}{
		{name: "an id that exists", id: "HA-2000002", wantFound: true},
		{name: "an id that does not exist", id: "XX-9999999", wantFound: false},
		{name: "an empty id", id: "", wantFound: false},
		{name: "an id with the wrong case", id: "ha-2000002", wantFound: false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, found := store.FindByID(testCase.id)

			if found != testCase.wantFound {
				t.Fatalf("found = %v, want %v", found, testCase.wantFound)
			}
			if found && got.ID != testCase.id {
				t.Errorf("ID = %q, want %q", got.ID, testCase.id)
			}
			if !found && got.ID != "" {
				t.Errorf("ID = %q, want an empty response when nothing is found", got.ID)
			}
		})
	}
}

func TestValidateSearchParams(t *testing.T) {
	tests := []struct {
		name      string
		options   models.FilterOptions
		wantError bool
	}{
		{name: "no options at all", options: models.FilterOptions{}, wantError: false},
		{name: "a sensible price range", options: models.FilterOptions{MinPrice: floatPtr(50), MaxPrice: floatPtr(150)}, wantError: false},
		{name: "an inverted price range", options: models.FilterOptions{MinPrice: floatPtr(150), MaxPrice: floatPtr(50)}, wantError: true},
		{name: "a negative minimum price", options: models.FilterOptions{MinPrice: floatPtr(-1)}, wantError: true},
		{name: "a negative maximum price", options: models.FilterOptions{MaxPrice: floatPtr(-1)}, wantError: true},
		{name: "a negative star rating", options: models.FilterOptions{MinStarRating: intPtr(-1)}, wantError: true},
		{name: "a negative review score", options: models.FilterOptions{MinReviewScore: floatPtr(-0.5)}, wantError: true},
		{name: "a negative review count", options: models.FilterOptions{MinReviews: intPtr(-5)}, wantError: true},
		{name: "a negative bedroom count", options: models.FilterOptions{MinBedroom: intPtr(-2)}, wantError: true},
		{name: "a limit of zero", options: models.FilterOptions{Limit: intPtr(0)}, wantError: true},
		{name: "a negative limit", options: models.FilterOptions{Limit: intPtr(-3)}, wantError: true},
		{name: "a limit of one", options: models.FilterOptions{Limit: intPtr(1)}, wantError: false},
		{name: "an allowed feed", options: models.FilterOptions{Feed: intPtr(22)}, wantError: false},
		{name: "a feed outside the allowed set", options: models.FilterOptions{Feed: intPtr(13)}, wantError: true},
		{name: "an allowed property type", options: models.FilterOptions{PropertyType: stringPtr("Hostel")}, wantError: false},
		{name: "a property type with the wrong case", options: models.FilterOptions{PropertyType: stringPtr("hotel")}, wantError: true},
		{name: "a property type that does not exist", options: models.FilterOptions{PropertyType: stringPtr("Castle")}, wantError: true},
		{name: "a star rating no property can reach", options: models.FilterOptions{MinStarRating: intPtr(9)}, wantError: false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := ValidateSearchParams(testCase.options)

			if testCase.wantError && err == nil {
				t.Error("Expected error, but found nil")
			} else if !testCase.wantError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}
