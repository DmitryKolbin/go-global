package client

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/DmitryKolbin/go-global/pkg/client/models"
	"github.com/dimchansky/utfbom"
	"github.com/gocarina/gocsv"
)

const (
	getDestinationsUrl = "https://static-data.tourismcloudservice.com/propsdata/Destinations/compress/true"
	getHotelsUrlFmt    = "https://static-data.tourismcloudservice.com/agency/hotels/%s"

	searchRequest           = goGlobalRequest("HOTEL_SEARCH_REQUEST")
	bookingValidation       = goGlobalRequest("BOOKING_VALUATION_REQUEST")
	bookingInsert           = goGlobalRequest("BOOKING_INSERT_REQUEST")
	bookingStatus           = goGlobalRequest("BOOKING_STATUS_REQUEST")
	bookingSearch           = goGlobalRequest("BOOKING_SEARCH_REQUEST")
	advBookingSearch        = goGlobalRequest("ADV_BOOKING_SEARCH_REQUEST")
	bookingCancel           = goGlobalRequest("BOOKING_CANCEL_REQUEST")
	voucherDetails          = goGlobalRequest("VOUCHER_DETAILS_REQUEST")
	bookingInfoForAmendment = goGlobalRequest("BOOKING_INFO_FOR_AMENDMENT_REQUEST")
	bookingAmendment        = goGlobalRequest("BOOKING_AMENDMENT_REQUEST")
	hotelInfo               = goGlobalRequest("HOTEL_INFO_REQUEST")
	priceBreakdown          = goGlobalRequest("PRICE_BREAKDOWN_REQUEST")
)

type goGlobalRequest string

var requestTypes = map[goGlobalRequest]int64{
	searchRequest:           11,
	bookingValidation:       9,
	bookingInsert:           2,
	bookingStatus:           5,
	bookingSearch:           4,
	advBookingSearch:        10,
	bookingCancel:           3,
	voucherDetails:          8,
	bookingInfoForAmendment: 15,
	bookingAmendment:        16,
	hotelInfo:               61,
	priceBreakdown:          14,
}
var defaultRequestVersion = map[goGlobalRequest]string{
	searchRequest:     "2.4",
	bookingValidation: "2.0",
	bookingInsert:     "2.3",
	bookingSearch:     "2.2",
	advBookingSearch:  "2.2",
	voucherDetails:    "2.3",
	hotelInfo:         "2.2",
	priceBreakdown:    "2.0",
}

type GoGlobalService interface {
	GetDestinations() ([]*Destination, error)
	GetHotels() ([]*Hotel, error)
	Search(models.HotelSearchRequest, RequestLogger, ResponseLogger) ([]models.HotelSearchResponseItem, error)
	BookingValuation(models.BookValuationRequest, RequestLogger, ResponseLogger) (models.BookValuationResponse, error)
	BookingInsert(models.BookingInsertRequest, RequestLogger, ResponseLogger) (models.BookingInsertResponse, error)
	BookingStatus(models.BookingStatusRequest, RequestLogger, ResponseLogger) (models.BookingStatusResponse, error)
	BookingSearch(models.BookingSearchRequest, RequestLogger, ResponseLogger) (models.BookingSearchResponse, error)
	AdvBookingSearch(models.AdvBookingSearchRequest, RequestLogger, ResponseLogger) (models.AdvBookingSearchResponse, error)
	BookingCancel(models.BookingCancelRequest, RequestLogger, ResponseLogger) (models.BookingCancelResponse, error)
	VoucherDetails(models.VoucherDetailsRequest, RequestLogger, ResponseLogger) (models.VoucherDetailsResponse, error)
	BookingInfoForAmendment(models.BookingInfoForAmendmentRequest, RequestLogger, ResponseLogger) (models.BookingInfoForAmendmentResponse, error)
	BookingAmendment(models.BookingAmendmentRequest, RequestLogger, ResponseLogger) error
	HotelInfo(models.HotelInfoRequest, RequestLogger, ResponseLogger) (models.HotelInfoResponse, error)
	PriceBreakdown(models.PriceBreakdownRequest, RequestLogger, ResponseLogger) (models.PriceBreakdownResponse, error)
}

type RequestLogger func(r *http.Request) error
type ResponseLogger func(r *http.Response) error

type Credentials struct {
	AgencyId string
	UserName string
	Password string
}

type goGlobalService struct {
	baseUrl         string
	baseCredentials Credentials
	client          *http.Client
}

func NewGoGlobalService(
	apiUrl string,
	baseCredentials Credentials,
) GoGlobalService {
	return &goGlobalService{
		baseUrl:         apiUrl,
		baseCredentials: baseCredentials,
		client:          &http.Client{},
	}
}

