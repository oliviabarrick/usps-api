package usps // import "github.com/justinbarrick/usps-api"

import (
    "encoding/xml"
    "github.com/stretchr/testify/assert"
    "io"
    "strings"
    "testing"
    "net/http"
    "net/http/httptest"
    "net/url"
    "os"
)

func testTrackFieldRequest(t *testing.T, parsed_url *url.URL, track_request TrackFieldRequest) {
    query, err := url.ParseQuery(parsed_url.RawQuery)
    if err != nil {
        t.Error(err)
    }

    var xml_request TrackFieldRequest
    decoder := xml.NewDecoder(strings.NewReader(query["XML"][0]))
    err = decoder.Decode(&xml_request)
    if err != nil {
        t.Error(err)
    }

    assert.Equal(t, query["API"][0], "TrackV2", "api mismatch")
    assert.Equal(t, xml_request, track_request, "request mismatch")
}

func TestConstructURL(t *testing.T) {
    var id TrackID
    id.Id = "abcd"

    var track_request TrackFieldRequest
    track_request.UserId = "1234"
    track_request.TrackID = id
    track_request.Revision = 1
    track_request.ClientIp = "111.0.0.1"
    track_request.SourceId = "hello"

    api_url := "https://secure.shippingapis.com/ShippingAPI.dll"

    usps_url, err := track_request.construct_url(api_url)
    if err != nil {
        t.Error(err)
    }

    parsed_url, err := url.Parse(usps_url)
    if err != nil {
        t.Error(err)
    }

    assert.Equal(t, parsed_url.Host, "secure.shippingapis.com", "hostname mismatch")
    assert.Equal(t, parsed_url.Path, "/ShippingAPI.dll", "path mismatch")

    testTrackFieldRequest(t, parsed_url, track_request)
}

func TestUsps(t *testing.T) {
    pt := &PackageTracker{}

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var id TrackID
        id.Id = "thisismytrackingnumber"

        var track_request TrackFieldRequest
        track_request.UserId = "1234"
        track_request.TrackID = id
        track_request.Revision = 1
        track_request.ClientIp = "111.0.0.1"
        track_request.SourceId = "hello"

        testTrackFieldRequest(t, r.URL, track_request)

        file, err := os.Open("test_data/usps.xml")
        if err != nil {
            t.Error(err)
        }

        io.Copy(w, file)
    }))
    defer ts.Close()

    pt.ApiUrl = ts.URL
    pt.UserId = "1234"

    track_response, err := pt.Fetch("thisismytrackingnumber")
    if err != nil {
        t.Error(err)
    }

    assert.Equal(t, track_response.TrackInfo.StatusCategory, "Delivered", "status category mismatch")
    assert.Equal(t, track_response.TrackInfo.Status[0], "Delivered, Individual Picked Up at Postal Facility", "status mismatch")

    assert.Equal(t, track_response.TrackInfo.OriginCity, "SANTA FE SPRINGS", "origin city mismatch")
    assert.Equal(t, track_response.TrackInfo.OriginState, "CA", "origin state mismatch")
    assert.Equal(t, track_response.TrackInfo.DestinationCity, "SAN FRANCISCO", "destination city mismatch")
    assert.Equal(t, track_response.TrackInfo.DestinationState, "CA", "destination state mismatch")

    detail := track_response.TrackInfo.TrackDetails[0]
    assert.Equal(t, detail.Event, "Business Closed", "event mismatch")
    assert.Equal(t, detail.EventDate, "April 7, 2017", "event date mismatch")
    assert.Equal(t, detail.EventCity, "SAN FRANCISCO", "event city mismatch")

    detail = track_response.TrackInfo.TrackDetails[1]
    assert.Equal(t, detail.Event, "Redelivery Scheduled", "event mismatch")
    assert.Equal(t, detail.EventDate, "April 5, 2017", "event date mismatch")
    assert.Equal(t, detail.EventCity, "SAN FRANCISCO", "event city mismatch")
}

func TestUspsAuthError(t *testing.T) {
    pt := &PackageTracker{}

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        file, err := os.Open("test_data/usps_error.xml")
        if err != nil {
            t.Error(err)
        }

        io.Copy(w, file)
    }))
    defer ts.Close()

    pt.ApiUrl = ts.URL
    pt.UserId = "1234"

    track_response, _ := pt.Fetch("thisismytrackingnumber")

    assert.Equal(t, track_response.Error.Number, "80040B1A", "error code mismatch")
    assert.Equal(t, track_response.Error.Description, "Authorization failure.  Perhaps username and/or password is incorrect.", "error description mismatch")
    assert.Equal(t, track_response.Error.Source, "USPSCOM::DoAuth", "error source mismatch")
}

func TestUspsTrackInfoError(t *testing.T) {
    pt := &PackageTracker{}

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        file, err := os.Open("test_data/usps_tracking_error.xml")
        if err != nil {
            t.Error(err)
        }

        io.Copy(w, file)
    }))
    defer ts.Close()

    pt.ApiUrl = ts.URL
    pt.UserId = "1234"

    track_response, _ := pt.Fetch("thisismytrackingnumber")

    assert.Equal(t, track_response.Error.Number, "-2147219302", "error code mismatch")
    assert.Equal(t, track_response.Error.Description, "The tracking number may be incorrect or the status update is not yet available. Please verify your tracking number and try again later.", "error description mismatch")
}
