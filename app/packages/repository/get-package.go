package repository

import (
	"fmt"
	"sample-service/app/common"
	"sample-service/app/infra"
	"sample-service/app/model"
	"sample-service/app/packages/dto"
	"strconv"
	"strings"
)

type PackageWithAreaCount struct {
	model.Package        // Include all Package props
	CountryList   string `gorm:"column:country_list"`
}

func GetPackagesWithAreaCount(query dto.GetPackageRequestDTO) ([]PackageWithAreaCount, error) {

	var idCondition string
	var hiddenCondition bool
	var regionCodeNameCondition string

	fmt.Printf("Query: %v", query)

	id := common.GetStr(query.Id)
	showHidden := common.GetBool(query.ShowHidden)
	regionCodeName := common.ClearQuery(common.GetStr(query.RegionCodeName))
	regionName := common.ClearQuery(common.GetStr(query.RegionName))
	vendorIds := common.GetArr(query.VendorIds)

	if id != "" {

		packageId, err := common.GetPackageFromGlobalId(id)
		fmt.Printf("PARAMS: %d", packageId)

		if err == nil {
			idCondition = "packages.vendor_id = " + strconv.Itoa(int(packageId))
		}
	}

	if showHidden == false {
		hiddenCondition = false
	}

	switch true {
	case regionCodeName != "":
		regionCodeNameCondition = " AND areas.iso = '" + regionCodeName + "'"
	case regionName != "":
		regionCodeNameCondition = "  AND areas.name = '" + regionName + "'"
	default:
		regionCodeNameCondition = ""
	}

	// WHERE
	where := ""

	conditions := []string{}

	// mock for empty conditions
	conditions = append(conditions, "1=1")
	//

	if idCondition != "" {
		conditions = append(conditions, idCondition)
	}
	if hiddenCondition {
		conditions = append(conditions, "hidden = false")
	}
	if len(vendorIds) > 0 {
		vendorIdsStr := ConvertSliceToString(vendorIds)
		conditions = append(conditions, "vendor_id IN("+vendorIdsStr+")")
	}

	if regionCodeNameCondition != "" {
		conditions = append(conditions, ` 
		(
			SELECT COUNT(areas.id) 
			FROM package_areas 
			INNER JOIN areas ON (areas.id = package_areas.area_id) 
			WHERE package_areas.package_id = packages.id `+regionCodeNameCondition+`
		) != 0`)
	}

	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ") + " "
	}

	queryStr := `
		SELECT 
		packages.*,
		(
			SELECT STRING_AGG(CONCAT(areas.name, ',', areas.region, ',', areas.iso), ';') 
			FROM package_areas 
			INNER JOIN areas ON (areas.id = package_areas.area_id) 
			WHERE package_areas.package_id = packages.id
		) AS country_list
		FROM packages ` + where + `
		ORDER BY packages.vendor_id ASC
		`

	db := infra.DB
	db = db.Debug()
	// Search
	var results []PackageWithAreaCount
	if err := db.Raw(queryStr).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// -----------------------------------------------------------------------------
func ConvertSliceToString(slice []int32) string {
	strSlice := make([]string, len(slice))
	for i, num := range slice {
		strSlice[i] = strconv.FormatInt(int64(num), 10)
	}
	return strings.Join(strSlice, ",")
}

func GetRegionAndCountries(countryList string) (dto.RegionDTO, error) {

	if countryList == "" {
		return dto.RegionDTO{}, fmt.Errorf("Empty country list")
	}

	var region dto.RegionDTO
	region.Name = ""

	var countries []dto.CountryDTO

	arrCountries := strings.Split(countryList, ";")

	for _, strCountry := range arrCountries {
		arrCountry := strings.Split(strCountry, ",")

		if arrCountry[0] == "" || arrCountry[0] == "null" {
			continue
		}

		if arrCountry[1] != "null" && arrCountry[2] != "null" {
			country := dto.CountryDTO{
				CodeName: arrCountry[2],
				Name:     arrCountry[0],
			}
			countries = append(countries, country)
		} else {
			region.Name = arrCountry[0]
		}
	}

	region.CodeName = ""
	region.IncludedCountriesAmount = int8(len(countries))
	region.Countries = countries

	if region.IncludedCountriesAmount == 1 && region.Name == "" {
		region.Name = countries[0].Name
		region.CodeName = countries[0].CodeName
	}

	return region, nil

}
