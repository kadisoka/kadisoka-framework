package app

import "sync"

type BuildInfo struct {
	RevisionID string
	Timestamp  string
}

var (
	buildRevisionID string = "unknown"
	buildTimestamp  string = "unknown"
)

var buildInfoSetOnce sync.Once

func SetBuildInfo(
	revisionID string,
	timestamp string,
) {
	buildInfoSetOnce.Do(func() {
		buildRevisionID = revisionID
		buildTimestamp = timestamp
	})
}

func GetBuildInfo() BuildInfo {
	return BuildInfo{
		RevisionID: buildRevisionID,
		Timestamp:  buildTimestamp,
	}
}
