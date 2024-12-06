package apis_test

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/subzerobo/go-apicurio-sdk/apis"
	"github.com/subzerobo/go-apicurio-sdk/client"
	"github.com/subzerobo/go-apicurio-sdk/models"
	"net/http"
	"testing"
	"time"
)

const (
	baseURL    = "http://localhost:9080/apis/registry/v3"
	groupID    = "test-group"
	artifactID = "test-artifact"
)

var (
	stubArtifactContent = `{"type": "record", "name": "Test", "fields": [{"name": "field1", "type": "string"}]}`
)

func setupClient() *apis.ArtifactsAPI {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	apiClient := client.NewClient(baseURL, client.WithHTTPClient(httpClient))
	return apis.NewArtifactsAPI(*apiClient)
}

func cleanup(t *testing.T, artifactsAPI *apis.ArtifactsAPI) {
	ctx := context.Background()
	err := artifactsAPI.DeleteArtifactsInGroup(ctx, groupID)
	if err != nil {
		var APIError *models.APIError
		if errors.As(err, &APIError) && APIError.Status == 404 {
			return
		}
		t.Fatalf("Failed to clean up artifacts: %v", err)
	}
}

func TestArtifactsAPI(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	artifactsAPI := setupClient()

	// Clean up before and after tests
	t.Cleanup(func() { cleanup(t, artifactsAPI) })
	cleanup(t, artifactsAPI)

	ctx := context.Background()

	// Test CreateArtifact
	t.Run("CreateArtifact", func(t *testing.T) {
		artifact := models.Artifact{
			Name:         artifactID,
			ArtifactID:   artifactID,
			ArtifactType: string(models.Json),
			Content:      stubArtifactContent,
		}
		params := models.CreateArtifactParams{
			IfExists: models.IfExistsFail,
		}

		resp, err := artifactsAPI.CreateArtifact(ctx, groupID, artifact, params)
		assert.NoError(t, err)
		assert.Equal(t, groupID, resp.Artifact.GroupID)
		assert.Equal(t, artifactID, resp.Artifact.Name)
	})

	// Test SearchArtifacts
	t.Run("SearchArtifacts", func(t *testing.T) {
		params := &models.SearchArtifactsParams{
			Name: artifactID,
		}
		resp, err := artifactsAPI.SearchArtifacts(ctx, params)
		fmt.Println(err)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(resp.Artifacts), 1)
	})

	// Test ListArtifactReferences
	t.Run("ListArtifactReferences", func(t *testing.T) {
		contentID := int64(12345) // Replace with a valid content ID for your tests
		_, err := artifactsAPI.ListArtifactReferences(ctx, contentID)
		assert.Error(t, err) // Expect an error since no content ID exists
	})

	// Test ListArtifactReferencesByGlobalID
	t.Run("ListArtifactReferencesByGlobalID", func(t *testing.T) {
		globalID := int64(12345) // Replace with a valid global ID for your tests
		params := &models.ListArtifactReferencesByGlobalIDParams{}
		_, err := artifactsAPI.ListArtifactReferencesByGlobalID(ctx, globalID, params)
		assert.Error(t, err) // Expect an error since no global ID exists
	})

	// Test ListArtifactReferencesByHash
	t.Run("ListArtifactReferencesByHash", func(t *testing.T) {
		contentHash := "invalidhash" // Replace with a valid content hash for your tests
		_, err := artifactsAPI.ListArtifactReferencesByHash(ctx, contentHash)
		assert.Error(t, err) // Expect an error since no hash exists
	})

	// Test ListArtifactsInGroup
	t.Run("ListArtifactsInGroup", func(t *testing.T) {
		params := &models.ListArtifactsInGroupParams{}
		resp, err := artifactsAPI.ListArtifactsInGroup(ctx, groupID, params)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(resp.Artifacts), 1)
	})

	// Test GetArtifactContentByHash
	t.Run("GetArtifactContentByHash", func(t *testing.T) {
		contentHash := "invalidhash" // Replace with a valid content hash for your tests
		_, err := artifactsAPI.GetArtifactContentByHash(ctx, contentHash)
		assert.Error(t, err) // Expect an error since no hash exists
	})

	// Test GetArtifactContentByID
	t.Run("GetArtifactContentByID", func(t *testing.T) {
		contentID := int64(12345) // Replace with a valid content ID for your tests
		_, err := artifactsAPI.GetArtifactContentByID(ctx, contentID)
		assert.Error(t, err) // Expect an error since no content ID exists
	})

	// Test DeleteArtifactsInGroup
	t.Run("DeleteArtifactsInGroup", func(t *testing.T) {
		err := artifactsAPI.DeleteArtifactsInGroup(ctx, groupID)
		assert.NoError(t, err)
	})

	// Test DeleteArtifact
	t.Run("DeleteArtifact", func(t *testing.T) {

		// Re-create the artifact
		artifact := models.Artifact{
			Name:         artifactID,
			ArtifactID:   artifactID,
			ArtifactType: string(models.Json),
			Content:      stubArtifactContent,
		}
		params := models.CreateArtifactParams{
			IfExists: models.IfExistsFail,
		}

		resp, err := artifactsAPI.CreateArtifact(ctx, groupID, artifact, params)
		assert.NoError(t, err)
		assert.Equal(t, groupID, resp.Artifact.GroupID)
		assert.Equal(t, artifactID, resp.Artifact.Name)

		// Delete the artifact
		err = artifactsAPI.DeleteArtifact(ctx, groupID, artifactID)
		assert.NoError(t, err)
	})
}
