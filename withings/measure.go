package withings

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// FormParam is for http request parameter.
type FormParam struct {
	key   string
	value string
}

// OffsetBase is used to check whether lastupdate or startdate/enddate is used in GetMeas
var OffsetBase time.Time = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

func createRequest(ctx context.Context, fp []FormParam, uri, method string) (*http.Request, error) {
	form := url.Values{}
	for _, v := range fp {
		form.Add(v.key, v.value)
	}

	body := strings.NewReader(form.Encode())
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), body)
	req = req.WithContext(ctx)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func parseResponse(resp *http.Response, result interface{}) error {
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(rbody, result)
}

func reqAndParse(c *Client, fp []FormParam, url, method string, result interface{}) error {
	ctx, cancel := getNewContext(c.Timeout)
	defer cancel()

	req, err := createRequest(ctx, fp, url, method)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = parseResponse(resp, result)
	if err != nil {
		return err
	}
	return nil
}

func createDataFields(fields interface{}) (string, error) {
	df := []string{}

	rv := reflect.ValueOf(fields)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			irv := rv.Index(i)
			switch irv.Kind() {
			case reflect.String:
				df = append(df, irv.String())
			case reflect.Int:
				df = append(df, fmt.Sprintf("%d", irv.Int()))
			default:
				return "", errors.Errorf("createDataFields allows string or int values.")
			}
		}
	default:
		return "", errors.Errorf("createDataFields Slice or Array.")
	}
	return strings.Join(df, ","), nil
}

// GetMeas call withings API Measure - GetMeas. (https://developer.withings.com/oauth2/#operation/measure-getmeas)
// cattype: category, 1 for real measures, 2 for user objectives
// startdate, enddate: Measure's start date, end date.
// lastupdate : Timestamp for requesting data that were updated or created after this date. Use this instead of startdate+endate.
//              If lastupdate is set to a timestamp other than Offsetbase, getMeas will use lastupdate in preference to startdate/enddate.
// offset: When a first call retuns more:1 and offset:XX, set value XX in this parameter to retrieve next available rows.
// isOldToNew: If true, results must be sorted by oldest to newest. If false, results must be sorted by newest to oldest.
// isSerialized: if true, results must be parsed to Measurement.SerializedData
// mtype: Measurement Type. Set the measurement type you want to get data. See MeasType in enum.go.
func (c *Client) GetMeas(cattype CatType, startdate, enddate, lastupdate time.Time, offset int, isOldToNew, isSerialized bool, mtype ...MeasType) (*Measurement, error) {

	if len(mtype) == 0 {
		return nil, errors.Errorf("Need least one param as MeasType.")
	}

	mym := new(Measurement)

	df, err := createDataFields(mtype)
	if err != nil {
		return nil, err
	}

	fp := []FormParam{
		{PPaction, MeasureA},
		{PPmeastype, df},
		{PPcategory, fmt.Sprintf("%d", cattype)},
	}

	if offset != 0 {
		fp = append(fp, FormParam{PPoffset, fmt.Sprintf("%d", offset)})
	}

	// The below sentence comes from Withings API document.
	// 	> Timestamp for requesting data that were updated or created after this date.
	// 	> Useful for data synchronization between systems.
	// 	> Use this instead of startdate + enddate.
	if lastupdate != OffsetBase {
		fp = append(fp, FormParam{PPlastupdate, strconv.FormatInt(lastupdate.Unix(), 10)})
	} else {
		fp = append(fp, FormParam{PPstartdate, strconv.FormatInt(startdate.Unix(), 10)})
		fp = append(fp, FormParam{PPenddate, strconv.FormatInt(enddate.Unix(), 10)})
	}

	err = reqAndParse(c, fp, c.MeasureURL, http.MethodPost, mym)
	if err != nil {
		return nil, err
	}

	if isOldToNew {
		sort.Slice(mym.Body.Measuregrps, func(i, j int) bool {
			return time.Unix(int64(mym.Body.Measuregrps[i].Date), 0).Before(time.Unix(int64(mym.Body.Measuregrps[j].Date), 0))
		})
	} else {
		sort.Slice(mym.Body.Measuregrps, func(i, j int) bool {
			return time.Unix(int64(mym.Body.Measuregrps[i].Date), 0).After(time.Unix(int64(mym.Body.Measuregrps[j].Date), 0))
		})
	}

	if isSerialized {
		mym.SerializedData, err = SerialMeas(mym)
		if err != nil {
			mym.SerializedData = nil
			return mym, err
		}
	}
	return mym, nil
}

