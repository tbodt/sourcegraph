package globalstatedb

import "github.com/sourcegraph/sourcegraph/internal/db/dbtesting"

func init() {
	dbtesting.DBNameSuffix = "globalstatedb"
}
