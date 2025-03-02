package voyage

import (
	"net/http"
	"os"
)

type (
	Voyage interface {
		Embed(EmbedRequest) (*EmbedResponse, error)
		Rerank(RerankRequest) (*RerankResponse, error)
	}

	VoyageConfig struct {
		APIKey string
		Client *http.Client
		Host   string
	}

	voyage struct {
		APIKey string
		Client *http.Client
		Host   string
	}
)

type VoyageModel string

const (
	// Recommended Models
	VoyageModel3Large   VoyageModel = "voyage-3-large"
	VoyageModel3        VoyageModel = "voyage-3"
	VoyageModel3Lite    VoyageModel = "voyage-3-lite"
	VoyageModelCode3    VoyageModel = "voyage-code-3"
	VoyageModelFinance2 VoyageModel = "voyage-finance-2"
	VoyageModelLaw2     VoyageModel = "voyage-law-2"

	// Older Models
	VoyageModelMultilingual2  VoyageModel = "voyage-multilingual-2"
	VoyageModelLarge2Instruct VoyageModel = "voyage-large-2-instruct"
	VoyageModelLarge2         VoyageModel = "voyage-large-2"
	VoyageModel2              VoyageModel = "voyage-2"
	VoyageModelLite02Instruct VoyageModel = "voyage-lite-02-instruct"
	VoyageModel02             VoyageModel = "voyage-02"
	VoyageModel01             VoyageModel = "voyage-01"
	VoyageModelLite01         VoyageModel = "voyage-lite-01"
	VoyageModelLite01Instruct VoyageModel = "voyage-lite-01-instruct"
)

type VoyageRerank string

const (
	VoyageRerank2     VoyageRerank = "rerank-2"
	VoyageRerank2Lite VoyageRerank = "rerank-2-lite"
)

type VoyageInputType string

const (
	VoyageInputTypeQuery    VoyageInputType = "query"
	VoyageInputTypeDocument VoyageInputType = "document"
)

type OutputDtype string

const (
	OutputDtypeFloat   OutputDtype = "float"
	OutputDtypeInt8    OutputDtype = "int8"
	OutputDtypeUint8   OutputDtype = "uint8"
	OutputDtypeBinary  OutputDtype = "binary"
	OutputDtypeUbinary OutputDtype = "ubinary"
)

type EncodingFormat string

const EncodingFormatBase64 EncodingFormat = "base64"

const BaseURL = "https://api.voyageai.com"

func NewClient(config *VoyageConfig) Voyage {
	v := &voyage{
		APIKey: os.Getenv("VOYAGE_API_KEY"),
		Client: &http.Client{},
		Host:   BaseURL,
	}

	if config != nil {
		if config.Client != nil {
			v.Client = config.Client
		}

		if config.Host != "" {
			v.Host = config.Host
		}

		if config.APIKey != "" {
			v.APIKey = config.APIKey
		}
	}

	return v
}
