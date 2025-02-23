package voyage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	EmbedRequest struct {
		// Input single text string, or a list of texts as a list of strings, such as
		// ["I like cats", "I also like dogs"].
		//
		// Currently, we have two constraints on the list:
		//  - The maximum length of the list is 128.
		//  - The total number of tokens in the list is at most 1M for voyage-3-lite;
		//    320K for voyage-3 and voyage-2; and 120K for voyage-3-large,
		//    voyage-code-3, voyage-large-2-instruct, voyage-finance-2,
		//    voyage-multilingual-2, voyage-law-2, and voyage-large-2.
		//
		// Check text embedding limits here: https://docs.voyageai.com/docs/embeddings
		//
		// Required
		Input []string `json:"input"`

		// Name of the model
		//
		// Recommend Options:
		//  - VoyageModel3Large
		//  - VoyageModel3
		//  - VoyageModel3Lite
		//  - VoyageModelCode3
		//  - VoyageModelFinance2
		//  - VoyageModelLaw2
		//
		// Older Models:
		//  - VoyageModelMultilingual2
		//  - VoyageModelLarge2Instruct
		//  - VoyageModelLarge2
		//  - VoyageModel2
		//  - VoyageModelLite02Instruct
		//  - VoyageModel02
		//  - VoyageModel01
		//  - VoyageModelLite01
		//  - VoyageModelLite01Instruct
		//
		// Required
		Model VoyageModel `json:"model"`

		// Type of the input text. Defaults to nil.
		//
		// Other options: query, document.
		//  - VoyageInputTypeQuery
		//  - VoyageInputTypeDocument
		//
		// When input_type is nil, the embedding model directly converts the inputs
		// (texts) into numerical vectors. For retrieval/search purposes, where a
		// "query" is used to search for relevant information among a collection of
		// data referred to as "documents," we recommend specifying whether your
		// inputs (texts) are intended as queries or documents by setting input_type
		// to query or document, respectively. In these cases, Voyage automatically
		// prepends a prompt to your inputs before vectorizing them, creating
		// vectors more tailored for retrieval/search tasks. Embeddings generated
		// with and without the input_type argument are compatible.
		//
		//  - For transparency, the following prompts are prepended to your input.
		//     - For VoyageInputTypeQuery, the prompt is "Represent the query for
		//       retrieving supporting documents: ".
		//     - For VoyageInputTypeDocument, the prompt is "Represent the document
		//       for retrieval: ".
		//
		// Optional
		InputType *VoyageInputType `json:"input_type,omitempty"`

		// Whether to truncate the input texts to fit within the context length.
		//
		// Defaults to true.
		//
		//  - If true, an over-length input texts will be truncated to fit within
		//    the context length, before vectorized by the embedding model.
		//  - If false, an error will be raised if any given text exceeds the
		//    context length.
		Truncate *bool `json:"truncation,omitempty"`

		// The number of dimensions for resulting output embeddings.
		//
		// Defaults to nil.
		//
		//  - Most models only support a single default dimension, used when
		//    output_dimension is set to null (see output embedding dimensions
		//    here: https://docs.voyageai.com/docs/embeddings).
		//  - voyage-3-large and voyage-code-3 support the following
		//    output_dimension, values: 2048, 1024 (default), 512, and 256.
		OutputDimension *int `json:"output_dimension,omitempty"`

		// The data type for the embeddings to be returned.
		//
		// Defaults to OutputDtypeFloat.
		//
		// Other options:
		//  - OutputDtypeInt8
		//  - OutputDtypeUint8
		//  - OutputDtypeBinary
		//  - OutputDtypeUbinary.
		//
		// OutputDtypeFloat is supported for all models. OutputDtypeInt8,
		// OutputDtypeUint8, OutputDtypeBinary, and OutputDtypeUbinary are
		// supported by voyage-3-large and voyage-code-3.
		//
		// Please see our guide for more details about output data types.
		//
		// https://docs.voyageai.com/docs/flexible-dimensions-and-quantization#quantization
		//
		// OutputDtypeFloat: Each returned embedding is a list of 32-bit (4-byte)
		// single-precision floating-point numbers. This is the default and provides
		// the highest precision / retrieval accuracy.
		//
		// OutputDtypeInt8 and OutputDtypeUint8: Each returned embedding is a list
		// of 8-bit (1-byte) integers ranging from -128 to 127 and 0 to 255,
		// respectively.
		//
		// OutputDtypeBinary and OutputDtypeUbinary: Each returned embedding is a
		// list of 8-bit integers that represent bit-packed, quantized single-bit
		// embedding values: int8 for OutputDtypeBinary and OutputDtypeUint8 for
		// OutputDtypeUbinary. The length of the returned list of integers is 1/8
		// of output_dimension (which is the actual dimension of the embedding).
		// The OutputDtypeBinary type uses the offset binary method.
		//
		// Please refer to our guide for details on offset binary (https://docs.voyageai.com/docs/flexible-dimensions-and-quantization#offset-binary)
		// and binary embeddings (https://docs.voyageai.com/docs/flexible-dimensions-and-quantization#quantization).
		OutputDtype *OutputDtype `json:"output_dtype,omitempty"`

		// Format in which the embeddings are encoded.
		//
		// Defaults to nil.
		//
		// Other options: base64.
		//
		// If nil, each embedding is an array of float numbers when OutputDtype
		// is set to float and as an array of integers for all other values of
		// OutputDtype (OutputDtypeInt8, OutputDtypeUint8, OutputDtypeBinary, and
		// OutputDtypeUbinary).
		//
		// If base64, the embeddings are represented as a Base64-encoded NumPy
		// array of:
		//
		//  - Floating-point numbers (numpy.float32) for OutputDtype set to
		//    OutputDtypeFloat.
		//  - Signed integers (numpy.int8) for OutputDtype set to OutputDtypeInt8
		//    or OutputDtypeBinary.
		//  - Unsigned integers (numpy.uint8) for OutputDtype set to OutputDtypeUint8
		//    or OutputDtypeUbinary.
		EncodingFormat *string `json:"encoding_format,omitempty"`
	}

	EmbedResponse struct {
		// The object type, which is always "list"
		Object string `json:"object"`

		// An array of embedding objects
		Data []EmbedResponseData `json:"data"`

		// Model is the name of the current model
		Model VoyageModel `json:"model"`

		Usage EmbedResponseUsage `json:"usage"`
	}

	EmbedResponseData struct {
		// The object type, which is always "embedding"
		Object string `json:"object"`

		// Each embedding is a vector represented as an array of float numbers when
		// OutputDtype is set to OutputDtypeFloat and as an array of integers for
		// all other values of OutputDtype (OutputDtypeInt8, OutputDtypeUint8,
		// OutputDtypeBinary, and OutputDtypeUbinary). The length of this vector
		// varies depending on the specific model, OutputDimension, and
		// OutputDtype.
		Embedding []float32 `json:"embedding"`

		// An integer representing the index of the embedding within the list
		// of embeddings
		Index uint `json:"index"`
	}

	EmbedResponseUsage struct {
		// The total number of tokens used for computing the embeddings
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
