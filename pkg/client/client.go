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
	"regexp"
	"strconv"

	"github.com/DmitryKolbin/go-global/pkg/client/models"
	"github.com/dimchansky/utfbom"
	"github.com/gocarina/gocsv"
)

const (
	getDestinationsUrl = "https://static-data.tourismcloudservice.com/propsdata/Destinations/compress/true"
	getHotelsUrlFmt    = "https://static-data.tourismcloudservice.com/agency/hotels/%d"

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
	GetDestinations(Credentials) ([]*Destination, error)
	GetHotels(Credentials) ([]*Hotel, error)
	Search(Credentials, models.HotelSearchRequest, RequestLogger, ResponseLogger) ([]models.HotelSearchResponseItem, error)
	BookingValuation(Credentials, models.BookValuationRequest, RequestLogger, ResponseLogger) (models.BookValuationResponse, error)
	BookingInsert(Credentials, models.BookingInsertRequest, RequestLogger, ResponseLogger) (models.BookingInsertResponse, error)
	BookingStatus(Credentials, models.BookingStatusRequest, RequestLogger, ResponseLogger) (models.BookingStatusResponse, error)
	BookingSearch(Credentials, models.BookingSearchRequest, RequestLogger, ResponseLogger) (models.BookingSearchResponse, error)
	AdvBookingSearch(Credentials, models.AdvBookingSearchRequest, RequestLogger, ResponseLogger) (models.AdvBookingSearchResponse, error)
	BookingCancel(Credentials, models.BookingCancelRequest, RequestLogger, ResponseLogger) (models.BookingCancelResponse, error)
	VoucherDetails(Credentials, models.VoucherDetailsRequest, RequestLogger, ResponseLogger) (models.VoucherDetailsResponse, error)
	BookingInfoForAmendment(Credentials, models.BookingInfoForAmendmentRequest, RequestLogger, ResponseLogger) (models.BookingInfoForAmendmentResponse, error)
	BookingAmendment(Credentials, models.BookingAmendmentRequest, RequestLogger, ResponseLogger) error
	HotelInfo(Credentials, models.HotelInfoRequest, RequestLogger, ResponseLogger) (models.HotelInfoResponse, error)
	PriceBreakdown(Credentials, models.PriceBreakdownRequest, RequestLogger, ResponseLogger) (models.PriceBreakdownResponse, error)
}

type RequestLogger func(r *http.Request) error
type ResponseLogger func(r *http.Response, responseBody []byte, isError bool) error

type Credentials struct {
	AgencyId int64
	UserName string
	Password string
}

type goGlobalService struct {
	baseUrl string
	client  *http.Client
}

func NewGoGlobalService(
	apiUrl string,
) GoGlobalService {
	return &goGlobalService{
		baseUrl: apiUrl,
		client:  &http.Client{},
	}
}

