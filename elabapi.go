package elabapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tealeg/xlsx"
)

// SampleMeta structure to match JSON request
type SampleMeta struct {
	SampleTypeMetaID int    `json:"sampleTypeMetaID"`
	SampleDataType   string `json:"sampleDataType"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	ChemicalFile     string `json:"chemicalFile"`
}

type SampleData struct {
	SampleMetas    []SampleMeta `json:"sampleMetas"`
	SampleTypeID   int          `json:"sampleTypeID"`
	CheckedOut     bool         `json:"checkedOut"`
	ParentSampleID int          `json:"parentSampleID"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Note           string       `json:"note"`
	AltID          string       `json:"altID"`
	StorageLayerID int          `json:"storageLayerID"`
	Position       int          `json:"position"`
}

type Metadata struct {
	SampleTypeMetaID int    `json:"sampleTypeMetaID"`
	Key              string `json:"key"`
	SampleDataType   string `json:"sampleDataType"`
}

func ApiTest() {
	fmt.Println("This is a test")
}

// getSamplesID retrieves samples from the eLabJournal API.
// If the sampleTypeID argument is not nil, the function will only retrieve samples with the specified sample type ID.
// If the sampleTypeID argument is nil, the function will retrieve all samples.

// Functions that work with samples
type APIError struct {
	Message string
	Errors  []string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %s", e.Message)
}

func GetSamples(apiToken string, sampleTypeID *string) ([]map[string]interface{}, error) {
	client := &http.Client{}
	var url string
	if sampleTypeID != nil {
		url = fmt.Sprintf("https://uio.elabjournal.com/api/v1/samples?sampleTypeID=%s", *sampleTypeID)
	} else {
		url = "https://uio.elabjournal.com/api/v1/samples"
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	data := result["data"].([]interface{})
	samples := make([]map[string]interface{}, len(data))
	for i, sample := range data {
		samples[i] = sample.(map[string]interface{})
	}
	//fmt.Printf("samples: %v\n", samples)
	return samples, nil
}

func GetSampleTypes(apiToken string) ([]map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://uio.elabjournal.com/api/v1/sampleTypes", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	data := result["data"].([]interface{})
	sampleTypes := make([]map[string]interface{}, len(data))
	for i, sampleType := range data {
		sampleTypes[i] = sampleType.(map[string]interface{})
	}

	return sampleTypes, nil
}

func GetSampleByID(apiToken string, sampleID int32) (map[string]interface{}, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://uio.elabjournal.com/api/v1/samples/%d", sampleID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetSampleMeta(apiToken string, sampleID int) ([]map[string]interface{}, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://uio.elabjournal.com/api/v1/samples/%d/meta", sampleID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	data := result["data"].([]interface{})
	metaFields := make([]map[string]interface{}, len(data))
	for i, metaField := range data {
		metaFields[i] = metaField.(map[string]interface{})
	}
	// fmt.Printf("metaFields: %v\n", metaFields)
	return metaFields, nil
}
func PostSample(apiToken string, sample SampleData) error {
	client := &http.Client{}
	payload, err := json.Marshal(sample)
	if err != nil {
		return err
	}

	// Add the autoCreateMetaDefaults=true query parameter
	url := "https://uio.elabjournal.com/api/v1/samples?autoCreateMetaDefaults=true"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", apiToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s, body: %s", resp.Status, string(body))
	}

	fmt.Printf("POST request successful: %s\n", string(body))
	return nil
}

// THis is now deprecated and replaced by PostSample above
/*
func PostSample(apiToken string, sample map[string]interface{}) (int32, error) {
	fmt.Println("Creating new sample...")
	client := &http.Client{}
	url := "https://uio.elabjournal.com/api/v1/samples?autoCreateMetaDefaults=true"
	sampleJSON, err := json.Marshal(sample)
	if err != nil {
		fmt.Println("Error marshaling sample:", err)
		return 0, err
	}
	fmt.Println("Sample JSON:", string(sampleJSON))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(sampleJSON))
	if err != nil {
		fmt.Println("Error creating new request:", err)
		return 0, err
	}
	req.Header.Add("Authorization", apiToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return 0, err
	}
	fmt.Println("Response body:", string(body))
	var result int32
	err = json.Unmarshal(body, &result)
	if err != nil {
		var apiErr APIError
		if json.Unmarshal(body, &apiErr) == nil {
			return 0, &apiErr
		}
		fmt.Println("Error unmarshaling response body:", err)
		return 0, err
	}
	fmt.Println("New sample ID:", result)
	return result, nil
}
*/

