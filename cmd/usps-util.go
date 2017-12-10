package main

import (
    "log"
    "os"
    "usps"
)

func main() {
    if len(os.Args) != 3 {
        log.Fatalf("Usage: %s <user id> <tracking code>", os.Args[0])
    }

    var pt usps.PackageTracker
    pt.UserId = os.Args[1]

    track_response, err := pt.Fetch(os.Args[2])
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Order status: ", track_response.TrackInfo.Status)
    log.Printf("%s, %s -> %s, %s\n", track_response.TrackInfo.OriginCity, track_response.TrackInfo.OriginState, track_response.TrackInfo.DestinationCity, track_response.TrackInfo.DestinationState)
}
