package models

import (
	"encoding/xml"
	"fmt"
)

type BookValuationRequest struct {
	XMLName xml.Name `xml:"Main"`
	//Attribute to define request version 2.0
	Version string `xml:"Version,attr"`
	//The Hotel Search Code
	HotelSearchCode string `xml:"HotelSearchCode"`
	//Check-in Date (yyyy-MM-dd)
	ArrivalDate string `xml:"ArrivalDate"`
}

type BookValuationRoot struct {
	XMLName xml.Name                  `xml:"Root"`
	Header  Header                    `xml:"Header"`
	Main    BookValuationMainResponse `xml:"Main"`
}

func (r BookValuationRoot) CheckError() error {
	if r.Header.OperationType == OperationTypeError || r.Header.OperationType == OperationTypeMessage {
		return fmt.Errorf("code: %d, message: %s", r.Main.Error.Code, r.Main.Error.Message)
	}

	return nil
}

func (r BookValuationRoot) GetResponse() BookValuationResponse {
	return r.Main.BookValuationResponse
}

type BookValuationMainResponse struct {
	XMLName xml.Name `xml:"Main"`
	BookValuationResponse

	ErrorResponse
}

type BookValuationResponse struct {
	HotelSearchCode      string            `xml:"HotelSearchCode"`
	ArrivalDate          string            `xml:"ArrivalDate"`
	CancellationDeadline string            `xml:"CancellationDeadline"`
	Remarks              string            `xml:"Remarks"`
	Rates                BookValuationRate `xml:"Rates"`
}

type BookValuationRate struct {
	XMLName  xml.Name `xml:"Rates"`
	Currency string   `xml:"Currency,attr"`
	Value    float64  `xml:",chardata"`
}
