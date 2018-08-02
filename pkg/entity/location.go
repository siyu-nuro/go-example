package entity

// Address entity describes a physical address
type Address struct {
	Street1          string
	Street2          *string
	City             string
	State            string
	PostalCode       string
	FormattedAddress string

	// Eventually, when we add this, we should use
	// https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3
	// (Because apparently it's better than Alpha-2)
	//CountryISOAlpha3 string
}

func (a Address) HasRequiredFields() bool {
	return true
}

func (a Address) StreetAddressString() string {
	if a.Street2 == nil || len(*a.Street2) == 0 {
		return a.Street1
	}
	return a.Street1 + ", " + *a.Street2
}

// This will break in languages other than English
func (a Address) AddressString() string {
	return a.StreetAddressString() + ", " + a.City + ", " + a.State + " " + a.PostalCode
}

// Geolocation entity describes a GPS coordinate
type Geolocation struct {
	Latitude  float64
	Longitude float64
}

func (gl Geolocation) HasRequiredFields() bool {
	return true
}

// Location entity
type Location struct {
	Geolocation *Geolocation
	Address     *Address
}

func (l Location) HasRequiredFields() bool {
	return l.Geolocation != nil && l.Geolocation.HasRequiredFields() &&
		l.Address != nil && l.Address.HasRequiredFields()
}
