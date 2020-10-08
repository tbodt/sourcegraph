package api

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	bundles "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/client"
	bundlemocks "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/client/mocks"
	clienttypes "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/client_types"
	commitmocks "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/commits/mocks"
	gitservermocks "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/gitserver/mocks"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/store"
	storemocks "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/store/mocks"
)

func TestRanges(t *testing.T) {
	mockStore := storemocks.NewMockStore()
	mockBundleManagerClient := bundlemocks.NewMockBundleManagerClient()
	mockBundleClient := bundlemocks.NewMockBundleClient()
	mockGitserverClient := gitservermocks.NewMockClient()
	mockCommitUpdater := commitmocks.NewMockUpdater()

	sourceRanges := []clienttypes.CodeIntelligenceRange{
		{
			Range:       clienttypes.Range{Start: clienttypes.Position{1, 2}, End: clienttypes.Position{3, 4}},
			Definitions: []clienttypes.Location{},
			References:  []clienttypes.Location{},
			HoverText:   "",
		},
		{
			Range:       clienttypes.Range{Start: clienttypes.Position{2, 3}, End: clienttypes.Position{4, 5}},
			Definitions: []clienttypes.Location{{Path: "foo.go", Range: clienttypes.Range{Start: clienttypes.Position{10, 20}, End: clienttypes.Position{30, 40}}}},
			References:  []clienttypes.Location{{Path: "bar.go", Range: clienttypes.Range{Start: clienttypes.Position{100, 200}, End: clienttypes.Position{300, 400}}}},
			HoverText:   "ht2",
		},
		{
			Range:       clienttypes.Range{Start: clienttypes.Position{3, 4}, End: clienttypes.Position{5, 6}},
			Definitions: []clienttypes.Location{{Path: "bar.go", Range: clienttypes.Range{Start: clienttypes.Position{11, 21}, End: clienttypes.Position{31, 41}}}},
			References:  []clienttypes.Location{{Path: "foo.go", Range: clienttypes.Range{Start: clienttypes.Position{101, 201}, End: clienttypes.Position{301, 401}}}},
			HoverText:   "ht3",
		},
	}

	setMockStoreGetDumpByID(t, mockStore, map[int]store.Dump{42: testDump1})
	setMockBundleManagerClientBundleClient(t, mockBundleManagerClient, map[int]bundles.BundleClient{42: mockBundleClient})
	setMockBundleClientRanges(t, mockBundleClient, "main.go", 10, 20, sourceRanges)

	api := testAPI(mockStore, mockBundleManagerClient, mockGitserverClient, mockCommitUpdater)
	ranges, err := api.Ranges(context.Background(), "sub1/main.go", 10, 20, 42)
	if err != nil {
		t.Fatalf("expected error getting ranges: %s", err)
	}

	expectedRanges := []ResolvedCodeIntelligenceRange{
		{
			Range:       clienttypes.Range{Start: clienttypes.Position{1, 2}, End: clienttypes.Position{3, 4}},
			Definitions: nil,
			References:  nil,
			HoverText:   "",
		},
		{
			Range:       clienttypes.Range{Start: clienttypes.Position{2, 3}, End: clienttypes.Position{4, 5}},
			Definitions: []ResolvedLocation{{Dump: testDump1, Path: "sub1/foo.go", Range: clienttypes.Range{Start: clienttypes.Position{10, 20}, End: clienttypes.Position{30, 40}}}},
			References:  []ResolvedLocation{{Dump: testDump1, Path: "sub1/bar.go", Range: clienttypes.Range{Start: clienttypes.Position{100, 200}, End: clienttypes.Position{300, 400}}}},
			HoverText:   "ht2",
		},
		{
			Range:       clienttypes.Range{Start: clienttypes.Position{3, 4}, End: clienttypes.Position{5, 6}},
			Definitions: []ResolvedLocation{{Dump: testDump1, Path: "sub1/bar.go", Range: clienttypes.Range{Start: clienttypes.Position{11, 21}, End: clienttypes.Position{31, 41}}}},
			References:  []ResolvedLocation{{Dump: testDump1, Path: "sub1/foo.go", Range: clienttypes.Range{Start: clienttypes.Position{101, 201}, End: clienttypes.Position{301, 401}}}},
			HoverText:   "ht3",
		},
	}
	if diff := cmp.Diff(expectedRanges, ranges); diff != "" {
		t.Errorf("unexpected range (-want +got):\n%s", diff)
	}
}

func TestRangesUnknownDump(t *testing.T) {
	mockStore := storemocks.NewMockStore()
	mockBundleManagerClient := bundlemocks.NewMockBundleManagerClient()
	mockGitserverClient := gitservermocks.NewMockClient()
	mockCommitUpdater := commitmocks.NewMockUpdater()
	setMockStoreGetDumpByID(t, mockStore, nil)

	api := testAPI(mockStore, mockBundleManagerClient, mockGitserverClient, mockCommitUpdater)
	if _, err := api.Ranges(context.Background(), "sub1", 42, 0, 10); err != ErrMissingDump {
		t.Fatalf("unexpected error getting ranges. want=%q have=%q", ErrMissingDump, err)
	}
}
