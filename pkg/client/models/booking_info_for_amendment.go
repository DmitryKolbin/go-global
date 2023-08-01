package models

import (
	"encoding/xml"
	"fmt"
)

type BookingInfoForAmendmentRequest struct {
	XMLName xml.Name `xml:"Main"`
	//The Reservation ref#/code
	GoBookingCode string `xml:"GoBookingCode"`
}

type BookingInfoForAmendmentRoot struct {
	XMLName xml.Name                            `xml:"Root"`
	Header  Header                              `xml:"Header"`
	Main    BookingInfoForAmendmentMainResponse `xml:"Main"`
}

func (r BookingInfoForAmendmentRoot) CheckError() error {
	if r.Header.OperationType == OperationTypeError || r.Header.OperationType == OperationTypeMessage {
		return fmt.Errorf("code: %d, message: %s", r.Main.Error.Code, r.Main.Error.Message)
	}

	return nil
}

func (r BookingInfoForAmendmentRoot) GetResponse() BookingInfoForAmendmentResponse {
	return r.Main.BookingInfoForAmendmentResponse
}

type BookingInfoForAmendmentMainResponse struct {
	XMLName xml.Name `xml:"Main"`
	BookingInfoForAmendmentResponse
	ErrorResponse
}

type BookingInfoForAmendmentResponse struct {
	//Check In Date	2013-10-08
	ArrivalDate string `xml:"ArrivalDate"`
	//Number of nights
	Nights int64 `xml:"Nights"`
	//Room List
	Rooms BookingInfoForAmendmenRoomsResponse `xml:"Rooms"`
	//Remark List
	Remarks BookingInfoForAmendmenRemarks `xml:"Remarks"`
}

type BookingInfoForAmendmenRoomsResponse struct {
	XMLName  xml.Name                                 `xml:"Rooms"`
	RoomType []BookingInfoForAmendmenRoomTypeResponse `xml:"RoomType"`
}

type BookingInfoForAmendmenRoomTypeResponse struct {
	XMLName xml.Name `xml:"RoomType"`
	//Attribute - Number of adults in a room
	Adults int64 `xml:"Adults,attr"`
	//Group to specify pax per room in group - can be more then one
	Room []BookingInfoForAmendmenRoomResponse `xml:"Room"`
}

type BookingInfoForAmendmenRoomResponse struct {
	XMLName xml.Name `xml:"Room"`
	//Attribute - A unique room ID for each type - incremental
	RoomId int64 `xml:"RoomId,attr"`
	//Attribute - booked room description
	Category string `xml:"Category,attr"`
	//Attribute - The number of cots for the given room type
	Cots     int64                        `xml:"Cots,attr,omitempty"`
	Person   BookingInfoForAmendmenPerson `xml:"Person"`
	ExtraBed ExtraBed                     `xml:"ExtraBed"`
}

type BookingInfoForAmendmenPerson struct {
	XMLName xml.Name `xml:"PersonName"`
	//Attribute - A unique Person ID for the booking - incremental
	PersonID int64 `xml:"PersonID,attr"`
	//Attribute - Pax Title - Version 2+ only
	Title string `xml:"Title,attr"`
	//Attribute - Adult First Name - Version 2+ only
	FirstName string `xml:"FirstName,attr"`
	//Attribute - Adult Last Name - Version 2+ only
	LastName string `xml:"LastName,attr"`
	//Pax Age (Required for child)
	Age int64 `xml:"Age,omitempty"`
}

type BookingInfoForAmendmenRemarks struct {
	XMLName xml.Name        `xml:"Remarks"`
	Remark  []BookingRemark `xml:"Remark"`
}

type BookingInfoForAmendmenRemark struct {
	XMLName xml.Name `xml:"Remark"`
	//Id of the remark
	Id int64 `xml:"Id,attr"`
	//Remark text
	Text []BookingRemark `xml:",chardata"`
}
