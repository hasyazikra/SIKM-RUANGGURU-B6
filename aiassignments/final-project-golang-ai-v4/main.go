package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type AIModelConnector struct {
	Client *http.Client
}

type Inputs struct {
	Table map[string][]string `json:"table"`
	Query string              `json:"query"`
}

type Response struct {
	Answer      string   `json:"answer"`
	Coordinates [][]int  `json:"coordinates"`
	Cells       []string `json:"cells"`
	Aggregator  string   `json:"aggregator"`
}

func CsvToSlice(data string) (map[string][]string, error) {
	r := csv.NewReader(strings.NewReader(data))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	headers := records[0]
	for _, header := range headers {
		result[header] = make([]string, 0)
	}

	for _, record := range records[1:] {
		for i, value := range record {
			result[headers[i]] = append(result[headers[i]], value)
		}
	}

	return result, nil
} //TO DO DONE

func (c *AIModelConnector) ConnectAIModel(payload interface{}, token string) (Response, error) {
	url := "http://your-ai-model-endpoint"
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return Response{}, err
	}

	return response, nil
} // TO DO DONE

func main() {
	// Read CSV file
	fileData, err := ioutil.ReadFile("data-series.csv")
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	// Convert CSV data to map of string slices
	inputs, err := CsvToSlice(string(fileData))
	if err != nil {
		fmt.Println("Error converting CSV:", err)
		return
	}

	// Example payload to AI model
	aiPayload := Inputs{
		Table: inputs,
		Query: "predict_energy_consumption",
	}

	// Initialize AIModelConnector with http.Client
	connector := AIModelConnector{
		Client: &http.Client{},
	}

	// Example token (replace with actual token retrieval)
	token := "your-auth-token"

	// Connect to AI model
	response, err := connector.ConnectAIModel(aiPayload, token)
	if err != nil {
		fmt.Println("Error connecting to AI model:", err)
		return
	}

	// Print AI model response
	fmt.Println("AI Model Response:")
	fmt.Println("Answer:", response.Answer)
	fmt.Println("Coordinates:", response.Coordinates)
	fmt.Println("Cells:", response.Cells)
	fmt.Println("Aggregator:", response.Aggregator)
} // TO DO DONE
