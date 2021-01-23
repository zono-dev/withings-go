package withings

// API endpoint
const (
	defaultMeasureURL   = "https://wbsapi.withings.net/measure"
	defaultMeasureURLv2 = "https://wbsapi.withings.net/v2/measure"
	defaultSleepURLv2   = "https://wbsapi.withings.net/v2/sleep"
)

// scopes
const (
	ScopeActivity string = "user.activity"
	ScopeMetrics  string = "user.metrics"
	ScopeInfo     string = "user.info"
)

// form keys
const (
	PPaction       string = "action"
	PPmeastype     string = "meastype"
	PPcategory     string = "category"
	PPstartdate    string = "startdate"
	PPenddate      string = "enddate"
	PPstartdateymd string = "startdateymd"
	PPenddateymd   string = "enddateymd"
	PPlastupdate   string = "lastupdate"
	PPoffset       string = "offset"
	PPdataFields   string = "data_fields"
)

// Service action
const (
	MeasureA  string = "getmeas"
	ActivityA string = "getactivity"
	SleepA    string = "get"
	SleepSA   string = "getsummary"
)

// MeasType is Measurement Type
type MeasType int

// Measurement Type
const (
	Weight        MeasType = 1   // Weight (kg)
	Height        MeasType = 4   // Height (meter)
	FatFreeMass   MeasType = 5   // Fat Free Mass (kg)
	FatRatio      MeasType = 6   // Fat Ratio (%)
	FatMassWeight MeasType = 8   // Fat Mass Weight (kg)
	DiastolicBP   MeasType = 9   // Diastolic Blood Pressure (mmHg)
	SystolicBP    MeasType = 10  // Systolic Blood Pressure (mmHg)
	HeartPulse    MeasType = 11  // Heart Pulse (bpm)
	Temp          MeasType = 12  // Temperature (celsius)
	SPO2          MeasType = 54  // SPO2 (%)
	BodyTemp      MeasType = 71  // Body Temperature (celsius)
	SkinTemp      MeasType = 73  // Skin temperature (celsius)
	MuscleMass    MeasType = 76  // Muscle Mass (kg)
	Hydration     MeasType = 77  // Hydration (kg)
	BoneMass      MeasType = 88  // Bone Mass (kg)
	PWaveVel      MeasType = 91  // Pulse Wave Velocity (m/s)
	VO2           MeasType = 123 // VO2 max is a numerical measurement of your body’s ability to consume oxygen (ml/min/kg).
)

// CatType is category type
type CatType int

// Category type
const (
	Real      CatType = 1 // Real is for real measures.
	Objective CatType = 2 // Objective is for user objectives.
)

// ActivityType is activity type
type ActivityType string

// Activity Type
const (
	Steps         ActivityType = "steps"         // Number of steps
	Distance      ActivityType = "distance"      // Distance travelled (in meters).
	Elevation     ActivityType = "elevation"     // Number of floors cliembed.
	Soft          ActivityType = "soft"          // Duration of soft activities (in seconds).
	Moderate      ActivityType = "moderate"      // Duration of moderate activities (in seconds).
	Intense       ActivityType = "intense"       // Duration of intense activities (in seconds).
	Active        ActivityType = "active"        // Sum of intense and moderate activity durations (in seconds).
	Calories      ActivityType = "calories"      // Active calories burned (in Kcal).
	TotalCalories ActivityType = "totalcalories" // Total calories burned (in Kcal).
	HrAverage     ActivityType = "hr_average"    // Average heart rate.
	HrMin         ActivityType = "hr_min"        // Minimal heart rate.
	HrMax         ActivityType = "hr_max"        // Maximal heart rate.
	HrZone0       ActivityType = "hr_zone_0"     // Duration in seconds when heart rate was in a light zone.
	HrZone1       ActivityType = "hr_zone_1"     // Duration in seconds when heart rate was in a moderate zone.
	HrZone2       ActivityType = "hr_zone_2"     // Duration in seconds when heart rate was in an intense zone.
	HrZone3       ActivityType = "hr_zone_3"     // Duration in seconds when heart rate was in maximal zone.
)

// SleepType is Sleep Type.
type SleepType string

// Sleep Type
const (
	HrSleep      SleepType = "hr"      // Heart Rate.
	RrSleep      SleepType = "rr"      // Respiration Rate.
	SnoringSleep SleepType = "snoring" // Total snoring time.
)

// SleepSummariesType is Sleep Summaries Type.
type SleepSummariesType string

// Sleep Summaries Types.
const (
	SSBdi   SleepSummariesType = "breathing_disturbances_intensity" // Intensity of breathing disturbances
	SSDsd   SleepSummariesType = "deepsleepduration"                // Duration in state deep sleep (in seconds).
	SSD2s   SleepSummariesType = "durationtosleep"                  // Time to sleep (in seconds).
	SSD2w   SleepSummariesType = "durationtowakeup"                 // Time to wake up (in seconds).
	SSHrAvr SleepSummariesType = "hr_average"                       // Average heart rate.
	SSHrMax SleepSummariesType = "hr_max"                           // Maximal heart rate.
	SSHrMin SleepSummariesType = "hr_min"                           // Minimal heart rate.
	SSLsd   SleepSummariesType = "lightsleepduration"               // Duration in state light sleep (in seconds).
	SSRsd   SleepSummariesType = "remsleepduration"                 // Duration in state REM sleep (in seconds).
	SSRRAvr SleepSummariesType = "rr_average"                       // Average respiration rate.
	SSRRMax SleepSummariesType = "rr_max"                           // Maximal respiration rate.
	SSRRMin SleepSummariesType = "rr_min"                           // Minimal respiration rate.
	SSSS    SleepSummariesType = "sleep_score"                      // Sleep score
	SSSng   SleepSummariesType = "snoring"                          // Total snoring time
	SSSngEC SleepSummariesType = "snoringepisodecount"              // Numbers of snoring episodes of at least one minute
	SSWupC  SleepSummariesType = "wakeupcount"                      // Number of times the user woke up.
	SSWupD  SleepSummariesType = "wakeupduration"                   // Time spent awake (in seconds).
)

// SleepState is Sleep state
type SleepState int

// Sleep state
const (
	Awake      SleepState = 0
	LightSleep SleepState = 1
	DeepSleep  SleepState = 2
	REM        SleepState = 3
)