// SerialMeas will parse measurement results.
func SerialMeas(mym *Measurement) (*SerialzedMeas, error) {

	sm := new(SerialzedMeas)

	for _, mGrp := range mym.Body.Measuregrps {
		for _, meas := range mGrp.Measures {
			val := MeasureData{
				GrpID:    mGrp.GrpID,
				Date:     time.Unix(int64(mGrp.Date), 0),
				Value:    float64(meas.Value) * math.Pow10(meas.Unit),
				Attrib:   mGrp.Attrib,
				Category: mGrp.Category,
				DeviceID: mGrp.DeviceID,
			}
			switch meas.Type {
			case int(Weight):
				sm.Weights = append(sm.Weights, val)
			case int(Height):
				sm.Heights = append(sm.Heights, val)
			case int(FatFreeMass):
				sm.FatFreeMass = append(sm.FatFreeMass, val)
			case int(FatRatio):
				sm.FatRatios = append(sm.FatRatios, val)
			case int(FatMassWeight):
				sm.FatMassWeights = append(sm.FatMassWeights, val)
			case int(DiastolicBP):
				sm.DiastolicBPs = append(sm.DiastolicBPs, val)
			case int(SystolicBP):
				sm.SystolicBPs = append(sm.SystolicBPs, val)
			case int(HeartPulse):
				sm.HeartPulses = append(sm.HeartPulses, val)
			case int(Temp):
				sm.Temps = append(sm.Temps, val)
			case int(SPO2):
				sm.SPO2s = append(sm.SPO2s, val)
			case int(BodyTemp):
				sm.BodyTemps = append(sm.BodyTemps, val)
			case int(SkinTemp):
				sm.SkinTemps = append(sm.SkinTemps, val)
			case int(MuscleMass):
				sm.MuscleMasses = append(sm.MuscleMasses, val)
			case int(Hydration):
				sm.Hydration = append(sm.Hydration, val)
			case int(BoneMass):
				sm.BoneMasses = append(sm.BoneMasses, val)
			case int(PWaveVel):
				sm.PWaveVel = append(sm.PWaveVel, val)
			case int(VO2):
				sm.VO2s = append(sm.VO2s, val)
			default:
				sm.UnknowVals = append(sm.UnknowVals, val)
			}
		}
	}
	return sm, nil
}

// GetActivity call withings API Measure v2 - Getactivity. (https://developer.withings.com/oauth2/#operation/measurev2-getactivity)
// startdate/enddate: Activity result start date, end date.
// lastupdate : Timestamp for requesting data that were updated or created after this date. Use this instead of startdate+endate.
//              If lastupdate is set to a timestamp other than Offsetbase, getMeas will use lastupdate in preference to startdate/enddate.
// offset: When a first call retuns more:1 and offset:XX, set value XX in this parameter to retrieve next available rows.
// atype: Acitivity Type. Set the activity type you want to get data. See ActivityType in enum.go.
func (c *Client) GetActivity(startdate, enddate string, lastupdate int, offset int, atype ...ActivityType) (*Activities, error) {
	if len(atype) == 0 {
		return nil, errors.Errorf("Need least one param as ActivityType.")
	}
	act := new(Activities)

	df, err := createDataFields(atype)

	var fp []FormParam = []FormParam{
		{PPaction, ActivityA},
		{PPoffset, fmt.Sprintf("%d", offset)},
		{PPdataFields, df},
	}

	if startdate == "" || enddate == "" {
		fp = append(fp, FormParam{PPlastupdate, fmt.Sprintf("%d", lastupdate)})
	} else {
		fp = append(fp, FormParam{PPstartdateymd, startdate}, FormParam{PPenddateymd, enddate})
	}

	err = reqAndParse(c, fp, c.MeasureURLv2, http.MethodPost, act)
	if err != nil {
		return nil, err
	}
	return act, nil
}

