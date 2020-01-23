/*******************************************************************************
*
* THIS FILE IS AUTOGENERATED.
*
* Edit `generated.go.in` instead and run `make generate` to produce this file.
*
*******************************************************************************/

package schwift

//AccountHeaders contains the headers for a schwift.Account instance.
//
//To read and write well-known headers, use the methods on this type.
//To read and write arbitary headers, use the methods on the Headers supertype.
type AccountHeaders struct {
	Headers
}

//NewAccountHeaders creates a new AccountHeaders instance. The return value
//will have the Headers attribute initialized to a non-nil map.
func NewAccountHeaders() AccountHeaders {
	return AccountHeaders{make(Headers)}
}

//Validate returns MalformedHeaderError if the value of any well-known header
//does not conform to its data type. This is called automatically by Schwift
//when preparing an AccountHeaders instance from a GET/HEAD response, so you
//usually do not need to do it yourself. You will get the validation error from
//the Account method doing the request, e.g. Headers().
func (h AccountHeaders) Validate() error {
	if err := h.BytesUsed().validate(); err != nil {
		return err
	}
	if err := h.ContainerCount().validate(); err != nil {
		return err
	}
	if err := h.Metadata().validate(); err != nil {
		return err
	}
	if err := h.BytesUsedQuota().validate(); err != nil {
		return err
	}
	if err := h.TempURLKey2().validate(); err != nil {
		return err
	}
	if err := h.TempURLKey().validate(); err != nil {
		return err
	}
	if err := h.ObjectCount().validate(); err != nil {
		return err
	}
	if err := h.CreatedAt().validate(); err != nil {
		return err
	}
	return evadeGolintComplaint1()
}

//BytesUsed provides type-safe access to X-Account-Bytes-Used headers.
func (h AccountHeaders) BytesUsed() FieldUint64Readonly {
	return FieldUint64Readonly{h.Headers, "X-Account-Bytes-Used"}
}

//ContainerCount provides type-safe access to X-Account-Container-Count headers.
func (h AccountHeaders) ContainerCount() FieldUint64Readonly {
	return FieldUint64Readonly{h.Headers, "X-Account-Container-Count"}
}

//Metadata provides type-safe access to X-Account-Meta- headers.
func (h AccountHeaders) Metadata() FieldMetadata {
	return FieldMetadata{h.Headers, "X-Account-Meta-"}
}

//BytesUsedQuota provides type-safe access to X-Account-Meta-Quota-Bytes headers.
func (h AccountHeaders) BytesUsedQuota() FieldUint64 {
	return FieldUint64{h.Headers, "X-Account-Meta-Quota-Bytes"}
}

//TempURLKey2 provides type-safe access to X-Account-Meta-Temp-URL-Key-2 headers.
func (h AccountHeaders) TempURLKey2() FieldString {
	return FieldString{h.Headers, "X-Account-Meta-Temp-URL-Key-2"}
}

//TempURLKey provides type-safe access to X-Account-Meta-Temp-URL-Key headers.
func (h AccountHeaders) TempURLKey() FieldString {
	return FieldString{h.Headers, "X-Account-Meta-Temp-URL-Key"}
}

//ObjectCount provides type-safe access to X-Account-Object-Count headers.
func (h AccountHeaders) ObjectCount() FieldUint64Readonly {
	return FieldUint64Readonly{h.Headers, "X-Account-Object-Count"}
}

//CreatedAt provides type-safe access to X-Timestamp headers.
func (h AccountHeaders) CreatedAt() FieldUnixTimeReadonly {
	return FieldUnixTimeReadonly{h.Headers, "X-Timestamp"}
}

//ContainerHeaders contains the headers for a schwift.Container instance.
//
//To read and write well-known headers, use the methods on this type.
//To read and write arbitary headers, use the methods on the Headers supertype.
type ContainerHeaders struct {
	Headers
}

//NewContainerHeaders creates a new ContainerHeaders instance. The return value
//will have the Headers attribute initialized to a non-nil map.
func NewContainerHeaders() ContainerHeaders {
	return ContainerHeaders{make(Headers)}
}

