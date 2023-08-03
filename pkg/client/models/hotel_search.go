package models

import "encoding/xml"

type HotelSearchRequest struct {
	XMLName xml.Name `xml:"Main"`
	//Attribute to define request version	2.3
	Version string `xml:"Version,attr"`
	//Attribute to return response with JSON
	ResponseFormat string `xml:"ResponseFormat,attr,omitempty"`
	//Attribute to request geo coordinates - default true
	IncludeGeo bool `xml:"IncludeGeo,attr,omitempty"`
	//Attribute to request Max Hotels - default all
	MaxHotels int64 `xml:"MaxHotels,attr,omitempty"`
	//Attribute to request Max offers per hotel - default all
	MaxOffers int64 `xml:"MaxOffers,attr,omitempty"`
	//Attribute to request currency (eg. EUR, GBP, USD) only if set in profile - default per profile/destination
	Currency string `xml:"Currency,attr,omitempty"`
	//Attribute to request commission in response - default false
	IncludeCommission bool `xml:"IncludeCommission,attr,omitempty"`
	//Attribute to request general hotel facilities in response - default false
	HotelFacilities string `xml:"HotelFacilities,attr,omitempty"`
	//Attribute to request general room facilities in response - default false
	RoomFacilities string `xml:"RoomFacilities,attr,omitempty"`
	//Sorts results of hotels search
	SortOrder string `xml:"SortOrder"`
	//Minimum price for hotels to include in the results
	FilterPriceMin *float64 `xml:"FilterPriceMin,omitempty"`
	//Maximum price for hotels to include in the results
	FilterPriceMax *float64 `xml:"FilterPriceMax,omitempty"`
	//Maximum wait time (in seconds) for response (using this may return less results than we have available on our system)
	MaximumWaitTime int64 `xml:"MaximumWaitTime,omitempty"`
	//How many search results to include (number) in response after sorting - def 1000 when not set	1000
	MaxResponses int64 `xml:"MaxResponses,omitempty"`
	//Group to filter only a specific room basis
	FilterRoomBasises SearchFilterRoomBasises `xml:"FilterRoomBasises,omitempty"`
	//Passport Nationality ISO Code of the Lead Pax	GB
	Nationality string `xml:"Nationality"`
	//Code of the city
	CityCode []int64 `xml:"CityCode,omitempty"`
	//Group to search for hotel by ID
	Hotels SearchHotels `xml:"Hotels,omitempty"`
	//Check-in Date (yyyy-MM-dd)
	ArrivalDate string `xml:"ArrivalDate"`
	//Number of nights
	Nights int64 `xml:"Nights"`
	//Single Stars Code/ID to filter results - from the list here Or Group for Stars Range Filter
	Stars SearchStars `xml:"Stars,omitempty"`
	//Search only for apartment hotels
	Apartments bool `xml:"Apartments,omitempty"`
	//Group to specify room search criteria
	Rooms SearchRooms `xml:"Rooms"`
}

type SearchFilterRoomBasises struct {
	//Code of Room basis - from RoomBasisList
	FilterRoomBasis []string `xml:"FilterRoomBasis"`
}

type SearchHotels struct {
	//Hotel ID - you can get it from the static data. Up to 500 Ids (can be multi-city)
	HotelId []int64 `xml:"HotelId"`
}

type SearchStars struct {
	//Attribute for the minimum Stars Code/Id to filter
	MinStar string `xml:"MinStar,attr,omitempty"`
	//Attribute for the maximum Stars Code/Id to filter
	MaxStar string `xml:"MaxStar,attr,omitempty"`
}

type SearchRooms struct {
	//Group to specify room search criteria - can be more then one
	Room []SearchRoom `xml:"Room"`
}

type SearchRoom struct {
	// Attribute for number of adults in a room
	Adults int64 `xml:"Adults,attr"`
	// Attribute for number of rooms for a given criteria
	RoomCount int64 `xml:"RoomCount,attr"`
	// Attribute for number of children in the room. Including Babies/Infants.
	ChildCount int64 `xml:"ChildCount,attr"`
	// Attribute for number of cots in the room. Only 1 cot is allowed per room
	CotCount int64 `xml:"CotCount,attr,omitempty"`
	// If specified - age of child in this room.
	ChildAge []int64 `xml:"ChildAge,omitempty"`
}

// region response
type HotelSearchResponse struct {
	Header Header                    `json:"Header"`
	Hotels []HotelSearchResponseItem `json:"Hotels"`
	//Fill only if error occurred
	Main ErrorResponse `json:"Main,omitempty"`
}

type HotelSearchResponseItem struct {
	//Name of Hotel
	HotelName string `json:"HotelName"`
	//Unique HotelID code for the hote
	HotelCode int `json:"HotelCode"`
	//Country Id of hotel
	CountryId int `json:"CountryId"`
	//City Id of hotel
	CityId int `json:"CityId"`
	//Text Location of the Hotel - City Centre, Airport, etc.
	Location string `json:"Location"`
	//Location Code
	LocationCode string `json:"LocationCode"`
	//Thumbnail images of hotel
	Thumbnail string `json:"Thumbnail"`
	//Longitude
	Longitude float64 `json:"Longitude"`
	//Latitude
	Latitude float64 `json:"Latitude"`
	//Hotel Rank in Destination
	BestSellerRank string `json:"BestSellerRank"`
	//Large High Quality Thumbnail
	HotelImage string `json:"HotelImage"`
	//array of HotelFacilities
	HotelFacilities []string `json:"HotelFacilities"`
	//array of RoomFacilities
	RoomFacilities []string           `json:"RoomFacilities"`
	Offers         []HotelSearchOffer `json:"Offers"`
}

type HotelSearchOffer struct {
	//Unique Code session code - used for subsequent requests
	HotelSearchCode string `json:"HotelSearchCode"`
	//Cancellation Deadline
	CxlDeadline string `json:"CxlDeadline"`
	//Indication of Refundability
	NonRef bool `json:"NonRef"`
	//array of roomNames
	Rooms []string `json:"Rooms"`
	//BoardBasis - BB, RO
	RoomBasis string `json:"RoomBasis"`
	//1- Hotel is Available , 0 - NotAvailable
	Availability int `json:"Availability"`
	//Total Price
	TotalPrice float64 `json:"TotalPrice"`
	//ISO Currency code
	Currency string `json:"Currency"`
	//The Comm flat value
	CommPercent *float64 `json:"CommPercent"`
	//The Comm % value
	CommValue *float64 `json:"CommValue"`
	//Star Rating of the Hotel
	Category string `json:"Category"`
	//Free text remark
	Remark string `json:"Remark"`
	//Special remarks
	Special string `json:"Special"`
	//Is Best buy
	Preferred            bool                 `json:"Preferred"`
	CancellationPolicies []CancellationPolicy `json:"CancellationPolicies"`
}

type CancellationPolicy struct {
	//Policy Index starting at 1
	Id int64 `json:"Id"`
	//Date when policy takes affect	(dd/mm/yyyy)
	Starting string `json:"Starting"`
	//How to Apply the penalty
	BasedOn string `json:"BasedOn"`
	//Is value %(PCT|FLAT)
	Mode string `json:"Mode"`
	//Penalty Value to apply
	Value string `json:"Value"`
}

//endregion
