package withings

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
)

const (
	tokenFile = "access_token.json"
	layout    = "2006-01-02"
	layout2   = "2006-01-02 15:04:05"
	isnotify  = false
)

const (
	testSettingsFile = ".test_settings.yaml"
	testMeasureFile  = "sample_measure.json"
	testActivityFile = "sample_activity.json"
	testSleepFile    = "sample_sleep.json"
)

var (
	jst        *time.Location
	tnow       time.Time
	adayago    time.Time
	lastupdate time.Time
	ed         string
	sd         string
	settings   map[string]string
	client     *(Client)
)

func setupForTest(settingsFile string, t *testing.T) {
	settings = ReadSettings(settingsFile)
	authWithTokenFile(settings, t)
	tokenFuncs()
	// to get data from a day ago to now
	jst = time.FixedZone("Asis/Tokyo", 9*60*60)
	tnow = time.Now()
	//adayago := t.Add(-96 * time.Hour)
	adayago = tnow.Add(-24 * time.Hour)
	ed = tnow.Format(layout)
	sd = adayago.Format(layout)
	//lastupdate = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	lastupdate = OffsetBase
	//lastupdate = time.Date(2020, 12, 20, 0, 0, 0, 0, time.UTC)
}

func authWithTokenFile(settings map[string]string, t *testing.T) {
	var err error
	client, err = New(settings["cid"], settings["secret"], settings["redirect_url"])

	if err != nil {
		t.Errorf("Failed to create New client:%v", err)
		return
	}

	if _, err := os.Open(tokenFile); err == nil {
		_, err = client.ReadToken(tokenFile)

		if err != nil {
			t.Fatalf("Failed to read token file:%v", err)
			return
		}
	} else {
		t.Fatalf("authWithTokenFile needs token file.")
	}

	return
}

func tokenFuncs() {
	// Show token
	client.PrintToken()

	// Refresh Token if you need
	//_, rf, err := client.RefreshToken()
	//if err != nil {
	//	fmt.Println("Failed to RefreshToken")
	//	fmt.Println(err)
	//	return
	//}
	//if rf {
	//	fmt.Println("You got new token!")
	//}

	//// Save Token if you need
	//err = client.SaveToken(tokenFile)
	//if err != nil {
	//	fmt.Println("Failed to RefreshToken")
	//	fmt.Println(err)
	//	return
	//}

	//client.PrintToken()
}

//func TestGetMeas(t *testing.T) {
//	setupForTest(testSettingsFile, t)
//	jsonBlob, err := ioutil.ReadFile(testMeasureFile)
//	if err != nil {
//		t.Fatalf("ioutil.ReadFile returns error(%v)", err)
//	}
//
//	ts := httptest.NewServer(http.HandlerFunc(
//		func(w http.ResponseWriter, r *http.Request) {
//			w.Header().Set("Content-Type", "application/json")
//			//d := TestStruct{"Bob", 18}
//			//slcBob, _ := json.Marshal(d)
//			//w.Write(slcBob)
//			w.Write(jsonBlob)
//		}))
//	defer ts.Close()
//
//	client.MeasureURL = ts.URL
//	res, err := client.GetMeas(Real, adayago, tnow, lastupdate, 0, Weight, Height, FatFreeMass)
//	if err != nil {
//		t.Errorf("client.GetMeas returned error:%v", err)
//	} else {
//		fmt.Println(res)
//	}
//}

func callCreateDataFields(correct string, mtype interface{}, t *testing.T) {
	df, err := createDataFields(mtype)
	if err != nil {
		t.Errorf("createDataFields, got error:%v", err)
	} else {
		fmt.Println(df)
		if df != correct {
			t.Errorf("createDataFields, = %v, want %v", df, correct)
		}
	}
}
func TestCreateDataFieldsInt(t *testing.T) {
	correct := "1,4,5,6,8,9,10,11,12,54,71,73,76,77,88,91,123"
	mtype := [...]MeasType{Weight, Height, FatFreeMass, FatRatio, FatMassWeight, DiastolicBP, SystolicBP, HeartPulse, Temp, SPO2, BodyTemp, SkinTemp, MuscleMass, Hydration, BoneMass, PWaveVel, VO2}
	fmt.Println(mtype)

	callCreateDataFields(correct, mtype, t)
}

