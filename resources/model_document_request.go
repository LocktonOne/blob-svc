/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type DocumentRequest struct {
	Key
	Attributes    DocumentRequestAttributes    `json:"attributes"`
	Relationships DocumentRequestRelationships `json:"relationships"`
}
type DocumentRequestResponse struct {
	Data     DocumentRequest `json:"data"`
	Included Included        `json:"included"`
}

type DocumentRequestListResponse struct {
	Data     []DocumentRequest `json:"data"`
	Included Included          `json:"included"`
	Links    *Links            `json:"links"`
}

// MustDocumentRequest - returns DocumentRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustDocumentRequest(key Key) *DocumentRequest {
	var documentRequest DocumentRequest
	if c.tryFindEntry(key, &documentRequest) {
		return &documentRequest
	}
	return nil
}
