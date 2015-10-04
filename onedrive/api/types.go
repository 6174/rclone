// Types passed and returned to and from the API

package api

import "time"

const (
	timeFormat = `"` + time.RFC3339 + `"`
)

// Error is returned from one drive when things go wrong
type Error struct {
	ErrorInfo struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		InnerError struct {
			Code string `json:"code"`
		} `json:"innererror"`
	} `json:"error"`
}

// Error returns a string for the error and statistifes the error interface
func (e *Error) Error() string {
	out := e.ErrorInfo.Code
	if e.ErrorInfo.InnerError.Code != "" {
		out += ": " + e.ErrorInfo.InnerError.Code
	}
	out += ": " + e.ErrorInfo.Message
	return out
}

// Check Error statisfies the error interface
var _ error = (*Error)(nil)

// Identity represents an identity of an actor. For example, and actor
// can be a user, device, or application.
type Identity struct {
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
}

// IdentitySet is a keyed collection of Identity objects. It is used
// to represent a set of identities associated with various events for
// an item, such as created by or last modified by.
type IdentitySet struct {
	User        Identity `json:"user"`
	Application Identity `json:"application"`
	Device      Identity `json:"device"`
}

// Quota groups storage space quota-related information on OneDrive into a single structure.
type Quota struct {
	Total     int    `json:"total"`
	Used      int    `json:"used"`
	Remaining int    `json:"remaining"`
	Deleted   int    `json:"deleted"`
	State     string `json:"state"` // normal | nearing | critical | exceeded
}

// Drive is a representation of a drive resource
type Drive struct {
	ID        string      `json:"id"`
	DriveType string      `json:"driveType"`
	Owner     IdentitySet `json:"owner"`
	Quota     Quota       `json:"quota"`
}

// Timestamp represents represents date and time information for the
// OneDrive API, by using ISO 8601 and is always in UTC time.
type Timestamp time.Time

// MarshalJSON turns a Timestamp into JSON (in UTC)
func (t *Timestamp) MarshalJSON() (out []byte, err error) {
	out = (*time.Time)(t).UTC().AppendFormat(out, timeFormat)
	return out, nil
}

// UnmarshalJSON turns JSON into a Timestamp
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	newT, err := time.Parse(timeFormat, string(data))
	if err != nil {
		return err
	}
	*t = Timestamp(newT)
	return nil
}

// ItemReference groups data needed to reference a OneDrive item
// across the service into a single structure.
type ItemReference struct {
	DriveID string `json:"driveId"` // Unique identifier for the Drive that contains the item.	Read-only.
	ID      string `json:"id"`      // Unique identifier for the item.	Read/Write.
	Path    string `json:"path"`    // Path that used to navigate to the item.	Read/Write.
}

// FolderFacet groups folder-related data on OneDrive into a single structure
type FolderFacet struct {
	ChildCount int64 `json:"childCount"` // Number of children contained immediately within this container.
}

// HashesType groups different types of hashes into a single structure, for an item on OneDrive.
type HashesType struct {
	Sha1Hash  string `json:"sha1Hash"`  // base64 encoded SHA1 hash for the contents of the file (if available)
	Crc32Hash string `json:"crc32Hash"` // base64 encoded CRC32 value of the file (if available)
}

// FileFacet groups file-related data on OneDrive into a single structure.
type FileFacet struct {
	MimeType string     `json:"mimeType"` // The MIME type for the file. This is determined by logic on the server and might not be the value provided when the file was uploaded.
	Hashes   HashesType `json:"hashes"`   // Hashes of the file's binary content, if available.
}

// FileSystemInfoFacet contains properties that are reported by the
// device's local file system for the local version of an item. This
// facet can be used to specify the last modified date or created date
// of the item as it was on the local device.
type FileSystemInfoFacet struct {
	CreatedDateTime      Timestamp `json:"createdDateTime"`      // The UTC date and time the file was created on a client.
	LastModifiedDateTime Timestamp `json:"lastModifiedDateTime"` // The UTC date and time the file was last modified on a client.
}

// DeletedFacet indicates that the item on OneDrive has been
// deleted. In this version of the API, the presence (non-null) of the
// facet value indicates that the file was deleted. A null (or
// missing) value indicates that the file is not deleted.
type DeletedFacet struct {
}

