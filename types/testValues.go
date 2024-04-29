package types

var (
	TestPhoneValue = PhoneNumber{
		PhoneInput:          "12018675309",
		IsValid:             true,
		E164Format:          "+12018675309",
		InternationalFormat: "+1 201-867-5309",
		NationalFormat:      "(201) 867-5309",
		Location:            "New Jersey",
		LineType:            "FIXED_LINE_OR_MOBILE",
		Country: Country{
			IsoAlpha2:        "US",
			IsoAlpha3:        "USA",
			M49Code:          840,
			Name:             "United States of America (the)",
			IsoName:          "United States of America (the)",
			IsoNameFull:      "the United States of America",
			UnRegion:         "Europe and Northern America/Northern America",
			CallingCode:      1,
			CountryFlagEmoji: "🇺🇸",
			WikidataID:       "Q30",
			GeonameID:        "6252001",
			IsIndependent:    true,
			IsoAdminLanguages: IsoAdminLanguages{{
				IsoAlpha3:  "eng",
				IsoAlpha2:  "en",
				IsoName:    "English",
				NativeName: "English",
			}},
			Currency: Currency{
				NumericCode: 840,
				Code:        "USD",
				Name:        "US Dollar",
				MinorUnits:  2,
			},
			WbRegion: WbRegion{
				ID:       "NAC",
				Iso2Code: "XU",
				Value:    "North America",
			},
			WbIncomeLevel: WbIncomeLevel{
				ID:       "HIC",
				Iso2Code: "XD",
				Value:    "High income",
			},
		},
	}
)
