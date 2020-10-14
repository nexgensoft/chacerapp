package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Pagination interface {
	GenerateNextPageToken(page PageInfo, pageSize int32) (string, error)
	ParsePageToken(token string) (PageInfo, error)
}

type paginator struct {
	secret []byte
}

func NewPaginator(secret []byte) Pagination {
	return &paginator{secret}
}

// PageInfo represents the data that is held within the page_token
// provided during pagination requests.
type PageInfo struct {
	// A unique Key that should identify the request. A typical value
	// for this field can be the request URI or the parent of the pagination
	// request if one is provided.
	RequestKey string
	// The end of the cursor from the previous page request. This
	// value is used to determine where the pagination should resume.
	EndCursor string
	// Filter is the filter that was used on the original
	// pagination request. When the filter is used in a request,
	// it must match the filter in then page token.
	Filter string
	// Order is the order that was used with the results of the
	// original pagination request. When the order option is used
	// in a request it must match the order in the page token.
	Order string
}

func (p *paginator) GenerateNextPageToken(pager PageInfo, pageSize int32) (string, error) {
	// default to setting to 0
	if pager.EndCursor == "" {
		pager.EndCursor = "0"
	}

	// Parse the old end cursor and add the page size to it
	endCursor, err := strconv.Atoi(pager.EndCursor)
	if err != nil {
		return "", err
	}

	newPageToken := PageInfo{
		RequestKey: pager.RequestKey,
		EndCursor:  strconv.Itoa(endCursor + int(pageSize)),
		Filter:     pager.Filter,
		Order:      pager.Order,
	}

	token, err := json.Marshal(newPageToken)
	if err != nil {
		return "", err
	}

	aesCipher, err := aes.NewCipher(p.secret)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, token, nil)), nil
}

// ParsePageToken will take a page token from a request and parse it
// into a PageInfo struct. This function does not validate that the
// contents of the token are valid for any request being made. When
// the page token was not encrypted with the configured secret then
// a gRPC friendly error will be returned. A nil struct will be returned
// when an empty page token is provided.
func (p *paginator) ParsePageToken(token string) (PageInfo, error) {
	if token == "" {
		return PageInfo{}, nil
	}

	aesCipher, err := aes.NewCipher(p.secret)
	if err != nil {
		return PageInfo{}, err
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return PageInfo{}, err
	}

	// The token should be base64 encoded, so decode before doing anything
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil || len(decodedToken) < gcm.NonceSize() {
		return PageInfo{}, status.Error(codes.InvalidArgument, "invalid page_token provided")
	}

	decypted, err := gcm.Open(nil, decodedToken[:gcm.NonceSize()], decodedToken[gcm.NonceSize():], nil)
	if err != nil {
		return PageInfo{}, status.Error(codes.InvalidArgument, "invalid page_token provided")
	}

	pageInfo := PageInfo{}
	if err := json.Unmarshal(decypted, &pageInfo); err != nil {
		return PageInfo{}, err
	}

	return pageInfo, nil
}

func paginateQuery(query string, info PageInfo, pageSize int32) string {
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	query += " LIMIT " + strconv.Itoa(int(pageSize))
	if info.EndCursor != "" {
		query += " OFFSET " + info.EndCursor
	}
	return query
}
