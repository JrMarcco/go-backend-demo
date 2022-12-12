package copier

import (
	"fmt"
	"reflect"
)

type RefCopier[S, D any] struct {
	rootFn *fieldNode
}

type fieldNode struct {
	name   string // 字段名
	fields []*fieldNode

	sIndex int // 源字段 index
	dIndex int // 目标字段 index
}

func NewRefCopier[S, D any]() (*RefCopier[S, D], error) {
	srcTyp := reflect.TypeOf(new(S)).Elem()
	dstTyp := reflect.TypeOf(new(D)).Elem()

	if srcTyp.Kind() != reflect.Struct || dstTyp.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid type")
	}

	rootFn := &fieldNode{
		fields: []*fieldNode{},
	}

	if err := createFieldNode(srcTyp, dstTyp, rootFn); err != nil {
		return nil, err
	}

	return &RefCopier[S, D]{
		rootFn: rootFn,
	}, nil
}

// createFieldNode 创建字段树节点
func createFieldNode(srcTyp, dstTyp reflect.Type, rootFn *fieldNode) error {
	// 记录原字段 index
	srcFdMap := map[string]int{}
	for i := 0; i < srcTyp.NumField(); i++ {
		fd := srcTyp.Field(i)
		if fd.IsExported() {
			srcFdMap[fd.Name] = i
		}
	}

	dstNum := dstTyp.NumField()
	for i := 0; i < dstNum; i++ {
		dfd := dstTyp.Field(i)
		if !dfd.IsExported() {
			continue
		}
		if sIndex, ok := srcFdMap[dfd.Name]; ok {
			sfd := srcTyp.Field(sIndex)

			if sfd.Type.Kind() != dfd.Type.Kind() {
				// 同名字段类型不一致忽略字段
				continue
			}

			if sfd.Type.Kind() == reflect.Pointer {
				if sfd.Type.Elem().Kind() == reflect.Pointer {
					// TODO 完善异常信息 不支持多重指针
					return fmt.Errorf("invalid field")
				}
				if sfd.Type.Elem().Kind() != dfd.Type.Elem().Kind() {
					// 同名指针指向的类型不一致，忽略字段
					continue
				}
			}
			fn := &fieldNode{
				name:   sfd.Name,
				fields: []*fieldNode{},
				sIndex: sIndex,
				dIndex: i,
			}

			sFdTyp := sfd.Type
			dFdTyp := dfd.Type

			if sFdTyp.Kind() == reflect.Pointer {
				sFdTyp = sFdTyp.Elem()
				dFdTyp = dFdTyp.Elem()
			}

			if isBuiltinType(sFdTyp.Kind()) {
				// 内建类型
				if sFdTyp != dFdTyp {
					// TODO 完善错误信息（字段不匹配）
					return fmt.Errorf("invalid field")
				}
			} else if sFdTyp.Kind() == reflect.Struct {
				// 字段类型为结构体，递归构建
				if err := createFieldNode(sFdTyp, dFdTyp, fn); err != nil {
					return err
				}
			} else {
				continue
			}
			rootFn.fields = append(rootFn.fields, fn)
		}
	}
	return nil
}

func isBuiltinType(kind reflect.Kind) bool {
	switch kind {
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.String,
		reflect.Slice,
		reflect.Map,
		reflect.Chan,
		reflect.Array:
		return true
	}
	return false
}

var _ Copier[any, any] = new(RefCopier[any, any])

func (r *RefCopier[S, D]) Copy(src *S) (*D, error) {
	dst := new(D)
	if err := r.CopyTo(src, dst); err != nil {
		return nil, err
	}
	return dst, nil
}

func (r *RefCopier[S, D]) CopyTo(src *S, dst *D) error {
	return r.execCopy(src, dst)
}

func (r *RefCopier[S, D]) execCopy(src *S, dst *D) error {
	srcTyp := reflect.TypeOf(src)
	dstTyp := reflect.TypeOf(dst)

	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	return r.copyNode(srcTyp, dstTyp, srcVal, dstVal, r.rootFn)
}

func (r *RefCopier[S, D]) copyNode(srcTyp, dstTyp reflect.Type, srcVal, dstVal reflect.Value, rootFn *fieldNode) error {
	if srcVal.Kind() == reflect.Pointer {
		// 当前字段类型为指针
		if srcVal.IsNil() {
			return nil
		}
		if dstVal.IsNil() {
			// 目标值为空，先初始化
			dstVal.Set(reflect.New(dstTyp.Elem()))
		}

		srcTyp, srcVal = srcTyp.Elem(), srcVal.Elem()
		dstTyp, dstVal = dstTyp.Elem(), dstVal.Elem()
	}

	if len(rootFn.fields) == 0 {
		// 叶子节点，拷贝值
		if dstVal.CanSet() {
			dstVal.Set(srcVal)
		}
		return nil
	}

	for _, fn := range rootFn.fields {
		srcFdTyp := srcTyp.Field(fn.sIndex)
		srcFdVal := srcVal.Field(fn.sIndex)

		dstFdTyp := dstTyp.Field(fn.dIndex)
		dstFdVal := dstVal.Field(fn.dIndex)

		if err := r.copyNode(srcFdTyp.Type, dstFdTyp.Type, srcFdVal, dstFdVal, fn); err != nil {
			return err
		}
	}
	return nil
}
