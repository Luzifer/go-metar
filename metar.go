package metar

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	apiSource = "https://www.aviationweather.gov/adds/dataserver_current/httpparam?dataSource=metars&requestType=retrieve&format=xml&stationString=%s&hoursBeforeNow=2&mostRecent=true"
)

var (
	// HTTPClient is used to make requests, you can insert your own
	HTTPClient = http.DefaultClient
)

// Result holds all the data from the METAR request
type Result struct {
	XMLName             xml.Name            `xml:"METAR"`
	RawText             string              `xml:"raw_text"`              // The raw METAR
	StationID           string              `xml:"station_id"`            // Station identifier; Always a four character alphanumeric( A-Z, 0-9)
	ObservationTime     time.Time           `xml:"observation_time"`      // Time this METAR was observed
	Latitude            float64             `xml:"latitude"`              // The latitude (in decimal degrees) of the station that reported this METAR
	Longitude           float64             `xml:"longitude"`             // The longitude (in decimal degrees) of the station that reported this METAR
	Temperature         float64             `xml:"temp_c"`                // Air temperature (celsius)
	Dewpoint            float64             `xml:"dewpoint_c"`            // Dewpoint temperature (celsius)
	WindDirDegrees      int64               `xml:"wind_dir_degrees"`      // Direction from which the wind is blowing. 0 degrees=variable wind direction.
	WindSpeed           int64               `xml:"wind_speed_kt"`         // Wind speed; 0 degree wdir and 0 wspd = calm winds (kts)
	WindGust            int64               `xml:"wind_gust_kt"`          // Wind gust
	VisibilityStatute   float64             `xml:"visibility_statute_mi"` // Horizontal visibility (statute miles)
	Altimeter           float64             `xml:"altim_in_hg"`           // Altimeter (inches of Hg)
	SeaLevelPressure    float64             `xml:"sea_level_pressure_mb"` // Sea-level pressure (mb)
	QualityControlFlags QualityControlFlags `xml:"quality_control_flags"` // Quality control flags provide useful information about the METAR station(s) that provide the data.
	WXString            string              `xml:"wx_string"`             // WX string descriptions (https://www.aviationweather.gov/static/adds/docs/metars/wxSymbols_anno2.pdf)
	SkyCondition        struct {
		SkyCover SkyCover `xml:"sky_cover,attr"` // Sky cover, up to four levels of sky cover can be reported ; OVX present when vert_vis_ft is reported
	} `xml:"sky_condition"`
	FlightCategory FlightCategory `xml:"flight_category"` // Flight category of this METAR
	// Fields 19 to 29 currently not implemented
	MetarType string  `xml:"metar_type"`  // METAR or SPECI
	Elevation float64 `xml:"elevation_m"` // The elevation of the station that reported this METAR (meters)
}

// QualityControlFlags provide useful information about the METAR station(s) that provide the data.
type QualityControlFlags struct {
	XMLName  xml.Name `xml:"quality_control_flags"`
	NoSignal bool     `xml:"no_signal"`
}

// SkyCover defines and explains possible sky coverage situations
type SkyCover string

// Common SkyCover situations
const (
	SkyCoverSKC   SkyCover = "SKC"   // "No cloud/Sky clear" used worldwide but in North America is used to indicate a human generated report
	SkyCoverCLR   SkyCover = "CLR"   // "No clouds below 12,000 ft (3,700 m) (U.S.) or 25,000 ft (7,600 m) (Canada)", used mainly within North America and indicates a station that is at least partly automated
	SkyCoverNSC   SkyCover = "NSC"   // "No (nil) significant cloud", i.e., none below 5,000 ft (1,500 m) and no TCU or CB. Not used in North America.
	SkyCoverFEW   SkyCover = "FEW"   // "Few" = 1–2 oktas
	SkyCoverSCT   SkyCover = "SCT"   // "Scattered" = 3–4 oktas
	SkyCoverBKN   SkyCover = "BKN"   // "Broken" = 5–7 oktas
	SkyCoverOVC   SkyCover = "OVC"   //	"Overcast" = 8 oktas, i.e., full cloud coverage
	SkyCoverCAVOK SkyCover = "CAVOK" // Ceiling And Visibility OKay, indicating no cloud below 5,000 ft (1,500 m) or the highest minimum sector altitude and no cumulonimbus or towering cumulus at any level, a visibility of 10 km (6 mi) or more and no significant weather change
)

// FlightCategory defines and explains possible flight category situations
type FlightCategory string

// Common FlightCategory situations
const (
	FlightCategoryVFR  FlightCategory = "VFR"  // Visual Flight Rules (Ceiling greater than 3,000 feet AGL and visibility greater than 5 miles)
	FlightCategoryMVFR FlightCategory = "MVFR" // Marginal Visual Flight Rules (Ceiling 1,000 to 3,000 feet AGL and/or visibility 3 to 5 miles)
	FlightCategoryIFR  FlightCategory = "IFR"  // Instrument Flight Rules (Ceiling 500 to below 1,000 feet AGL and/or visibility 1 mile to less than 3 miles)
	FlightCategoryLIFR FlightCategory = "LIFR" // Low Instrument Flight Rules (Ceiling below 500 feet AGL and/or visibility less than 1 mile)
)

type response struct {
	XMLName xml.Name `xml:"response"`
	Data    struct {
		NumResults int      `xml:"num_results,attr"`
		Results    []Result `xml:"METAR"`
	} `xml:"data"`
}

// FetchCurrentStationWeather fetches the last result from the specified station if it was reported during last 2 hours
func FetchCurrentStationWeather(station string) (*Result, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf(apiSource, station), nil)
	res, err := HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	r := &response{}
	if err = xml.NewDecoder(res.Body).Decode(r); err != nil {
		return nil, err
	}

	if r.Data.NumResults != len(r.Data.Results) {
		return nil, errors.New("Got inconsistent number of results")
	}

	if r.Data.NumResults == 0 {
		return nil, errors.New("Did not find any data for your station")
	}

	return &r.Data.Results[0], nil
}

// InHgTohPa converts "inch of mercury" to "hectopascal"
func InHgTohPa(inHg float64) float64 {
	return inHg * 33.8638866667
}

// KtsToMs converts "knots" to "meters per second"
func KtsToMs(kts float64) float64 {
	return kts * 0.514444
}

// StatMileToKm converts "statute miles" to "kilometers"
func StatMileToKm(sm float64) float64 {
	return sm * 1.60934
}

// MbTohPa converts "millibar" to "hectopascal"
func MbTohPa(mb float64) float64 {
	return mb * 0.1
}
