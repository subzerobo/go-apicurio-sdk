package apis_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/subzerobo/go-apicurio-sdk/apis"
	"github.com/subzerobo/go-apicurio-sdk/client"
	"github.com/subzerobo/go-apicurio-sdk/models"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupVersionAPIClient() *apis.ArtifactsAPI {
	apiClient := setupHTTPClient()
	return apis.NewArtifactsAPI(apiClient)
}

func TestVersionsAPI_DeleteArtifactVersion(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate the request
			assert.Equal(t, "/groups/test-group/artifacts/test-artifact/versions/1.0.0", r.URL.Path)
			assert.Equal(t, http.MethodDelete, r.Method)

			// Respond with a successful status code
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		// Create a mock client and API instance
		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		// Call the method
		err := api.DeleteArtifactVersion(context.Background(), "test-group", "test-artifact", "1.0.0")

		// Assertions
		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/test-group/artifacts/test-artifact/versions/1.0.0", r.URL.Path)
			assert.Equal(t, http.MethodDelete, r.Method)

			// Simulate a 404 Not Found response
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{
				Detail: "Artifact version not found",
				Status: http.StatusNotFound,
			})
			if err != nil {
				t.Error(err)
			}
		}))
		defer server.Close()

		// Create a mock client and API instance
		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		// Call the method
		err := api.DeleteArtifactVersion(context.Background(), "test-group", "test-artifact", "1.0.0")

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/test-group/artifacts/test-artifact/versions/1.0.0", r.URL.Path)
			assert.Equal(t, http.MethodDelete, r.Method)

			// Simulate a 405 Method Not Allowed response
			w.WriteHeader(http.StatusMethodNotAllowed)
			err := json.NewEncoder(w).Encode(models.APIError{
				Detail: "Method not allowed",
				Status: http.StatusMethodNotAllowed,
			})
			if err != nil {
				t.Error(err)
			}
		}))
		defer server.Close()

		// Create a mock client and API instance
		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		// Call the method
		err := api.DeleteArtifactVersion(context.Background(), "test-group", "test-artifact", "1.0.0")

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Method not allowed")
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/test-group/artifacts/test-artifact/versions/1.0.0", r.URL.Path)
			assert.Equal(t, http.MethodDelete, r.Method)

			// Simulate a 500 Internal Server Error response
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{
				Detail: "Internal Server Error",
				Status: http.StatusInternalServerError,
			})
			if err != nil {
				t.Error(err)
			}
		}))
		defer server.Close()

		// Create a mock client and API instance
		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		// Call the method
		err := api.DeleteArtifactVersion(context.Background(), "test-group", "test-artifact", "1.0.0")

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Internal Server Error")
	})
}

