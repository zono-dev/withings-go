# withings-go

withings-go is UNOFFICIAL go client to acess Withings API easily.
More Withings API document can be found in [Withings developer documentation](https://developer.withings.com/oauth2/#).

Also withings-go document can be found in [withings · pkg.go.dev](https://pkg.go.dev/github.com/zono-dev/withings-go/withings)


## Supported Resources

- Offline Authorization
- [Measure - Getmeas](https://developer.withings.com/api-reference/#operation/measure-getmeas)
- [Measure v2 - Getactivity](https://developer.withings.com/api-reference/#operation/measurev2-getactivity)
- [Measure v2 - Getworkouts](https://developer.withings.com/api-reference/#operation/measurev2-getworkouts)
- [Sleep v2 - Get](https://developer.withings.com/api-reference/#operation/sleepv2-get)
- [Sleep v2 - Getsummary](https://developer.withings.com/api-reference/#operation/sleepv2-getsummary)

## Requirements

This library requires Go 1.15 or later.

## Installation

`go get github.com/zono-dev/withings-go/withings`

## Getting started

### Create your Withings account

If you have any Withings gear, you should already have your Withings account.

### Register your app

[Register your app](https://account.withings.com/partner/add_oauth2) and get your client ID and consumer secret.

### Setup your settings file

Create `.test_settings.yaml` in your working directory.

Tip: In the `cmd` directory, there is examples of how to use withings-go.

`.test_settings.yaml` sample is here.

```
CID: "YOUR CLIENT ID HERE"
Secret: "YOUR CONSUMER SECRET HERE"
RedirectURL: "https://example.com/"
```

### Run

Copy `cmd/auth/main.go` to your working directory and `go build -o auth`.
When you run `auth`, some messages will be displayed in the console as follows.

```
map[CID:YOUR-CLIENT-ID RedirectURL:https://example.com/ Secret:YOUR-SECRET]
[user.activity,user.metrics,user.info]
URL to authorize:http://account.withings.com/oauth2_user/authorize2?access_type=offline&client_id=yourclientid&redirect_uri=https%3A%2F%2Fexample.com&response_type=code&scope=user.activity%2Cuser.metrics%2Cuser.info&state=state
Open url your browser and Enter your grant code here.
 Grant Code:
```

### Get token

Go to the `URL to authorize` URL in your browser and allow your app connect to your withings account.

Then your browser will be redirect to your `RedirectURL` automatically and you will find your `Grant code` in redirected URL.

Type this into the console followed by the `Grant Code:`.

Lastly, your `access_token.json` will be generated in `cmd/auth`.

## Get Measurements 

Create `.test_settings.yaml` in your working directory.

`.test_settings.yaml` sample is here.

```
CID: "YOUR CLIENT ID HERE"
Secret: "YOUR CONSUMER SECRET HERE"
RedirectURL: "https://example.com/"
```

Copy `cmd/getMeasurements/main.go` to your working directory and `go build -o getMeasurements`.

When you run `getMeasurements`, your measurements will be displayed in the console.

## Usage

### Configuration

```Go
client, err = withings.New("YourConsumerID", "YourConsumerSecret", "RedirectURL")
if err != nil {
    fmt.Println(err)
    return
}
```

### Authorize

```Go
client, err = withings.New("YourConsumerID", "YourConsumerSecret", "RedirectURL")
if err != nil {
    fmt.Println(err)
    return
}

// First time authorization
client.Token, e = withings.AuthorizeOffline(client.Conf)
client.Client = withings.GetClient(client.Conf, client.Token)

// Readtoken file and Refresh
token, err = client.ReadToken("path_to_token_file")

// Refresh token
token, rf, err := client.RefreshToken()

// Save Token if you need
err = client.SaveToken("path_to_token_file")
if err != nil {
	fmt.Println("Failed to SaveToken")
	fmt.Println(err)
	return
}

```

### Get Measurements

```Go
// GetMeas call withings API Measure - GetMeas. (https://developer.withings.com/oauth2/#operation/measure-getmeas)
//
// cattype: category, 1 for real measures, 2 for user objectives
// startdate, enddate: Measure's start date, end date.
// lastupdate : Timestamp for requesting data that were updated or created after this date. Use this instead of startdate+endate.If lastupdate is set to a timestamp other than Offsetbase, getMeas will use lastupdate in preference to startdate/enddate.
// offset: When a first call retuns more:1 and offset:XX, set value XX in this parameter to retrieve next available rows.
// isOldToNew: If true, results must be sorted by oldest to newest. If false, results must be sorted by newest to oldest.
// isSerialized: if true, results must be parsed to Measurement.SerializedData
// mtype: Measurement Type. Set the measurement type you want to get data. See MeasType in enum.go.
mym, err := client.GetMeas(withings.Real, adayago, t, lastupdate, 0, false, true, withings.Weight, withings.Height, withings.FatFreeMass, withings.BoneMass, withings.FatRatio, withings.FatMassWeight, withings.Temp, withings.HeartPulse, withings.Hydration)

if err != nil {
	fmt.Println(err)
	return
}

// Status codes information can be found in the following URL.
// https://developer.withings.com/oauth2/#section/Response-status
fmt.Printf("Status: %d\n", mym.Status)

// If "isSerialized" was true, results must be parsed to Measurement.SerializedData
for _, v := range mym.SerializedData.Weights {
	fmt.Printf("Weight(Grpid:%v, Category:%v, Attrib: %v, DeviceID:%v)\n", v.GrpID, v.Category, v.Attrib, v.DeviceID)
	fmt.Printf("%v, %.1f Kg\n", v.Date.In(jst).Format(layout2), v.Value)
}

// Raw data should be provided from mym.Body.Measuregrps
for _, v := range mym.Body.Measuregrps {
	weight := float64(v.Measures[0].Value) * math.Pow10(v.Measures[0].Unit)
	fmt.Printf("Weight:%.1f Kgs\n", weight)
}
```

### Get Activity

```Go

// GetActivity call withings API Measure v2 - Getactivity. (https://developer.withings.com/oauth2/#operation/measurev2-getactivity)
// startdate/enddate: Activity result start date, end date.
// lastupdate : Timestamp for requesting data that were updated or created after this date. Use this instead of startdate+endate. If lastupdate is set to a timestamp other than Offsetbase, getMeas will use lastupdate in preference to startdate/enddate.
// offset: When a first call retuns more:1 and offset:XX, set value XX in this parameter to retrieve next available rows.
// atype: Acitivity Type. Set the activity type you want to get data. See ActivityType in enum.go.
act, err := client.GetActivity(sd, ed, 0, 0, withings.Steps, withings.Calories, withings.HrAverage, withings.HrMin, withings.HrMax)

if err != nil {
	fmt.Println("getActivity Error.")
	fmt.Println(err)
	return
}

for _, v := range act.Body.Activities {
	fmt.Printf("Date:%s, Steps:%d, BurnedCalories: %g, HRAverage: %d, HRMinimum: %d, HRMax:%d \n", v.Date, v.Steps, v.Calories, v.HrAverage, v.HrMin, v.HrMax)
}

```

### Get Workouts

```Go
// GetWorkouts call withings API Measure v2 - Getworkouts. (https://developer.withings.com/api-reference#operation/measurev2-getworkouts)
// startdate/enddate: Workouts result start date, end date.
// lastupdate : Timestamp for requesting data that were updated or created after this date. Use this instead of startdate+endate. If lastupdate is set to a timestamp other than Offsetbase, GetWorkouts will use lastupdate in preference to startdate/enddate.
// offset: When a first call retuns more:1 and offset:XX, set value XX in this parameter to retrieve next available rows.
// wtype: Workout Type. Set the workout type you want to get data. See WorkoutType in enum.go.
workouts, err := client.GetWorkouts(sd, ed, 0, 0, withings.WTCalories, withings.WTEffduration, withings.WTSteps, withings.WTDistance)

if err != nil {
	fmt.Println("getWorkouts Error.")
	fmt.Println(err)
	return
}

for _, v := range workouts.Body.Series {
	fmt.Printf("Date:%s, Category: %d, Duration: %d, Steps:%d, Distance:%.1f, Calories: %.1f\n", v.Date, v.Category, v.Data.Effduration, v.Data.Steps, v.Data.Distance, v.Data.Calories)
}
```

### Get Sleep

```Go
// GetSleep cal withings API Sleep v2 - Get. (https://developer.withings.com/oauth2/#operation/sleepv2-get)
// startdate/enddate: Measures' start date, end date.
// stype: Sleep Type. Set the sleep type you want to get data. See SleepType in enum.go.
slp, err := client.GetSleep(adayago, t, withings.HrSleep, withings.RrSleep, withings.SnoringSleep)
if err != nil {
	fmt.Println("getSleep Error!")
	fmt.Println(err)
	return
}
for _, v := range slp.Body.Series {
	st := ""
	switch v.State {
	case int(withings.Awake):
		st = "Awake"
	case int(withings.LightSleep):
		st = "LightSleep"
	case int(withings.DeepSleep):
		st = "DeepSleep"
	case int(withings.REM):
		st = "REM"
	default:
		st = "Unknown"
	}
	stimeUnix := time.Unix(v.Startdate, 0)
	etimeUnix := time.Unix(v.Enddate, 0)

	stime := (stimeUnix.In(jst)).Format(layout2)
	etime := (etimeUnix.In(jst)).Format(layout2)
	message := fmt.Sprintf("%s to %s: %s, Hr:%d, Rr:%d, Snoring:%d\n", stime, etime, st, v.Hr, v.Rr, v.Snoring)
	fmt.Printf(message)
}
```

### Get Sleep Summary

```Go
// GetSleepSummary call withings API Sleep v2 - Getsummary. (https://developer.withings.com/oauth2/#operation/sleepv2-getsummary)
// startdate/enddate: Measurement result start date, end date.
// lastupdate : Timestamp for requesting data that were updated or created after this date. Use this instead of startdate+endate. If lastupdate is set to a timestamp other than Offsetbase, getMeas will use lastupdate in preference to startdate/enddate.
// stype: Sleep Summaries Type. Set the sleep summaries data you want to get. See SleepSummariesType in enum.go.
slpsum, err := client.GetSleepSummary(sd, ed, 0, withings.SSBdi, withings.SSDsd, withings.SSD2s, withings.SSD2w, withings.SSHrAvr, withings.SSHrMax, withings.SSHrMin, withings.SSLsd, withings.SSRsd, withings.SSRRAvr, withings.SSRRMax, withings.SSRRMin, withings.SSSS,withings.SSSng, withings.SSSngEC, withings.SSWupC, withings.SSWupD)


if err != nil {
	fmt.Println("getSleepSummary Error!")
	fmt.Println(err)
	return
}

for _, v := range slpsum.Body.Series {
	stimeUnix := time.Unix(v.Startdate, 0)
	etimeUnix := time.Unix(v.Enddate, 0)

	stime := (stimeUnix.In(jst)).Format(layout2)
	etime := (etimeUnix.In(jst)).Format(layout2)
	message := fmt.Sprintf(
		"%s-%s: BDI:%d, duration to deep sleep(sec):%d, duration to sleep(sec):%d, duration to wakeup(sec):%d, HrAverage:%d, Max:%d, Min:%d, WakeupCounts:%d",
		stime, etime, v.Data.BreathingDisturbancesIntensity, v.Data.Deepsleepduration, v.Data.Durationtosleep, v.Data.Durationtowakeup, v.Data.HrAverage, v.Data.HrMax, v.Data.HrMin, v.Data.Wakeupcount)
	fmt.Println(message)
}
```
