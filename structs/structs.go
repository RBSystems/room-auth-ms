package structs

type WSO2CredentialResponse struct {
	Values []WSO2CredentialPerson `json:"values"`
}

type WSO2Field struct {
	Value       string `json:"value"`
	Description string `json:"description"`
	APIType     string `json:"api_type"`
}

type WSO2CredentialPerson struct {
	Basic struct {
		UpdatedDatetime      WSO2Field `json:"updated_datetime"`
		UpdatedByByuID       WSO2Field `json:"updated_by_byu_id"`
		CreatedDatetime      WSO2Field `json:"created_datetime"`
		CreatedByByuID       WSO2Field `json:"created_by_byu_id"`
		ByuID                WSO2Field `json:"byu_id"`
		PersonID             WSO2Field `json:"person_id"`
		NetID                WSO2Field `json:"net_id"`
		PersonalEmailAddress WSO2Field `json:"personal_email_address"`
		PrimaryPhoneNumber   WSO2Field `json:"primary_phone_number"`
		//Deceased             WSO2Field `json:"deceased"`
		Sex                WSO2Field `json:"sex"`
		FirstName          WSO2Field `json:"first_name"`
		MiddleName         WSO2Field `json:"middle_name"`
		Surname            WSO2Field `json:"surname"`
		Suffix             WSO2Field `json:"suffix"`
		PreferredFirstName WSO2Field `json:"preferred_first_name"`
		PreferredSurname   WSO2Field `json:"preferred_surname"`
		RestOfName         WSO2Field `json:"rest_of_name"`
		NameLnf            WSO2Field `json:"name_lnf"`
		NameFnf            WSO2Field `json:"name_fnf"`
		PreferredName      WSO2Field `json:"preferred_name"`
		HomeTown           WSO2Field `json:"home_town"`
		HomeStateCode      WSO2Field `json:"home_state_code"`
		HomeCountryCode    WSO2Field `json:"home_country_code"`
		//MergeInProcess       WSO2Field `json:"merge_in_process"`
	} `json:"basic"`
}