func (c *goGlobalService) GetDestinations(credentials Credentials) ([]*Destination, error) {
	req, err := http.NewRequest(http.MethodGet, getDestinationsUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(credentials.UserName, credentials.Password)

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

func (c *goGlobalService) GetHotels(credentials Credentials) ([]*Hotel, error) {
	url := fmt.Sprintf(getHotelsUrlFmt, credentials.AgencyId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(credentials.UserName, credentials.Password)

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
	credentials Credentials,
	request models.HotelSearchRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) ([]models.HotelSearchResponseItem, error) {
	if request.Version == "" {
		request.Version = defaultRequestVersion[searchRequest]
	}
	results := models.HotelSearchResponse{}

	var resp *http.Response
	var respBody []byte
	response, err := c.doRequest(credentials, searchRequest, request, requestLogger, func(r *http.Response, body []byte, isError bool) error {
		resp = r
		respBody = body
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &results)
	if err != nil {
		return nil, err
	}

	if results.Header.OperationType == models.OperationTypeError || results.Header.OperationType == models.OperationTypeMessage {
		if responseLogger != nil {
			err = responseLogger(resp, respBody, true)
			if err != nil {
				return nil, fmt.Errorf("cant't save responsee log with error code: %d, message: %s; actual error %v", results.Main.Error.Code, results.Main.Error.Message, err)
			}
		}
		return nil, fmt.Errorf("code: %d, message: %s", results.Main.Error.Code, results.Main.Error.Message)
	}
	if responseLogger != nil {
		err = responseLogger(resp, respBody, false)
		if err != nil {
			return nil, fmt.Errorf("cant't save responsee log")
		}
	}

	return results.Hotels, nil
}

func (c *goGlobalService) BookingValuation(
	credentials Credentials,
	request models.BookValuationRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.BookValuationResponse, error) {
	if request.Version == "" {
		request.Version = defaultRequestVersion[bookingValidation]
	}
	return genericDoRequest[models.BookValuationRequest, models.BookValuationRoot, models.BookValuationResponse](
		credentials,
		c,
		bookingValidation,
		request,
		requestLogger,
		responseLogger,
	)
}

func (c *goGlobalService) BookingInsert(
	credentials Credentials,
	request models.BookingInsertRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.BookingInsertResponse, error) {
	if request.Version == "" {
		request.Version = defaultRequestVersion[bookingInsert]
	}
	return genericDoRequest[models.BookingInsertRequest, models.BookingInsertRoot, models.BookingInsertResponse](
		credentials,
		c,
		bookingInsert,
		request,
		requestLogger,
		responseLogger,
	)

}

func (c *goGlobalService) BookingStatus(
	credentials Credentials,
	request models.BookingStatusRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.BookingStatusResponse, error) {
	return genericDoRequest[models.BookingStatusRequest, models.BookingStatusRoot, models.BookingStatusResponse](
		credentials,
		c,
		bookingStatus,
		request,
		requestLogger,
		responseLogger,
	)
}

func (c *goGlobalService) BookingSearch(
	credentials Credentials,
	request models.BookingSearchRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.BookingSearchResponse, error) {
	return genericDoRequest[models.BookingSearchRequest, models.BookingSearchRoot, models.BookingSearchResponse](
		credentials,
		c,
		bookingSearch,
		request,
		requestLogger,
		responseLogger,
	)
}

func (c *goGlobalService) AdvBookingSearch(
	credentials Credentials,
	request models.AdvBookingSearchRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.AdvBookingSearchResponse, error) {
	if request.Version == "" {
		request.Version = defaultRequestVersion[advBookingSearch]
	}
	return genericDoRequest[models.AdvBookingSearchRequest, models.AdvBookingSearchRoot, models.AdvBookingSearchResponse](
		credentials,
		c,
		advBookingSearch,
		request,
		requestLogger,
		responseLogger,
	)
}

func (c *goGlobalService) BookingCancel(
	credentials Credentials,
	request models.BookingCancelRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.BookingCancelResponse, error) {
	return genericDoRequest[models.BookingCancelRequest, models.BookingCancelRoot, models.BookingCancelResponse](
		credentials,
		c,
		bookingCancel,
		request,
		requestLogger,
		responseLogger,
	)
}

func (c *goGlobalService) VoucherDetails(
	credentials Credentials,
	request models.VoucherDetailsRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.VoucherDetailsResponse, error) {
	return genericDoRequest[models.VoucherDetailsRequest, models.VoucherDetailsRoot, models.VoucherDetailsResponse](
		credentials,
		c,
		voucherDetails,
		request,
		requestLogger,
		responseLogger,
	)
}

func (c *goGlobalService) BookingInfoForAmendment(
	credentials Credentials,
	request models.BookingInfoForAmendmentRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.BookingInfoForAmendmentResponse, error) {
	return genericDoRequest[models.BookingInfoForAmendmentRequest, models.BookingInfoForAmendmentRoot, models.BookingInfoForAmendmentResponse](
		credentials,
		c,
		bookingInfoForAmendment,
		request,
		requestLogger,
		responseLogger,
	)
}

func (c *goGlobalService) BookingAmendment(
	credentials Credentials,
	request models.BookingAmendmentRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) error {
	_, err := genericDoRequest[models.BookingAmendmentRequest, models.BookingAmendmentRoot, models.BookingAmendmentResponse](
		credentials,
		c,
		bookingAmendment,
		request,
		requestLogger,
		responseLogger,
	)

	return err
}

func (c *goGlobalService) HotelInfo(
	credentials Credentials,
	request models.HotelInfoRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.HotelInfoResponse, error) {
	if request.Version == "" {
		request.Version = defaultRequestVersion[hotelInfo]
	}
	return genericDoRequest[models.HotelInfoRequest, models.HotelInfoRoot, models.HotelInfoResponse](
		credentials,
		c,
		hotelInfo,
		request,
		requestLogger,
		responseLogger,
	)
}

func (c *goGlobalService) PriceBreakdown(
	credentials Credentials,
	request models.PriceBreakdownRequest,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (models.PriceBreakdownResponse, error) {

	return genericDoRequest[models.PriceBreakdownRequest, models.PriceBreakdownRoot, models.PriceBreakdownResponse](
		credentials,
		c,
		priceBreakdown,
		request,
		requestLogger,
		responseLogger,
	)
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
	credentials Credentials,
	operation goGlobalRequest,
	request any,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) ([]byte, error) {
	encoded, err := xml.Marshal(request)
	if err != nil {
		return nil, err
	}
	requestRoot := models.RequestRoot{
		Header: models.Header{
			Agency:        json.Number(strconv.FormatInt(credentials.AgencyId, 10)),
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
	req.Header.Add("API-AgencyID", strconv.FormatInt(credentials.AgencyId, 10))
	req.Header.Add("API-Operation", string(operation))
	req.Header.Add("Accept", "application/json")

	if requestLogger != nil {
		err2 := requestLogger(req)
		if err2 != nil {
			log.Printf("save request error: %v \n", err2)
		}
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
	if err != nil {
		return nil, err
	}

	if responseLogger != nil {
		err2 := responseLogger(resp, body, false)
		if err2 != nil {
			log.Printf("save response error: %v \n", err2)
		}
	}
	response := models.EnvelopeResponse{}

	//assume that data in response contain &#x0000 characters only when it's typed by mistake
	//and remove them cause it break go xml decoder
	re := regexp.MustCompile(`&#x[\da-fA-F]+;`)
	body = []byte(re.ReplaceAllString(string(body), ""))
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("can't parse response: %w", err)
	}

	return response.Body.MakeRequestResponse.MakeRequestResult.Data, err
}

func genericDoRequest[REQ any, ROOT models.ResponseRoot[RES], RES any](
	credentials Credentials,
	service *goGlobalService,
	operation goGlobalRequest,
	req REQ,
	requestLogger RequestLogger,
	responseLogger ResponseLogger,
) (RES, error) {
	var response RES

	var resp *http.Response
	var respBody []byte

	xmlResponse, err := service.doRequest(credentials, operation, req, requestLogger, func(r *http.Response, body []byte, isError bool) error {
		resp = r
		respBody = body
		return nil
	})
	if err != nil {
		return response, err
	}

	var root ROOT
	err = xml.Unmarshal(xmlResponse, &root)
	if err != nil {
		return response, err
	}

	if err = root.CheckError(); err != nil {
		if responseLogger != nil {
			err2 := responseLogger(resp, respBody, true)
			if err2 != nil {
				return response, fmt.Errorf("cant't save responsee log with error %v, actual eror %v", err, err2)
			}
		}
		return response, err
	}

	if responseLogger != nil {
		err = responseLogger(resp, respBody, false)
		if err != nil {
			return response, fmt.Errorf("cant't save responsee log")
		}
	}

	response = root.GetResponse()
	return response, nil
}
