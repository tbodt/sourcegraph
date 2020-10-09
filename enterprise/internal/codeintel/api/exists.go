package api

import (
	"context"
	"strings"

	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	bundles "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/bundles/client"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/gitserver"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/store"
)

// FindClosestDumps returns the set of dumps that can most accurately answer code intelligence
// queries for the given path. If exactPath is true, then only dumps that definitely contain the
// exact document path are returned. Otherwise, dumps containing any document for which the given
// path is a prefix are returned. These dump IDs should be subsequently passed to invocations of
// Definitions, References, and Hover.
func (api *codeIntelAPI) FindClosestDumps(ctx context.Context, repositoryID int, commit, path string, exactPath bool, indexer string) ([]store.Dump, error) {
	candidates, err := api.updateCommitGraph(ctx, repositoryID, commit, path, exactPath, indexer)
	if err != nil {
		return nil, err
	}

	var dumps []store.Dump
	for _, dump := range candidates {
		// TODO(efritz) - ensure there's a valid document path
		// for the other condition. This should probably look like
		// an additional parameter on the following exists query.
		if exactPath {
			exists, err := api.bundleManagerClient.BundleClient(dump.ID).Exists(ctx, strings.TrimPrefix(path, dump.Root))
			if err != nil {
				if err == bundles.ErrNotFound {
					log15.Warn("Bundle does not exist")
					return nil, nil
				}
				return nil, errors.Wrap(err, "bundleManagerClient.BundleClient")
			}
			if !exists {
				continue
			}
		}

		dumps = append(dumps, dump)
	}

	return dumps, nil
}

// TODO - rename
// TODO - redocument
// TODO - update tests
// updateCommitGraph will perform an update of the given repository's commit graph if it appears to be out
// of date. If we know already know about this commit, we do not perform an update. Otherwise, it is likely
// that a user is browsing a commit that was pushed after commit that owns the last index we processed. If
// the repository has no index data at all, we skip the update and return false as there would be no useful
// information available in the commit graph.
func (api *codeIntelAPI) updateCommitGraph(ctx context.Context, repositoryID int, commit, path string, exactPath bool, indexer string) ([]store.Dump, error) {
	commitExists, err := api.store.HasCommit(ctx, repositoryID, commit)
	if err != nil {
		return nil, errors.Wrap(err, "store.HasCommit")
	}
	if commitExists {
		// The parameters exactPath and rootMustEnclosePath align here: if we're looking for dumps
		// that can answer queries for a directory (e.g. diagnostics), we want any dump that happens
		// to intersect the target directory. If we're looking for dumps that can answer queries for
		// a single file, then we need a dump with a root that properly encloses that file.
		dumps, err := api.store.FindClosestDumps(ctx, repositoryID, commit, path, exactPath, indexer)
		if err != nil {
			return nil, errors.Wrap(err, "store.FindClosestDumps")
		}

		return dumps, nil
	}

	repositoryExists, err := api.store.HasRepository(ctx, repositoryID)
	if err != nil {
		return nil, errors.Wrap(err, "store.HasRepository")
	}
	if !repositoryExists {
		// TODO - differentiate this from error with disjoint history from
		// gitserver + what's in postgres? Allow the extension to retry or
		// display a message in that case so it's not just failing randomly
		// on first page loads sometimes.
		return nil, nil
	}

	// TODO - need to do smaller subgraph
	graph, err := api.gitserverClient.CommitGraph(ctx, api.store, repositoryID,gitserver.CommitGraphOptions{
		// TODO
	})
	if err != nil {
		return nil, err
	}

	dumps, err := api.store.FindClosestDumpsFromGraphFragment(ctx, repositoryID, commit, path, exactPath, indexer, graph)
	if err != nil {
		return nil, err
	}

	if err := api.store.MarkRepositoryAsDirty(ctx, repositoryID); err != nil {
		return nil, errors.Wrap(err, "store.MarkRepositoryAsDirty")
	}

	return dumps, nil
}
