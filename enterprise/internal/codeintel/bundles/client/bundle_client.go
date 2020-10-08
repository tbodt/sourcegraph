package client

import (
	"context"
	"encoding/json"
	"fmt"

	clienttypes "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/client_types"
)

// BundleClient is the interface to the precise-code-intel-bundle-manager service scoped to a particular dump.
type BundleClient interface {
	// ID gets the identifier of the target bundle.
	ID() int

	// Exists determines if the given path exists in the dump.
	Exists(ctx context.Context, path string) (bool, error)

	// Ranges returns definition, reference, and hover data for each range within the given span of lines.
	Ranges(ctx context.Context, path string, startLine, endLine int) ([]clienttypes.CodeIntelligenceRange, error)

	// Definitions retrieves a list of definition locations for the symbol under the given location.
	Definitions(ctx context.Context, path string, line, character int) ([]clienttypes.Location, error)

	// Definitions retrieves a list of reference locations for the symbol under the given location.
	References(ctx context.Context, path string, line, character int) ([]clienttypes.Location, error)

	// Hover retrieves the hover text for the symbol under the given location.
	Hover(ctx context.Context, path string, line, character int) (string, clienttypes.Range, bool, error)

	// Diagnostics retrieves the diagnostics and total count of diagnostics for the documents that have the given path prefix.
	Diagnostics(ctx context.Context, prefix string, skip, take int) ([]clienttypes.Diagnostic, int, error)

	// MonikersByPosition retrieves a list of monikers attached to the symbol under the given location. There may
	// be multiple ranges enclosing this point. The returned monikers are partitioned such that inner ranges occur
	// first in the result, and outer ranges occur later.
	MonikersByPosition(ctx context.Context, path string, line, character int) ([][]clienttypes.MonikerData, error)

	// MonikerResults retrieves a page of locations attached to a moniker and a total count of such locations.
	MonikerResults(ctx context.Context, modelType, scheme, identifier string, skip, take int) ([]clienttypes.Location, int, error)

	// PackageInformation retrieves package information data by its identifier.
	PackageInformation(ctx context.Context, path, packageInformationID string) (clienttypes.PackageInformationData, error)
}

type bundleClientImpl struct {
	base     baseClient
	bundleID int
}

var _ BundleClient = &bundleClientImpl{}

// ID gets the identifier of the target bundle.
func (c *bundleClientImpl) ID() int {
	return c.bundleID
}

// Exists determines if the given path exists in the dump.
func (c *bundleClientImpl) Exists(ctx context.Context, path string) (exists bool, err error) {
	err = c.request(ctx, "exists", map[string]interface{}{"path": path}, &exists, func() (interface{}, error) { return nil, fmt.Errorf("unimplemented") })
	return exists, err
}

// Ranges returns definition, reference, and hover data for each range within the given span of lines.
func (c *bundleClientImpl) Ranges(ctx context.Context, path string, startLine, endLine int) (codeintelRanges []clienttypes.CodeIntelligenceRange, err error) {
	args := map[string]interface{}{
		"path":      path,
		"startLine": startLine,
		"endLine":   endLine,
	}

	err = c.request(ctx, "ranges", args, &codeintelRanges, func() (interface{}, error) { return nil, fmt.Errorf("unimplemented") })
	return codeintelRanges, err
}

// Definitions retrieves a list of definition locations for the symbol under the given location.
func (c *bundleClientImpl) Definitions(ctx context.Context, path string, line, character int) (locations []clienttypes.Location, err error) {
	args := map[string]interface{}{
		"path":      path,
		"line":      line,
		"character": character,
	}

	err = c.request(ctx, "definitions", args, &locations, func() (interface{}, error) { return nil, fmt.Errorf("unimplemented") })
	c.addBundleIDToLocations(locations)
	return locations, err
}

// Definitions retrieves a list of reference locations for the symbol under the given location.
func (c *bundleClientImpl) References(ctx context.Context, path string, line, character int) (locations []clienttypes.Location, err error) {
	args := map[string]interface{}{
		"path":      path,
		"line":      line,
		"character": character,
	}

	err = c.request(ctx, "references", args, &locations, func() (interface{}, error) { return nil, fmt.Errorf("unimplemented") })
	c.addBundleIDToLocations(locations)
	return locations, err
}

