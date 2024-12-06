package models

// BaseMetadata contains common fields shared by both artifact and artifact version metadata.
type BaseMetadata struct {
	GroupID      string            `json:"groupId"`
	ArtifactID   string            `json:"artifactId"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	ArtifactType string            `json:"artifactType"`
	Owner        string            `json:"owner"`
	CreatedOn    string            `json:"createdOn"`
	Labels       map[string]string `json:"labels"`
}

// ArtifactVersionMetadata represents metadata for a single artifact version.
type ArtifactVersionMetadata struct {
	BaseMetadata
	Version   string `json:"version"`
	GlobalID  int64  `json:"globalId"`
	ContentID int64  `json:"contentId"`
}

// ArtifactMetadata represents metadata for an artifact.
type ArtifactMetadata struct {
	BaseMetadata
	ModifiedBy string `json:"modifiedBy"`
	ModifiedOn string `json:"modifiedOn"`
}

// UpdateArtifactMetadataRequest represents the metadata update request.
type UpdateArtifactMetadataRequest struct {
	Name        string            `json:"name,omitempty"`        // Editable name
	Description string            `json:"description,omitempty"` // Editable description
	Labels      map[string]string `json:"labels,omitempty"`      // Editable labels
	Owner       string            `json:"owner,omitempty"`       // Editable owner
}
