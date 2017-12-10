package usps

import (
    "fmt"
    "net/http"
    "bytes"
    "encoding/xml"
    "net/url"
)

type TrackResponse struct {
    TrackInfo TrackInfo
}

type TrackInfo struct {
    ID string `xml:"ID,attr"`
    Class string
    ClassOfMailCode string
    DestinationCity string
    DestinationState string
    DestinationZip string
    EmailEnabled bool
    KahalaIndicator bool
    MailTypeCode string
    MPDATE string
    MPSUFFIX string
    OriginCity string
    OriginState string
    OriginZip string
    PodEnabled bool
    RestoreEnabled bool
    RramEnabled bool
    RreEnabled bool
    Service []string
    ServiceTypeCode int
    Status []string
    StatusCategory string
    StatusSummary string
    TABLECODE string
    TrackSummary TrackSummary
    TrackDetails []TrackDetail `xml:"TrackDetail"`
}

type TrackSummary struct {
    EventTime string
    EventDate string
    Event string
    EventCity string
    EventState string
    EventZIPCode string
    EventCount string
    FirmName string
    Name string
    AuthorizedAgent bool
    EventCode int
    DeliveryAttributeCode int
}

type TrackDetail struct {
    EventTime string
    EventDate string
    Event string
    EventCity string
    EventState string
    EventZIPCode string
    EventCountry string
    FirmName string
    Name string
    AuthorizedAgent bool
    EventCode string
}

type TrackFieldRequest struct {
    UserId string `xml:"USERID,attr"`
    Revision int
    ClientIp string
    SourceId string
    TrackID TrackID
}

type TrackID struct {
    Id string `xml:"ID,attr"`
}

func (tr TrackFieldRequest) construct_url(base_url string) (usps_url string, err error) {
    proc_inst := xml.ProcInst{
        Target: "xml",
        Inst:   []byte("version=\"1.0\" encoding=\"UTF-8\""),
    }

    encoded := &bytes.Buffer{}

    enc := xml.NewEncoder(encoded)
    enc.Indent(" ", "    ")

    if err = enc.EncodeToken(proc_inst); err != nil {
        return
    }

    if err = enc.Encode(tr); err != nil {
        return
    }

    parsed_url, err := url.Parse(base_url)
    if err != nil {
        return
    }

    parameters := url.Values{}
    parameters.Add("XML", encoded.String())
    parameters.Add("API", "TrackV2")
    parsed_url.RawQuery = parameters.Encode()

    usps_url = parsed_url.String()
    return
}

type PackageTracker struct {
    client *http.Client
    ApiUrl string
    UserId string
}

func (pt *PackageTracker) Fetch(track_id string) (track_response TrackResponse, err error) {
    if pt.ApiUrl == "" {
        pt.ApiUrl = "https://secure.shippingapis.com/ShippingAPI.dll"
    }

    var id TrackID
    id.Id = track_id

    var track_request TrackFieldRequest
    track_request.UserId = pt.UserId
    track_request.TrackID = id
    track_request.Revision = 1
    track_request.ClientIp = "111.0.0.1"
    track_request.SourceId = "hello"

    usps_url, err := track_request.construct_url(pt.ApiUrl)
    if err != nil {
        return
    }

    if pt.client == nil {
        pt.client = &http.Client{}
    }

    resp, err := pt.client.Get(usps_url)
    if err != nil {
        return
    }

    if resp.StatusCode != 200 {
        err = fmt.Errorf("Response status %d != 200\n", resp.StatusCode)
        return
    }

    /*
    buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)
    fmt.Println(buf.String())
    */

    if err = xml.NewDecoder(resp.Body).Decode(&track_response); err != nil {
        return
    }

    return
}
