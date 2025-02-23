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

	VoyageModel string

	VoyageRerank string

	VoyageInputType string

	OutputDtype string

	EncodingFormat string
)

const (
	BaseURL = "https://api.voyageai.com"

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

	VoyageRerank2     VoyageRerank = "rerank-2"
	VoyageRerank2Lite VoyageRerank = "rerank-2-lite"

	VoyageInputTypeQuery    VoyageInputType = "query"
	VoyageInputTypeDocument VoyageInputType = "document"

	OutputDtypeFloat   OutputDtype = "float"
	OutputDtypeInt8    OutputDtype = "int8"
	OutputDtypeUint8   OutputDtype = "uint8"
	OutputDtypeBinary  OutputDtype = "binary"
	OutputDtypeUbinary OutputDtype = "ubinary"

	EncodingFormatBase64 EncodingFormat = "base64"
)

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
