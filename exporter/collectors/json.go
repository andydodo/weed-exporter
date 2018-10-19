package collectors

type Info struct {
	Version string  `json:"Version"`
	Volumes Volumes `json:"Volumes"`
}

type Volumes struct {
	DataCenters DataCenters `json:"DataCenters"`
	Free        int64       `json:"Free"`
	Max         int64       `json:"Max"`
}

type DataCenters struct {
	DefaultDataCenter map[string]map[string][]DefaultDataCenter `json:"DefaultDataCenter"`
}

type DefaultDataCenter struct {
	ID               int64            `json:"Id"`
	Size             int64            `json:"Size"`
	ReplicaPlacement ReplicaPlacement `json:"ReplicaPlacement"`
	TTL              TTL              `json:"Ttl"`
	Collection       string           `json:"Collection"`
	Version          int64            `json:"Version"`
	FileCount        int64            `json:"FileCount"`
	DeleteCount      int64            `json:"DeleteCount"`
	DeletedByteCount int64            `json:"DeletedByteCount"`
	ReadOnly         bool             `json:"ReadOnly"`
}

type ReplicaPlacement struct {
	SameRackCount       int64 `json:"SameRackCount"`
	DiffRackCount       int64 `json:"DiffRackCount"`
	DiffDataCenterCount int64 `json:"DiffDataCenterCount"`
}

type TTL struct {
}
