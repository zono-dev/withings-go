package withings

import "time"

// MeasureData is used for parsed Measurement.
type MeasureData struct {
	GrpID    int64
	Date     time.Time
	Value    float64
	Attrib   int
	Category int
	DeviceID string
}

// SerialzedMeas has parsed Measurements.
type SerialzedMeas struct {
	Weights        []MeasureData
	Heights        []MeasureData
	FatFreeMass    []MeasureData
	FatRatios      []MeasureData
	FatMassWeights []MeasureData
	DiastolicBPs   []MeasureData
	SystolicBPs    []MeasureData
	HeartPulses    []MeasureData
	Temps          []MeasureData
	SPO2s          []MeasureData
	BodyTemps      []MeasureData
	SkinTemps      []MeasureData
	MuscleMasses   []MeasureData
	Hydration      []MeasureData
	BoneMasses     []MeasureData
	PWaveVel       []MeasureData
	VO2s           []MeasureData
	UnknowVals     []MeasureData
}

// Measurement is raw data from Measure API.
// See https://developer.withings.com/oauth2/#operation/measure-getmeas .
type Measurement struct {
	Status int `json:"status"`
	Body   struct {
		Updatetime  int    `json:"updatetime"`
		Timezone    string `json:"timezone"`
		Measuregrps []struct {
			GrpID        int64  `json:"grpid"`
			Attrib       int    `json:"attrib"`
			Date         int    `json:"date"`
			Created      int    `json:"created"`
			Category     int    `json:"category"`
			DeviceID     string `json:"deviceid"`
			HashDeviceID string `json:"hash_deviceid"`
			Measures     []struct {
				Value int `json:"value"`
				Type  int `json:"type"`
				Unit  int `json:"unit"`
				Algo  int `json:"algo"`
				Fm    int `json:"fm"`
			} `json:"measures"`
			Comment string `json:"comment"`
		} `json:"measuregrps"`
		More   int `json:"more"`
		Offset int `json:"offset"`
	} `json:"body"`
	SerializedData *SerialzedMeas
}

// Activities is raw data from Measure Activity API.
// See https://developer.withings.com/oauth2/#operation/measurev2-getactivity .
type Activities struct {
	Status int `json:"status"`
	Body   struct {
		Activities []struct {
			Date          string  `json:"date"`
			Timezone      string  `json:"timezone"`
			Deviceid      string  `json:"deviceid"`
			Brand         int     `json:"brand"`
			IsTracker     bool    `json:"is_tracker"`
			Steps         int     `json:"steps"`
			Distance      int     `json:"distance"`
			Elevation     int     `json:"elevation"`
			Soft          int     `json:"soft"`
			Moderate      int     `json:"moderate"`
			Intense       int     `json:"intense"`
			Active        int     `json:"active"`
			Calories      float64 `json:"calories"`
			Totalcalories int     `json:"totalcalories"`
			HrAverage     int     `json:"hr_average"`
			HrMin         int     `json:"hr_min"`
			HrMax         int     `json:"hr_max"`
			HrZone0       int     `json:"hr_zone_0"`
			HrZone1       int     `json:"hr_zone_1"`
			HrZone2       int     `json:"hr_zone_2"`
			HrZone3       int     `json:"hr_zone_3"`
		} `json:"activities"`
		More   bool `json:"more"`
		Offset int  `json:"offset"`
	} `json:"body"`
}

