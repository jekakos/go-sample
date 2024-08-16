package internal_vendor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	esimDto "sample-service/app/esim/dto"
)

func GetEsim(requestQuery esimDto.VendorGetEsimRequestDTO) (esimDto.VendorCreateEsimResponseDTO, error) {
	println("Get Vendor Esim")

	// Prepare request to Vendor
	token := os.Getenv("VENDOR_TOKEN")
	baseURL := os.Getenv("VENDOR_API_BASE")

	url := fmt.Sprintf("%s/wholesale/esims/%s", baseURL, requestQuery.Iccid)
	println(url)

	var response esimDto.VendorCreateEsimResponseDTO

	jsonBody, err := json.Marshal(requestQuery)
	if err != nil {
		return esimDto.VendorCreateEsimResponseDTO{}, fmt.Errorf("Error marshalling JSON: %v\n", err)
	}

	bytes := bytes.NewReader(jsonBody)
	req, err := http.NewRequest("GET", url, bytes)
	if err != nil {
		return esimDto.VendorCreateEsimResponseDTO{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Apikey %s", token))
	req.Header.Add("Content-Type", "application/json")

	fmt.Printf("Vendor request: %+v\n", req)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return esimDto.VendorCreateEsimResponseDTO{}, fmt.Errorf("GetEsim: Error getting response from Vendor")
	}
	fmt.Printf("Vendor response: %+v\n", resp)
	fmt.Printf("SatusCode: %d\n", resp.StatusCode)
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return esimDto.VendorCreateEsimResponseDTO{}, fmt.Errorf("GetEsim: Vendor server returned an internal error")
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return esimDto.VendorCreateEsimResponseDTO{}, fmt.Errorf("GetEsim: Error decoding response from Vendor")
	}

	return response, nil

}
