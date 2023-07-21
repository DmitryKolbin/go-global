package models

import (
	"encoding/xml"
	"fmt"
)

type PriceBreakdownRequest struct {
	XMLName xml.Name `xml:"Main"`
	//The Hotel Search Code
	HotelSearchCode string `xml:"HotelSearchCode"`
}

type PriceBreakdownRoot struct {
	XMLName xml.Name                   `xml:"Root"`
	Header  Header                     `xml:"Header"`
	Main    PriceBreakdownMainResponse `xml:"Main"`
}

func (r PriceBreakdownRoot) CheckError() error {
	if r.Header.OperationType == OperationTypeError || r.Header.OperationType == OperationTypeMessage {
		return fmt.Errorf("code: %d, message: %s", r.Main.Error.Code, r.Main.Error.Message)
	}

	return nil
}

func (r PriceBreakdownRoot) GetResponse() PriceBreakdownResponse {
	return r.Main.PriceBreakdownResponse
}

type PriceBreakdownMainResponse struct {
	XMLName xml.Name `xml:"Main"`

	PriceBreakdownResponse
	ErrorResponse
}

type PriceBreakdownResponse struct {
	//The selected hotel
	HotelName string             `xml:"HotelName"`
	Room      PriceBreakdownRoom `xml:"Room"`
}

type PriceBreakdownRoom struct {
	XMLName xml.Name `xml:"Room"`
	//Search Type of this room Twin, Room for 2 Adults, etc
	RoomType string `xml:"RoomType"`
	//Number of children
	Children int64 `xml:"Children"`
	//Includes an Infant	0 or 1
	Cots           int64          `xml:"Cots"`
	PriceBreakdown PriceBreakdown `xml:"PriceBreakdown"`
}

type PriceBreakdown struct {
	XMLName xml.Name `xml:"PriceBreakdown"`
	//Break-down starting date (YYYY-mm-dd)
	FromDate string `xml:"FromDate"`
	//Break-down ending date(YYYY-mm-dd)
	ToDate string `xml:"ToDate"`
	//Price per 1 night	234.5
	Price float64 `xml:"Price"`
	//
	Currency string
}
