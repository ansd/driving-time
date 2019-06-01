package maps

import (
	"context"

	"googlemaps.github.io/maps"
)

type Client interface {
	DistanceMatrix(context.Context, *maps.DistanceMatrixRequest) (*maps.DistanceMatrixResponse, error)
}
