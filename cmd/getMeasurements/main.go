package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/zono-dev/withings-go/withings"
)

const (
	tokenFile = "access_token.json"
	layout    = "2006-01-02"
	layout2   = "2006-01-02 15:04:05"
	isnotify  = false
)

var (
	jst        *time.Location
	t          time.Time
	adayago    time.Time
	lastupdate time.Time
	ed         string
	sd         string
	client     *(withings.Client)
	settings   map[string]string
)

func auth(settings map[string]string) {
	var err error
	client, err = withings.New(settings["CID"], settings["Secret"], settings["RedirectURL"])

	if err != nil {
		fmt.Println("Failed to create New client")
		fmt.Println(err)
		return
	}

	if _, err := os.Open(tokenFile); err != nil {
		var e error

		client.Token, e = withings.AuthorizeOffline(client.Conf)
		client.Client = withings.GetClient(client.Conf, client.Token)

		if e != nil {
			fmt.Println("Failed to authorize offline.")
		}
		fmt.Println("~~ authorized. Let's check the token file!")
	} else {
		_, err = client.ReadToken(tokenFile)

		if err != nil {
			fmt.Println("Failed to read token file.")
			fmt.Println(err)
			return
		}
	}
}

func tokenFuncs() {
	// Show token
	client.PrintToken()

	// Refresh Token if you need
	_, rf, err := client.RefreshToken()
	if err != nil {
		fmt.Println("Failed to RefreshToken")
		fmt.Println(err)
		return
	}
	if rf {
		fmt.Println("You got new token!")
		client.PrintToken()
	}

	// Save Token if you need
	err = client.SaveToken(tokenFile)
	if err != nil {
		fmt.Println("Failed to RefreshToken")
		fmt.Println(err)
		return
	}
}

func mainSetup() {
	jst = time.FixedZone("Asis/Tokyo", 9*60*60)
	t = time.Now()
	// to get sample data from 2 days ago to now
	adayago = t.Add(-48 * time.Hour)
	ed = t.Format(layout)
	sd = adayago.Format(layout)
	lastupdate = withings.OffsetBase
	//lastupdate = time.Date(2020, 12, 20, 0, 0, 0, 0, time.UTC)
}

func printMeas(v withings.MeasureData, name, unit string) {
	fmt.Printf("%s(Grpid:%v, Category:%v, Attrib: %v, DeviceID:%v)\n", name, v.GrpID, v.Category, v.Attrib, v.DeviceID)
	fmt.Printf("%v, %.1f %s\n", v.Date.In(jst).Format(layout2), v.Value, unit)
}

func testGetmeas() {

	fmt.Println("========== Getmeas[START] ========== ")
	mym, err := client.GetMeas(withings.Real, adayago, t, lastupdate, 0, false, true, withings.Weight, withings.Height, withings.FatFreeMass, withings.BoneMass, withings.FatRatio, withings.FatMassWeight, withings.Temp, withings.HeartPulse, withings.Hydration)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Status: %d\n", mym.Status)

	for _, v := range mym.SerializedData.Weights {
		printMeas(v, "Weight", "Kg")
	}
	for _, v := range mym.SerializedData.FatFreeMass {
		printMeas(v, "FatFreeMass", "Kg")
	}
	for _, v := range mym.SerializedData.FatRatios {
		printMeas(v, "FatRatio", "%%")
	}
	for _, v := range mym.SerializedData.FatMassWeights {
		printMeas(v, "FatMassWeight", "Kg")
	}
	for _, v := range mym.SerializedData.BoneMasses {
		printMeas(v, "BoneMass", "Kg")
	}

	for _, v := range mym.SerializedData.UnknowVals {
		printMeas(v, "UnknownVal", "N/A")
	}

	// Raw data should be provided from mym.Body.Measuregrps
	for _, v := range mym.Body.Measuregrps {
		weight := float64(v.Measures[0].Value) * math.Pow10(v.Measures[0].Unit)
		fmt.Printf("Weight:%.1f Kgs\n", weight)
	}
	fmt.Printf("More:%d, Offset:%d\n", mym.Body.More, mym.Body.Offset)

	fmt.Println("========== Getmeas[END] ========== ")
}

func testGetactivity() {

	fmt.Println("========== Getactivity[START] ========== ")

	act, err := client.GetActivity(sd, ed, 0, 0, withings.Steps, withings.Calories, withings.HrAverage, withings.HrMin, withings.HrMax)

	if err != nil {
		fmt.Println("getActivity Error.")
		fmt.Println(err)
		return
	}

	//fmt.Println(act)
	for _, v := range act.Body.Activities {
		fmt.Printf("Date:%s, Steps:%d, BurnedCalories: %g, HRAverage: %d, HRMinimum: %d, HRMax:%d \n", v.Date, v.Steps, v.Calories, v.HrAverage, v.HrMin, v.HrMax)
	}
	fmt.Println("========== Getactivity[END] ========== ")
}

func testGetsleep() {
	fmt.Println("========== Getsleep[START] ========== ")

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
	//fmt.Println(slp)
	fmt.Println("========== Getsleep[END] ========== ")

}

func testGetsleepsummary() {
	fmt.Println("========== Getsleepsummary[START] ========== ")

	slpsum, err := client.GetSleepSummary(sd, ed, 0, withings.SSBdi, withings.SSDsd, withings.SSD2s, withings.SSD2w, withings.SSHrAvr, withings.SSHrMax, withings.SSHrMin, withings.SSLsd, withings.SSRsd, withings.SSRRAvr, withings.SSRRMax, withings.SSRRMin, withings.SSSS,
		withings.SSSng, withings.SSSngEC, withings.SSWupC, withings.SSWupD)

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
	//fmt.Println(slpsum)
	fmt.Println("========== Getsleepsummary[END] ========== ")
}

func main() {

	settings = withings.ReadSettings(".test_settings.yaml")

	auth(settings)
	tokenFuncs()
	mainSetup()

	testGetmeas()
	testGetactivity()
	testGetsleep()
	testGetsleepsummary()
}
