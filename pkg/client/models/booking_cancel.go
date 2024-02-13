package models

import (
	"encoding/xml"
)

type BookingCancelRequest struct {
	XMLName xml.Name `xml:"Main"`
	//The GOBookingCode
	GoBookingCode string `xml:"GoBookingCode"`
}

type BookingCancelRoot struct {
	XMLName xml.Name                  `xml:"Root"`
	Header  Header                    `xml:"Header"`
	Main    BookingCancelMainResponse `xml:"Main"`
}

func (r BookingCancelRoot) CheckError() error {
	if r.Header.OperationType == OperationTypeError || r.Header.OperationType == OperationTypeMessage {
		return r.Main.ErrorResponse.Error
	}

	return nil
}

func (r BookingCancelRoot) GetResponse() BookingCancelResponse {
	return r.Main.BookingCancelResponse
}

type BookingCancelMainResponse struct {
	XMLName xml.Name `xml:"Main"`
	BookingCancelResponse
	ErrorResponse
}

type BookingCancelResponse struct {
	//The Go booking code
	GoBookingCode string `xml:"GoBookingCode"`
	//The status of the booking RX, X, C etc.
	BookingStatus string `xml:"BookingStatus"`
}
