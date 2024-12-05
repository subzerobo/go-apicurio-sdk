package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/subzerobo/go-apicurio-sdk/client"
	"github.com/subzerobo/go-apicurio-sdk/models"
	"io"
	"net/http"
	"regexp"
)

type ArtifactsAPI struct {
	Client *client.Client
}

func NewArtifactsAPI(client client.Client) *ArtifactsAPI {
	return &ArtifactsAPI{
		Client: &client,
	}
}

var (
	ErrArtifactNotFound = errors.New("artifact not found")
	ErrUnexpectedError  = errors.New("unexpected server error")
	ErrMethodNotAllowed = errors.New("method not allowed or disabled on the server")
	ErrConflict         = errors.New("artifact conflicts with existing data")
	ErrInvalidInput     = errors.New("input must be between 1 and 512 characters")

	regexValidation = regexp.MustCompile(`^.{1,512}$`)
)

// SearchArtifacts - Search for artifacts using the given filter parameters.
// Search for artifacts using the given filter parameters.
// See:
func (api *ArtifactsAPI) SearchArtifacts(ctx context.Context, params *models.SearchArtifactsParams) (*models.SearchArtifactsResponse, error) {
	queryString := ""
	if params != nil {
		queryString = params.ToQuery().Encode()
	}

	url := fmt.Sprintf("%s/artifacts?%s", api.Client.BaseURL, queryString)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, fmt.Errorf("unexpected error: %s", parseErr)
		}
		return nil, apiError
	}

	var result models.SearchArtifactsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil

}

// SearchArtifactsByContent searches for artifacts that match the provided content.
// Returns a paginated list of all artifacts with at least one version that matches the posted content.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/searchArtifactsByContent
func (api *ArtifactsAPI) SearchArtifactsByContent(ctx context.Context, content []byte, params *models.SearchArtifactsByContentParams) (*models.SearchArtifactsResponse, error) {
	// Convert params to query string
	queryString := ""
	if params != nil {
		queryString = params.ToQuery().Encode()
	}

	url := fmt.Sprintf("%s/search/artifacts?%s", api.Client.BaseURL, queryString)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, errors.Wrap(parseErr, "unexpected error")
		}
		return nil, apiError
	}
	if resp.StatusCode != http.StatusOK {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, errors.Wrap(parseErr, "unexpected error")
		}
		return nil, apiError
	}

	var result models.SearchArtifactsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListArtifactReferences Returns a list containing all the artifact references using the artifact content ID.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/referencesByContentId
func (api *ArtifactsAPI) ListArtifactReferences(ctx context.Context, contentID int64) (*[]models.ArtifactReference, error) {
	url := fmt.Sprintf("%s/ids/contentId/%d/references", api.Client.BaseURL, contentID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, errors.Wrap(parseErr, "unexpected error")
		}
		return nil, apiError
	}

	var references []models.ArtifactReference
	if err := json.NewDecoder(resp.Body).Decode(&references); err != nil {
		return nil, err
	}

	return &references, nil
}

// ListArtifactReferencesByGlobalID Returns a list containing all the artifact references using the artifact global ID.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/referencesByContentHash
func (api *ArtifactsAPI) ListArtifactReferencesByGlobalID(ctx context.Context, globalID int64, params *models.ListArtifactReferencesByGlobalIDParams) (*[]models.ArtifactReference, error) {
	query := ""
	if params != nil {
		query = params.ToQuery()
	}

	url := fmt.Sprintf("%s/ids/globalIds/%d/references%s", api.Client.BaseURL, globalID, query)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, errors.Wrap(parseErr, "unexpected error")
		}
		return nil, apiError
	}

	var references []models.ArtifactReference
	if err := json.NewDecoder(resp.Body).Decode(&references); err != nil {
		return nil, err
	}

	return &references, nil
}

// ListArtifactReferencesByHash Returns a list containing all the artifact references using the artifact content hash.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/referencesByContentHash
func (api *ArtifactsAPI) ListArtifactReferencesByHash(ctx context.Context, contentHash string) (*[]models.ArtifactReference, error) {
	url := fmt.Sprintf("%s/ids/contentHashes/%s/references", api.Client.BaseURL, contentHash)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, fmt.Errorf("unexpected error: %s", parseErr)
		}
		return nil, apiError
	}

	var references []models.ArtifactReference
	if err := json.NewDecoder(resp.Body).Decode(&references); err != nil {
		return nil, err
	}

	return &references, nil
}

// ListArtifactsInGroup lists all artifacts in a specified group.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/referencesByContentHash
func (api *ArtifactsAPI) ListArtifactsInGroup(ctx context.Context, groupID string, params *models.ListArtifactsInGroupParams) (*models.ListArtifactsResponse, error) {
	err := validateInput(groupID, "Group ID")
	if err != nil {
		return nil, err
	}

	query := ""
	if params != nil {
		query = params.ToQuery()
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts%s", api.Client.BaseURL, groupID, query)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, fmt.Errorf("unexpected error: %s", parseErr)
		}
		return nil, apiError
	}

	var result models.ListArtifactsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetArtifactContentByHash Gets the content for an artifact version in the registry using the SHA-256 hash of the content