//Validate returns MalformedHeaderError if the value of any well-known header
//does not conform to its data type. This is called automatically by Schwift
//when preparing an ContainerHeaders instance from a GET/HEAD response, so you
//usually do not need to do it yourself. You will get the validation error from
//the Container method doing the request, e.g. Headers().
func (h ContainerHeaders) Validate() error {
	if err := h.BytesUsed().validate(); err != nil {
		return err
	}
	if err := h.Metadata().validate(); err != nil {
		return err
	}
	if err := h.BytesUsedQuota().validate(); err != nil {
		return err
	}
	if err := h.ObjectCountQuota().validate(); err != nil {
		return err
	}
	if err := h.TempURLKey2().validate(); err != nil {
		return err
	}
	if err := h.TempURLKey().validate(); err != nil {
		return err
	}
	if err := h.ObjectCount().validate(); err != nil {
		return err
	}
	if err := h.ReadACL().validate(); err != nil {
		return err
	}
	if err := h.SyncKey().validate(); err != nil {
		return err
	}
	if err := h.SyncTo().validate(); err != nil {
		return err
	}
	if err := h.WriteACL().validate(); err != nil {
		return err
	}
	if err := h.HistoryLocation().validate(); err != nil {
		return err
	}
	if err := h.StoragePolicy().validate(); err != nil {
		return err
	}
	if err := h.CreatedAt().validate(); err != nil {
		return err
	}
	if err := h.VersionsLocation().validate(); err != nil {
		return err
	}
	return evadeGolintComplaint1()
}

//BytesUsed provides type-safe access to X-Container-Bytes-Used headers.
func (h ContainerHeaders) BytesUsed() FieldUint64Readonly {
	return FieldUint64Readonly{h.Headers, "X-Container-Bytes-Used"}
}

//Metadata provides type-safe access to X-Container-Meta- headers.
func (h ContainerHeaders) Metadata() FieldMetadata {
	return FieldMetadata{h.Headers, "X-Container-Meta-"}
}

//BytesUsedQuota provides type-safe access to X-Container-Meta-Quota-Bytes headers.
func (h ContainerHeaders) BytesUsedQuota() FieldUint64 {
	return FieldUint64{h.Headers, "X-Container-Meta-Quota-Bytes"}
}

//ObjectCountQuota provides type-safe access to X-Container-Meta-Quota-Count headers.
func (h ContainerHeaders) ObjectCountQuota() FieldUint64 {
	return FieldUint64{h.Headers, "X-Container-Meta-Quota-Count"}
}

//TempURLKey2 provides type-safe access to X-Container-Meta-Temp-URL-Key-2 headers.
func (h ContainerHeaders) TempURLKey2() FieldString {
	return FieldString{h.Headers, "X-Container-Meta-Temp-URL-Key-2"}
}

//TempURLKey provides type-safe access to X-Container-Meta-Temp-URL-Key headers.
func (h ContainerHeaders) TempURLKey() FieldString {
	return FieldString{h.Headers, "X-Container-Meta-Temp-URL-Key"}
}

//ObjectCount provides type-safe access to X-Container-Object-Count headers.
func (h ContainerHeaders) ObjectCount() FieldUint64Readonly {
	return FieldUint64Readonly{h.Headers, "X-Container-Object-Count"}
}

//ReadACL provides type-safe access to X-Container-Read headers.
func (h ContainerHeaders) ReadACL() FieldString {
	return FieldString{h.Headers, "X-Container-Read"}
}

//SyncKey provides type-safe access to X-Container-Sync-Key headers.
func (h ContainerHeaders) SyncKey() FieldString {
	return FieldString{h.Headers, "X-Container-Sync-Key"}
}

//SyncTo provides type-safe access to X-Container-Sync-To headers.
func (h ContainerHeaders) SyncTo() FieldString {
	return FieldString{h.Headers, "X-Container-Sync-To"}
}

//WriteACL provides type-safe access to X-Container-Write headers.
func (h ContainerHeaders) WriteACL() FieldString {
	return FieldString{h.Headers, "X-Container-Write"}
}

//HistoryLocation provides type-safe access to X-History-Location headers.
func (h ContainerHeaders) HistoryLocation() FieldString {
	return FieldString{h.Headers, "X-History-Location"}
}

//StoragePolicy provides type-safe access to X-Storage-Policy headers.
func (h ContainerHeaders) StoragePolicy() FieldString {
	return FieldString{h.Headers, "X-Storage-Policy"}
}

//CreatedAt provides type-safe access to X-Timestamp headers.
func (h ContainerHeaders) CreatedAt() FieldUnixTimeReadonly {
	return FieldUnixTimeReadonly{h.Headers, "X-Timestamp"}
}