// Hover retrieves the hover text for the symbol under the given location.
func (c *bundleClientImpl) Hover(ctx context.Context, path string, line, character int) (string, clienttypes.Range, bool, error) {
	args := map[string]interface{}{
		"path":      path,
		"line":      line,
		"character": character,
	}

	var target *json.RawMessage
	if err := c.request(ctx, "hover", args, &target, func() (interface{}, error) { return nil, fmt.Errorf("unimplemented") }); err != nil {
		return "", clienttypes.Range{}, false, err
	}

	if target == nil {
		return "", clienttypes.Range{}, false, nil
	}

	payload := struct {
		Text  string            `json:"text"`
		Range clienttypes.Range `json:"range"`
	}{}

	if err := json.Unmarshal(*target, &payload); err != nil {
		return "", clienttypes.Range{}, false, err
	}

	return payload.Text, payload.Range, true, nil
}

// Diagnostics retrieves the diagnostics and total count of diagnostics for the documents that have the given path prefix.
func (c *bundleClientImpl) Diagnostics(ctx context.Context, prefix string, skip, take int) (diagnostics []clienttypes.Diagnostic, count int, err error) {
	args := map[string]interface{}{
		"prefix": prefix,
	}
	if skip != 0 {
		args["skip"] = skip
	}
	if take != 0 {
		args["take"] = take
	}

	target := struct {
		Diagnostics []clienttypes.Diagnostic `json:"diagnostics"`
		Count       int                      `json:"count"`
	}{}

	err = c.request(ctx, "diagnostics", args, &target, func() (interface{}, error) { return nil, fmt.Errorf("unimplemented") })
	diagnostics = target.Diagnostics
	count = target.Count
	c.addBundleIDToDiagnostics(diagnostics)
	return diagnostics, count, err
}

// MonikersByPosition retrieves a list of monikers attached to the symbol under the given location. There may
// be multiple ranges enclosing this point. The returned monikers are partitioned such that inner ranges occur
// first in the result, and outer ranges occur later.
func (c *bundleClientImpl) MonikersByPosition(ctx context.Context, path string, line, character int) (target [][]clienttypes.MonikerData, err error) {
	args := map[string]interface{}{
		"path":      path,
		"line":      line,
		"character": character,
	}

	err = c.request(ctx, "monikersByPosition", args, &target, func() (interface{}, error) { return nil, fmt.Errorf("unimplemented") })
	return target, err
}

// MonikerResults retrieves a page of locations attached to a moniker and a total count of such locations.
func (c *bundleClientImpl) MonikerResults(ctx context.Context, modelType, scheme, identifier string, skip, take int) (locations []clienttypes.Location, count int, err error) {
	args := map[string]interface{}{
		"modelType":  modelType,
		"scheme":     scheme,
		"identifier": identifier,
	}
	if skip != 0 {
		args["skip"] = skip
	}
	if take != 0 {
		args["take"] = take
	}

	target := struct {
		Locations []clienttypes.Location `json:"locations"`
		Count     int                    `json:"count"`
	}{}

	err = c.request(ctx, "monikerResults", args, &target, func() (interface{}, error) { return nil, fmt.Errorf("unimplemented") })
	locations = target.Locations
	count = target.Count
	c.addBundleIDToLocations(locations)
	return locations, count, err
}

// PackageInformation retrieves package information data by its identifier.
func (c *bundleClientImpl) PackageInformation(ctx context.Context, path, packageInformationID string) (target clienttypes.PackageInformationData, err error) {
	args := map[string]interface{}{
		"path":                 path,
		"packageInformationId": packageInformationID,
	}

	err = c.request(ctx, "packageInformation", args, &target, func() (interface{}, error) { return nil, fmt.Errorf("unimplemented") })
	return target, err
}

func (c *bundleClientImpl) request(ctx context.Context, path string, qs map[string]interface{}, target interface{}, handler func() (interface{}, error)) error {
	// TODO - short-circuit here

	return c.base.QueryBundle(ctx, c.bundleID, path, qs, &target)
}

func (c *bundleClientImpl) addBundleIDToLocations(locations []clienttypes.Location) {
	for i := range locations {
		locations[i].DumpID = c.bundleID
	}
}

func (c *bundleClientImpl) addBundleIDToDiagnostics(diagnostics []clienttypes.Diagnostic) {
	for i := range diagnostics {
		diagnostics[i].DumpID = c.bundleID
	}
}