func TestVersionsAPI_GetArtifactVersionReferences(t *testing.T) {
	t.Run("Success with Parameters", func(t *testing.T) {
		mockResponse := []models.ArtifactReference{
			{GroupID: "test-group", ArtifactID: "artifact-1", Version: "1", Name: "Reference 1"},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/test-group/artifacts/artifact-1/versions/1/references?refType=INBOUND", r.URL.String())
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		params := &models.ArtifactVersionReferencesParams{RefType: "INBOUND"}
		result, err := api.GetArtifactVersionReferences(context.Background(), "test-group", "artifact-1", "1", params)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(*result))
		assert.Equal(t, "Reference 1", (*result)[0].Name)
	})

	t.Run("Success without Parameters", func(t *testing.T) {
		mockResponse := []models.ArtifactReference{
			{GroupID: "test-group", ArtifactID: "artifact-1", Version: "1", Name: "Reference 1"},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/test-group/artifacts/artifact-1/versions/1/references", r.URL.String())
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionReferences(context.Background(), "test-group", "artifact-1", "1", nil)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(*result))
		assert.Equal(t, "Reference 1", (*result)[0].Name)
	})

	t.Run("Bad Request (400)", func(t *testing.T) {
		mockError := models.APIError{Title: "Bad Request", Detail: "Invalid refType parameter"}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(mockError)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		params := &models.ArtifactVersionReferencesParams{RefType: "INVALID"}
		result, err := api.GetArtifactVersionReferences(context.Background(), "test-group", "artifact-1", "1", params)
		assert.Error(t, err)
		assert.Nil(t, result)
		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, "Bad Request", apiErr.Title)
		assert.Equal(t, "Invalid refType parameter", apiErr.Detail)
	})

	t.Run("Not Found (404)", func(t *testing.T) {
		mockError := models.APIError{Title: "Not Found", Detail: "Artifact not found"}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(mockError)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		params := &models.ArtifactVersionReferencesParams{}
		result, err := api.GetArtifactVersionReferences(context.Background(), "test-group", "artifact-1", "1", params)
		assert.Error(t, err)
		assert.Nil(t, result)
		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, "Not Found", apiErr.Title)
		assert.Equal(t, "Artifact not found", apiErr.Detail)
	})

	t.Run("Method Not Allowed (405)", func(t *testing.T) {
		mockError := models.APIError{Title: "Method Not Allowed", Detail: "This method is not allowed"}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			err := json.NewEncoder(w).Encode(mockError)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		params := &models.ArtifactVersionReferencesParams{}
		result, err := api.GetArtifactVersionReferences(context.Background(), "test-group", "artifact-1", "1", params)
		assert.Error(t, err)
		assert.Nil(t, result)
		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, "Method Not Allowed", apiErr.Title)
		assert.Equal(t, "This method is not allowed", apiErr.Detail)
	})

	t.Run("Internal Server Error (500)", func(t *testing.T) {
		mockError := models.APIError{Title: "Internal Server Error", Detail: "An unexpected error occurred"}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(mockError)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		params := &models.ArtifactVersionReferencesParams{}
		result, err := api.GetArtifactVersionReferences(context.Background(), "test-group", "artifact-1", "1", params)
		assert.Error(t, err)
		assert.Nil(t, result)
		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, "Internal Server Error", apiErr.Title)
		assert.Equal(t, "An unexpected error occurred", apiErr.Detail)
	})
}

func TestVersionsAPI_GetArtifactVersionComments(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := []models.ArtifactComment{
			{CommentID: "12345", Value: "This is a comment on an artifact version.", Owner: "dwayne", CreatedOn: "2023-07-01T15:22:01Z"},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/test-group/artifacts/artifact-1/versions/1/comments", r.URL.String())
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionComments(context.Background(), "test-group", "artifact-1", "1")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, len(*result))
		assert.Equal(t, "This is a comment on an artifact version.", (*result)[0].Value)
	})

	t.Run("Bad Request (400)", func(t *testing.T) {
		mockError := models.APIError{Title: "Bad Request", Detail: "Invalid version expression"}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(mockError)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionComments(context.Background(), "test-group", "artifact-1", "invalid-version")
		assert.Error(t, err)
		assert.Nil(t, result)
		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, "Bad Request", apiErr.Title)
		assert.Equal(t, "Invalid version expression", apiErr.Detail)
	})

	t.Run("Not Found (404)", func(t *testing.T) {
		mockError := models.APIError{Title: "Not Found", Detail: "Artifact not found"}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(mockError)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionComments(context.Background(), "non-existent-group", "non-existent-artifact", "1")
		assert.Error(t, err)
		assert.Nil(t, result)
		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, "Not Found", apiErr.Title)
		assert.Equal(t, "Artifact not found", apiErr.Detail)
	})

	t.Run("Internal Server Error (500)", func(t *testing.T) {
		mockError := models.APIError{Title: "Internal Server Error", Detail: "An unexpected error occurred"}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(mockError)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		result, err := api.GetArtifactVersionComments(context.Background(), "test-group", "artifact-1", "1")
		assert.Error(t, err)
		assert.Nil(t, result)
		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, "Internal Server Error", apiErr.Title)
		assert.Equal(t, "An unexpected error occurred", apiErr.Detail)
	})
}

