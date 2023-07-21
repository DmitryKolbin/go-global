package client

const ResponseFormatJson = "JSON"
const SortOrderDefault = "Default"

type Destination struct {
	CityId    int64  `csv:"CityId"`
	City      string `csv:"City"`
	CountryId int64  `csv:"CountryId"`
	Country   string `csv:"Country"`
	IsoCode   string `csv:"IsoCode"`
}

type Hotel struct {
	CountryId   int64   `csv:"CountryId"`
	Country     string  `csv:"Country"`
	IsoCode     string  `csv:"IsoCode"`
	CityId      int64   `csv:"CityId"`
	City        string  `csv:"City"`
	HotelID     int64   `csv:"HotelID"`
	Name        string  `csv:"Name"`
	Address     string  `csv:"Address"`
	Phone       string  `csv:"Phone"`
	Fax         string  `csv:"Fax"`
	Stars       string  `csv:"Stars"`
	StarsID     string  `csv:"StarsID"`
	Longitude   float64 `csv:"Longitude"`
	Latitude    float64 `csv:"Latitude"`
	IsApartment bool    `csv:"IsApartment"`
	GiataCode   string  `csv:"GiataCode"`
}
