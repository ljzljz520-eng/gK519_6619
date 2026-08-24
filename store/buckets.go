package store

var bucketNames = struct {
	Bridges      []byte
	Scenarios    []byte
	Trajectories []byte
	Views        []byte
	Events       []byte
}{
	Bridges:      []byte("BridgeRecord"),
	Scenarios:    []byte("WindScenarioRecord"),
	Trajectories: []byte("TrajectoryRecord"),
	Views:        []byte("ViewPreferenceRecord"),
	Events:       []byte("AuditEvent"),
}

func allBuckets() [][]byte {
	return [][]byte{bucketNames.Bridges, bucketNames.Scenarios, bucketNames.Trajectories, bucketNames.Views, bucketNames.Events}
}
