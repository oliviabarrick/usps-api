package usps

type PackageStatus struct {
    StatusDelivered bool
}

type Location struct {
    City string
    State string
    Zip string
}

type CarrierEvent struct {
    TrackID string
    EventLocation Location
    Status PackageStatus
    Description string
}

type Package struct {
    TrackID string
    Origin Location
    Destination Location
    Status PackageStatus
    Events []CarrierEvent
}

func (tr *TrackResponse) Normalize() (p Package, err error) {
    if tr.TrackInfo.StatusCategory == "Delivered" {
        p.Status.StatusDelivered = true
    }

    return
}
