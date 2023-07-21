package models

import (
	"encoding/xml"
	"fmt"
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
		return fmt.Errorf("code: %d, message: %s", r.Main.Error.Code, r.Main.Error.Message)
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
