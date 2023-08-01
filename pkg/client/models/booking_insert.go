package models

import (
	"encoding/xml"
	"fmt"
)

type BookingInsertRequest struct {
	XMLName xml.Name `xml:"Main"`
	//Attribute to define request version	2.3
	Version string `xml:"Version,attr"`
	//Attribute to request payments in response - default false
	IncludePayments bool `xml:"IncludePayments,attr,omitempty"`
	//Attribute to request commission in response - default false
	IncludeCommission bool `xml:"IncludeCommission,attr,omitempty"`
	//Agent Reference
	AgentReference string `xml:"AgentReference"`
	//The Hotel search code of the chosen hotel
	HotelSearchCode string `xml:"HotelSearchCode"`
	//Check-in Date (yyyy-MM-dd)
	ArrivalDate string `xml:"ArrivalDate"`
	//Number of nights
	Nights int64 `xml:"Nights"`
	//Indicates whether client wants to be advised on alternates if his room is unavailable for booking (1=no alternatives)
	NoAlternativeHotel int64 `xml:"NoAlternativeHotel"`
	//The lead pax of the booking
	Leader Leader `xml:"Leader"`
	//Group to specify room search criteria
	Rooms RoomsRequest `xml:"Rooms"`
	//Additional Preferences Requested
	Preferences Preferences `xml:"Preferences,omitempty"`
	//Free Text remarks to pass in the booking, like special requests
	Remark string `xml:"Remark,omitempty"`
	//Base64 Encoded details for the credit card to use for payment}
	PaymentCreditCard string `xml:"PaymentCreditCard,omitempty"`
}

type Leader struct {
	XMLName        xml.Name `xml:"Leader"`
	LeaderPersonID int64    `xml:"LeaderPersonID,attr"`
}

type RoomsRequest struct {
	XMLName  xml.Name          `xml:"Rooms"`
	RoomType []RoomTypeRequest `xml:"RoomType"`
}

type RoomTypeRequest struct {
	XMLName xml.Name `xml:"RoomType"`
	//Attribute - Number of adults in a room
	Adults int64 `xml:"Adults,attr"`
	//Attribute - The number of cots for the given room type
	Cots int64 `xml:"Cots,attr,omitempty"`
	//Group to specify pax per room in group - can be more then one
	Room []RoomRequest `xml:"Room"`
}

type RoomRequest struct {
	XMLName xml.Name `xml:"Room"`
	//Attribute - A unique room ID for each type - incremental
	RoomId     int64        `xml:"RoomId,attr"`
	PersonName []PersonName `xml:"PersonName"`
	ExtraBed   []ExtraBed   `xml:"ExtraBed"`
}

type PersonName struct {
	XMLName xml.Name `xml:"PersonName"`
	//Attribute - A unique Person ID for the booking - incremental
	PersonID int64 `xml:"PersonID,attr"`
	//Attribute - Pax Title - Version 2+ only
	Title string `xml:"Title,attr"`
	//Attribute - Adult First Name - Version 2+ only
	FirstName string `xml:"FirstName,attr"`
	//Attribute - Adult Last Name - Version 2+ only
	LastName string `xml:"LastName,attr"`
}

type ExtraBed struct {
	XMLName xml.Name `xml:"ExtraBed"`
	//Attribute - A unique Person ID for the booking - incremental - with adults
	PersonID int64 `xml:"PersonID,attr"`
	//Attribute - Child First Name - Version 2+ only
	FirstName string `xml:"FirstName,attr"`
	//Attribute - Child Last Name - Version 2+ only
	LastName string `xml:"LastName,attr"`
	//Attribute - The age of the child. Version 2.2 and above, supporting child age 1 - 18. Before 2.2, ages 2 - 10.
	ChildAge int64 `xml:"ChildAge,attr"`
}

type Preferences struct {
	XMLName xml.Name `xml:"Preferences"`
	//Adjoining Rooms
	AdjoiningRooms int64 `xml:"AdjoiningRooms,omitempty"`
	//Connecting Rooms
	ConnectingRooms int64 `xml:"ConnectingRooms,omitempty"`
	//Non-Smoking Rooms
	NonSmokingRooms int64 `xml:"NonSmokingRooms,omitempty"`
	//Late Arrival Time (HH:mm)
	LateArrival string `xml:"LateArrival,omitempty"`
}

type BookingInsertRoot struct {
	XMLName xml.Name               `xml:"Root"`
	Header  Header                 `xml:"Header"`
	Main    BookInsertMainResponse `xml:"Main"`
}

func (r BookingInsertRoot) CheckError() error {
	if r.Header.OperationType == OperationTypeError || r.Header.OperationType == OperationTypeMessage {
		return fmt.Errorf("code: %d, message: %s", r.Main.Error.Code, r.Main.Error.Message)
	}

	return nil
}

func (r BookingInsertRoot) GetResponse() BookingInsertResponse {
	return r.Main.BookingInsertResponse
}

type BookInsertMainResponse struct {
	XMLName xml.Name `xml:"Main"`
	BookingInsertResponse
	ErrorResponse
}

