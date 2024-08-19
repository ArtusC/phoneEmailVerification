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

	TestPhoneValueUpsert_1 = PhoneNumber{
		ID:                  "0dab7e5e343206634713474e42af8fe3",
		PhoneInput:          "12018675309",
		IsValid:             true,
		E164Format:          "+12018675309",
		InternationalFormat: "+1 201-867-5309",
		NationalFormat:      "(201) 867-5309",
		Location:            "Texas Man",
		LineType:            "FIXED_LINE_OR_MOBILE",
		Country: Country{
			IsoAlpha2:        "US",
			IsoAlpha3:        "USA",
			M49Code:          841,
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
				IsoAlpha3:  "eng2",
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

	TestPhoneValueUpsertResult = PhoneNumberResults{
		PhoneNumber{
			ID:                  "0dab7e5e343206634713474e42af8fe3",
			PhoneInput:          "12018675309",
			IsValid:             true,
			E164Format:          "+12018675309",
			InternationalFormat: "+1 201-867-5309",
			NationalFormat:      "(201) 867-5309",
			Location:            "Texas Man",
			LineType:            "FIXED_LINE_OR_MOBILE",
			Country: Country{
				IsoAlpha2:        "US",
				IsoAlpha3:        "USA",
				M49Code:          841,
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
					IsoAlpha3:  "eng2",
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
		},
	}

	TestPhoneValueUpsert_2 = PhoneNumber{
		ID:         "0dab7e5e343206634713474e42af8fe3",
		PhoneInput: "12018675309",
	}

	TestPhoneValueUpsertResult_2 = PhoneNumberResults{
		TestPhoneValueUpsert_2,
	}

	TestPhoneValueUpsert_3 = PhoneNumber{
		ID:         "0dab7e5e343206634713474e42af8fe3",
		PhoneInput: "12018675309",
		Location:   "OLD Texas",
		Country: Country{
			WbRegion: WbRegion{
				ID: "1",
			},
			WbIncomeLevel: WbIncomeLevel{
				ID: "2",
			},
			Currency: Currency{
				NumericCode: 842,
			},
			IsoAdminLanguages: IsoAdminLanguages{
				{NativeName: "Brazilian"},
			},
		},
	}

	TestPhoneValueUpsertResult_3 = PhoneNumberResults{
		TestPhoneValueUpsert_3,
	}

	TestPhoneValue_2 = PhoneNumber{
		PhoneInput: "12018675310",
		IsValid:    true,
		E164Format: "+12018675310",
	}

	TestPhoneValue_3 = PhoneNumber{
		PhoneInput: "12018675355",
		IsValid:    true,
		E164Format: "+12018675355",
	}

	TestPhoneValue_Null = PhoneNumber{
		ID:         "",
		PhoneInput: "",
	}
)
