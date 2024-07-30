package models

import (
	"encoding/xml"
)

type BookValuationRequest struct {
	XMLName xml.Name `xml:"Main"`
	//Attribute to define request version 2.0
	Version string `xml:"Version,attr"`
	//The Hotel Search Code
	HotelSearchCode string `xml:"HotelSearchCode"`
	//Check-in Date (yyyy-MM-dd)
	ArrivalDate string `xml:"ArrivalDate"`
	//Attribute to request TotalTax and RoomRate when available - default false
	ReturnTaxData bool `xml:"ReturnTaxData,attr,omitempty"`
}

type BookValuationRoot struct {
	XMLName xml.Name                  `xml:"Root"`
	Header  Header                    `xml:"Header"`
	Main    BookValuationMainResponse `xml:"Main"`
}

func (r BookValuationRoot) CheckError() error {
	if r.Header.OperationType == OperationTypeError || r.Header.OperationType == OperationTypeMessage {
		return r.Main.ErrorResponse.Error
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

type CancellationPolicies struct {
	XMLName xml.Name             `xml:"CancellationPolicies"`
	Policy  []CancellationPolicy `xml:"Policy"`
}

type BookValuationResponse struct {
	HotelSearchCode      string               `xml:"HotelSearchCode"`
	ArrivalDate          string               `xml:"ArrivalDate"`
	CancellationDeadline string               `xml:"CancellationDeadline"`
	Remarks              string               `xml:"Remarks"`
	Rates                BookValuationRate    `xml:"Rates"`
	TotalTax             float64              `xml:"TotalTax"`
	RoomRate             float64              `xml:"RoomRate"`
	CancellationPolicies CancellationPolicies `xml:"CancellationPolicies"`
}

type BookValuationRate struct {
	XMLName       xml.Name `xml:"Rates"`
	Currency      string   `xml:"currency,attr"`
	CurrencyUpper string   `xml:"Currency,attr"` // in docs currency attribute is lower case. this is fallback for some cases
	Value         float64  `xml:",chardata"`
}
