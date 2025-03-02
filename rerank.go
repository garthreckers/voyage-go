package voyage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	RerankRequest struct {
		Query           string       `json:"query"`
		Documents       []string     `json:"documents"`
		Model           VoyageRerank `json:"model"`
		TopK            *uint        `json:"top_k,omitempty"`
		ReturnDocuments bool         `json:"return_documents"`
		Truncate        *bool        `json:"truncate,omitempty"`
	}

	RerankResponse struct {
		Object string               `json:"object"`
		Data   []RerankResponseData `json:"data"`
		Model  VoyageModel          `json:"model"`
		Usage  RerankResponseUsage  `json:"usage"`
	}

	RerankResponseData struct {
		Document       string  `json:"document"`
		RelevanceScore float64 `json:"relevance_score"`
		Index          uint    `json:"index"`
	}

	RerankResponseUsage struct {
		TotalTokens uint `json:"total_tokens"`
	}
)

// IsValid checks if the Rerank struct is valid
func (r *RerankRequest) IsValid() error {
	if len(r.Documents) == 0 {
		return ErrDocumentsRequired
	}

	if len(r.Documents) > 1000 {
		return ErrDocumentsTooLarge
	}

	if r.Query == "" {
		return ErrQueryRequired
	}

	if r.Model == "" {
		return ErrModelRequired
	}

	return nil
}

// Rerank reranks the given documents based on the query
func (v *voyage) Rerank(rerankRequest RerankRequest) (*RerankResponse, error) {
	if err := rerankRequest.IsValid(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v1/rerank", v.Host)

	jsonReq, err := json.Marshal(rerankRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", v.APIKey))

	resp, err := v.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("resp: %+v\n", resp)

	var rerankResponse RerankResponse
	if err := json.NewDecoder(resp.Body).Decode(&rerankResponse); err != nil {
		return nil, err
	}

	return &rerankResponse, nil
}