// GetWorkouts call withings API Measure v2 - Getworkouts. (https://developer.withings.com/api-reference#operation/measurev2-getworkouts)
// startdate/enddate: Workouts result start date, end date.
// lastupdate : Timestamp for requesting data that were updated or created after this date. Use this instead of startdate+endate.
//              If lastupdate is set to a timestamp other than Offsetbase, GetWorkouts will use lastupdate in preference to startdate/enddate.
// offset: When a first call retuns more:1 and offset:XX, set value XX in this parameter to retrieve next available rows.
// wtype: Workout Type. Set the workout type you want to get data. See WorkoutType in enum.go.
func (c *Client) GetWorkouts(startdate, enddate string, lastupdate int, offset int, wtype ...WorkoutType) (*Workouts, error) {
	if len(wtype) == 0 {
		return nil, errors.Errorf("Need least one param as WorkoutType.")
	}
	workouts := new(Workouts)

	df, err := createDataFields(wtype)

	var fp []FormParam = []FormParam{
		{PPaction, WorkoutsA},
		{PPoffset, fmt.Sprintf("%d", offset)},
		{PPdataFields, df},
	}

	if startdate == "" || enddate == "" {
		fp = append(fp, FormParam{PPlastupdate, fmt.Sprintf("%d", lastupdate)})
	} else {
		fp = append(fp, FormParam{PPstartdateymd, startdate}, FormParam{PPenddateymd, enddate})
	}

	err = reqAndParse(c, fp, c.MeasureURLv2, http.MethodPost, workouts)
	if err != nil {
		return nil, err
	}
	return workouts, nil
}

// GetSleep cal withings API Sleep v2 - Get. (https://developer.withings.com/oauth2/#operation/sleepv2-get)
// startdate/enddate: Measures' start date, end date.
// stype: Sleep Type. Set the sleep type you want to get data. See SleepType in enum.go.
func (c *Client) GetSleep(startdate, enddate time.Time, stype ...SleepType) (*Sleeps, error) {
	if len(stype) == 0 {
		return nil, errors.Errorf("Need least one param as SleepType.")
	}

	df, err := createDataFields(stype)
	if err != nil {
		return nil, err
	}

	var fp []FormParam = []FormParam{
		{PPaction, SleepA},
		{PPstartdate, strconv.FormatInt(startdate.Unix(), 10)},
		{PPenddate, strconv.FormatInt(enddate.Unix(), 10)},
		{PPdataFields, df},
	}

	slp := new(Sleeps)
	err = reqAndParse(c, fp, c.SleepURLv2, http.MethodPost, slp)
	if err != nil {
		return nil, err
	}

	// sort by startdate
	sort.Slice(slp.Body.Series, func(i, j int) bool {
		return slp.Body.Series[i].Startdate < slp.Body.Series[j].Startdate
	})
	return slp, nil
}

// GetSleepSummary call withings API Sleep v2 - Getsummary. (https://developer.withings.com/oauth2/#operation/sleepv2-getsummary)
// startdate/enddate: Measurement result start date, end date.
// lastupdate : Timestamp for requesting data that were updated or created after this date. Use this instead of startdate+endate.
//              If lastupdate is set to a timestamp other than Offsetbase, getMeas will use lastupdate in preference to startdate/enddate.
// stype: Sleep Summaries Type. Set the sleep summaries data you want to get. See SleepSummariesType in enum.go.
func (c *Client) GetSleepSummary(startdate, enddate string, lastupdate int, sstype ...SleepSummariesType) (*SleepSummaries, error) {
	if len(sstype) == 0 {
		return nil, errors.Errorf("Need least one param as SleepSummariesType.")
	}
	df, err := createDataFields(sstype)
	if err != nil {
		return nil, err
	}

	var fp []FormParam = []FormParam{
		{PPaction, SleepSA},
		{PPdataFields, df},
	}

	if startdate == "" || enddate == "" {
		fp = append(fp, FormParam{PPlastupdate, fmt.Sprintf("%d", lastupdate)})
	} else {
		fp = append(fp, FormParam{PPstartdateymd, startdate}, FormParam{PPenddateymd, enddate})
	}
	slpss := new(SleepSummaries)
	err = reqAndParse(c, fp, c.SleepURLv2, http.MethodPost, slpss)
	if err != nil {
		return nil, err
	}

	// sort by startdate
	sort.Slice(slpss.Body.Series, func(i, j int) bool {
		return slpss.Body.Series[i].Startdate < slpss.Body.Series[j].Startdate
	})

	return slpss, nil
}
