package internal_vendor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sample-service/app/esim/dto"
	"strconv"
)

func CreateEsim(requestBody dto.VendorCreateEsimRequestDTO, isProd bool) (dto.VendorCreateEsimResponseDTO, error) {

	// Prepare request to Vendor
	token := os.Getenv("VENDOR_TOKEN")
	baseURL := os.Getenv("VENDOR_API_BASE")

	url := fmt.Sprintf("%s/wholesale/orders", baseURL)
	println(url)

	var response dto.VendorCreateEsimResponseDTO

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return dto.VendorCreateEsimResponseDTO{}, fmt.Errorf("Error marshalling JSON: %v\n", err)
	}

	fmt.Printf("Vendor request JSON body: %+v\n", requestBody)

	bytes := bytes.NewReader(jsonBody)
	req, err := http.NewRequest("POST", url, bytes)
	if err != nil {
		return dto.VendorCreateEsimResponseDTO{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Apikey %s", token))
	req.Header.Add("Content-Type", "application/json")

	fmt.Printf("Vendor request: %+v\n", req)

	//This is Prod E-SIM creation
	if isProd {

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return dto.VendorCreateEsimResponseDTO{}, fmt.Errorf("CreateEsim: Error getting response from Vendor")
		}
		fmt.Printf("Vendor response: %+v\n", resp)
		fmt.Printf("SatusCode: %d\n", resp.StatusCode)
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			return dto.VendorCreateEsimResponseDTO{}, fmt.Errorf("CreateEsim: Vendor server returned an internal error")
		}

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return dto.VendorCreateEsimResponseDTO{}, fmt.Errorf("CreateEsim: Error decoding response from Vendor")
		}
	} else {

		// Test simulation

		// Random ICCID
		min := 1000000
		max := 9999999
		randomNumber := rand.Intn(max-min+1) + min
		randomNumberStr := strconv.Itoa(randomNumber)

		jsonString := `{
    "iosSetup": "https://esimsetup.apple.com/esim_qrcode_provisioning?carddata=LPA:1$smdp.io$K2-00TEST-1KGTEST",
    "iccid": "893720401616` + randomNumberStr + `",
    "iosmatchingIdSetup": "",
    "smdpAddress": "smdp.io",
    "profileStatus": "RELEASED",
    "installedAt": "",
    "assignedPlans": [
					{
							"initialQuantityInBytes": 1000000000,
							"remainingQuantityInBytes": 1000000000,
							"startTime": "",
							"endTime": "",
							"isExpired": false,
							"areas": [
									{
											"id": 93,
											"name": "Portugal",
											"region": "Europe",
											"iso": "PT"
									}
							]
						}
				]
		}`

		err = json.Unmarshal([]byte(jsonString), &response)
		if err != nil {
			return dto.VendorCreateEsimResponseDTO{}, fmt.Errorf("Error decoding response from Vendor: %v", err)
		}
	}

	return response, nil
}
