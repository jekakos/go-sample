package internal_vendor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sample-service/app/esim/dto"
)

func EsimTopup(iccid string, data dto.VendorCreateEsimRequestDTO) (dto.VendorCreateEsimResponseDTO, error) {

	// Check if possible to add Plan
	availibalePlans, err := GetPackagesForTopup(iccid)
	if err != nil {
		return dto.VendorCreateEsimResponseDTO{}, fmt.Errorf("Error getting availibalePlans: %v\n", err)
	}

	found := false
	for _, plan := range availibalePlans {
		if plan.Id == data.PlanId {
			fmt.Printf("Possible to add Plan %d\n", data.PlanId)
			found = true
			break
		}
	}

	if found == false {
		return dto.VendorCreateEsimResponseDTO{}, fmt.Errorf("Cannot find availibale plan (%d) in  %+v\n", data.PlanId, availibalePlans)
	}

	// Add Plan
	var resultEsim dto.VendorCreateEsimResponseDTO
	token := os.Getenv("VENDOR_TOKEN")
	baseURL := os.Getenv("VENDOR_API_BASE")

	url := fmt.Sprintf("%s/wholesale/esims/%s/top-up", baseURL, iccid)

	println(url)

	jsonBody, err := json.Marshal(data)
	if err != nil {
		return dto.VendorCreateEsimResponseDTO{}, fmt.Errorf("Error marshalling JSON: %v\n", err)
	}

	fmt.Printf("Vendor request JSON body: %+v\n", data)

	bytes := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", url, bytes)
	if err != nil {
		return dto.VendorCreateEsimResponseDTO{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Apikey %s", token))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return dto.VendorCreateEsimResponseDTO{}, fmt.Errorf("TopupEsim: Error getting response from Vendor")
	}
	fmt.Printf("Vendor response: %+v\n", resp)
	fmt.Printf("SatusCode: %d\n", resp.StatusCode)
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return dto.VendorCreateEsimResponseDTO{}, fmt.Errorf("TopupEsim: Vendor server returned an internal error")
	}

	err = json.NewDecoder(resp.Body).Decode(&resultEsim)
	if err != nil {
		return dto.VendorCreateEsimResponseDTO{}, fmt.Errorf("TopupEsim: Error decoding response from Vendor")
	}

	return resultEsim, nil
}
