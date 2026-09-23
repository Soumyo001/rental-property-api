# Rental Property API

A read only REST API over a JSON file of rental properties, built with Go and Beego.
The data file is loaded into memory once at startup and every request is served
from there.

## Requirements

- Go 1.22 or higher
- The bee CLI

```
go install github.com/beego/bee/v2@latest
```
- In `.bashrc` or `.zshrc` add:
```
export PATH=$PATH:$(go env GOPATH)/bin
```

## Setup

```
git clone https://github.com/Soumyo001/rental-property-api.git
cd rental-property-api
go mod download
```

The data file is already at `data/rental_properties.json`. The path is read from
`conf/app.conf` under the `datapath` key.

> This project needs no credentials of any kind. If any were ever required they would be read from environment variables, not from `conf/app.conf`.

## Running

```
bee run
```

The server listens on the port set by `httpport` in `conf/app.conf`, which is
8080 by default.

To regenerate the route table and the Swagger spec:

```
bee generate routers
bee generate docs
bee run -downdoc=true
```

`bee generate routers` must be re-run after any change to a `@router` comment.

Swagger UI: http://localhost:8080/swagger/

## Tests

```
go test ./... -v
go vet ./...
```

Coverage for the service layer:

```
go test ./services/... -cover
```

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /v1/properties | List properties with optional filters and a limit |
| GET | /v1/properties/:id | Get one property by id |

### Query parameters

| Parameter | Type | Rule |
|-----------|------|------|
| min_price | float | usd_price >= min_price |
| max_price | float | usd_price <= max_price |
| min_star_rating | int | star_rating >= min_star_rating |
| min_review_score | float | review_score_general >= min_review_score |
| min_reviews | int | number_of_review >= min_reviews |
| min_bedroom | int | bedroom_count >= min_bedroom |
| published | bool | true or false |
| property_type | string | Hotel, House, Apartment, Villa, Resort or Hostel |
| feed | int | 11, 12, 22 or 24 |
| amenities | string | Comma separated, case sensitive, matched with OR |
| limit | int | Caps how many results are returned |

Scalar filters combine with AND. The amenities list is matched with OR among
itself, and that result is then one more AND condition.

## Sample requests

List everything:

```
curl "http://localhost:8080/v1/properties"
```

Cap the number of results:

```
curl "http://localhost:8080/v1/properties?limit=5"
```

Two AND filters:

```
curl "http://localhost:8080/v1/properties?feed=11&published=false"
```

A price range with a property type:

```
curl "http://localhost:8080/v1/properties?min_price=50&max_price=150&property_type=Hotel"
```

Amenities with OR logic:

```
curl "http://localhost:8080/v1/properties?amenities=Internet,Parking"
```

Amenities OR combined with an AND filter:

```
curl "http://localhost:8080/v1/properties?feed=11&amenities=Internet,Parking"
```

Highly rated properties with several bedrooms:

```
curl "http://localhost:8080/v1/properties?min_star_rating=4&min_review_score=8&min_bedroom=3"
```

A single property:

```
curl "http://localhost:8080/v1/properties/BC-1000001"
```

An id that does not exist, which returns 404:

```
curl -i "http://localhost:8080/v1/properties/XX-9999999"
```

An invalid parameter, which returns 400:

```
curl -i "http://localhost:8080/v1/properties?feed=13"
```