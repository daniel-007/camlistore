/*
Copyright 2016 The Camlistore Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"net/url"
	"time"

	"camlistore.org/pkg/blob"
	"camlistore.org/pkg/types/camtypes"
)

// Duplicating the search pkg types in here - since we only use them for json
// decoding - , instead of importing them through the search package, which would
// bring in more dependencies, and hence a larger js file.
// To give an idea, the generated publisher.js is ~3.5MB, whereas if we instead import
// camlistore.org/pkg/search to use its types instead of the ones below, we grow to
// ~5.7MB.

// TODO(mpl): keep these types in sync with camlistore.org/pkg/search.
// Brad suggested using go:generate

// A MetaMap is a map from blobref to a DescribedBlob.
type MetaMap map[string]*DescribedBlob

type DescribedBlob struct {
	BlobRef   blob.Ref `json:"blobRef"`
	CamliType string   `json:"camliType,omitempty"`
	Size      int64    `json:"size,"`

	// if camliType "permanode"
	Permanode *DescribedPermanode `json:"permanode,omitempty"`

	// if camliType "file"
	File *camtypes.FileInfo `json:"file,omitempty"`
	// if camliType "directory"
	Dir *camtypes.FileInfo `json:"dir,omitempty"`
	// if camliType "file", and File.IsImage()
	Image *camtypes.ImageInfo `json:"image,omitempty"`
	// if camliType "file" and media file
	MediaTags map[string]string `json:"mediaTags,omitempty"`

	// if camliType "directory"
	DirChildren []blob.Ref `json:"dirChildren,omitempty"`

	// Stub is set if this is not loaded, but referenced.
	Stub bool `json:"-"`
}

type DescribedPermanode struct {
	Attr    url.Values `json:"attr"` // a map[string][]string
	ModTime time.Time  `json:"modtime,omitempty"`
}

// SearchResult is the result of the Search method for a given SearchQuery.
type SearchResult struct {
	Blobs    []*SearchResultBlob `json:"blobs"`
	Describe *DescribeResponse   `json:"description"`

	// Continue optionally specifies the continuation token to to
	// continue fetching results in this result set, if interrupted
	// by a Limit.
	Continue string `json:"continue,omitempty"`
}

type SearchResultBlob struct {
	Blob blob.Ref `json:"blob"`
	// ... file info, permanode info, blob info ... ?
}

func (r *SearchResultBlob) String() string {
	return fmt.Sprintf("[blob: %s]", r.Blob)
}

// DescribeResponse is the JSON response from $searchRoot/camli/search/describe.
type DescribeResponse struct {
	Meta MetaMap `json:"meta"`
}
