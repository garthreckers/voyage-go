package voyage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	RerankRequest struct {
		// The query as a string. The query can contain a maximum of 4000
		// tokens for rerank-2, 2000 tokens for rerank-2-lite and rerank-1,
		// and 1000 tokens for rerank-lite-1.
		//
		// Required
		Query string `json:"query"`

		// The documents to be reranked as a list of strings.
		//
		// - The number of documents cannot exceed 1000.
		// - The sum of the number of tokens in the query and the number of
		//   tokens in any single document cannot exceed 16000 for rerank-2;
		//   8000 for rerank-2-lite and rerank-1; and 4000 for rerank-lite-1.
		// - The total number of tokens, defined as "the number of query tokens
		//   × the number of documents + sum of the number of tokens in all
		//   documents", cannot exceed 600K for rerank-2 and rerank-2-lite,
		//   and 300K for rerank-1 and rerank-lite-1. Please see our FAQ.
		//
		// Required
		Documents []string `json:"documents"`

		// Name of the model
		//
		//  - VoyageRerank2
		//  - VoyageRerank2Lite
		//
		// Required
		Model VoyageRerank `json:"model"`

		// The number of most relevant documents to return. If not specified,
		// the reranking results of all documents will be returned.
		//
		// Optional
		TopK *uint `json:"top_k,omitempty"`

		// Whether to return the documents in the response. Defaults to false.
		//
		// - If false, the API will return a list of {"index", "relevance_score"}
		//   where "index" refers to the index of a document within the input list.
		//
		// - If true, the API will return a list of {"index", "document",
		//   "relevance_score"} where "document" is the corresponding document
		//   from the input list.
		//
		// Optional
		ReturnDocuments bool `json:"return_documents"`

		// Whether to truncate the input to satisfy the "context length limit" on
		// the query and the documents. Defaults to true.
		//
		// If true, the query and documents will be truncated to fit within the
		// context length limit, before processed by the reranker model.
		//
		// If false, an error will be raised when the query exceeds 4000 tokens
		// for rerank-2; 2000 tokens rerank-2-lite and rerank-1; and 1000 tokens
		// for rerank-lite-1, or the sum of the number of tokens in the query and
		// the number of tokens in any single document exceeds 16000 for rerank-2;
		// 8000 for rerank-2-lite and rerank-1; and 4000 for rerank-lite-1.
		//
		// Optional
		Truncate *bool `json:"truncate,omitempty"`
	}

	RerankResponse struct {
		// The object type, which is always "list"
		Object string `json:"object"`

		// An array of embedding objects
		Data []RerankResponseData `json:"data"`

		// Model is the name of the current model
		Model VoyageModel `json:"model"`

		Usage RerankResponseUsage `json:"usage"`
	}

	RerankResponseData struct {
		// The document string. Only returned when return_documents is set to true.
		Document string `json:"document"`

		// The relevance score of the document with respect to the query.
		RelevanceScore float64 `json:"relevance_score"`

		// The index of the document in the input list.
		Index uint `json:"index"`
	}

	RerankResponseUsage struct {
		// The total number of tokens used for computing the reranking
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
