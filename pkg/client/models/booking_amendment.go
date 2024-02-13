package models

import (
	"encoding/xml"
)

type BookingAmendmentRequest struct {
	XMLName xml.Name `xml:"Main"`
	//Check In Date	2013-10-08
	ArrivalDate string `xml:"ArrivalDate"`
	//Number of nights
	Nights int64 `xml:"Nights"`
	//Room List
	Rooms BookingInfoForAmendmenRoomsResponse `xml:"Rooms"`
	//Remark List
	Remarks BookingInfoForAmendmenRemarks `xml:"Remarks"`
}

type BookingAmendmentRoot struct {
	XMLName xml.Name                     `xml:"Root"`
	Header  Header                       `xml:"Header"`
	Main    BookingAmendmentMainResponse `xml:"Main"`
}

func (r BookingAmendmentRoot) CheckError() error {
	if r.Header.OperationType == OperationTypeError || r.Header.OperationType == OperationTypeMessage {
		return r.Main.ErrorResponse.Error
	}

	return nil
}

func (r BookingAmendmentRoot) GetResponse() BookingAmendmentResponse {
	return r.Main.BookingAmendmentResponse
}

// BookingAmendmentMainResponse
// If the Main tag is empty, the amendment request was RECEIVED correctly, else, an exception message will be returned.
type BookingAmendmentMainResponse struct {
	XMLName xml.Name `xml:"Main"`

	ErrorResponse
	BookingAmendmentResponse
}

type BookingAmendmentResponse struct {
}
