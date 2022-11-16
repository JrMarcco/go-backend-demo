package copier

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRefCopier_Copy(t *testing.T) {
	tcs := []struct {
		name     string
		copyFunc func() (any, error)
		wantRes  any
		wantErr  error
	}{
		{
			name: "basic",
			copyFunc: func() (any, error) {
				copier, err := NewRefCopier[basicSrc, basicDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&basicSrc{
					Name: "Foo",
					Age:  18,
				})
			},
			wantErr: nil,
			wantRes: &basicDst{
				Name: "Foo",
				Age:  18,
			},
		},
		{
			name: "diff type",
			copyFunc: func() (any, error) {
				copier, err := NewRefCopier[diffTypeSrc, diffTypeDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&diffTypeSrc{
					Name: "Foo",
					Sex:  "male",
				})
			},
			wantErr: nil,
			wantRes: &diffTypeDst{
				Name: "Foo",
			},
		},
		{
			name: "diff ptr type",
			copyFunc: func() (any, error) {
				copier, err := NewRefCopier[diffTypeSrc, diffTypeDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&diffTypeSrc{
					Name:   "Foo",
					PtrSex: toPtr("male"),
				})
			},
			wantErr: nil,
			wantRes: &diffTypeDst{
				Name: "Foo",
			},
		},
		{
			name: "struct field",
			copyFunc: func() (any, error) {
				copier, err := NewRefCopier[basicSrc, basicDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&basicSrc{
					Name: "Foo",
					Age:  18,
					StructAddr: address{
						Province: "Fujian",
						City:     "Xiamen",
					},
				})
			},
			wantErr: nil,
			wantRes: &basicDst{
				Name: "Foo",
				Age:  18,
				StructAddr: address{
					Province: "Fujian",
					City:     "Xiamen",
				},
			},
		},
		{
			name: "ptr field",
			copyFunc: func() (any, error) {
				copier, err := NewRefCopier[basicSrc, basicDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&basicSrc{
					Name: "Foo",
					Age:  18,
					PtrAddr: &address{
						Province: "Fujian",
						City:     "Xiamen",
						Code:     toPtr(361000),
					},
				})
			},
			wantErr: nil,
			wantRes: &basicDst{
				Name: "Foo",
				Age:  18,
				PtrAddr: &address{
					Province: "Fujian",
					City:     "Xiamen",
					Code:     toPtr(361000),
				},
			},
		},
		{
			name: "str map field",
			copyFunc: func() (any, error) {
				copier, err := NewRefCopier[basicSrc, basicDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&basicSrc{
					Name: "Foo",
					StrMap: map[string]string{
						"first":  "abc",
						"second": "efg",
					},
				})
			},
			wantErr: nil,
			wantRes: &basicDst{
				Name: "Foo",
				StrMap: map[string]string{
					"first":  "abc",
					"second": "efg",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.copyFunc()
			assert.Equal(t, tc.wantErr, err)

			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

type basicSrc struct {
	Name       string
	Age        int
	StructAddr address
	PtrAddr    *address

	StrMap map[string]string
}

type basicDst struct {
	Name     string
	Nickname string
	Age      int

	StructAddr address
	PtrAddr    *address

	StrMap map[string]string
}

type diffTypeSrc struct {
	Name   string
	Sex    string
	PtrSex *string
}

type diffTypeDst struct {
	Name   string
	Sex    int
	PtrSex *int
}

type address struct {
	Province string
	City     string

	Code *int
}

func toPtr[T any](t T) *T {
	return &t
}
