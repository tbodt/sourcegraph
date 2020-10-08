package graphqlbackend

import (
	"context"
	"testing"

	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/types"
	"github.com/sourcegraph/sourcegraph/internal/db"
)

func TestRepositories(t *testing.T) {
	resetMocks()
	repos := []*types.Repo{
		{ID: 0, Name: "repo1"},
		{ID: 1, Name: "repo2"},
		{
			ID:   2,
			Name: "repo3",
			RepoFields: &types.RepoFields{
				Cloned: true,
			},
		},
	}
	db.Mocks.Repos.List = func(ctx context.Context, opt db.ReposListOptions) ([]*types.Repo, error) {
		if opt.NoCloned {
			return repos[0:2], nil
		}
		if opt.OnlyCloned {
			return repos[2:], nil
		}
		if opt.LimitOffset != nil {
			if opt.After > 0 {
				return repos[opt.After : opt.LimitOffset.Limit+1], nil
			}
			return repos[:opt.LimitOffset.Limit], nil
		}
		return repos, nil
	}

	db.Mocks.Repos.Count = func(context.Context, db.ReposListOptions) (int, error) { return 3, nil }
	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Schema: mustParseGraphQLSchema(t),
			Query: `
				{
					repositories {
						nodes { name }
						totalCount
						pageInfo { hasNextPage }
					}
				}
			`,
			ExpectedResult: `
				{
					"repositories": {
						"nodes": [
							{ "name": "repo1" },
							{ "name": "repo2" },
							{ "name": "repo3" }
						],
						"totalCount": null,
						"pageInfo": {"hasNextPage": false}
					}
				}
			`,
		},
		{
			Schema: mustParseGraphQLSchema(t),
			// cloned and notCloned are true by default
			// this test ensures the behavior is the same
			// when setting them explicitly
			Query: `
				{
					repositories(cloned: true, notCloned: true) {
						nodes { name }
						totalCount
						pageInfo { hasNextPage }
					}
				}
			`,
			ExpectedResult: `
				{
					"repositories": {
						"nodes": [
							{ "name": "repo1" },
							{ "name": "repo2" },
							{ "name": "repo3" }
						],
						"totalCount": null,
						"pageInfo": {"hasNextPage": false}
					}
				}
			`,
		},
		{
			Schema: mustParseGraphQLSchema(t),
			Query: `
				{
					repositories(first: 2) {
						nodes { name }
						pageInfo { hasNextPage }
					}
				}
			`,
			ExpectedResult: `
				{
					"repositories": {
						"nodes": [
							{ "name": "repo1" },
							{ "name": "repo2" }
						],
						"pageInfo": {"hasNextPage": true}
					}
				}
			`,
		},
		{
			Schema: mustParseGraphQLSchema(t),
			Query: `
				{
					repositories(cloned: false) {
						nodes { name }
						pageInfo { hasNextPage }
					}
				}
			`,
			ExpectedResult: `
				{
					"repositories": {
						"nodes": [
							{ "name": "repo1" },
							{ "name": "repo2" }
						],
						"pageInfo": {"hasNextPage": false}
					}
				}
			`,
		},
		{
			Schema: mustParseGraphQLSchema(t),
			Query: `
				{
					repositories(notCloned: false) {
						nodes { name }
						pageInfo { hasNextPage }
					}
				}
			`,
			ExpectedResult: `
				{
					"repositories": {
						"nodes": [
							{ "name": "repo3" }
						],
						"pageInfo": {"hasNextPage": false}
					}
				}
			`,
		},
		{
			Schema: mustParseGraphQLSchema(t),
			Query: `
				{
					repositories(notCloned: false, cloned: false) {
						nodes { name }
						pageInfo { hasNextPage }
					}
				}
			`,
			ExpectedResult: `
				{
					"repositories": {
						"nodes": [
							{ "name": "repo1" },
							{ "name": "repo2" }
						],
						"pageInfo": {"hasNextPage": false}
					}
				}
			`,
		},
		{
			Schema: mustParseGraphQLSchema(t),
			Query: `
				{
					repositories(first: 1) {
						nodes {
							name
						}
						pageInfo {
							endCursor
						}
					}
				}
			`,
			ExpectedResult: `
				{
					"repositories": {
						"nodes": [
							{
								"name": "repo1"
							}
						],
						"pageInfo": {
							"endCursor": "UmVwb3NpdG9yeUN1cnNvcjp7IkFmdGVyIjoxfQ=="
						}
					}
				}
			`,
		},
		{
			Schema: mustParseGraphQLSchema(t),
			Query: `
				{
					repositories(first: 1, after: "UmVwb3NpdG9yeUN1cnNvcjp7IkFmdGVyIjoxfQ==") {
						nodes {
							name
						}
						pageInfo {
							endCursor
						}
					}
				}
			`,
			ExpectedResult: `
				{
					"repositories": {
						"nodes": [
							{
								"name": "repo2"
							}
						],
						"pageInfo": {
							"endCursor": "UmVwb3NpdG9yeUN1cnNvcjp7IkFmdGVyIjoyfQ=="
						}
					}
				}
			`,
		},
		{
			Schema: mustParseGraphQLSchema(t),
			Query: `
				{
					repositories(first: 1, after: "UmVwb3NpdG9yeUN1cnNvcjp7IkFmdGVyIjoyfQ==") {
						nodes {
							name
						}
						pageInfo {
							endCursor
						}
					}
				}
			`,
			ExpectedResult: `
				{
					"repositories": {
						"nodes": [
							{
								"name": "repo3"
							}
						],
						"pageInfo": {
							"endCursor": null
						}
					}
				}
			`,
		},
	})
}