// Workouts is raw data from Measure Getworkouts API.
// See https://developer.withings.com/api-reference#operation/measurev2-getworkouts .
type Workouts struct {
	Status int `json:"status"`
	Body   struct {
		Series []struct {
			ID        int64           `json:"id"`
			Category  WorkoutCategory `json:"category"`
			Timezone  string          `json:"timezone"`
			Model     int             `json:"model"`
			Attrib    int             `json:"attrib"`
			Startdate int64           `json:"startdate"`
			Enddate   int64           `json:"enddate"`
			Date      string          `json:"date"`
			Modified  int64           `json:"modified"`
			DeviceID  string          `json:"deviceid"`
			Data      struct {
				AlgoPauseDuration int     `json:"algo_pause_duration"`
				Calories          float64 `json:"calories"`
				Distance          float64 `json:"distance"`
				Effduration       int     `json:"effduration"`
				Elevation         int     `json:"elevation"`
				HrAverage         int     `json:"hr_average"`
				HrMax             int     `json:"hr_max"`
				HrMin             int     `json:"hr_min"`
				HrZone0           int     `json:"hr_zone_0"`
				HrZone1           int     `json:"hr_zone_1"`
				HrZone2           int     `json:"hr_zone_2"`
				HrZone3           int     `json:"hr_zone_3"`
				Intensity         int     `json:"intensity"`
				ManualCalories    int     `json:"manual_calories"`
				ManualDistance    int     `json:"manual_distance"`
				PauseDuration     int     `json:"pause_duration"`
				PoolLaps          int     `json:"pool_laps"`
				PoolLength        int     `json:"pool_length"`
				Spo2Average       int     `json:"spo2_average"`
				Steps             int     `json:"steps"`
				Strokes           int     `json:"strokes"`
			} `json:"data"`
		} `json:"series"`
		More   bool `json:"more"`
		Offset int  `json:"offset"`
	} `json:"body"`
}

// Sleeps is raw data from Sleep API.
// See https://developer.withings.com/oauth2/#tag/sleep .
type Sleeps struct {
	Status int `json:"status"`
	Body   struct {
		Series []struct {
			Startdate int64 `json:"startdate"`
			Enddate   int64 `json:"enddate"`
			State     int   `json:"state"`
			Hr        struct {
				Timestamp int `json:"timestamp"`
			} `json:"hr"`
			Rr struct {
				Timestamp int `json:"timestamp"`
			} `json:"rr"`
			Snoring struct {
				Timestamp int `json:"timestamp"`
			} `json:"snoring"`
		} `json:"series"`
		Model   int `json:"model"`
		ModelID int `json:"model_id"`
	} `json:"body"`
}

// SleepSummaries is raw data from Sleep Summaries API.
// See https://developer.withings.com/oauth2/#operation/sleepv2-getsummary .
type SleepSummaries struct {
	Status int `json:"status"`
	Body   struct {
		Series []struct {
			Timezone  string `json:"timezone"`
			Model     int    `json:"model"`
			ModelID   int    `json:"model_id"`
			Startdate int64  `json:"startdate"`
			Enddate   int64  `json:"enddate"`
			Date      string `json:"date"`
			Created   int64  `json:"created"`
			Modified  int64  `json:"modified"`
			Data      struct {
				BreathingDisturbancesIntensity int `json:"breathing_disturbances_intensity"`
				Deepsleepduration              int `json:"deepsleepduration"`
				Durationtosleep                int `json:"durationtosleep"`
				Durationtowakeup               int `json:"durationtowakeup"`
				HrAverage                      int `json:"hr_average"`
				HrMax                          int `json:"hr_max"`
				HrMin                          int `json:"hr_min"`
				Lightsleepduration             int `json:"lightsleepduration"`
				Remsleepduration               int `json:"remsleepduration"`
				RrAverage                      int `json:"rr_average"`
				RrMax                          int `json:"rr_max"`
				RrMin                          int `json:"rr_min"`
				SleepScore                     int `json:"sleep_score"`
				Snoring                        int `json:"snoring"`
				Snoringepisodecount            int `json:"snoringepisodecount"`
				Wakeupcount                    int `json:"wakeupcount"`
				Wakeupduration                 int `json:"wakeupduration"`
			} `json:"data"`
		} `json:"series"`
		More   bool `json:"more"`
		Offset int  `json:"offset"`
	} `json:"body"`
}
