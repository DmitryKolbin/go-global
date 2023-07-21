package models

import (
	"encoding/xml"
	"fmt"
)

const (
	BookingRemarksAgent  = "Agent"
	BookingRemarksTariff = "Tariff"
)

type VoucherDetailsRequest struct {
	XMLName xml.Name `xml:"Main"`
	//The GOBookingCode
	GoBookingCode string `xml:"GoBookingCode"`
	//Return emergency contact phone
	GetEmergencyPhone bool `xml:"GetEmergencyPhone"`
}

type VoucherDetailsRoot struct {
	XMLName xml.Name                   `xml:"Root"`
	Header  Header                     `xml:"Header"`
	Main    VoucherDetailsMainResponse `xml:"Main"`
}

func (r VoucherDetailsRoot) CheckError() error {
	if r.Header.OperationType == OperationTypeError || r.Header.OperationType == OperationTypeMessage {
		return fmt.Errorf("code: %d, message: %s", r.Main.Error.Code, r.Main.Error.Message)
	}

	return nil
}

func (r VoucherDetailsRoot) GetResponse() VoucherDetailsResponse {
	return r.Main.VoucherDetailsResponse
}

type VoucherDetailsMainResponse struct {
	XMLName xml.Name `xml:"Main"`

	VoucherDetailsResponse
	ErrorResponse
}

type VoucherDetailsResponse struct {
	//The Go booking code	162345
	GoBookingCode string `xml:"GoBookingCode"`
	//Hotel name
	HotelName string `xml:"HotelName"`
	//Hotel address
	Address string `xml:"Address"`
	//Hotel phome
	Phone string `xml:"Phone"`
	//Hotel fax
	Fax string `xml:"Fax"`
	//checkin date (03/Mar/11)
	CheckInDate string `xml:"CheckInDate"`
	//Room bassis
	RoomBasis string `xml:"RoomBasis"`
	//Number of nights
	Nights int64 `xml:"Nights"`
	//Room and Pax details
	Rooms string `xml:"Rooms"`
	//Remark
	Remarks string `xml:"Remarks"`
	//Url for a PDF version of the vouche
	VoucherDownloadURL string `xml:"VoucherDownloadURL"`
	//Group of booking remarks per type - Version 2+ only
	BookingRemarks BookingRemarks `xml:"BookingRemarks,omitempty"`
	//Hotel supplier
	BookedAndPayableBy string `xml:"BookedAndPayableBy"`
	//Supplier booking code
	SupplierReferenceNumber string `xml:"SupplierReferenceNumber"`
	//Contact phone in case of Emergency
	EmergencyPhone string `xml:"EmergencyPhone"`
}

type BookingRemarks struct {
	XMLName xml.Name `xml:"BookingRemarks"`
	//Attribute - Type of remark group - Version 2+ only
	Type string `xml:"Type,attr"`
	//Remarks data - list - Version 2+ only
	Remark []BookingRemark `xml:"Remark"`
}

type BookingRemark struct {
	XMLName xml.Name `xml:"Remark"`
	Value   string   `xml:",cdata"`
}