//VersionsLocation provides type-safe access to X-Versions-Location headers.
func (h ContainerHeaders) VersionsLocation() FieldString {
	return FieldString{h.Headers, "X-Versions-Location"}
}

//ObjectHeaders contains the headers for a schwift.Object instance.
//
//To read and write well-known headers, use the methods on this type.
//To read and write arbitary headers, use the methods on the Headers supertype.
type ObjectHeaders struct {
	Headers
}

//NewObjectHeaders creates a new ObjectHeaders instance. The return value
//will have the Headers attribute initialized to a non-nil map.
func NewObjectHeaders() ObjectHeaders {
	return ObjectHeaders{make(Headers)}
}

//Validate returns MalformedHeaderError if the value of any well-known header
//does not conform to its data type. This is called automatically by Schwift
//when preparing an ObjectHeaders instance from a GET/HEAD response, so you
//usually do not need to do it yourself. You will get the validation error from
//the Object method doing the request, e.g. Headers().
func (h ObjectHeaders) Validate() error {
	if err := h.ContentDisposition().validate(); err != nil {
		return err
	}
	if err := h.ContentEncoding().validate(); err != nil {
		return err
	}
	if err := h.SizeBytes().validate(); err != nil {
		return err
	}
	if err := h.ContentType().validate(); err != nil {
		return err
	}
	if err := h.Etag().validate(); err != nil {
		return err
	}
	if err := h.UpdatedAt().validate(); err != nil {
		return err
	}
	if err := h.ExpiresAt().validate(); err != nil {
		return err
	}
	if err := h.Metadata().validate(); err != nil {
		return err
	}
	if err := h.SymlinkTargetAccount().validate(); err != nil {
		return err
	}
	if err := h.SymlinkTarget().validate(); err != nil {
		return err
	}
	if err := h.CreatedAt().validate(); err != nil {
		return err
	}
	return evadeGolintComplaint1()
}

//ContentDisposition provides type-safe access to Content-Disposition headers.
func (h ObjectHeaders) ContentDisposition() FieldString {
	return FieldString{h.Headers, "Content-Disposition"}
}

//ContentEncoding provides type-safe access to Content-Encoding headers.
func (h ObjectHeaders) ContentEncoding() FieldString {
	return FieldString{h.Headers, "Content-Encoding"}
}

//SizeBytes provides type-safe access to Content-Length headers.
func (h ObjectHeaders) SizeBytes() FieldUint64 {
	return FieldUint64{h.Headers, "Content-Length"}
}

//ContentType provides type-safe access to Content-Type headers.
func (h ObjectHeaders) ContentType() FieldString {
	return FieldString{h.Headers, "Content-Type"}
}

//Etag provides type-safe access to Etag headers.
func (h ObjectHeaders) Etag() FieldString {
	return FieldString{h.Headers, "Etag"}
}

//UpdatedAt provides type-safe access to Last-Modified headers.
func (h ObjectHeaders) UpdatedAt() FieldHTTPTimeReadonly {
	return FieldHTTPTimeReadonly{h.Headers, "Last-Modified"}
}

//ExpiresAt provides type-safe access to X-Delete-At headers.
func (h ObjectHeaders) ExpiresAt() FieldUnixTime {
	return FieldUnixTime{h.Headers, "X-Delete-At"}
}

//Metadata provides type-safe access to X-Object-Meta- headers.
func (h ObjectHeaders) Metadata() FieldMetadata {
	return FieldMetadata{h.Headers, "X-Object-Meta-"}
}

//SymlinkTargetAccount provides type-safe access to X-Symlink-Target-Account headers.
func (h ObjectHeaders) SymlinkTargetAccount() FieldString {
	return FieldString{h.Headers, "X-Symlink-Target-Account"}
}

//SymlinkTarget provides type-safe access to X-Symlink-Target headers.
func (h ObjectHeaders) SymlinkTarget() FieldString {
	return FieldString{h.Headers, "X-Symlink-Target"}
}

//CreatedAt provides type-safe access to X-Timestamp headers.
func (h ObjectHeaders) CreatedAt() FieldUnixTimeReadonly {
	return FieldUnixTimeReadonly{h.Headers, "X-Timestamp"}
}

func evadeGolintComplaint1() error {
	return nil
}
