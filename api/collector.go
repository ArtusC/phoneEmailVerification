package api

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"

// 	"golang.org/x/net/context"
// )

// func CollectAndStoreData() error {
// 	// Collect data from the API
// 	ctx := context.Background()

// 	data, err := collectData(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("Data collected: %v\n", data)
// 	// Store the data in the database
// 	// err = db.StoreData(data)
// 	// if err != nil {
// 	//     return err
// 	// }

// 	return nil
// }

// func collectData(ctx context.Context) (phoneNumberResponse PhoneNumberValidationCollector, err error) {
// 	// Implement logic to collect data from the API
// 	baseURL := "https://api-bdc.net/data/phone-number-validate"
// 	apiKey := "bdc_7467afbb1e5a4790b12598f6a85d35da"

// 	// Create a map to store query parameters
// 	queryParams := url.Values{}
// 	queryParams.Set("number", "201 867-5309")
// 	queryParams.Set("countryCode", "us")
// 	queryParams.Set("localityLanguage", "en")
// 	queryParams.Set("key", apiKey)

// 	// Construct the URL with query parameters
// 	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

// 	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()
// 	bodyText, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	json.Unmarshal(bodyText, &phoneNumberResponse)

// 	fmt.Printf("%v\n", bodyText)

// 	return phoneNumberResponse, nil
// }