type BookingInsertResponse struct {
	//The GoGlobal booking code
	GoBookingCode string `xml:"GoBookingCode"`
	//The GoGlobal Reference
	GoReference string `xml:"GoReference"`
	//The client booking code
	ClientBookingCode string `xml:"ClientBookingCode"`
	//The status of the booking RQ, X, C etc.
	BookingStatus string `xml:"BookingStatus"`
	//The total booking price
	TotalPrice float64 `xml:"TotalPrice"`
	//Currency
	Currency string `xml:"Currency"`
	//The Comm flat value - with IncludeCommission
	Commission Commission `xml:"Commission"`
	//Id of the Hotel - Version 2+ Only
	HotelId int64 `xml:"HotelId"`
	//Name of the Hotel
	HotelName string `xml:"HotelName"`
	//Hotel Search code
	HotelSearchCode string `xml:"HotelSearchCode"`
	//DEPRECATED Ignore it
	RoomType string `xml:"RoomType"`
	//RoomBasis
	RoomBasis string `xml:"RoomBasis"`
	//Check-in Date (yyyy-MM-dd)
	ArrivalDate string `xml:"ArrivalDate"`
	//CXL deadline date
	CancellationDeadline string `xml:"CancellationDeadline"`
	//Number of nights
	Nights int64 `xml:"Nights"`
	//Don't return alternative
	NoAlternativeHotel string `xml:"NoAlternativeHotel"`
	//The LeadPax of the booking
	Leader Leader `xml:"Leader"`
	//
	PaymentTransactions string `xml:"PaymentTransactions"`
	//Group to specify room search criteria
	Rooms RoomsResponse `xml:"Rooms"`
	//Additional Preferences Requested
	Preferences Preferences `xml:"Preferences"`
	//Free Text remarks to pass in the booking, like special requests	Req. Wine in Room	Y
	Remark string `xml:"Remark"`
	//Payment Details if opted to pay on booking
	PaymentInfo PaymentInfo `xml:"PaymentInfo"`
}

type Commission struct {
	XMLName xml.Name `xml:"Commission"`
	//Attribute - The Comm % value - with IncludeCommission
	Pct float64 `xml:"pct,attr"`
}

type BookInsertRate struct {
	XMLName  xml.Name `xml:"Rates"`
	Currency string   `xml:"Currency,attr"`
	Value    float64  `xml:",chardata"`
}

type RoomsResponse struct {
	XMLName  xml.Name           `xml:"Rooms"`
	RoomType []RoomTypeResponse `xml:"RoomType"`
}

type RoomTypeResponse struct {
	XMLName xml.Name `xml:"RoomType"`
	//Attribute - Number of adults in a room
	Adults int64 `xml:"Adults,attr"`
	//Attribute - The number of cots for the given room type
	Cots int64 `xml:"Cots,attr,omitempty"`
	//Group to specify pax per room in group - can be more then one
	Room []RoomResponse `xml:"Room"`
}

type RoomResponse struct {
	XMLName xml.Name `xml:"Room"`
	//Attribute - A unique room ID for each type - incremental
	RoomId int64 `xml:"RoomId,attr"`
	//Attribute - booked room description
	Category   string       `xml:"Category,attr"`
	PersonName []PersonName `xml:"PersonName"`
	ExtraBed   []ExtraBed   `xml:"ExtraBed"`
}

type PaymentTransactions struct {
	XMLName     xml.Name      `xml:"PaymentTransactions"`
	Transaction []Transaction `xml:"Transaction"`
}

type Transaction struct {
	XMLName xml.Name `xml:"Transaction"`
	//Attribute - Date of the transaction (format: 2022-11-01 00:00:00)
	Date string `xml:"Date,attr"`
	//Attribute - Type of the transaction: Payment|Refund
	Type string `xml:"Type,attr"`
	//Attribute - Transaction Category: Reservation|Refund
	Category string `xml:"Category,attr"`
	//Attribute - Form of payment: BANK_TRANSFER, CREDITCARD, etc.
	Method string `xml:"Method,attr"`
	//Attribute - Paid amount in local currency
	PaidAmount float64 `xml:"PaidAmount,attr"`
	//Attribute - Currency of the paid amount - ISO Code
	PaidCurrency string `xml:"PaidCurrency,attr"`
	//Attribute - Paid amount in Booking Currency
	BookingAmount float64 `xml:"BookingAmount"`
	//Attribute - Currency of the booking - ISO Code
	BookingCurrency string `xml:"BookingCurrency"`
	//Attribute - Local to Booking exchange rate at time of transaction
	ExchangeRate float64 `xml:"ExchangeRate"`
}

type PaymentInfo struct {
	XMLName xml.Name `xml:"PaymentInfo"`
	//Payment transaction details
	PaymentResult PaymentResult `xml:"PaymentResult"`
	//Refund transaction details
	RefundResult RefundResult `xml:"RefundResult"`
}

type PaymentResult struct {
	XMLName xml.Name `xml:"PaymentResult"`
	//Indication if the payment completed correctly
	Successful bool `xml:"Successful"`
	//Amount charged for the booking
	Amount float64 `xml:"Amount"`
	//Currency used to pay for the booking - booking currency
	Currency string `xml:"Currency"`
	//Confirmation code for the transaction
	ApprovalCode string `xml:"ApprovalCode"`
	//If payment was not successful, will hold the reason for failure
	ErrorMessage string `xml:"ErrorMessage"`
}

type RefundResult struct {
	XMLName xml.Name `xml:"RefundResult"`
	//Indication if the refund completed correctly
	Successful bool `xml:"Successful"`
	//Confirmation code for the transaction
	ApprovalCode string `xml:"ApprovalCode"`
	//If refund was not successful, will hold the reason for failure
	ErrorMessage string `xml:"ErrorMessage"`
}
