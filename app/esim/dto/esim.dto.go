package dto

type CreateEsimRequestDTO struct {
	PackageId string `json:"package_id" validate:"required"`
	VendorId  int8   `json:"vendor_id" validate:"required"`
	UserUUID  string `json:"user_uuid" validate:"required"`
	Prod      *bool  `json:"prod,omitempty"`
}

type CreateEsimResponseDTO struct {
	Iccid       string  `json:"iccid"`
	Lpa         string  `json:"lpa"`
	Imsi        *string `json:"imsi"`
	Msisdn      *string `json:"msisdn"`
	ExtensionId *int64  `json:"extension_id"`
	UserUUID    string  `json:"user_uuid"`
	VendorId    int8    `json:"vendor_id"`
	PackageId   string  `json:"package_id"`
}

/*pub struct CreateEsimResponseInterfaceDTO {
    pub iccid: String,
    pub lpa: String,
    pub imsi: Option<String>,
    pub msisdn: Option<String>,

    pub extension_id: Option<i32>,

    pub user_uuid: String,
    pub vendor_id: i16,
    pub package_id: Option<String>,
}
*/

type VendorCreateEsimRequestDTO struct {
	PlanId int64 `json:"planId" validate:"required"`
}

type UpdateEsimRequestDTO struct {
	ProfileStatus *string `json:"profileStatus,omitempty"`
	InstalledAt   *string `json:"installedAt,omitempty"`
}

type VendorCreateEsimResponseDTO struct {
	IosSetup      string             `json:"iosSetup"`
	Iccid         string             `json:"iccid"`
	MatchingId    string             `json:"iosmatchingIdSetup"`
	SmdpAddress   string             `json:"smdpAddress"`
	ProfileStatus string             `json:"profileStatus"`
	InstalledAt   string             `json:"installedAt"`
	AssignedPlans []AssignedPlansDTO `json:"assignedPlans"`
}

type VendorGetEsimRequestDTO struct {
	Iccid string `json:"iccid"`
}

type AssignedPlansDTO struct {
	Id                       int64     `json:"id"`
	PlanId                   int64     `json:"planId"`
	InitialQuantityInBytes   int64     `json:"initialQuantityInBytes"`
	RemainingQuantityInBytes int64     `json:"remainingQuantityInBytes"`
	StartTime                string    `json:"startTime"`
	EndTime                  string    `json:"endTime"`
	IsExpired                bool      `json:"isExpired"`
	Areas                    []AreaDTO `json:"areas"`
}

type AreaDTO struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
	Iso    string `json:"iso"`
}

/*
{
  "iosSetup": "https://esimsetup.apple.com/esim_qrcode_provisioning?carddata=LPA:1$address.test$X2-2XXXX7-LCIUXQ",
  "iccid": "89372040161688888888",
  "matchingId": "X2-2XXXX7-LCIUXQ",
  "smdpAddress": "address.test",
  "profileStatus": "RELEASED",
  "installedAt": null,
  "assignedPlans": [
    {
			"id": 89,
      "planId": 5,
      "initialQuantityInBytes": 1000000000,
      "remainingQuantityInBytes": 1000000000,
      "startTime": null,
      "endTime": null,
      "isExpired": false,
      "areas": [
        {
          "id": 49,
          "name": "Greece",
          "region": "Europe",
          "iso": "GR"
        }
      ]
    }
  ]
}
*/

type GetEsimRequestDTO struct {
	Iccid        *string `json:"iccid,omitempty"  query:"iccid"`
	UserUUID     *string `json:"user_uuid,omitempty" query:"user_uuid"`
	ActualStatus *bool   `json:"actual_status,omitempty" query:"actual_status"`
}

type GetEsimResponseDTO struct {
	Iccid  string `json:"iccid"`
	Lpa    string `json:"lpa"`
	Imsi   string `json:"imsi"`
	Msisdn string `json:"msisdn"`

	UserUUID string `json:"user_uuid"`
	VendorId int8   `json:"vendor_id"`
}

type GetEsimDbDTO struct {
	ID            int64  `json:"id"`
	Iccid         string `json:"iccid"`
	Lpa           string `json:"lpa"`
	ProfileStatus string `json:"profile_status"`
	SmdpAddress   string `json:"smdp_address"`
	UserUUID      string `json:"user_uuid"`
}

/*
pub struct GetEsimRequestInterfaceDTO {
    pub iccid: Option<String>,
    pub user_uuid: Option<String>,
}

pub struct GetEsimResponseInterfaceDTO {
    pub iccid: String,
    pub lpa: String,
    pub imsi: Option<String>,
    pub msisdn: Option<String>,


    pub user_uuid: String,
    pub vendor_id: i16,
}
*/
