//go:build integration
// +build integration

// package repository

package repository_test

import (
	tp "github.com/ArtusC/phoneEmailVerification/types"
)

var (
	TestPhoneValue = tp.PhoneNumber{
		PhoneInput:          "12018675309",
		IsValid:             true,
		E164Format:          "+12018675309",
		InternationalFormat: "+1 201-867-5309",
		NationalFormat:      "(201) 867-5309",
		Location:            "New Jersey",
		LineType:            "FIXED_LINE_OR_MOBILE",
		Country: tp.Country{
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
			IsoAdminLanguages: tp.IsoAdminLanguages{{
				IsoAlpha3:  "eng",
				IsoAlpha2:  "en",
				IsoName:    "English",
				NativeName: "English",
			}},
			Currency: tp.Currency{
				NumericCode: 840,
				Code:        "USD",
				Name:        "US Dollar",
				MinorUnits:  2,
			},
			WbRegion: tp.WbRegion{
				ID:       "NAC",
				Iso2Code: "XU",
				Value:    "North America",
			},
			WbIncomeLevel: tp.WbIncomeLevel{
				ID:       "HIC",
				Iso2Code: "XD",
				Value:    "High income",
			},
		},
	}

	TestPhoneValue_2 = tp.PhoneNumber{
		PhoneInput: "12018675310",
		IsValid:    true,
		E164Format: "+12018675310",
	}

	TestPhoneValue_3 = tp.PhoneNumber{
		PhoneInput: "12018675355",
		IsValid:    true,
		E164Format: "+12018675355",
	}
)
