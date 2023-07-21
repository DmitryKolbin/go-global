package models

import (
	"encoding/xml"
	"fmt"
)

type BookingSearchRequest struct {
	XMLName xml.Name `xml:"Main"`
	//Attribute to request payments in response - default false
	IncludePayments bool `xml:"IncludePayments,attr,omitempty"`
	//Attribute to request commission in response - default false
	IncludeCommission bool `xml:"IncludeCommission,attr,omitempty"`
	//The GOBookingCode or GORef
	GoBookingCode string `xml:"GoBookingCode"`
}

type BookingSearchRoot struct {
	XMLName xml.Name                  `xml:"Root"`
	Header  Header                    `xml:"Header"`
	Main    BookingSearchMainResponse `xml:"Main"`
}

func (r BookingSearchRoot) CheckError() error {
	if r.Header.OperationType == OperationTypeError || r.Header.OperationType == OperationTypeMessage {
		return fmt.Errorf("code: %d, message: %s", r.Main.Error.Code, r.Main.Error.Message)
	}

	return nil
}

func (r BookingSearchRoot) GetResponse() BookingSearchResponse {
	return r.Main.BookingSearchResponse
}

type BookingSearchMainResponse struct {
	XMLName xml.Name `xml:"Main"`

	BookingSearchResponse
	ErrorResponse
}

type BookingSearchResponse struct {
	//The Go booking code
	GoBookingCode string `xml:"GoBookingCode"`
	//The Go Reference
	GoReference string `xml:"GoReference"`
	//The client booking code
	ClientBookingCode string `xml:"ClientBookingCode"`
	//The status of the booking RQ, X, C etc.
	BookingStatus string `xml:"BookingStatus"`
	//The total client price
	TotalPrice float64 `xml:"TotalPrice"`
	//Currency
	Currency string `xml:"Currency"`
	//The total agent price
	GrossPrice GrossPrice `xml:"GrossPrice,omitempty"`
	//The Comm flat value - with IncludeCommission
	Commission Commission `xml:"Commission"`
	//Name of the Hotel
	HotelName string `xml:"HotelName,omitempty"`
	//Hotel search code
	HotelSearchCode string `xml:"HotelSearchCode,omitempty"`
	//Hotel city code
	CityCode string `xml:"CityCode,omitempty"`
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

type GrossPrice struct {
	XMLName  xml.Name `xml:"GrossPrice"`
	Currency string   `xml:"Currency,attr"`
	Value    float64  `xml:",chardata"`
}

type BookingSearchRoomsResponse struct {
	XMLName  xml.Name                        `xml:"Rooms"`
	RoomType []BookingSearchRoomTypeResponse `xml:"RoomType"`
}

type BookingSearchRoomTypeResponse struct {
	XMLName xml.Name `xml:"RoomType"`
	//Attribute - Number of adults in a room
	Adults int64 `xml:"Adults,attr"`
	//Group to specify pax per room in group - can be more then one
	Room []BookingSearchRoomResponse `xml:"Room"`
}

type BookingSearchRoomResponse struct {
	XMLName xml.Name `xml:"Room"`
	//Attribute - A unique room ID for each type - incremental
	RoomId int64 `xml:"RoomId,attr"`
	//Attribute - booked room description
	Category string `xml:"Category,attr"`
	//Attribute - The number of cots for the given room type
	Cots       int64      `xml:"Cots,attr,omitempty"`
	PersonName PersonName `xml:"PersonName"`
	ExtraBed   ExtraBed   `xml:"ExtraBed"`
}

type Vehicle struct {
	XMLName xml.Name `xml:"Vehicle"`
	//Code of Transfer Vehicle - may be empty
	VehicleCode string `xml:"VehicleCode,omitempty"`
	//Type of Vehicle
	VehicleName string `xml:"VehicleName,omitempty"`
	//Max # of pax
	MaximumPassengers int64 `xml:"MaximumPassengers,omitempty"`
	//Actual # of pax
	NumberOfPassengers int64 `xml:"NumberOfPassengers"`
}
