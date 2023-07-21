# Go library for Go Global Travel

documentation: https://data.goglobal.travel/specs/XmlSpecs_v3.14.3.html

Example of usage:
```
service := client.NewGoGlobalService(
		API_URL,
		AGENCY_ID,
		USERNAME,
		PASSWORD,
	)

	response, err := service.Search(models.HotelSearchRequest{
		Version:        "2.3",
		ResponseFormat: client.ResponseFormatJson,
		IncludeGeo:     true,
		Currency:       "EUR",
		FilterPriceMin: nil,
		FilterPriceMax: nil,
		CityCode:       []int64{75},
		ArrivalDate:    "2023-10-01",
		Nights:         2,
		Apartments:     false,
		Nationality:    "RU",
		Rooms: models.SearchRooms{
			Room: []models.SearchRoom{
				{
					Adults:     2,
					RoomCount:  1,
					ChildCount: 0,
					CotCount:   0,
					ChildAge:   nil,
				},
			},
		},
	}, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range response {
	    ///handle offers
	    ...
	}
```