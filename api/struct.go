package api

type PhoneNumber struct {
	IsValid             bool   `json:"isValid"`
	E164Format          string `json:"e164Format"`
	InternationalFormat string `json:"internationalFormat"`
	NationalFormat      string `json:"nationalFormat"`
	Location            string `json:"location"`
	LineType            string `json:"lineType"`
	Country             struct {
		IsoAlpha2         string `json:"isoAlpha2"`
		IsoAlpha3         string `json:"isoAlpha3"`
		M49Code           int    `json:"m49Code"`
		Name              string `json:"name"`
		IsoName           string `json:"isoName"`
		IsoNameFull       string `json:"isoNameFull"`
		IsoAdminLanguages []struct {
			IsoAlpha3  string `json:"isoAlpha3"`
			IsoAlpha2  string `json:"isoAlpha2"`
			IsoName    string `json:"isoName"`
			NativeName string `json:"nativeName"`
		} `json:"isoAdminLanguages"`
		UnRegion string `json:"unRegion"`
		Currency struct {
			NumericCode int    `json:"numericCode"`
			Code        string `json:"code"`
			Name        string `json:"name"`
			MinorUnits  int    `json:"minorUnits"`
		} `json:"currency"`
		WbRegion struct {
			ID       string `json:"id"`
			Iso2Code string `json:"iso2Code"`
			Value    string `json:"value"`
		} `json:"wbRegion"`
		WbIncomeLevel struct {
			ID       string `json:"id"`
			Iso2Code string `json:"iso2Code"`
			Value    string `json:"value"`
		} `json:"wbIncomeLevel"`
		CallingCode      string `json:"callingCode"`
		CountryFlagEmoji string `json:"countryFlagEmoji"`
		WikidataID       string `json:"wikidataId"`
		GeonameID        int    `json:"geonameId"`
		IsIndependent    bool   `json:"isIndependent"`
	} `json:"country"`
}
