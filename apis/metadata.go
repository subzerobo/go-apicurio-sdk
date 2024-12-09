package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/subzerobo/go-apicurio-sdk/client"
	"github.com/subzerobo/go-apicurio-sdk/models"
	"net/http"
)

// MetadataAPI handles metadata-related operations for artifacts.
type MetadataAPI struct {
	Client *client.Client
}

// NewMetadataAPI creates a new MetadataAPI instance.
func NewMetadataAPI(client *client.Client) *MetadataAPI {
	return &MetadataAPI{
		Client: client,
	}
}

// GetArtifactVersionMetadata retrieves metadata for a single artifact version.
func (api *MetadataAPI) GetArtifactVersionMetadata(ctx context.Context, groupId, artifactId, versionExpression string) (*models.ArtifactVersionMetadata, error) {
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions/%s", api.Client.BaseURL, groupId, artifactId, versionExpression)

	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var metadata models.ArtifactVersionMetadata
	if err := handleResponse(resp, http.StatusOK, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// UpdateArtifactVersionMetadata updates the user-editable metadata of an artifact version.
func (api *MetadataAPI) UpdateArtifactVersionMetadata(ctx context.Context, groupId, artifactId, versionExpression string, metadata models.UpdateArtifactMetadataRequest) error {
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s/versions/%s", api.Client.BaseURL, groupId, artifactId, versionExpression)

	resp, err := api.executeRequest(ctx, http.MethodPut, url, metadata)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// GetArtifactMetadata retrieves metadata for an artifact based on the latest version or the next available non-disabled version.
func (api *MetadataAPI) GetArtifactMetadata(ctx context.Context, groupId, artifactId string) (*models.ArtifactMetadata, error) {
	// Construct the URL
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s", api.Client.BaseURL, groupId, artifactId)

	resp, err := api.executeRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var metadata models.ArtifactMetadata
	if err := handleResponse(resp, http.StatusOK, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// UpdateArtifactMetadata updates the editable parts of an artifact's metadata.
func (api *MetadataAPI) UpdateArtifactMetadata(ctx context.Context, groupId, artifactId string, metadata models.UpdateArtifactMetadataRequest) error {
	// Construct the URL
	url := fmt.Sprintf("%s/groups/%s/artifacts/%s", api.Client.BaseURL, groupId, artifactId)

	resp, err := api.executeRequest(ctx, http.MethodPut, url, metadata)
	if err != nil {
		return err
	}

	return handleResponse(resp, http.StatusNoContent, nil)
}

// executeRequest executes an HTTP request with the given method, URL, and body.
func (api *MetadataAPI) executeRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal request body")
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request")
	}

	return resp, nil
}

// handleResponse reads the response body and checks the status code.
func handleResponse(resp *http.Response, expectedStatus int, result interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return errors.Wrap(parseErr, "unexpected server error")
		}
		return apiError
	}

	if result != nil && resp.StatusCode == expectedStatus {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return errors.Wrap(err, "failed to parse response body")
		}
	}

	return nil
}