func (c *goGlobalService) GetDestinations() ([]*Destination, error) {
	req, err := http.NewRequest(http.MethodGet, getDestinationsUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.baseCredentials.UserName, c.baseCredentials.Password)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("GetDestinations: close connection: %s \n", err)
		}
	}()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("do request: %v", resp.Status)
	}

	var destinations []*Destination
	err = c.getDumpContent(resp.Body, &destinations)
	if err != nil {
		return nil, err
	}

	return destinations, nil
}

func (c *goGlobalService) GetHotels() ([]*Hotel, error) {
	url := fmt.Sprintf(getHotelsUrlFmt, c.baseCredentials.AgencyId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.baseCredentials.UserName, c.baseCredentials.Password)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("GetHotels: close connection: %s \n", err)
		}
	}()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("do request: %v", resp.Status)
	}
	var hotels []*Hotel
	err = c.getDumpContent(resp.Body, &hotels)
	if err != nil {
		return nil, err
	}

	return hotels, nil
}

func (c *goGlobalService) Search(
	request models.HotelSearchRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) ([]models.HotelSearchResponseItem, error) {
	if request.Version == "" {
		request.Version = defaultRequestVersion[searchRequest]
	}
	results := models.HotelSearchResponse{}
	response, err := c.doRequest(searchRequest, request, requestLogger, responseLogger)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &results)
	if err != nil {
		return nil, err
	}

	if results.Header.OperationType == models.OperationTypeError || results.Header.OperationType == models.OperationTypeMessage {
		return nil, fmt.Errorf("code: %d, message: %s", results.Main.Error.Code, results.Main.Error.Message)
	}

	return results.Hotels, nil
}

func (c *goGlobalService) BookingValuation(
	request models.BookValuationRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.BookValuationResponse, error) {
	if request.Version == "" {
		request.Version = defaultRequestVersion[bookingValidation]
	}
	return genericDoRequest[models.BookValuationRequest, models.BookValuationRoot, models.BookValuationResponse](
		c,
		bookingValidation, request, requestLogger, responseLogger)
}

func (c *goGlobalService) BookingInsert(
	request models.BookingInsertRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.BookingInsertResponse, error) {
	if request.Version == "" {
		request.Version = defaultRequestVersion[bookingInsert]
	}
	return genericDoRequest[models.BookingInsertRequest, models.BookingInsertRoot, models.BookingInsertResponse](
		c,
		bookingInsert, request, requestLogger, responseLogger)

}

func (c *goGlobalService) BookingStatus(
	request models.BookingStatusRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.BookingStatusResponse, error) {
	return genericDoRequest[models.BookingStatusRequest, models.BookingStatusRoot, models.BookingStatusResponse](
		c,
		bookingStatus, request, requestLogger, responseLogger)
}

func (c *goGlobalService) BookingSearch(
	request models.BookingSearchRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.BookingSearchResponse, error) {
	return genericDoRequest[models.BookingSearchRequest, models.BookingSearchRoot, models.BookingSearchResponse](
		c,
		bookingSearch, request, requestLogger, responseLogger)
}

func (c *goGlobalService) AdvBookingSearch(
	request models.AdvBookingSearchRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.AdvBookingSearchResponse, error) {
	if request.Version == "" {
		request.Version = defaultRequestVersion[advBookingSearch]
	}
	return genericDoRequest[models.AdvBookingSearchRequest, models.AdvBookingSearchRoot, models.AdvBookingSearchResponse](
		c,
		advBookingSearch, request, requestLogger, responseLogger)
}

func (c *goGlobalService) BookingCancel(
	request models.BookingCancelRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.BookingCancelResponse, error) {
	return genericDoRequest[models.BookingCancelRequest, models.BookingCancelRoot, models.BookingCancelResponse](
		c,
		bookingCancel, request, requestLogger, responseLogger)
}

func (c *goGlobalService) VoucherDetails(
	request models.VoucherDetailsRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.VoucherDetailsResponse, error) {
	return genericDoRequest[models.VoucherDetailsRequest, models.VoucherDetailsRoot, models.VoucherDetailsResponse](
		c,
		voucherDetails, request, requestLogger, responseLogger)
}

func (c *goGlobalService) BookingInfoForAmendment(
	request models.BookingInfoForAmendmentRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.BookingInfoForAmendmentResponse, error) {
	return genericDoRequest[models.BookingInfoForAmendmentRequest, models.BookingInfoForAmendmentRoot, models.BookingInfoForAmendmentResponse](
		c,
		bookingInfoForAmendment, request, requestLogger, responseLogger)
}

func (c *goGlobalService) BookingAmendment(
	request models.BookingAmendmentRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) error {
	_, err := genericDoRequest[models.BookingAmendmentRequest, models.BookingAmendmentRoot, models.BookingAmendmentResponse](
		c,
		bookingAmendment, request, requestLogger, responseLogger)

	return err
}

