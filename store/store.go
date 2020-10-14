package store

import (
	"context"
	"database/sql"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	fieldmask "github.com/mennanov/fieldmask-utils"
	"google.golang.org/genproto/protobuf/field_mask"
)

type Storage interface {
	Account
	Location
	Pagination
	Room
}

type store struct {
	Pagination
	db *sql.DB
}

func New(db *sql.DB, paginator Pagination) Storage {
	return &store{paginator, db}
}

const defaultPageSize = 25
const serviceName = "//apis.chacerapp.com/"

var protoMarshaller = jsonpb.Marshaler{}
var protoUnmarshaller = jsonpb.Unmarshaler{}

type ListOption func(*listOptions)

type UpdateOption func(*updateOptions)

func WithUpdateMask(mask *field_mask.FieldMask) UpdateOption {
	return func(opts *updateOptions) {
		opts.fieldMask = mask
	}
}

func WithPageSize(size int32) ListOption {
	return func(opts *listOptions) {
		opts.pageSize = size
	}
}

func WithPageInfo(info PageInfo) ListOption {
	return func(opts *listOptions) {
		opts.pageInfo = PageInfo{
			EndCursor: info.EndCursor,
			Filter:    info.Filter,
			Order:     info.Order,
		}
	}
}

func getListOptions(opts ...ListOption) *listOptions {
	options := &listOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func getUpdateOptions(opts ...UpdateOption) *updateOptions {
	options := &updateOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

type listOptions struct {
	pageSize int32
	pageInfo PageInfo
}

type updateOptions struct {
	fieldMask *field_mask.FieldMask
}

type inserter interface {
	ExecContext(context.Context, string, ...interface{}) (*sql.Result, error)
}

type retriever interface {
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
}

type scanner interface {
	Scan(dst ...interface{}) error
}

func applyUpdateMask(dst, src proto.Message, fieldMask *field_mask.FieldMask) (merged proto.Message, err error) {
	var paths []string
	if fieldMask != nil {
		paths = fieldMask.Paths
	}

	mask, err := fieldmask.MaskFromPaths(paths, generator.CamelCase)
	if err != nil {
		return nil, err
	}

	// Clone the destination struct
	merged = proto.Clone(dst)
	// Apply the mask to the cloned version with the new updates
	if err := fieldmask.StructToStruct(mask, src, merged); err != nil {
		return nil, err
	}

	return merged, nil
}

func doTransaction(ctx context.Context, db *sql.DB, callback func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := callback(ctx, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
