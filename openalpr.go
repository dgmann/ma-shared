package shared

type OpenAlprResponse struct {
	EpochTime int64 `json:"epoch_time"`
	ImageWidth int `json:"img_width"`
	ImageHeight int `json:"img_height"`
	TotalProcessingTimeMs float64 `json:"processing_time_ms"`
	Results []OpenAlprLicensePlateResult `json:"results"`
	RegionsOfInterest []string `json:"regions_of_interest"`
}

type OpenAlprLicensePlateResult struct {
	RequestedTopN int `json:"requested_topn"`
	Plate string `json:"plate"`
	ProcessingTimeMs float64 `json:"processing_time_ms"`
	PlateIndex int `json:"plate_index"`
	RegionConfidence float64 `json:"region_confidence"`
	Region string `json:"region"`
	Coordinates []Coordinates `json:"coordinates"`
	Plates []Plate `json:"candidates"`
}

type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Plate struct {
	Plate string `json:"plate"`
	Confidence float64 `json:"confidence"`
	MatchesTemplate bool
}

