package models

import (
	"encoding/xml"
	"fmt"
)

type AdvBookingSearchRequest struct {
	XMLName xml.Name `xml:"Main"`
	//Attribute to define request version
	Version string `xml:"Version,attr"`
	//Attribute to request payments in response - default false
	IncludePayments bool `xml:"IncludePayments,attr,omitempty"`
	//Search also hierarchy agencies
	IncludeSubAgencies bool `xml:"IncludeSubAgencies,omitempty"`
	//Return full (default) or short (no rooms/pax) details
	DetailLevel string `xml:"DetailLevel,omitempty"`
	//Any passenger name (or part of)
	PaxName string `xml:"PaxName,omitempty"`
	//Code of the city
	CityCode int64 `xml:"CityCode,omitempty"`
	//Earliest arrival date (YYYY-mm-dd)
	ArrivalDateRangeFrom string `xml:"ArrivalDateRangeFrom,omitempty"`
	//Latest arrival date (YYYY-mm-dd)
	ArrivalDateRangeTo string `xml:"ArrivalDateRangeTo,omitempty"`
	//Match exact date (deprecated)	(YYYY-mm-dd)
	ArrivalDate string `xml:"ArrivalDate,omitempty"`
	//Single Date of booking creation (YYYY-mm-dd)
	CreatedDate string `xml:"CreatedDate,omitempty"`
	//Earliest Date of booking creation	(YYYY-mm-dd)
	CreatedDateRangeFrom string `xml:"CreatedDateRangeFrom,omitempty"`
	//Latest Date of booking creation (YYYY-mm-dd)
	CreatedDateRangeTo string `xml:"CreatedDateRangeTo,omitempty"`
	//Agent Ref as provided in the booking request
	ClientBookingCode string `xml:"ClientBookingCode,omitempty"`
	//No. of nights
	Nights int64 `xml:"Nights,omitempty"`
	//The Hotel Search Code
	HotelSearchCode string `xml:"HotelSearchCode,omitempty"`
	//Hotel name
	HotelName string `xml:"HotelName,omitempty"`
}

type AdvBookingSearchRoot struct {
	XMLName xml.Name                     `xml:"Root"`
	Header  Header                       `xml:"Header"`
	Main    AdvBookingSearchMainResponse `xml:"Main"`
}

func (r AdvBookingSearchRoot) CheckError() error {
	if r.Header.OperationType == OperationTypeError || r.Header.OperationType == OperationTypeMessage {
		return fmt.Errorf("code: %d, message: %s", r.Main.Error.Code, r.Main.Error.Message)
	}

	return nil
}

func (r AdvBookingSearchRoot) GetResponse() AdvBookingSearchResponse {
	return r.Main.Bookings
}

type AdvBookingSearchMainResponse struct {
	XMLName  xml.Name                 `xml:"Main"`
	Bookings AdvBookingSearchResponse `xml:"Bookings"`

	ErrorResponse
}

type AdvBookingSearchResponse struct {
	XMLName xml.Name                  `xml:"Bookings"`
	Booking []AdvBookingSearchBooking `xml:"Booking"`
}

type AdvBookingSearchBooking struct {
	XMLName xml.Name `xml:"Booking"`
	//The Go booking code
	GoBookingCode string `xml:"GoBookingCode"`
	//The Go Reference
	GoReference string `xml:"GoReference"`
	//The client booking code
	ClientBookingCode string `xml:"ClientBookingCode"`
	//Date of creation (yyyy-MM-dd HH:mm)
	CreatedDate string `xml:"CreatedDate"`
	//Id of Agency - usually same as credentials
	AgencyID string `xml:"AgencyID"`
	//Name of agency - usually belonging to the credentials
	AgencyName string `xml:"AgencyName"`
	//The status of the booking RQ, X, C etc.
	BookingStatus string `xml:"BookingStatus"`
	//The total client price
	TotalPrice float64 `xml:"TotalPrice"`
	//Currency
	Currency string `xml:"Currency"`
	//The total agent price
	GrossPrice GrossPrice `xml:"GrossPrice,omitempty"`
	//Name of the Hotel
	HotelName string `xml:"HotelName,omitempty"`
	//Hotel city code
	CityCode string `xml:"CityCode,omitempty"`
	//Hotel search code
	HotelSearchCode string `xml:"HotelSearchCode,omitempty"`
	//DEPRECATED
	RoomType string `xml:"RoomType,omitempty"`
	//BoardBasis
	RoomBasis string `xml:"RoomBasis,omitempty"`
	//Check-in Date (yyyy-MM-dd)
	ArrivalDate string `xml:"ArrivalDate,omitempty"`
	//Transfer Country
	Country string `xml:"Country,omitempty"`
	//Description of Transfer
	TransferName string `xml:"TransferName,omitempty"`
	//Location of Pickup
	PickupLocation string `xml:"PickupLocation,omitempty"`
	//Location of Dropoff
	DropOffLocation string `xml:"DropOffLocation,omitempty"`
	//Pickup Date and time (yyyy-MM-dd HH:mm)
	PickupDate string `xml:"PickupDate,omitempty"`
	//Cancellaton deadline date
	CancellationDeadline string `xml:"CancellationDeadline"`
	//Number of nights
	Nights int64 `xml:"Nights,omitempty"`
	//No Alternative Hotel Returned
	NoAlternativeHotel int64 `xml:"NoAlternativeHotel,omitempty"`
	//The LeadPax of the booking
	Leader Leader `xml:"Leader"`
	//The Pax Nationality as provided on avail step
	Nationality string `xml:"Nationality,omitempty"`
	//Rooms
	Rooms BookingSearchRoomsResponse `xml:"Rooms"`
	//PaymentTransactions
	PaymentTransactions PaymentTransactions `xml:"PaymentTransactions"`
	//Additional Preferences Requested
	Preferences Preferences `xml:"Preferences"`
	// Vehicle
	Vehicle Vehicle `xml:"Vehicle"`
	//Free Text remarks to pass in the booking, like special requests	Req. Wine in Room
	Remark string `xml:"Remark"`
}
