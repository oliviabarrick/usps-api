package main

import (
    "log"
    "os"
    "github.com/justinbarrick/usps-api"
    "github.com/jessevdk/go-flags"
)

func main() {
    var opts struct {
        Debug bool `short:"d" long:"debug" description:"Whether or not to enable debug mode."`
        Test bool `short:"t" long:"test" descrition:"Whether or not to use the USPS test server."`
        Args struct {
            UserId string `positional-arg-name:"user_id" description:"User ID to use with USPS API."`
            TrackingCode string `positional-arg-name:"tracking_code" description:"Tracking code to lookup."`
        } `positional-args:"true" required:"2"`
    }

     _, err := flags.Parse(&opts)
    if err != nil {
        if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
            os.Exit(0)
        } else {
          log.Fatal(err)
           os.Exit(1)
        }
    }

    var pt usps.PackageTracker
    pt.UserId = opts.Args.UserId
    pt.Debug = opts.Debug

    if opts.Test == true {
        pt.ApiUrl = "https://stg-secure.shippingapis.com/ShippingAPI.dll"
    }

    track_response, err := pt.Fetch(opts.Args.TrackingCode)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Order status: ", track_response.TrackInfo.Status)
    log.Printf("%s, %s -> %s, %s\n", track_response.TrackInfo.OriginCity, track_response.TrackInfo.OriginState, track_response.TrackInfo.DestinationCity, track_response.TrackInfo.DestinationState)
}