// GetExperiments retrieves a list of experiments from the ELAB journal API with optional filters
func GetExperiments(apiToken string, filters map[string]string) ([]map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://uio.elabjournal.com/api/v1/experiments", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", apiToken)

	// Add optional filters as query parameters
	q := req.URL.Query()
	for k, v := range filters {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	data := result["data"].([]interface{})
	experiments := make([]map[string]interface{}, len(data))
	for i, experiment := range data {
		experiments[i] = experiment.(map[string]interface{})
	}

	return experiments, nil
}

func PostExperiment(apiToken string, experiment map[string]interface{}) (int32, error) {
	client := &http.Client{}
	payload, err := json.Marshal(experiment)
	// fmt.Printf("payload: %v\n", payload)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", "https://uio.elabjournal.com/api/v1/experiments", bytes.NewBuffer(payload))
	if err != nil {
		return 0, err
	}

	req.Header.Add("Authorization", apiToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result int32
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func PostSection(apiToken string, experimentID int32, section map[string]interface{}) (int32, error) {
	client := &http.Client{}
	payload, err := json.Marshal(section)
	if err != nil {
		return 0, err
	}

	url := fmt.Sprintf("https://uio.elabjournal.com/api/v1/experiments/%d/sections", experimentID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return 0, err
	}

	req.Header.Add("Authorization", apiToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result int32 // This is the expJournalID value
	err = json.Unmarshal(body, &result)
	if err != nil {
		// Check if the error is due to the response body not being an int32 value
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return 0, fmt.Errorf("unexpected response format: %s", string(body))
		}
		return 0, err
	}

	return result, nil
}
func GetExperimentSections(apiToken string, experimentID int32, filters map[string]string) ([]map[string]interface{}, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://uio.elabjournal.com/api/v1/experiments/%d/sections", experimentID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", apiToken)

	// Add optional filters as query parameters
	q := req.URL.Query()
	for k, v := range filters {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	data := result["data"].([]interface{})
	sections := make([]map[string]interface{}, len(data))
	for i, section := range data {
		sections[i] = section.(map[string]interface{})
	}
	return sections, nil
}
func GetExperimentSectionExcel(apiToken string, expJournalID int32) ([][]string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://uio.elabjournal.com/api/v1/experiments/sections/%d/excel", expJournalID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", apiToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected server response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	file, err := xlsx.OpenBinary(body)
	if err != nil {
		return nil, err
	}

	var excelData [][]string
	for _, sheet := range file.Sheets {
		for _, row := range sheet.Rows {
			var vals []string
			for _, cell := range row.Cells {
				vals = append(vals, cell.String())
			}
			excelData = append(excelData, vals)
		}
	}

	return excelData, nil
}

// GetExperimentSectionHTML retrieves the HTML of an experiment section from the ELAB journal API
func GetExperimentSectionHTML(apiToken string, expJournalID string) (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://uio.elabjournal.com/api/v1/experiments/sections/%s/html", expJournalID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", apiToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetExpTextSectionContent retrieves the content of an experiment text section from the ELAB journal
func GetExpTextSectionContent(apiToken string, expJournalID int32) (map[string]interface{}, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://uio.elabjournal.com/api/v1/experiments/sections/%d/content", expJournalID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetExperimentSamples(apiToken string, expJournalID int32, filters map[string]string) ([]map[string]interface{}, error) {
	/*
		This function will send a GET request to the specified URL and return a list of samples from a SAMPLESIN or SAMPLESOUT section.
		The filters parameter allows you to add optional query parameters to the request.
		The function will return an error if it fails to send the request or parse the response.
	*/

	client := &http.Client{}
	url := fmt.Sprintf("https://uio.elabjournal.com/api/v1/experiments/sections/%d/samples", expJournalID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", apiToken)

	// Add optional filters as query parameters
	q := req.URL.Query()
	for k, v := range filters {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	data := result["data"].([]interface{})
	samples := make([]map[string]interface{}, len(data))
	for i, sample := range data {
		samples[i] = sample.(map[string]interface{})
	}
	return samples, nil
}

// UpdateExpTextSectionContent updates the content of an experiment text section in the ELAB journal
func UpdateExperimentSection(apiToken string, expJournalID int32, data map[string]interface{}) error {
	client := &http.Client{}
	url := fmt.Sprintf("https://uio.elabjournal.com/api/v1/experiments/sections/%d/content", expJournalID)
	payloadJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", apiToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// UpdateExperimentSectionHTML updates the HTML content of an experiment text section in the ELAB journal
func UpdateExperimentSectionHTML(apiToken string, expJournalID int32, htmlContent string) error {
	client := &http.Client{}
	url := fmt.Sprintf("https://uio.elabjournal.com/api/v1/experiments/sections/%d/html", expJournalID)
	req, err := http.NewRequest("PUT", url, bytes.NewBufferString(htmlContent))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// ListFiles retrieves a list of files from the ELAB journal API with optional filters
func ListFiles(apiToken string, filters map[string]string) (map[string]interface{}, error) {
	client := &http.Client{}
	url := "https://uio.elabjournal.com/api/v1/files"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", apiToken)

	// Add optional filters as query parameters
	q := req.URL.Query()
	for k, v := range filters {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Studies
// ListStudies retrieves a list of studies from the ELAB journal API with optional filters
func ListStudies(apiToken string, filters map[string]string) (map[string]interface{}, error) {
	client := &http.Client{}
	url := "https://uio.elabjournal.com/api/v1/studies"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", apiToken)

	// Add optional filters as query parameters
	q := req.URL.Query()
	for k, v := range filters {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/* Example use
# Use without filters
files, err := ListFiles(apiToken, map[string]string{})
if err != nil {
    // handle error
}


# Use with filters
filters := map[string]string{
    "fileName": "example.txt",
    "userID":   "12345",
}
files, err := ListFiles(apiToken, filters)
if err != nil {
    // handle error
}
*/

// ListSampleTypes retrieves a list of sample types from the eLab Journal API
func ListSampleTypes(apiToken string, filters map[string]string) (map[string]interface{}, error) {
	client := &http.Client{}
	url := "https://uio.elabjournal.com/api/v1/sampleTypes"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", apiToken)

	// Add optional filters as query parameters
	q := req.URL.Query()
	for k, v := range filters {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %v", err)
	}

	return result, nil
}

func FetchSampleTypeMetaDetails(apiToken, sampleTypeID string) (map[string]interface{}, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://uio.elabjournal.com/api/v1/sampleTypes/%s/meta", sampleTypeID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %v", err)
	}

	return result, nil
}
