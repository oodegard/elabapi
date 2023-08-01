package elabapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ApiTest() {
	fmt.Println("This is a test")
}

// getSamplesID retrieves samples from the eLabJournal API.
// If the sampleTypeID argument is not nil, the function will only retrieve samples with the specified sample type ID.
// If the sampleTypeID argument is nil, the function will retrieve all samples.

// Functions that work with samples

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
	body, err := ioutil.ReadAll(resp.Body)
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
	body, err := ioutil.ReadAll(resp.Body)
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
	body, err := ioutil.ReadAll(resp.Body)
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
	body, err := ioutil.ReadAll(resp.Body)
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return 0, err
	}
	fmt.Println("Response body:", string(body))
	var result int32
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error unmarshaling response body:", err)
		return 0, err
	}
	fmt.Println("New sample ID:", result)
	return result, nil
}

// Functions that work with the ELAB journal

func GetExperiments(apiToken string) ([]map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://uio.elabjournal.com/api/v1/experiments", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
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

func PostExperiment(apiToken string) error {
	client := &http.Client{}
	data := map[string]string{"name": "New Experiment"}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://uio.elabjournal.com/api/v1/experiments", bytes.NewBuffer(payload))
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

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	return nil
}
