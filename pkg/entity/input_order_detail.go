package entity

import "time"

type InputOrderDetail struct {
	MerchantID          int16
	MerchantOrderID     string
	PickupInfo          *InputPickupInfo
	DropoffInfo         *InputDropoffInfo
	MerchantNotesOption *string
}

type InputPickupInfo struct {
	Location           *Location
	PhoneNumber        string
	EarliestPickupTime *time.Time // ISO8601 UTC
}

type InputDropoffInfo struct {
	Location            *Location
	FirstName           string
	LastName            string
	Email               string
	PhoneNumber         string
	EarliestDropoffTime *time.Time // ISO8601 UTC
	LatestDropoffTime   *time.Time // ISO8601 UTC
}
