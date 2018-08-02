package model

import "time"

type OrderDetail struct {
	MerchantID      int16        `json:"merchantId"`
	MerchantOrderID string       `json:"merchantOrderId"`
	Capacity        *Capacity    `json:"capacity,omitempty"`
	PickupInfo      *PickupInfo  `json:"pickupInfo"`
	DropoffInfo     *DropoffInfo `json:"dropoffInfo"`
	MerchantNotes   *string      `json:"merchantNotes,omitempty"`
}

func (o *OrderDetail) HasRequiredFields() bool {
	return o.Capacity != nil && o.Capacity.HasRequiredFields() &&
		o.PickupInfo != nil && o.PickupInfo.HasRequiredFields() &&
		o.DropoffInfo != nil && o.DropoffInfo.HasRequiredFields()
}

type Capacity struct {
	GroceryToteCount int16 `json:"groceryToteCount"`
	LargeItemCount   int16 `json:"largeItemCount"`
}

func (c Capacity) HasRequiredFields() bool { return true }

type PickupInfo struct {
	Location           *Location      `json:"location"`
	PhoneNumber        string         `json:"phoneNumber"`
	EarliestPickupTime *time.Time     `json:"earliestPickupTime"` // ISO8601 UTC
	PickupAccessCode   string         `json:"pickupAccessCode"`
	RouteGoals         []*Geolocation `json:"routeGoals,omitempty"`
}

func (p *PickupInfo) HasRequiredFields() bool {
	return p.Location != nil && p.Location.HasRequiredFields()
}

type DropoffInfo struct {
	Location            *Location      `json:"location"`
	FirstName           string         `json:"firstName"`
	LastName            string         `json:"lastName"`
	Email               string         `json:"email"`
	PhoneNumber         string         `json:"phoneNumber"`
	EarliestDropoffTime *time.Time     `json:"earliestDropoffTime"` // ISO8601 UTC
	LatestDropoffTime   *time.Time     `json:"latestDropoffTime"`   // ISO8601 UTC
	RouteGoals          []*Geolocation `json:"routeGoals,omitempty"`
}

func (d *DropoffInfo) HasRequiredFields() bool {
	return d.Location != nil && d.Location.HasRequiredFields()
}

type Location struct {
	Geolocation *Geolocation `json:"geolocation"`
	Address     *Address     `json:"address"`
}

func (l *Location) HasRequiredFields() bool {
	return l.Geolocation != nil && l.Address != nil
}

type Geolocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Address struct {
	Street1          string  `json:"street1"`
	Street2          *string `json:"street2,omitempty"`
	City             string  `json:"city"`
	State            string  `json:"state"`
	PostalCode       string  `json:"postalCode"`
	FormattedAddress string  `json:"formattedAddress"`
}
