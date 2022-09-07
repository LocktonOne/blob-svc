/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Document struct {
	Key
	Attributes    DocumentAttributes    `json:"attributes"`
	Relationships DocumentRelationships `json:"relationships"`
}
type DocumentResponse struct {
	Data     Document `json:"data"`
	Included Included `json:"included"`
}

type DocumentListResponse struct {
	Data     []Document `json:"data"`
	Included Included   `json:"included"`
	Links    *Links     `json:"links"`
}

// MustDocument - returns Document from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustDocument(key Key) *Document {
	var document Document
	if c.tryFindEntry(key, &document) {
		return &document
	}
	return nil
}
