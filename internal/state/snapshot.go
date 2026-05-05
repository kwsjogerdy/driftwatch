package state

// Snapshot represents the recorded state of a single environment.
type Snapshot struct {
	// Environment is the name of the environment (e.g. "staging", "production").
	Environment string `json:"environment"`

	// Resources is a flat map of resource key to its current value or hash.
	// Keys are typically dot-separated paths such as "aws.s3.bucket.name".
	Resources map[string]string `json:"resources"`
}