func TestVersionsAPI_AddArtifactVersionComment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.ArtifactComment{
			CommentID: "12345",
			Value:     "This is a new comment on an artifact version.",
			Owner:     "dwayne",
			CreatedOn: "2023-07-01T15:22:01Z",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/my-group/artifacts/example-artifact/versions/v1/comments", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			// Check if request body matches
			var requestBody models.ArtifactComment
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			assert.NoError(t, err)
			assert.Equal(t, "This is a new comment on an artifact version.", requestBody.Value)

			// Write the response
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		comment := "This is a new comment on an artifact version."
		result, err := api.AddArtifactVersionComment(context.Background(), "my-group", "example-artifact", "v1", comment)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockResponse, *result)
	})

	t.Run("BadRequest", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 400, Title: "Invalid input"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		comment := ""
		result, err := api.AddArtifactVersionComment(context.Background(), "my-group", "example-artifact", "v1", comment)

		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 400, apiErr.Status)
		assert.Equal(t, "Invalid input", apiErr.Title)
	})

	t.Run("NotFound", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 404, Title: "Artifact not found"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		comment := "This is a new comment"
		result, err := api.AddArtifactVersionComment(context.Background(), "invalid-group", "example-artifact", "v1", comment)

		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 404, apiErr.Status)
		assert.Equal(t, "Artifact not found", apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 500, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		comment := "This is a new comment"
		result, err := api.AddArtifactVersionComment(context.Background(), "my-group", "example-artifact", "v1", comment)

		assert.Error(t, err)
		assert.Nil(t, result)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 500, apiErr.Status)
		assert.Equal(t, "Internal server error", apiErr.Title)
	})
}

func TestVersionsAPI_UpdateArtifactVersionComment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/my-group/artifacts/example-artifact/versions/v1/comments/12345", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)
			// Return success response
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		comment := "Updated comment value"
		err := api.UpdateArtifactVersionComment(context.Background(), "my-group", "example-artifact", "v1", "12345", comment)
		assert.NoError(t, err)
	})

	t.Run("BadRequest", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 400, Title: "Invalid input"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		comment := ""
		err := api.UpdateArtifactVersionComment(context.Background(), "my-group", "example-artifact", "v1", "12345", comment)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 400, apiErr.Status)
		assert.Equal(t, "Invalid input", apiErr.Title)
	})

	t.Run("NotFound", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 404, Title: "Comment not found"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		comment := ""
		err := api.UpdateArtifactVersionComment(context.Background(), "my-group", "example-artifact", "v1", "invalid-comment-id", comment)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 404, apiErr.Status)
		assert.Equal(t, "Comment not found", apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 500, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		comment := ""
		err := api.UpdateArtifactVersionComment(context.Background(), "my-group", "example-artifact", "v1", "12345", comment)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 500, apiErr.Status)
		assert.Equal(t, "Internal server error", apiErr.Title)
	})
}

func TestVersionsAPI_DeleteArtifactVersionComment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/my-group/artifacts/example-artifact/versions/v1/comments/12345", r.URL.Path)
			assert.Equal(t, http.MethodDelete, r.Method)
			// Return success response
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersionComment(context.Background(), "my-group", "example-artifact", "v1", "12345")
		assert.NoError(t, err)
	})

	t.Run("BadRequest", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 400, Title: "Invalid input"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersionComment(context.Background(), "my-group", "example-artifact", "v1", "12345")
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 400, apiErr.Status)
		assert.Equal(t, "Invalid input", apiErr.Title)
	})

	t.Run("NotFound", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 404, Title: "Comment not found"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersionComment(context.Background(), "my-group", "example-artifact", "v1", "invalid-comment-id")
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 404, apiErr.Status)
		assert.Equal(t, "Comment not found", apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 500, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.DeleteArtifactVersionComment(context.Background(), "my-group", "example-artifact", "v1", "12345")
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 500, apiErr.Status)
		assert.Equal(t, "Internal server error", apiErr.Title)
	})
}

