package voyage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	EmbedRequest struct {
		Input           []string         `json:"input"`
		Model           VoyageModel      `json:"model"`
		InputType       *VoyageInputType `json:"input_type,omitempty"`
		Truncate        *bool            `json:"truncation,omitempty"`
		OutputDimension *int             `json:"output_dimension,omitempty"`
		OutputDtype     *OutputDtype     `json:"output_dtype,omitempty"`
		EncodingFormat  *string          `json:"encoding_format,omitempty"`
	}

	EmbedResponse struct {
		Object string              `json:"object"`
		Data   []EmbedResponseData `json:"data"`
		Model  VoyageModel         `json:"model"`
		Usage  EmbedResponseUsage  `json:"usage"`
	}

	EmbedResponseData struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     uint      `json:"index"`
	}

	EmbedResponseUsage struct {
		TotalTokens uint `json:"total_tokens"`
	}
)

// IsValid checks if the Embed struct is valid
func (e *EmbedRequest) IsValid() error {
	if len(e.Input) == 0 {
		return ErrInputRequired
	}

	if len(e.Input) > 128 {
		return ErrInputTooLarge
	}

	if e.Model == "" {
		return ErrModelRequired
	}

	return nil
}

// Embed generates embeddings for the given input text
func (v *voyage) Embed(embedRequest EmbedRequest) (*EmbedResponse, error) {
	if err := embedRequest.IsValid(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v1/embeddings", v.Host)

	jsonReq, err := json.Marshal(embedRequest)
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

	var embedResponse EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResponse); err != nil {
		return nil, err
	}

	return &embedResponse, nil
}
