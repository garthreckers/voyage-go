# Voyage Go

This is a work-in-progress library for Voyage, written in Go.

# Client
`voyage.NewClient` accepts a `voyage.VoyageConfig{}` parameter

```go
type VoyageConfig struct {
    // By default, APIKey uses the VOYAGE_API_KEY environment variable
    APIKey string
    
    // If you need to modify or mock for testing, you can set the client and/or host in the VoyageConfig.
    Client *http.Client
    Host   string
}

// Create a new client
voy := voyage.NewClient(&voyage.VoyageConfig{})
```

# Basic Examples

## Create an embedding
```go
voy := voyage.NewClient(&voyage.VoyageConfig{})

// Prepare the embedding request
data := []string{
    "Document to embed",
    "Another document to embed",
}

req := voyage.EmbedRequest{
    Input:     data,
    Model:     voyage.VoyageModel3,
    InputType: voyage.P(voyage.VoyageInputTypeDocument),
}

// Get the embedding response
res, err := voy.Embed(req)
if err != nil {
    // Handle error
}
```

## Query an embedding

```go
voy := voyage.NewClient(&voyage.VoyageConfig{})

// Prepare the embedding request
data := []string{
    "Search query",
}

req := voyage.EmbedRequest{
    Input:     data,
    Model:     voyage.VoyageModel3,
    InputType: voyage.P(voyage.VoyageInputTypeQuery),
}

// Get the embedding response
res, err := voy.Embed(req)
if err != nil {
    // Handle error
}
```