func TestVersionsAPI_ListArtifactVersions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.ArtifactVersionListResponse{
			Count: 2,
			Versions: []models.ArtifactVersion{
				{
					CreatedOn:    "2024-12-10T08:56:40Z",
					ArtifactType: models.Json,
					State:        models.StateEnabled,
					GlobalID:     47,
					Version:      "2.0.0",
					ContentID:    47,
					ArtifactID:   "example-artifact",
					GroupID:      "my-group",
					ModifiedOn:   "2024-12-10T08:56:40Z",
				},
				{
					CreatedOn:    "2024-12-10T08:56:17Z",
					ArtifactType: models.Json,
					State:        models.StateEnabled,
					GlobalID:     46,
					Version:      "1.0.0",
					ContentID:    46,
					ArtifactID:   "example-artifact",
					GroupID:      "my-group",
					ModifiedOn:   "2024-12-10T08:56:17Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/my-group/artifacts/example-artifact/versions", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
			// Write the response
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)

		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		versions, err := api.ListArtifactVersions(context.Background(), "my-group", "example-artifact", nil)
		assert.NoError(t, err)
		assert.NotNil(t, versions)
		assert.Equal(t, 2, len(*versions))
		assert.Equal(t, "2.0.0", (*versions)[0].Version)
		assert.Equal(t, "1.0.0", (*versions)[1].Version)
	})

	t.Run("NotFound", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 404, Title: "not found"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		versions, err := api.ListArtifactVersions(context.Background(), "my-group", "example-artifact", nil)
		assert.Error(t, err)
		assert.Nil(t, versions)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 404, apiErr.Status)
		assert.Equal(t, "not found", apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 500, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		versions, err := api.ListArtifactVersions(context.Background(), "my-group", "example-artifact", nil)
		assert.Error(t, err)
		assert.Nil(t, versions)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 500, apiErr.Status)
		assert.Equal(t, "Internal server error", apiErr.Title)
	})
}

