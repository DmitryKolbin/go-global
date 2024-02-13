package models

import (
	"encoding/xml"
)

type BookingStatusRequest struct {
	XMLName xml.Name `xml:"Main"`
	//The GOBookingCode or GORef
	GoBookingCode string `xml:"GoBookingCode"`
}

type BookingStatusRoot struct {
	XMLName xml.Name                  `xml:"Root"`
	Header  Header                    `xml:"Header"`
	Main    BookingStatusMainResponse `xml:"Main"`
}

func (r BookingStatusRoot) CheckError() error {
	if r.Header.OperationType == OperationTypeError || r.Header.OperationType == OperationTypeMessage {
		return r.Main.ErrorResponse.Error
	}

	return nil
}

func (r BookingStatusRoot) GetResponse() BookingStatusResponse {
	return r.Main.BookingStatusResponse
}

type BookingStatusMainResponse struct {
	XMLName xml.Name `xml:"Main"`

	BookingStatusResponse
	ErrorResponse
}

type BookingStatusResponse struct {
	GoBookingCode GoBookingCode `xml:"GoBookingCode"`
}

type GoBookingCode struct {
	XMLName xml.Name `xml:"GoBookingCode"`
	//	The status of the booking RQ, X, C etc.	RQ
	Status string `xml:"Status,attr"`
	//	The Go Reference
	GoReference string `xml:"GoReference,attr"`
	//	The total booking price
	TotalPrice float64 `xml:"TotalPrice,attr"`
	//	Currency
	Currency string `xml:"Currency,attr"`
	//The Go booking code
	Code string `xml:",chardata"`
}