func (c *goGlobalService) HotelInfo(
	request models.HotelInfoRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.HotelInfoResponse, error) {
	if request.Version == "" {
		request.Version = defaultRequestVersion[hotelInfo]
	}
	return genericDoRequest[models.HotelInfoRequest, models.HotelInfoRoot, models.HotelInfoResponse](
		c,
		hotelInfo, request, requestLogger, responseLogger)
}

func (c *goGlobalService) PriceBreakdown(
	request models.PriceBreakdownRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.PriceBreakdownResponse, error) {

	return genericDoRequest[models.PriceBreakdownRequest, models.PriceBreakdownRoot, models.PriceBreakdownResponse](
		c,
		priceBreakdown, request, requestLogger, responseLogger)
}

func (c *goGlobalService) getDumpContent(compressedDump io.ReadCloser, out any) error {
	//т.к. дампы небольшие - десяток мегабайт в zip'е и ~ в 3 раза больше в распакованном, то просто загружаем в память
	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, compressedDump)
	if err != nil {
		return fmt.Errorf("getDumpContent: iocopy %w", err)
	}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '|'
		r.LazyQuotes = true
		r.TrimLeadingSpace = true
		r.ReuseRecord = true
		return r // Allows use pipe as delimiter
	})

	reader := bytes.NewReader(buff.Bytes())

	zipReader, err := zip.NewReader(reader, size)
	//у них в дампах всегда один текстовый файл, так что просто читаем и возвращаем содержимое первого
	for _, file := range zipReader.File {
		f, err := file.Open()
		if err != nil {
			return fmt.Errorf("getDumpContent: file open %w", err)
		}
		defer func() {
			if err = f.Close(); err != nil {
				log.Printf("getDumpContent: %s \n", err)
			}
		}()
		err = gocsv.Unmarshal(utfbom.SkipOnly(f), out)
		if err != nil {
			return err
		}

		return nil
	}
	return errors.New("getDumpContent: missing files in dump")
}

func (c *goGlobalService) SetBaseUrl(url string) {
	c.baseUrl = url
}

func (c *goGlobalService) doRequest(
	operation goGlobalRequest,
	request any,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) ([]byte, error) {
	encoded, err := xml.Marshal(request)
	if err != nil {
		return nil, err
	}
	credentials := c.baseCredentials

	requestRoot := models.RequestRoot{
		Header: models.Header{
			Agency:        credentials.AgencyId,
			User:          credentials.UserName,
			Password:      credentials.Password,
			Operation:     string(operation),
			OperationType: models.OperationTypeRequest,
		},
		Main: encoded,
	}

	xmlRequest, err := xml.Marshal(requestRoot)
	if err != nil {
		return nil, err
	}

	envelope := models.EnvelopeRequest{
		XsiType: "http://www.w3.org/2001/XMLSchema-instance",
		Body: models.Body{
			MakeRequest: models.MakeRequest{
				RequestType: requestTypes[operation],
				XmlRequest: models.XmlRequestCdata{
					Text: string(xmlRequest),
				},
			},
		},
	}

	payload, err := xml.Marshal(envelope)
	if err != nil {
		return nil, err
	}

	payload = append([]byte(xml.Header), payload...)

	req, err := http.NewRequest("POST", c.baseUrl, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.ContentLength = int64(len(payload))

	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("API-AgencyID", credentials.AgencyId)
	req.Header.Add("API-Operation", string(operation))
	req.Header.Add("Accept", "application/json")

	if requestLogger != nil {
		defer func() {
			go func() {
				err2 := requestLogger(req)
				if err2 != nil {
					log.Printf("save request error: %v \n", err)
				}

			}()
		}()
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("make request: close connection: %s \n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err == nil && responseLogger != nil {
		err2 := responseLogger(resp)
		if err2 != nil {
			log.Printf("save response error: %v \n", err)
		}
	}
	response := models.EnvelopeResponse{}

	err = xml.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("can't parse response: %w", err)
	}

	return response.Body.MakeRequestResponse.MakeRequestResult.Data, err
}

func genericDoRequest[REQ any, ROOT models.ResponseRoot[RES], RES any](
	service *goGlobalService,
	operation goGlobalRequest,
	req REQ,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (RES, error) {
	var response RES
	xmlResponse, err := service.doRequest(operation, req, requestLogger, responseLogger)
	if err != nil {
		return response, err
	}

	var root ROOT
	err = xml.Unmarshal(xmlResponse, &root)
	if err != nil {
		return response, err
	}

	if err = root.CheckError(); err != nil {
		return response, err
	}
	response = root.GetResponse()
	return response, nil
}