func TestVersionsAPI_CreateArtifactVersion(t *testing.T) {

	t.Run("Success", func(t *testing.T) {

		mockResponse := models.ArtifactVersionDetailed{
			ArtifactVersion: models.ArtifactVersion{
				Version:      "1.0.0",
				CreatedOn:    "2024-12-10T08:56:40Z",
				ArtifactType: models.Json,
				GlobalID:     40,
				State:        models.StateEnabled,
				ContentID:    10,
				ArtifactID:   "my-artifact-id",
				GroupID:      "my-group",
				ModifiedOn:   "2024-12-10T08:56:40Z",
			},
			Name:        "Artifact Name",
			Description: "Artifact Description",
			Labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/my-group/artifacts/example-artifact/versions", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
			// Return success response
			// Write the response
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createVersion := &models.CreateVersionRequest{
			Version: "1.0.0",
			Content: models.CreateContentRequest{
				Content:     `{"a": "1"}`,
				ContentType: "application/json",
			},
			Name:        "Artifact Name",
			Description: "Artifact Description",
			Labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			IsDraft: false,
		}
		res, err := api.CreateArtifactVersion(context.Background(), "my-group", "example-artifact", createVersion, false)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "1.0.0", res.Version)
		assert.Equal(t, "Artifact Name", res.Name)
		assert.Equal(t, "Artifact Description", res.Description)
		assert.Equal(t, 2, len(res.Labels))
		assert.Equal(t, models.Json, res.ArtifactType)
	})

	t.Run("BadRequest", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 400, Title: "Invalid input"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createVersion := &models.CreateVersionRequest{
			Version: "1.0.0",
			Content: models.CreateContentRequest{
				Content:     `{"a": "1"}`,
				ContentType: "application/json",
			},
			Name:        "Artifact Name",
			Description: "Artifact Description",
			Labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			IsDraft: false,
		}
		res, err := api.CreateArtifactVersion(context.Background(), "my-group", "example-artifact", createVersion, false)
		assert.Error(t, err)
		assert.Nil(t, res)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 400, apiErr.Status)
		assert.Equal(t, "Invalid input", apiErr.Title)
	})

	t.Run("NotFound", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 404, Title: "Comment not found"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createVersion := &models.CreateVersionRequest{
			Version: "1.0.0",
			Content: models.CreateContentRequest{
				Content:     `{"a": "1"}`,
				ContentType: "application/json",
			},
			Name:        "Artifact Name",
			Description: "Artifact Description",
			Labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			IsDraft: false,
		}
		res, err := api.CreateArtifactVersion(context.Background(), "my-group", "example-artifact", createVersion, false)
		assert.Error(t, err)
		assert.Nil(t, res)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 404, apiErr.Status)
		assert.Equal(t, "Comment not found", apiErr.Title)
	})

	t.Run("Conflict", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 409, Title: "Conflict"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createVersion := &models.CreateVersionRequest{
			Version: "1.0.0",
			Content: models.CreateContentRequest{
				Content:     `{"a": "1"}`,
				ContentType: "application/json",
			},
			Name:        "Artifact Name",
			Description: "Artifact Description",
			Labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			IsDraft: false,
		}
		res, err := api.CreateArtifactVersion(context.Background(), "my-group", "example-artifact", createVersion, false)
		assert.Error(t, err)
		assert.Nil(t, res)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 409, apiErr.Status)
		assert.Equal(t, "Conflict", apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 500, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createVersion := &models.CreateVersionRequest{
			Version: "1.0.0",
			Content: models.CreateContentRequest{
				Content:     `{"a": "1"}`,
				ContentType: "application/json",
			},
			Name:        "Artifact Name",
			Description: "Artifact Description",
			Labels: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			IsDraft: false,
		}
		res, err := api.CreateArtifactVersion(context.Background(), "my-group", "example-artifact", createVersion, false)
		assert.Error(t, err)
		assert.Nil(t, res)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 500, apiErr.Status)
		assert.Equal(t, "Internal server error", apiErr.Title)
	})

}