func TestCreateDataFieldsString(t *testing.T) {
	correct := "steps,distance,elevation,soft,moderate,intense,active,calories,totalcalories,hr_average,hr_min,hr_max,hr_zone_0,hr_zone_1,hr_zone_2,hr_zone_3"

	mtype := [...]ActivityType{Steps, Distance, Elevation, Soft, Moderate, Intense, Active, Calories, TotalCalories, HrAverage, HrMin, HrMax, HrZone0, HrZone1, HrZone2, HrZone3}
	fmt.Println(mtype)

	callCreateDataFields(correct, mtype, t)
}

func compareRequests(t *testing.T, wURL *url.URL, wMethod, wProto, wBody string, res *http.Request) bool {
	if res.URL.String() != wURL.String() {
		t.Errorf("createRequest URL = %s, want %s", res.URL.String(), wURL.String())
	}
	if res.Method != wMethod {
		t.Errorf("createRequest = %s, want %s", res.Method, wMethod)
		return false
	}
	if res.Proto != wProto {
		t.Errorf("createRequest = %s, want %s", res.Proto, wProto)
		return false
	}
	b, e := res.GetBody()
	if e != nil {
		t.Errorf("res.GetBody return error(%v)", e)
		return false
	} else {
		bbuf, e := ioutil.ReadAll(b)
		if e != nil {
			t.Errorf("ioutil.ReadAll(b) return error(%v)", e)
			return false
		} else {
			bstr := string(bbuf)
			if bstr != wBody {
				t.Errorf("createRequest = %s, want %s", bstr, wBody)
				return false
			} else {
				return true
			}
		}
	}
}

func TestCreateRequest(t *testing.T) {
	URL := "https://example.com"
	wBody := "action=getmeas&category=1&meastype=1%2C4%2C5%2C6%2C8%2C9%2C10%2C11%2C12%2C54%2C71%2C73%2C76%2C77%2C88%2C91%2C123"
	wMethod := http.MethodPost
	wURL, _ := url.ParseRequestURI(URL)
	wProto := "HTTP/1.1"

	df := "1,4,5,6,8,9,10,11,12,54,71,73,76,77,88,91,123"

	fp := []FormParam{
		{PPaction, MeasureA},
		{PPmeastype, df},
		{PPcategory, fmt.Sprintf("%d", Real)},
	}

	ctx, cancel := getNewContext(5)
	defer cancel()

	res, err := createRequest(ctx, fp, URL, http.MethodPost)
	if err != nil {
		t.Errorf("createRequest returns error(%v)", err)
	} else {
		b, e := httputil.DumpRequest(res, true)
		if e != nil {
			t.Errorf("DumpRequest returns error(%v)", e)
		} else {
			fmt.Println(string(b))
			compareRequests(t, wURL, wMethod, wProto, wBody, res)
		}
	}
}

type TestStruct struct {
	Name string `json:"name"`
	Age  int    `json:"int"`
}

// TODO: add fail pattern
func TestParseResponse(t *testing.T) {

	jsonBlob, err := ioutil.ReadFile(testMeasureFile)
	if err != nil {
		t.Fatalf("ioutil.ReadFile returns error(%v)", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonBlob)
		}))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Errorf("http.Get returns error(%v)", err)
	}
	defer resp.Body.Close()

	mym := new(Measurement)

	if err == nil {
		bytes, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("\tDumpResponse Results\n%s\n", string(bytes))
		e := parseResponse(resp, mym)
		if e != nil {
			t.Errorf("parseResponse returns error(%v)", e)
		}

	} else {
		t.Errorf("http.Get returns error(%v)", err)
		fmt.Println(err)
	}
}
