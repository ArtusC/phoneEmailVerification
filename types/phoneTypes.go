package types

type IsoAdminLanguages []struct {
	IsoAlpha3  string `bson:"isoAlpha3" json:"isoAlpha3,omitempty"`
	IsoAlpha2  string `bson:"isoAlpha2" json:"isoAlpha2,omitempty"`
	IsoName    string `bson:"isoName" json:"isoName,omitempty"`
	NativeName string `bson:"nativeName" json:"nativeName,omitempty"`
}

type Currency struct {
	NumericCode int32  `bson:"numericCode" json:"numericCode,omitempty"`
	Code        string `bson:"code" json:"code,omitempty"`
	Name        string `bson:"name" json:"name,omitempty"`
	MinorUnits  int32  `bson:"minorUnits" json:"minorUnits,omitempty"`
}

type WbRegion struct {
	ID       string `bson:"id" json:"id,omitempty"`
	Iso2Code string `bson:"iso2Code" json:"iso2Code,omitempty"`
	Value    string `bson:"value" json:"value,omitempty"`
}

type WbIncomeLevel struct {
	ID       string `bson:"id" json:"id,omitempty"`
	Iso2Code string `bson:"iso2Code" json:"iso2Code,omitempty"`
	Value    string `bson:"value" json:"value,omitempty"`
}

type Country struct {
	IsoAlpha2         string            `bson:"isoAlpha2" json:"isoAlpha2,omitempty"`
	IsoAlpha3         string            `bson:"isoAlpha3" json:"isoAlpha3,omitempty"`
	M49Code           int               `bson:"m49Code" json:"m49Code,omitempty"`
	Name              string            `bson:"name" json:"name,omitempty"`
	IsoName           string            `bson:"isoName" json:"isoName,omitempty"`
	IsoNameFull       string            `bson:"isoNameFull" json:"isoNameFull,omitempty"`
	IsoAdminLanguages IsoAdminLanguages `bson:"isoAdminLanguages" json:"isoAdminLanguages,omitempty"`
	UnRegion          string            `bson:"unRegion" json:"unRegion,omitempty"`
	Currency          Currency          `bson:"currency" json:"currency,omitempty"`
	WbRegion          WbRegion          `bson:"wbRegion" json:"wbRegion,omitempty"`
	WbIncomeLevel     WbIncomeLevel     `bson:"wbIncomeLevel" json:"wbIncomeLevel,omitempty"`
	CallingCode       int32             `bson:"callingCode" json:"callingCode,omitempty"`
	CountryFlagEmoji  string            `bson:"countryFlagEmoji" json:"countryFlagEmoji,omitempty"`
	WikidataID        string            `bson:"wikidataId" json:"wikidataId,omitempty"`
	GeonameID         string            `bson:"geonameId" json:"geonameId,omitempty"`
	IsIndependent     bool              `bson:"isIndependent" json:"isIndependent,omitempty"`
}

type PhoneNumber struct {
	ID                  string  `bson:"_id" json:"_id"`
	PhoneInput          string  `bson:"phoneInput" json:"phoneInput,omitempty"`
	IsValid             bool    `bson:"isValid" json:"isValid,omitempty"`
	E164Format          string  `bson:"e164Format" json:"e164Format,omitempty"`
	InternationalFormat string  `bson:"internationalFormat" json:"internationalFormat,omitempty"`
	NationalFormat      string  `bson:"nationalFormat" json:"nationalFormat,omitempty"`
	Location            string  `bson:"location" json:"location,omitempty"`
	LineType            string  `bson:"lineType" json:"lineType,omitempty"`
	Country             Country `bson:"country" json:"country,omitempty"`
}

type PhoneNumberResults []PhoneNumber