func TestVersionsAPI_GetArtifactVersionContent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := `{"a": "1"}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/my-group/artifacts/example-artifact/versions/1.0.0/content", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
			// Write the response
			w.Header().Set("X-Registry-ArtifactType", string(models.Json))
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(mockResponse))
			assert.NoError(t, err)

		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		content, err := api.GetArtifactVersionContent(context.Background(), "my-group", "example-artifact", "1.0.0", nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
		assert.Equal(t, `{"a": "1"}`, content.Content)
	})

	t.Run("BadRequest", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 400, Title: "bad request"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		version, err := api.GetArtifactVersionContent(context.Background(), "my-group", "example-artifact", "1.0.0", nil)
		assert.Error(t, err)
		assert.Nil(t, version)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 400, apiErr.Status)
		assert.Equal(t, "bad request", apiErr.Title)
	})

	t.Run("NotFound", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 404, Title: "not found"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		version, err := api.GetArtifactVersionContent(context.Background(), "my-group", "example-artifact", "1.0.0", nil)
		assert.Error(t, err)
		assert.Nil(t, version)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 404, apiErr.Status)
		assert.Equal(t, "not found", apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 500, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		version, err := api.GetArtifactVersionContent(context.Background(), "my-group", "example-artifact", "1.0.0", nil)
		assert.Error(t, err)
		assert.Nil(t, version)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 500, apiErr.Status)
		assert.Equal(t, "Internal server error", apiErr.Title)
	})
}

func TestVersionsAPI_UpdateArtifactVersionContent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/my-group/artifacts/example-artifact/versions/1.0.0/content", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createVersion := &models.CreateContentRequest{
			Content:     `{"a": "1"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(context.Background(), "my-group", "example-artifact", "1.0.0", createVersion)
		assert.NoError(t, err)
	})

	t.Run("BadRequest", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/my-group/artifacts/example-artifact/versions/1.0.0/content", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 400, Title: "Invalid input"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createVersion := &models.CreateContentRequest{
			Content:     `{"a": "1"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(context.Background(), "my-group", "example-artifact", "1.0.0", createVersion)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 400, apiErr.Status)
		assert.Equal(t, "Invalid input", apiErr.Title)
	})

	t.Run("NotFound", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 404, Title: "Comment not found"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createVersion := &models.CreateContentRequest{
			Content:     `{"a": "1"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(context.Background(), "my-group", "example-artifact", "1.0.0", createVersion)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 404, apiErr.Status)
		assert.Equal(t, "Comment not found", apiErr.Title)
	})

	t.Run("Conflict", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 409, Title: "Conflict"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createVersion := &models.CreateContentRequest{
			Content:     `{"a": "1"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(context.Background(), "my-group", "example-artifact", "1.0.0", createVersion)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 409, apiErr.Status)
		assert.Equal(t, "Conflict", apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 500, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		createVersion := &models.CreateContentRequest{
			Content:     `{"a": "1"}`,
			ContentType: "application/json",
		}

		err := api.UpdateArtifactVersionContent(context.Background(), "my-group", "example-artifact", "1.0.0", createVersion)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 500, apiErr.Status)
		assert.Equal(t, "Internal server error", apiErr.Title)
	})
}

func TestVersionsAPI_SearchForArtifactVersions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.ArtifactVersionListResponse{
			Count: 2,
			Versions: []models.ArtifactVersion{
				{
					CreatedOn:    "2024-12-10T08:56:40Z",
					ArtifactType: models.Json,
					State:        models.StateEnabled,
					GlobalID:     47,
					Version:      "2.0.0",
					ContentID:    47,
					ArtifactID:   "example-artifact",
					GroupID:      "my-group",
					ModifiedOn:   "2024-12-10T08:56:40Z",
				},
				{
					CreatedOn:    "2024-12-10T08:56:17Z",
					ArtifactType: models.Json,
					State:        models.StateEnabled,
					GlobalID:     46,
					Version:      "1.0.0",
					ContentID:    46,
					ArtifactID:   "example-artifact",
					GroupID:      "my-group",
					ModifiedOn:   "2024-12-10T08:56:17Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/search/versions", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)

			// Write the response
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		// Search for json artifact and enabled state
		params := &models.SearchVersionParams{
			ArtifactType: models.Json,
			State:        models.StateEnabled,
		}
		versions, err := api.SearchForArtifactVersions(context.Background(), params)
		assert.NoError(t, err)
		assert.NotNil(t, versions)
		assert.Equal(t, 2, len(*versions))
		assert.Equal(t, "2.0.0", (*versions)[0].Version)
		assert.Equal(t, "1.0.0", (*versions)[1].Version)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 500, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		// Search for json artifact and enabled state
		params := &models.SearchVersionParams{
			ArtifactType: models.Json,
			State:        models.StateEnabled,
		}
		versions, err := api.SearchForArtifactVersions(context.Background(), params)
		assert.Error(t, err)
		assert.Nil(t, versions)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 500, apiErr.Status)
		assert.Equal(t, "Internal server error", apiErr.Title)
	})
}

func TestVersionsAPI_SearchForArtifactVersionByContent(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.ArtifactVersionListResponse{
			Count: 2,
			Versions: []models.ArtifactVersion{
				{
					CreatedOn:    "2024-12-10T08:56:40Z",
					ArtifactType: models.Json,
					State:        models.StateEnabled,
					GlobalID:     47,
					Version:      "2.0.0",
					ContentID:    47,
					ArtifactID:   "example-artifact",
					GroupID:      "my-group",
					ModifiedOn:   "2024-12-10T08:56:40Z",
				},
				{
					CreatedOn:    "2024-12-10T08:56:17Z",
					ArtifactType: models.Json,
					State:        models.StateEnabled,
					GlobalID:     46,
					Version:      "1.0.0",
					ContentID:    46,
					ArtifactID:   "example-artifact",
					GroupID:      "my-group",
					ModifiedOn:   "2024-12-10T08:56:17Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/search/versions", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Equal(t, "test-content", string(body))

			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		params := &models.SearchVersionByContentParams{Limit: 10, Offset: 0}
		versions, err := api.SearchForArtifactVersionByContent(context.Background(), "test-content", params)
		assert.NoError(t, err)
		assert.NotNil(t, versions)
		assert.Equal(t, 2, len(*versions))
		assert.Equal(t, "2.0.0", (*versions)[0].Version)
		assert.Equal(t, "1.0.0", (*versions)[1].Version)
	})

	t.Run("BadRequest - Empty Content", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 400, Title: "content cannot be empty"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		params := &models.SearchVersionByContentParams{Limit: 10, Offset: 0}
		versions, err := api.SearchForArtifactVersionByContent(context.Background(), "", params)
		assert.Error(t, err)
		assert.Nil(t, versions)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 400, apiErr.Status)
		assert.Equal(t, "content cannot be empty", apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 500, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		params := &models.SearchVersionByContentParams{Limit: 10, Offset: 0}
		versions, err := api.SearchForArtifactVersionByContent(context.Background(), "test-content", params)
		assert.Error(t, err)
		assert.Nil(t, versions)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 500, apiErr.Status)
		assert.Equal(t, "Internal server error", apiErr.Title)
	})
}

func TestVersionsAPI_GetArtifactVersionState(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockResponse := models.StateResponse{
			State: models.StateEnabled,
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/my-group/artifacts/example-artifact/versions/1.0/state", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)

			// Write the response
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(mockResponse)
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		state, err := api.GetArtifactVersionState(context.Background(), "my-group", "example-artifact", "1.0")
		assert.NoError(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, models.StateEnabled, *state)
	})

	t.Run("NotFound", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 404, Title: "not found"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		state, err := api.GetArtifactVersionState(context.Background(), "my-group", "example-artifact", "1.0")
		assert.Error(t, err)
		assert.Nil(t, state)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 404, apiErr.Status)
		assert.Equal(t, "not found", apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 500, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		state, err := api.GetArtifactVersionState(context.Background(), "my-group", "example-artifact", "1.0")
		assert.Error(t, err)
		assert.Nil(t, state)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 500, apiErr.Status)
		assert.Equal(t, "Internal server error", apiErr.Title)
	})
}

func TestVersionsAPI_UpdateArtifactVersionState(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/groups/my-group/artifacts/example-artifact/versions/1.0/state", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			// Validate request body
			var requestBody map[string]string
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			assert.NoError(t, err)
			assert.Equal(t, "ENABLED", requestBody["state"])

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionState(context.Background(), "my-group", "example-artifact", "1.0", models.StateEnabled, false)
		assert.NoError(t, err)
	})

	t.Run("BadRequest", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 400, Title: "Invalid state"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionState(context.Background(), "my-group", "example-artifact", "1.0", "INVALID_STATE", false)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 400, apiErr.Status)
		assert.Equal(t, "Invalid state", apiErr.Title)
	})

	t.Run("Conflict", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 409, Title: "Conflict"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionState(context.Background(), "my-group", "example-artifact", "1.0", models.StateDraft, false)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 409, apiErr.Status)
		assert.Equal(t, "Conflict", apiErr.Title)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(models.APIError{Status: 500, Title: "Internal server error"})
			assert.NoError(t, err)
		}))
		defer server.Close()

		mockClient := &client.Client{BaseURL: server.URL, HTTPClient: server.Client()}
		api := apis.NewVersionsAPI(mockClient)

		err := api.UpdateArtifactVersionState(context.Background(), "my-group", "example-artifact", "1.0", models.StateEnabled, false)
		assert.Error(t, err)

		var apiErr *models.APIError
		ok := errors.As(err, &apiErr)
		assert.True(t, ok)
		assert.Equal(t, 500, apiErr.Status)
		assert.Equal(t, "Internal server error", apiErr.Title)
	})
}
