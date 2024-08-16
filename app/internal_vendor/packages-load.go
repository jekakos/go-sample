package internal_vendor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sample-service/app/internal_vendor/dto"
)

func LoadPackages() ([]dto.VendorPackage, error) {

	var resultPackages []dto.VendorPackage
	maxCount := 2000
	limit := 100

	start := 0
	token := os.Getenv("VENDOR_TOKEN")
	baseURL := os.Getenv("VENDOR_API_BASE")
	client := &http.Client{}
	all_count := 0

	for start < maxCount {

		// get response from start with offset
		url := fmt.Sprintf("%s/wholesale/plans?limit=%d&offset=%d", baseURL, limit, start)

		println(url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", fmt.Sprintf("Apikey %s", token))
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Vendor error: unexpected status code: %d", resp.StatusCode)
		}

		var response dto.VendorPackageResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("Decore JSON error: %s", err.Error()))
		}

		// Add results to common resultPackages
		resultPackages = append(resultPackages, response.Rows...)

		// if reached the end
		count := len(response.Rows)
		all_count += count

		if count < limit {
			println("break")
			break
		}

		start += limit
	}

	fmt.Printf("All %d\n", all_count)

	return resultPackages, nil

}
