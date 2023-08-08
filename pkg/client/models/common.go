package models

import (
	"encoding/json"
	"encoding/xml"
)

const (
	OperationTypeRequest  = "Request"
	OperationTypeResponse = "Response"
	OperationTypeError    = "Error"
	OperationTypeMessage  = "Message"

	LanguageUs = "us"
	LanguageEs = "es"
	LanguageHr = "hr"
	LanguageSk = "sk"
	LanguageBg = "bg"
	LanguagePl = "pl"
	LanguageRu = "ru"
	LanguageUa = "ua"
	LanguageFr = "fr"
	LanguageIt = "it"
	LanguageLv = "lv"
	LanguageRo = "ro"
	LanguageCz = "cz"
	LanguageHu = "hu"
	LanguageDe = "de"
	LanguageSi = "si"
	LanguageTr = "tr"
	LanguageRs = "rs"
	LanguageIl = "il"
	LanguagePt = "pt"
	LanguageZh = "zh"
	LanguageKr = "kr"
	LanguageBr = "br"

	//RoomBasisBb BED AND BREAKFAST
	RoomBasisBb = "BB"
	//RoomBasisCb CONTINENTAL BREAKFAST
	RoomBasisCb = "CB"
	//RoomBasisAi ALL-INCLUSIVE
	RoomBasisAi = "AI"
	//RoomBasisFb FULL-BOARD
	RoomBasisFb = "FB"
	//RoomBasisHb HALF-BOARD
	RoomBasisHb = "HB"
	//RoomBasisRo ROOM ONLY
	RoomBasisRo = "RO"
	//RoomBasisBd BED AND DINNER
	RoomBasisBd = "BD"

	//Start1 1 star
	Start1 = 1
	//Start1point5 1.5 stars
	Start1point5 = 2
	//Start2 2 stars
	Start2 = 3
	//Start2point5 2.5 stars
	Start2point5 = 4
	//Start3 3 stars
	Start3 = 5
	//Start3point5 3.5 stars
	Start3point5 = 6
	//Start4 4 stars
	Start4 = 7
	//Start4point5 4.5 star
	Start4point5 = 8
	//Start5 5 stars
	Start5 = 9
	//Start5point5 5.5 stars
	Start5point5 = 10
	//Start6 6 stars
	Start6 = 11

	//SortByPriceAsc first lowest prices
	SortByPriceAsc = "1"
	//SortByPriceDesc first high prices
	SortByPriceDesc = "2" //S
	//SortByCxlAsc first smallest dates
	SortByCxlAsc = "3"
	//SortByCxlDesc first latest dates
	SortByCxlDesc = "4"

	//StatusRequested Request for Booking was received - status not final, pending confirmation (C) or RJ
	StatusRequested = "RQ"
	//StatusConfirmed Booking is finalized and active
	StatusConfirmed = "C"
	//StatusReqCancellation Cancellation Request was received - status not final, expect either X, C, VI or XF
	StatusReqCancellation = "RX"
	//StatusCancelled Booking canceled FOC
	StatusCancelled = "X"
	//StatusRejected Booking was rejected
	StatusRejected = "RJ"
	//StatusVoucherIssued Issued	Booking is finalized, confirmed and the voucher was issued
	StatusVoucherIssued = "VCH"
	//StatusVoucherReq	Booking is finalized, confirmed and request for voucher was issued
	StatusVoucherReq = "VRQ"

	AmendmentCategoryStandard    = "STANDARD"
	AmendmentCategorySuperior    = "SUPERIOR"
	AmendmentCategoryDeluxe      = "DELUXE"
	AmendmentCategoryLuxury      = "LUXURY"
	AmendmentCategoryPremium     = "PREMIUM"
	AmendmentCategoryJuniorSuite = "JUNIOR SUITE"
	AmendmentCategorySuite       = "SUITE"
	AmendmentCategoryMiniSuite   = "MINI SUITE"
	AmendmentCategoryStudio      = "STUDIO"
	AmendmentCategoryExecutive   = "EXECUTIVE"

	CancellationPolicyModePercent = "PCT"
	CancellationPolicyModeFix     = "FLAT"
)

type ResponseRoot[T any] interface {
	CheckError() error
	GetResponse() T
}

type EnvelopeRequest struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	XsiType  string   `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr,omitempty"`
	XmlnsXsi string   `xml:"xmlns xsi,attr,omitempty"`

	Body Body
}

type Body struct {
	XMLName     xml.Name `xml:"Body"`
	MakeRequest MakeRequest
}

type MakeRequest struct {
	XMLName     xml.Name `xml:"http://www.goglobal.travel/ MakeRequest"`
	RequestType int64    `xml:"requestType"`
	XmlRequest  XmlRequestCdata
}

type XmlRequestCdata struct {
	XMLName xml.Name `xml:"xmlRequest"`
	Text    string   `xml:",cdata"`
}

type RequestRoot struct {
	XMLName xml.Name `xml:"Root"`

	Header Header
	Main   []byte `xml:",innerxml"`
}

type Header struct {
	XMLName xml.Name `xml:"Header"`
	//Agency/Client Identification
	Agency json.Number `xml:"Agency" json:"Agency"`
	//UserName
	User string `xml:"User" json:"User"`
	//User password
	Password string `xml:"Password" json:"Password"`
	//Operation Name (Header) based of the RequestType
	Operation string `xml:"Operation" json:"Operation"`
	//Type of operation
	OperationType string `xml:"OperationType" json:"OperationType"`
	// In response - search stats
	Stats SearchStats `xml:"-" json:"Stats"`
}

type SearchStats struct {
	//Number of hotels returned
	HotelQty int `json:"HotelQty"`
	//Number of results returned
	ResultsQty int `json:"ResultsQty"`
}

type EnvelopeResponse struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	XsiType  string   `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr,omitempty"`
	XmlnsXsi string   `xml:"xmlns xsi,attr,omitempty"`

	Body BodyResponse
}

type BodyResponse struct {
	XMLName             xml.Name `xml:"Body"`
	MakeRequestResponse MakeRequestResponse
}

type MakeRequestResponse struct {
	XMLName           xml.Name `xml:"http://www.goglobal.travel/ MakeRequestResponse"`
	MakeRequestResult MakeRequestResult
}
type MakeRequestResult struct {
	XMLName xml.Name `xml:"MakeRequestResult"`
	Data    []byte   `xml:",chardata"`
}

type ErrorResponse struct {
	Error      GoGlobalError      `json:"Error"`
	DebugError GoGlobalDebugError `json:"DebugError"`
}

type GoGlobalError struct {
	XMLName xml.Name `xml:"Error" json:"-"`
	Code    int64    `xml:"code,attr" json:"Code"`
	Message string   `xml:",cdata" json:"Message"`
}
type GoGlobalDebugError struct {
	XMLName     xml.Name `xml:"DebugError" json:"-"`
	Incident    int64    `xml:"incident,attr" json:"Incident"`
	TimeStamp   string   `xml:"timestamp,attr" json:"TimeStamp"`
	Message     string   `xml:",cdata" json:"Message"`
	FinalAction string   `json:"finalAction"`
}
