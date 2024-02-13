package models

import (
	"encoding/xml"
)

// HotelInfoRequest receive addition information about hotel
// Only 1 of either InfoHotelId or HotelSearchCode must be specified in each request.
type HotelInfoRequest struct {
	XMLName xml.Name `xml:"Main"`
	//Attribute to define request version
	Version string `xml:"Version,attr"`
	//Unique Hotel Id
	InfoHotelId int64 `xml:"InfoHotelId,omitempty"`
	//Hotel search code
	HotelSearchCode string `xml:"HotelSearchCode,omitempty"`
	//Lower Case 2 Letter ISO Language Code
	InfoLanguage string `xml:"InfoLanguage,omitempty"`
}

type HotelInfoRoot struct {
	XMLName xml.Name              `xml:"Root"`
	Header  Header                `xml:"Header"`
	Main    HotelInfoMainResponse `xml:"Main"`
}

func (r HotelInfoRoot) CheckError() error {
	if r.Header.OperationType == OperationTypeError || r.Header.OperationType == OperationTypeMessage || r.Main.Error.Code > 0 {
		return r.Main.ErrorResponse.Error
	}

	return nil
}

func (r HotelInfoRoot) GetResponse() HotelInfoResponse {
	return r.Main.HotelInfoResponse
}

type HotelInfoMainResponse struct {
	XMLName xml.Name `xml:"Main"`

	HotelInfoResponse
	ErrorResponse
}

type HotelInfoResponse struct {
	//Hotel code	12345/3243212345/53
	HotelSearchCode string `xml:"HotelSearchCode"`
	//Name of the Hotel	AC
	HotelName string `xml:"HotelName"`
	//Hotel Code/Id (Version 2.2+)
	HotelId int64 `xml:"HotelId"`
	//Address
	Address string `xml:"Address"`
	//City code
	CityCode int64 `xml:"CityCode"`
	//Holds Geocodes. (in req 61 only)
	GeoCodes GeoCodes `xml:"GeoCodes"`
	//Phone
	Phone string `xml:"Phone"`
	//Fax
	Fax string `xml:"Fax"`
	//Category
	Category string `xml:"Category"`
	//Description
	Description string `xml:"Description"`
	//Hotel Facilities
	HotelFacilities string `xml:"HotelFacilities"`
	//Room Facilities
	RoomFacilities string `xml:"RoomFacilities"`
	//Number of rooms
	RoomCount int64    `xml:"RoomCount"`
	Pictures  Pictures `xml:"Pictures"`
}

type GeoCodes struct {
	XMLName xml.Name `xml:"GeoCodes"`
	//Longitude (in req 61 only)
	Longitude float64 `xml:"Longitude"`
	//Latitude (in req 61 only)
	Latitude float64 `xml:"Latitude"`
}

type Pictures struct {
	XMLName xml.Name  `xml:"Pictures"`
	Picture []Picture `xml:"Picture"`
}
type Picture struct {
	XMLName xml.Name `xml:"Picture"`
	//Image Description, When Exists
	Description string `xml:"Description,omitempty"`
	//Image URL in CDATA
	Value string `xml:",cdata"`
}