// Item represents metadata for an item in OneDrive
type Item struct {
	ID                   string               `json:"id"`                   // The unique identifier of the item within the Drive. Read-only.
	Name                 string               `json:"name"`                 // The name of the item (filename and extension). Read-write.
	ETag                 string               `json:"eTag"`                 // eTag for the entire item (metadata + content). Read-only.
	CTag                 string               `json:"cTag"`                 // An eTag for the content of the item. This eTag is not changed if only the metadata is changed. Read-only.
	CreatedBy            IdentitySet          `json:"createdBy"`            // Identity of the user, device, and application which created the item. Read-only.
	LastModifiedBy       IdentitySet          `json:"lastModifiedBy"`       // Identity of the user, device, and application which last modified the item. Read-only.
	CreatedDateTime      Timestamp            `json:"createdDateTime"`      // Date and time of item creation. Read-only.
	LastModifiedDateTime Timestamp            `json:"lastModifiedDateTime"` // Date and time the item was last modified. Read-only.
	Size                 int64                `json:"size"`                 // Size of the item in bytes. Read-only.
	ParentReference      *ItemReference       `json:"parentReference"`      // Parent information, if the item has a parent. Read-write.
	WebURL               string               `json:"webUrl"`               // URL that displays the resource in the browser. Read-only.
	Description          string               `json:"description"`          // Provide a user-visible description of the item. Read-write.
	Folder               *FolderFacet         `json:"folder"`               // Folder metadata, if the item is a folder. Read-only.
	File                 *FileFacet           `json:"file"`                 // File metadata, if the item is a file. Read-only.
	FileSystemInfo       *FileSystemInfoFacet `json:"fileSystemInfo"`       // File system information on client. Read-write.
	//	Image                *ImageFacet          `json:"image"`                // Image metadata, if the item is an image. Read-only.
	//	Photo                *PhotoFacet          `json:"photo"`                // Photo metadata, if the item is a photo. Read-only.
	//	Audio                *AudioFacet          `json:"audio"`                // Audio metadata, if the item is an audio file. Read-only.
	//	Video                *VideoFacet          `json:"video"`                // Video metadata, if the item is a video. Read-only.
	//	Location             *LocationFacet       `json:"location"`             // Location metadata, if the item has location data. Read-only.
	Deleted *DeletedFacet `json:"deleted"` // Information about the deleted state of the item. Read-only.
}

// ViewDeltaResponse is the response to the view delta method
type ViewDeltaResponse struct {
	Value      []Item `json:"value"`            // An array of Item objects which have been created, modified, or deleted.
	NextLink   string `json:"@odata.nextLink"`  // A URL to retrieve the next available page of changes.
	DeltaLink  string `json:"@odata.deltaLink"` // A URL returned instead of @odata.nextLink after all current changes have been returned. Used to read the next set of changes in the future.
	DeltaToken string `json:"@delta.token"`     // A token value that can be used in the query string on manually-crafted calls to view.delta. Not needed if you're using nextLink and deltaLink.
}

// ListChildrenResponse is the response to the list children method
type ListChildrenResponse struct {
	Value    []Item `json:"value"`           // An array of Item objects
	NextLink string `json:"@odata.nextLink"` // A URL to retrieve the next available page of items.
}

// CreateItemRequest is the request to create an item object
type CreateItemRequest struct {
	Name             string      `json:"name"`                   // Name of the folder to be created.
	Folder           FolderFacet `json:"folder"`                 // Empty Folder facet to indicate that folder is the type of resource to be created.
	ConflictBehavior string      `json:"@name.conflictBehavior"` // Determines what to do if an item with a matching name already exists in this folder. Accepted values are: rename, replace, and fail (the default).
}

// SetFileSystemInfo is used to Update an object's FileSystemInfo.
type SetFileSystemInfo struct {
	FileSystemInfo FileSystemInfoFacet `json:"fileSystemInfo"` // File system information on client. Read-write.
}

// CreateUploadResponse is the response from creating an upload session
type CreateUploadResponse struct {
	UploadURL          string    `json:"uploadUrl"`          // "https://sn3302.up.1drv.com/up/fe6987415ace7X4e1eF866337",
	ExpirationDateTime Timestamp `json:"expirationDateTime"` // "2015-01-29T09:21:55.523Z",
	NextExpectedRanges []string  `json:"nextExpectedRanges"` // ["0-"]
}

// UploadFragmentResponse is the response from uploading a fragment
type UploadFragmentResponse struct {
	ExpirationDateTime Timestamp `json:"expirationDateTime"` // "2015-01-29T09:21:55.523Z",
	NextExpectedRanges []string  `json:"nextExpectedRanges"` // ["0-"]
}