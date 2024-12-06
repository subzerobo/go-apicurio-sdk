package models

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var (
	// ErrUnknownArtifactType is returned when an unknown artifact type is encountered.
	ErrUnknownArtifactType = fmt.Errorf("unknown artifact type")
)

type IfExistsType string

const (
	IfExistsFail                IfExistsType = "FAIL"                   // (default) - server rejects the content with a 409 error
	IfExistsCreate              IfExistsType = "CREATE_VERSION"         // server creates a new version of the existing artifact and returns it
	IfExistsFindOrCreateVersion IfExistsType = "FIND_OR_CREATE_VERSION" // server returns an existing version that matches the provided content if such a version exists, otherwise a new version is created
)

type SearchedArtifact struct {
	GroupId      string `json:"groupId"`
	ArtifactId   string `json:"artifactId"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ArtifactType string `json:"artifactType"`
	Owner        string `json:"owner"`
	CreatedOn    string `json:"createdOn"`
	ModifiedBy   string `json:"modifiedBy"`
	ModifiedOn   string `json:"modifiedOn"`
}

// SearchArtifactsParams represents the optional parameters for searching artifacts.
type SearchArtifactsParams struct {
	Name         string   // Filter by artifact name
	Offset       int      // Default: 0
	Limit        int      // Default: 20
	Order        string   // Default: "asc", Enum: "asc", "desc"
	OrderBy      string   // Field to sort by, e.g., "name", "createdOn"
	Labels       []string // Filter by one or more name/value labels
	Description  string   // Filter by description
	GroupID      string   // Filter by artifact group
	GlobalID     int64    // Filter by globalId
	ContentID    int64    // Filter by contentId
	ArtifactID   string   // Filter by artifactId
	ArtifactType string   // Filter by artifact type (e.g., AVRO, JSON)
}

// ToQuery converts the SearchArtifactsParams struct to URL query parameters.
func (p *SearchArtifactsParams) ToQuery() url.Values {
	query := url.Values{}

	if p.Name != "" {
		query.Set("name", p.Name)
	}
	if p.Offset != 0 {
		query.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Limit != 0 {
		query.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Order != "" {
		query.Set("order", p.Order)
	}
	if p.OrderBy != "" {
		query.Set("orderby", p.OrderBy)
	}
	if len(p.Labels) > 0 {
		query.Set("labels", strings.Join(p.Labels, ","))
	}
	if p.Description != "" {
		query.Set("description", p.Description)
	}
	if p.GroupID != "" {
		query.Set("groupId", p.GroupID)
	}
	if p.GlobalID != 0 {
		query.Set("globalId", strconv.FormatInt(p.GlobalID, 10))
	}
	if p.ContentID != 0 {
		query.Set("contentId", strconv.FormatInt(p.ContentID, 10))
	}
	if p.ArtifactID != "" {
		query.Set("artifactId", p.ArtifactID)
	}
	if p.ArtifactType != "" {
		query.Set("artifactType", p.ArtifactType)
	}

	return query
}

// SearchArtifactsByContentParams represents the query parameters for the search by content API.
type SearchArtifactsByContentParams struct {
	Canonical    bool   // Canonicalize the content
	ArtifactType string // Artifact type (e.g., AVRO, JSON)
	GroupID      string // Filter by group ID
	Offset       int    // Number of artifacts to skip
	Limit        int    // Number of artifacts to return
	Order        string // Sort order (asc, desc)
	OrderBy      string // Field to sort by
}

// ToQuery converts the SearchArtifactsByContentParams struct to query parameters.
func (p *SearchArtifactsByContentParams) ToQuery() url.Values {
	query := url.Values{}

	if p.Canonical {
		query.Set("canonical", "true")
	}
	if p.ArtifactType != "" {
		query.Set("artifactType", p.ArtifactType)
	}
	if p.GroupID != "" {
		query.Set("groupId", p.GroupID)
	}
	if p.Offset != 0 {
		query.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Limit != 0 {
		query.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Order != "" {
		query.Set("order", p.Order)
	}
	if p.OrderBy != "" {
		query.Set("orderby", p.OrderBy)
	}

	return query
}

type SearchArtifactsResponse struct {
	Artifacts []SearchedArtifact `json:"artifacts"`
	Count     int                `json:"count"`
}

// ArtifactReference represents a reference to an artifact.
type ArtifactReference struct {
	GroupID    string `json:"groupId"`
	ArtifactID string `json:"artifactId"`
	Version    string `json:"version"`
	Name       string `json:"name"`
}

type RefType string

const (
	OutBound RefType = "OUTBOUND"
	InBound  RefType = "INBOUND"
)

type ArtifactType string

const (
	Avro     ArtifactType = "AVRO"     // Avro artifact type
	Protobuf ArtifactType = "PROTOBUF" // Protobuf artifact type
	Json     ArtifactType = "JSON"     // JSON artifact type
	KConnect ArtifactType = "KCONNECT" // Kafka Connect artifact type
	OpenAPI  ArtifactType = "OPENAPI"  // OpenAPI artifact type
	AsyncAPI ArtifactType = "ASYNCAPI" // AsyncAPI artifact type
	GraphQL  ArtifactType = "GRAPHQL"  // GraphQL artifact type
	WSDL     ArtifactType = "WSDL"     // WSDL artifact type
	XSD      ArtifactType = "XSD"      // XSD artifact type
)

// ParseArtifactType parses a string and returns the corresponding ArtifactType.
func ParseArtifactType(artifactType string) (ArtifactType, error) {
	switch artifactType {
	case string(Avro):
		return Avro, nil
	case string(Protobuf):
		return Protobuf, nil
	case string(Json):
		return Json, nil
	case string(KConnect):
		return KConnect, nil
	case string(OpenAPI):
		return OpenAPI, nil
	case string(AsyncAPI):
		return AsyncAPI, nil
	case string(GraphQL):
		return GraphQL, nil
	case string(WSDL):
		return WSDL, nil
	case string(XSD):
		return XSD, nil
	default:
		return "", ErrUnknownArtifactType
	}
}

// ListArtifactReferencesByGlobalIDParams represents the optional parameters for listing references by global ID.
type ListArtifactReferencesByGlobalIDParams struct {
	RefType RefType
}

// ToQuery converts the params struct to URL query parameters.
func (p *ListArtifactReferencesByGlobalIDParams) ToQuery() string {
	if p != nil && p.RefType != "" {
		return fmt.Sprintf("?refType=%s", p.RefType)
	}
	return ""
}

// ListArtifactsInGroupParams represents the query parameters for listing artifacts in a group.
type ListArtifactsInGroupParams struct {
	Limit   int    // Number of artifacts to return (default: 20)
	Offset  int    // Number of artifacts to skip (default: 0)
	Order   string // Enum: "asc", "desc"
	OrderBy string // Enum: "groupId", "artifactId", "createdOn", "modifiedOn", "artifactType", "name"
}

// ToQuery converts the ListArtifactsInGroupParams struct to query parameters.
func (p *ListArtifactsInGroupParams) ToQuery() string {
	query := ""
	if p != nil {
		query = "?"
		if p.Limit != 0 {
			query += fmt.Sprintf("limit=%d&", p.Limit)
		}
		if p.Offset != 0 {
			query += fmt.Sprintf("offset=%d&", p.Offset)
		}
		if p.Order != "" {
			query += fmt.Sprintf("order=%s&", p.Order)
		}
		if p.OrderBy != "" {
			query += fmt.Sprintf("orderby=%s&", p.OrderBy)
		}
		// Remove the trailing "&" if present
		query = strings.TrimSuffix(query, "&")
	}
	return query
}

type ListArtifactsResponse struct {
	Artifacts []SearchedArtifact `json:"artifacts"`
	Count     int                `json:"count"`
}

type ArtifactContent struct {
	Content      string       `json:"content"`
	ArtifactType ArtifactType `json:"artifactType"`
}

// Artifact represents the artifact metadata and content.
type Artifact struct {
	ArtifactID   string            `json:"artifactId,omitempty"`
	ArtifactType string            `json:"artifactType"`
	Name         string            `json:"name,omitempty"`
	Description  string            `json:"description,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	Content      string            `json:"content"`
}

type ArtifactDetail struct {
	GroupID     string            `json:"groupId"`
	ArtifactID  string            `json:"artifactId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	CreatedOn   string            `json:"createdOn"`
	ModifiedOn  string            `json:"modifiedOn"`
	ContentID   int64             `json:"contentId"`
	Labels      map[string]string `json:"labels"`
}

type ArtifactResponse struct {
	Artifact ArtifactDetail `json:"artifact"`
}

type CreateArtifactParams struct {
	IfExists  IfExistsType // IfExists behavior @See IfExistsType
	Canonical bool         // Indicates whether to canonicalize the artifact content.
	DryRun    bool         // If true, no changes are made, only checks are performed.
}

// ToQuery converts the parameters into a query string.
func (p CreateArtifactParams) ToQuery() string {
	query := "?"
	if p.IfExists != "" {
		query += fmt.Sprintf("ifExists=%s&", p.IfExists)
	}
	if p.IfExists == IfExistsFindOrCreateVersion && p.Canonical {
		query += "canonical=true&"
	}
	if p.DryRun {
		query += "dryRun=true&"
	}
	return strings.TrimSuffix(query, "&") // Remove trailing "&"
}
