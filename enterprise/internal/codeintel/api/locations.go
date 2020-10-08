package api

import (
	clienttypes "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/client_types"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/store"
)

type ResolvedLocation struct {
	Dump  store.Dump
	Path  string
	Range clienttypes.Range
}

func sliceLocations(locations []clienttypes.Location, lo, hi int) []clienttypes.Location {
	if lo >= len(locations) {
		return nil
	}
	if hi >= len(locations) {
		hi = len(locations)
	}
	return locations[lo:hi]
}

func resolveLocationsWithDump(dump store.Dump, locations []clienttypes.Location) []ResolvedLocation {
	var resolvedLocations []ResolvedLocation
	for _, location := range locations {
		resolvedLocations = append(resolvedLocations, ResolvedLocation{
			Dump:  dump,
			Path:  dump.Root + location.Path,
			Range: location.Range,
		})
	}

	return resolvedLocations
}
