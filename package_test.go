package usps // import "github.com/justinbarrick/usps-api"

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNormalizePackage(t *testing.T) {
    track_response := TrackResponse{
        TrackInfo: TrackInfo{
            StatusCategory: "Delivered",
            OriginCity: "SANTA FE SPRINGS",
            OriginState: "CA",
            DestinationCity: "SAN FRANCISCO",
            DestinationState: "CA",
        },
    }

    pstatus, err := track_response.Normalize()
    if err != nil {
        t.Fatal(err)
    }

    assert.Equal(t, pstatus.Status.StatusDelivered, true, "Package should be delivered")
}