// This content hash may be shared by multiple artifact versions in the case where the artifact versions have identical content.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/getContentByHash
func (api *ArtifactsAPI) GetArtifactContentByHash(ctx context.Context, contentHash string) (*models.ArtifactContent, error) {
	url := fmt.Sprintf("%s/ids/contentHashes/%s", api.Client.BaseURL, contentHash)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.Wrapf(ErrArtifactNotFound, "content hash: %s", contentHash)
	}

	if resp.StatusCode != http.StatusOK {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, errors.Wrap(parseErr, "unexpected error")
		}
		return nil, apiError
	}

	artifactTypeHeader := resp.Header.Get("X-Registry-ArtifactType")
	artifactType, err := models.ParseArtifactType(artifactTypeHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid artifact type in response: %s", artifactType)
	}
	// Read the content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return &models.ArtifactContent{
		Content:      string(content),
		ArtifactType: artifactType,
	}, nil
}

// GetArtifactContentByID Gets the content for an artifact version in the registry using the unique content identifier for that content
// This content ID may be shared by multiple artifact versions in the case where the artifact versions are identical.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/getContentById
func (api *ArtifactsAPI) GetArtifactContentByID(ctx context.Context, contentID int64) (*models.ArtifactContent, error) {
	url := fmt.Sprintf("%s/ids/contentIds/%d", api.Client.BaseURL, contentID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.Wrapf(ErrArtifactNotFound, "content ID: %d", contentID)
	}

	if resp.StatusCode != http.StatusOK {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, errors.Wrap(parseErr, "unexpected error")
		}
		return nil, apiError
	}

	// Get and parse the artifact type
	artifactTypeHeader := resp.Header.Get("X-Registry-ArtifactType")
	artifactType, err := models.ParseArtifactType(artifactTypeHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "artifact type: %s", artifactTypeHeader)
	}

	// Read the content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return &models.ArtifactContent{
		Content:      string(content),
		ArtifactType: artifactType,
	}, nil
}

// DeleteArtifactsInGroup deletes all artifacts in a given group.
// Deletes all the artifacts that exist in a given group.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/deleteArtifactsInGroup
func (api *ArtifactsAPI) DeleteArtifactsInGroup(ctx context.Context, groupID string) error {
	err := validateInput(groupID, "Group ID")
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/groups/%s/artifacts", api.Client.BaseURL, groupID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create HTTP request")
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to execute HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusInternalServerError {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return errors.Wrap(parseErr, "unexpected server error")
		}
		return apiError
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(ErrUnexpectedError, "unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// DeleteArtifact deletes a specific artifact identified by groupId and artifactId.
// Deletes an artifact completely, resulting in all versions of the artifact also being deleted. This may fail for one of the following reasons:
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/deleteArtifact
func (api *ArtifactsAPI) DeleteArtifact(ctx context.Context, groupID, artifactId string) error {
	err := validateInput(artifactId, "Artifact ID")
	if err != nil {
		return err
	}

	err = validateInput(groupID, "Group ID")
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/groups/%s/artifacts/%s", api.Client.BaseURL, groupID, artifactId)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create HTTP request")
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to execute HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errors.Wrapf(ErrArtifactNotFound, "artifact ID: %s in group: %s", artifactId, groupID)
	}

	if resp.StatusCode == http.StatusMethodNotAllowed {
		return errors.Wrapf(ErrMethodNotAllowed, "method not allowed for artifact ID: %s in group: %s", artifactId, groupID)
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return errors.Wrap(parseErr, "unexpected server error")
		}
		return apiError
	}

	return nil
}

// CreateArtifact Creates a new artifact.
// See: https://www.apicur.io/registry/docs/apicurio-registry/3.0.x/assets-attachments/registry-rest-api.htm#tag/Artifacts/operation/createArtifact
func (api *ArtifactsAPI) CreateArtifact(ctx context.Context, groupId string, artifact models.Artifact, params models.CreateArtifactParams) (*models.ArtifactResponse, error) {
	// Build query string
	query := params.ToQuery()

	err := validateInput(artifact.Name, "Artifact ID")
	if err != nil {
		return nil, err
	}

	err = validateInput(groupId, "Group ID")
	if err != nil {
		return nil, err
	}

	// Construct URL
	url := fmt.Sprintf("%s/groups/%s/artifacts%s", api.Client.BaseURL, groupId, query)
	body, err := json.Marshal(artifact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal artifact request body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return nil, ErrConflict
	}

	if resp.StatusCode == http.StatusBadRequest {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, errors.Wrap(parseErr, "invalid request data")
		}
		return nil, apiError
	}

	if resp.StatusCode != http.StatusOK {
		apiError, parseErr := parseAPIError(resp)
		if parseErr != nil {
			return nil, errors.Wrap(parseErr, "unexpected server error")
		}
		return nil, apiError
	}

	// Parse response
	var response models.ArtifactResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to parse response body")
	}

	return &response, nil
}

func validateInput(input, name string) error {
	if match := regexValidation.MatchString(input); !match {
		return errors.Wrapf(ErrInvalidInput, "%s: %s", name, input)
	}
	return nil
}

// parseAPIError parses an API error response and returns an APIError struct.
func parseAPIError(resp *http.Response) (*models.APIError, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read error response body: %w", err)
	}

	var apiError models.APIError
	if err := json.Unmarshal(body, &apiError); err != nil {
		return nil, fmt.Errorf("failed to parse error response: %w", err)
	}

	return &apiError, nil
}