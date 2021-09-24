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
	WorkoutsA string = "getworkouts"
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

// WorkoutType is workout type
type WorkoutType string

// Workout Type
const (
	WTCalories          WorkoutType = "calories"            // Active calories burned (in Kcal).
	WTEffduration       WorkoutType = "effduration"         // Effective duration.
	WTIntensity         WorkoutType = "intensity"           // Intensity.
	WTManualDistance    WorkoutType = "manual_distance"     // Distance travelled manually entered by user (in meters).
	WTManualCalories    WorkoutType = "manual_calories"     // Active calories burned manually entered by user (in Kcal).
	WTHrAverage         WorkoutType = "hr_average"          // Average heart rate.
	WTHrMin             WorkoutType = "hr_min"              // Minimal heart rate.
	WTHrMax             WorkoutType = "hr_max"              // Maximal heart rate.
	WTHrZone0           WorkoutType = "hr_zone_0"           // Duration in seconds when heart rate was in a light zone (cf. Glossary).
	WTHrZone1           WorkoutType = "hr_zone_1"           // Duration in seconds when heart rate was in a moderate zone (cf. Glossary).
	WTHrZone2           WorkoutType = "hr_zone_2"           // Duration in seconds when heart rate was in an intense zone (cf. Glossary).
	WTHrZone3           WorkoutType = "hr_zone_3"           // Duration in seconds when heart rate was in maximal zone (cf. Glossary).
	WTPauseDuration     WorkoutType = "pause_duration"      // Total pause time in second filled by user
	WTAlgoPauseDuration WorkoutType = "algo_pause_duration" // Total pause time in seconds detected by Withings device (swim only)
	WTSpo2Average       WorkoutType = "spo2_average"        // Average percent of SpO2 percent value during a workout
	WTSteps             WorkoutType = "steps"               // Number of steps.
	WTDistance          WorkoutType = "distance"            // Distance travelled (in meters).
	WTElevation         WorkoutType = "elevation"           // Number of floors climbed.
	WTPoolLaps          WorkoutType = "pool_laps"           // Number of pool laps.
	WTStrokes           WorkoutType = "strokes"             // Number of strokes.
	WTPoolLength        WorkoutType = "pool_length"         // Length of the pool.
)

// WorkoutCategory is category of workout
type WorkoutCategory int

// Workout Category
const (
	WCWalk          WorkoutCategory = 1
	WCRun           WorkoutCategory = 2
	WCHiking        WorkoutCategory = 3
	WCSkating       WorkoutCategory = 4
	WCBMX           WorkoutCategory = 5
	WCBicycling     WorkoutCategory = 6
	WCSwimming      WorkoutCategory = 7
	WCSurfing       WorkoutCategory = 8
	WCKitesurfing   WorkoutCategory = 9
	WCWindsurfing   WorkoutCategory = 10
	WCBodyboard     WorkoutCategory = 11
	WCTennis        WorkoutCategory = 12
	WCTableTennis   WorkoutCategory = 13
	WCSquash        WorkoutCategory = 14
	WCBadminton     WorkoutCategory = 15
	WCLiftWeights   WorkoutCategory = 16
	WCCalisthenics  WorkoutCategory = 17
	WCElliptical    WorkoutCategory = 18
	WCPilates       WorkoutCategory = 19
	WCBasketBall    WorkoutCategory = 20
	WCSoccer        WorkoutCategory = 21
	WCFootball      WorkoutCategory = 22
	WCRugby         WorkoutCategory = 23
	WCVolleyBall    WorkoutCategory = 24
	WCWaterpolo     WorkoutCategory = 25
	WCHorseRiding   WorkoutCategory = 26
	WCGolf          WorkoutCategory = 27
	WCYoga          WorkoutCategory = 28
	WCDancing       WorkoutCategory = 29
	WCBoxing        WorkoutCategory = 30
	WCFencing       WorkoutCategory = 31
	WCWrestling     WorkoutCategory = 32
	WCMartialArts   WorkoutCategory = 33
	WCSkiing        WorkoutCategory = 34
	WCSnowboarding  WorkoutCategory = 35
	WCOther         WorkoutCategory = 36
	WCNoActivity    WorkoutCategory = 128
	WCRowing        WorkoutCategory = 187
	WCZumba         WorkoutCategory = 188
	WCBaseball      WorkoutCategory = 191
	WCHandball      WorkoutCategory = 192
	WCHockey        WorkoutCategory = 193
	WCIceHockey     WorkoutCategory = 194
	WCClimbing      WorkoutCategory = 195
	WCIceSkating    WorkoutCategory = 196
	WCMultiSport    WorkoutCategory = 272
	WCIndoorRunning WorkoutCategory = 307
	WCIndoorCycling WorkoutCategory = 308
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
