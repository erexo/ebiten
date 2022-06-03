// Copyright 2022 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package directx

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/hajimehoshi/ebiten/v2/internal/microsoftgdk"
)

type (
	_BOOL int32
)

func boolToUintptr(v bool) uintptr {
	if v {
		return 1
	}
	return 0
}

func syscallN(ptr uintptr, args ...uintptr) (r1, r2 uintptr, e windows.Errno) {
	if ptr == 0 {
		panic("directx: function pointer must not be 0")
	}
	switch len(args) {
	case 0:
		return syscall.Syscall(ptr, 0, 0, 0, 0)
	case 1:
		return syscall.Syscall(ptr, 1, args[0], 0, 0)
	case 2:
		return syscall.Syscall(ptr, 2, args[0], args[1], 0)
	case 3:
		return syscall.Syscall(ptr, 3, args[0], args[1], args[2])
	case 4:
		return syscall.Syscall6(ptr, 4, args[0], args[1], args[2], args[3], 0, 0)
	case 5:
		return syscall.Syscall6(ptr, 5, args[0], args[1], args[2], args[3], args[4], 0)
	case 6:
		return syscall.Syscall6(ptr, 6, args[0], args[1], args[2], args[3], args[4], args[5])
	case 7:
		return syscall.Syscall9(ptr, 7, args[0], args[1], args[2], args[3], args[4], args[5], args[6], 0, 0)
	case 8:
		return syscall.Syscall9(ptr, 8, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], 0)
	case 9:
		return syscall.Syscall9(ptr, 9, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8])
	default:
		panic(fmt.Sprintf("directx: too many arguments: %d", len(args)))
	}
}

// Reference:
// * https://github.com/microsoft/DirectX-Headers
// * https://github.com/microsoft/win32metadata
// * https://raw.githubusercontent.com/microsoft/win32metadata/master/generation/WinSDK/RecompiledIdlHeaders/um/d3d12.h

const (
	_D3D12_APPEND_ALIGNED_ELEMENT            = 0xffffffff
	_D3D12_DEFAULT_DEPTH_BIAS                = 0
	_D3D12_DEFAULT_DEPTH_BIAS_CLAMP          = 0.0
	_D3D12_DEFAULT_STENCIL_READ_MASK         = 0xff
	_D3D12_DEFAULT_STENCIL_WRITE_MASK        = 0xff
	_D3D12_DEFAULT_SLOPE_SCALED_DEPTH_BIAS   = 0.0
	_D3D12_DESCRIPTOR_RANGE_OFFSET_APPEND    = 0xffffffff
	_D3D12_MAX_DEPTH                         = 1.0
	_D3D12_MIN_DEPTH                         = 0.0
	_D3D12_REQ_TEXTURE2D_U_OR_V_DIMENSION    = 16384
	_D3D12_RESOURCE_BARRIER_ALL_SUBRESOURCES = 0xffffffff
	_D3D12_SDK_VERSION                       = (_D3D12_SDK_VERSION_MAJOR << 16) | _D3D12_SDK_VERSION_MINOR
	_D3D12_SDK_VERSION_MAJOR                 = 2
	_D3D12_SDK_VERSION_MINOR                 = 4
	_D3D12XBOX_DEFAULT_SIZE_BYTES            = 0xffffffff
)

type _D3D_FEATURE_LEVEL int32

const (
	_D3D_FEATURE_LEVEL_11_0 _D3D_FEATURE_LEVEL = 0xb000
)

type _D3D_PRIMITIVE_TOPOLOGY int32

const (
	_D3D_PRIMITIVE_TOPOLOGY_TRIANGLELIST _D3D_PRIMITIVE_TOPOLOGY = 4
)

type _D3D_ROOT_SIGNATURE_VERSION int32

const (
	_D3D_ROOT_SIGNATURE_VERSION_1_0 _D3D_ROOT_SIGNATURE_VERSION = 0x1
)

type _D3D12_BLEND int32

const (
	_D3D12_BLEND_ZERO             _D3D12_BLEND = 1
	_D3D12_BLEND_ONE              _D3D12_BLEND = 2
	_D3D12_BLEND_SRC_COLOR        _D3D12_BLEND = 3
	_D3D12_BLEND_INV_SRC_COLOR    _D3D12_BLEND = 4
	_D3D12_BLEND_SRC_ALPHA        _D3D12_BLEND = 5
	_D3D12_BLEND_INV_SRC_ALPHA    _D3D12_BLEND = 6
	_D3D12_BLEND_DEST_ALPHA       _D3D12_BLEND = 7
	_D3D12_BLEND_INV_DEST_ALPHA   _D3D12_BLEND = 8
	_D3D12_BLEND_DEST_COLOR       _D3D12_BLEND = 9
	_D3D12_BLEND_INV_DEST_COLOR   _D3D12_BLEND = 10
	_D3D12_BLEND_SRC_ALPHA_SAT    _D3D12_BLEND = 11
	_D3D12_BLEND_BLEND_FACTOR     _D3D12_BLEND = 14
	_D3D12_BLEND_INV_BLEND_FACTOR _D3D12_BLEND = 15
	_D3D12_BLEND_SRC1_COLOR       _D3D12_BLEND = 16
	_D3D12_BLEND_INV_SRC1_COLOR   _D3D12_BLEND = 17
	_D3D12_BLEND_SRC1_ALPHA       _D3D12_BLEND = 18
	_D3D12_BLEND_INV_SRC1_ALPHA   _D3D12_BLEND = 19
)

type _D3D12_BLEND_OP int32

const (
	_D3D12_BLEND_OP_ADD          _D3D12_BLEND_OP = 1
	_D3D12_BLEND_OP_SUBTRACT     _D3D12_BLEND_OP = 2
	_D3D12_BLEND_OP_REV_SUBTRACT _D3D12_BLEND_OP = 3
	_D3D12_BLEND_OP_MIN          _D3D12_BLEND_OP = 4
	_D3D12_BLEND_OP_MAX          _D3D12_BLEND_OP = 5
)

type _D3D12_CLEAR_FLAGS int32

const (
	_D3D12_CLEAR_FLAG_DEPTH   _D3D12_CLEAR_FLAGS = 0x1
	_D3D12_CLEAR_FLAG_STENCIL _D3D12_CLEAR_FLAGS = 0x2
)

type _D3D12_COLOR_WRITE_ENABLE int32

const (
	_D3D12_COLOR_WRITE_ENABLE_RED   _D3D12_COLOR_WRITE_ENABLE = 1
	_D3D12_COLOR_WRITE_ENABLE_GREEN _D3D12_COLOR_WRITE_ENABLE = 2
	_D3D12_COLOR_WRITE_ENABLE_BLUE  _D3D12_COLOR_WRITE_ENABLE = 4
	_D3D12_COLOR_WRITE_ENABLE_ALPHA _D3D12_COLOR_WRITE_ENABLE = 8
	_D3D12_COLOR_WRITE_ENABLE_ALL   _D3D12_COLOR_WRITE_ENABLE = _D3D12_COLOR_WRITE_ENABLE_RED | _D3D12_COLOR_WRITE_ENABLE_GREEN | _D3D12_COLOR_WRITE_ENABLE_BLUE | _D3D12_COLOR_WRITE_ENABLE_ALPHA
)

type _D3D12_COMMAND_LIST_TYPE int32

const (
	_D3D12_COMMAND_LIST_TYPE_DIRECT _D3D12_COMMAND_LIST_TYPE = 0
)

type _D3D12_COMMAND_QUEUE_FLAGS int32

const (
	_D3D12_COMMAND_QUEUE_FLAG_NONE _D3D12_COMMAND_QUEUE_FLAGS = 0
)

type _D3D12_COMPARISON_FUNC int32

const (
	_D3D12_COMPARISON_FUNC_NEVER         _D3D12_COMPARISON_FUNC = 1
	_D3D12_COMPARISON_FUNC_LESS          _D3D12_COMPARISON_FUNC = 2
	_D3D12_COMPARISON_FUNC_EQUAL         _D3D12_COMPARISON_FUNC = 3
	_D3D12_COMPARISON_FUNC_LESS_EQUAL    _D3D12_COMPARISON_FUNC = 4
	_D3D12_COMPARISON_FUNC_GREATER       _D3D12_COMPARISON_FUNC = 5
	_D3D12_COMPARISON_FUNC_NOT_EQUAL     _D3D12_COMPARISON_FUNC = 6
	_D3D12_COMPARISON_FUNC_GREATER_EQUAL _D3D12_COMPARISON_FUNC = 7
	_D3D12_COMPARISON_FUNC_ALWAYS        _D3D12_COMPARISON_FUNC = 8
)

type _D3D12_CONSERVATIVE_RASTERIZATION_MODE int32

const (
	_D3D12_CONSERVATIVE_RASTERIZATION_MODE_OFF _D3D12_CONSERVATIVE_RASTERIZATION_MODE = 0
	_D3D12_CONSERVATIVE_RASTERIZATION_MODE_ON  _D3D12_CONSERVATIVE_RASTERIZATION_MODE = 1
)

type _D3D12_CPU_PAGE_PROPERTY int32

const (
	_D3D12_CPU_PAGE_PROPERTY_UNKNOWN _D3D12_CPU_PAGE_PROPERTY = 0
)

type _D3D12_CULL_MODE int32

const (
	_D3D12_CULL_MODE_NONE  _D3D12_CULL_MODE = 1
	_D3D12_CULL_MODE_FRONT _D3D12_CULL_MODE = 2
	_D3D12_CULL_MODE_BACK  _D3D12_CULL_MODE = 3
)

type _D3D12_DEBUG_FEATURE int32

const (
	_D3D12_DEBUG_FEATURE_NONE                                   _D3D12_DEBUG_FEATURE = 0
	_D3D12_DEBUG_FEATURE_ALLOW_BEHAVIOR_CHANGING_DEBUG_AIDS     _D3D12_DEBUG_FEATURE = 0x1
	_D3D12_DEBUG_FEATURE_CONSERVATIVE_RESOURCE_STATE_TRACKING   _D3D12_DEBUG_FEATURE = 0x2
	_D3D12_DEBUG_FEATURE_DISABLE_VIRTUALIZED_BUNDLES_VALIDATION _D3D12_DEBUG_FEATURE = 0x4
)

type _D3D12_DEPTH_WRITE_MASK int32

const (
	_D3D12_DEPTH_WRITE_MASK_ZERO _D3D12_DEPTH_WRITE_MASK = 0
	_D3D12_DEPTH_WRITE_MASK_ALL  _D3D12_DEPTH_WRITE_MASK = 1
)

type _D3D12_DESCRIPTOR_HEAP_TYPE int32

const (
	_D3D12_DESCRIPTOR_HEAP_TYPE_CBV_SRV_UAV _D3D12_DESCRIPTOR_HEAP_TYPE = iota
	_D3D12_DESCRIPTOR_HEAP_TYPE_SAMPLER
	_D3D12_DESCRIPTOR_HEAP_TYPE_RTV
	_D3D12_DESCRIPTOR_HEAP_TYPE_DSV
	_D3D12_DESCRIPTOR_HEAP_TYPE_NUM_TYPES
)

type _D3D12_DESCRIPTOR_HEAP_FLAGS int32

const (
	_D3D12_DESCRIPTOR_HEAP_FLAG_NONE           _D3D12_DESCRIPTOR_HEAP_FLAGS = 0
	_D3D12_DESCRIPTOR_HEAP_FLAG_SHADER_VISIBLE _D3D12_DESCRIPTOR_HEAP_FLAGS = 0x1
)

type _D3D12_DESCRIPTOR_RANGE_TYPE int32

const (
	_D3D12_DESCRIPTOR_RANGE_TYPE_SRV _D3D12_DESCRIPTOR_RANGE_TYPE = iota
	_D3D12_DESCRIPTOR_RANGE_TYPE_UAV
	_D3D12_DESCRIPTOR_RANGE_TYPE_CBV
	_D3D12_DESCRIPTOR_RANGE_TYPE_SAMPLER
)

type _D3D12_DSV_DIMENSION int32

const (
	_D3D12_DSV_DIMENSION_UNKNOWN          _D3D12_DSV_DIMENSION = 0
	_D3D12_DSV_DIMENSION_TEXTURE1D        _D3D12_DSV_DIMENSION = 1
	_D3D12_DSV_DIMENSION_TEXTURE1DARRAY   _D3D12_DSV_DIMENSION = 2
	_D3D12_DSV_DIMENSION_TEXTURE2D        _D3D12_DSV_DIMENSION = 3
	_D3D12_DSV_DIMENSION_TEXTURE2DARRAY   _D3D12_DSV_DIMENSION = 4
	_D3D12_DSV_DIMENSION_TEXTURE2DMS      _D3D12_DSV_DIMENSION = 5
	_D3D12_DSV_DIMENSION_TEXTURE2DMSARRAY _D3D12_DSV_DIMENSION = 6
)

type _D3D12_DSV_FLAGS int32

const (
	_D3D12_DSV_FLAG_NONE              _D3D12_DSV_FLAGS = 0
	_D3D12_DSV_FLAG_READ_ONLY_DEPTH   _D3D12_DSV_FLAGS = 0x1
	_D3D12_DSV_FLAG_READ_ONLY_STENCIL _D3D12_DSV_FLAGS = 0x2
)

type _D3D12_FENCE_FLAGS int32

const (
	_D3D12_FENCE_FLAG_NONE _D3D12_FENCE_FLAGS = 0
)

type _D3D12_FILL_MODE int32

const (
	_D3D12_FILL_MODE_WIREFRAME _D3D12_FILL_MODE = 2
	_D3D12_FILL_MODE_SOLID     _D3D12_FILL_MODE = 3
)

type _D3D12_FILTER int32

const (
	_D3D12_FILTER_MIN_MAG_MIP_POINT _D3D12_FILTER = 0
)

type _D3D12_HEAP_FLAGS int32

const (
	_D3D12_HEAP_FLAG_NONE _D3D12_HEAP_FLAGS = 0
)

type _D3D12_HEAP_TYPE int32

const (
	_D3D12_HEAP_TYPE_DEFAULT  _D3D12_HEAP_TYPE = 1
	_D3D12_HEAP_TYPE_UPLOAD   _D3D12_HEAP_TYPE = 2
	_D3D12_HEAP_TYPE_READBACK _D3D12_HEAP_TYPE = 3
	_D3D12_HEAP_TYPE_CUSTOM   _D3D12_HEAP_TYPE = 4
)

type _D3D12_INDEX_BUFFER_STRIP_CUT_VALUE int32

const (
	_D3D12_INDEX_BUFFER_STRIP_CUT_VALUE_DISABLED   _D3D12_INDEX_BUFFER_STRIP_CUT_VALUE = 0
	_D3D12_INDEX_BUFFER_STRIP_CUT_VALUE_0xFFFF     _D3D12_INDEX_BUFFER_STRIP_CUT_VALUE = 1
	_D3D12_INDEX_BUFFER_STRIP_CUT_VALUE_0xFFFFFFFF _D3D12_INDEX_BUFFER_STRIP_CUT_VALUE = 2
)

type _D3D12_INPUT_CLASSIFICATION int32

const (
	_D3D12_INPUT_CLASSIFICATION_PER_VERTEX_DATA   _D3D12_INPUT_CLASSIFICATION = 0
	_D3D12_INPUT_CLASSIFICATION_PER_INSTANCE_DATA _D3D12_INPUT_CLASSIFICATION = 1
)

type _D3D12_LOGIC_OP int32

const (
	_D3D12_LOGIC_OP_CLEAR _D3D12_LOGIC_OP = iota
	_D3D12_LOGIC_OP_SET
	_D3D12_LOGIC_OP_COPY
	_D3D12_LOGIC_OP_COPY_INVERTED
	_D3D12_LOGIC_OP_NOOP
	_D3D12_LOGIC_OP_INVERT
	_D3D12_LOGIC_OP_AND
	_D3D12_LOGIC_OP_NAND
	_D3D12_LOGIC_OP_OR
	_D3D12_LOGIC_OP_NOR
	_D3D12_LOGIC_OP_XOR
	_D3D12_LOGIC_OP_EQUIV
	_D3D12_LOGIC_OP_AND_REVERSE
	_D3D12_LOGIC_OP_AND_INVERTED
	_D3D12_LOGIC_OP_OR_REVERSE
	_D3D12_LOGIC_OP_OR_INVERTED
)

type _D3D12_MEMORY_POOL int32

const (
	_D3D12_MEMORY_POOL_UNKNOWN _D3D12_MEMORY_POOL = 0
)

type _D3D12_PIPELINE_STATE_FLAGS int32

const (
	_D3D12_PIPELINE_STATE_FLAG_NONE       _D3D12_PIPELINE_STATE_FLAGS = 0
	_D3D12_PIPELINE_STATE_FLAG_TOOL_DEBUG _D3D12_PIPELINE_STATE_FLAGS = 0x1
)

type _D3D12_PRIMITIVE_TOPOLOGY_TYPE int32

const (
	_D3D12_PRIMITIVE_TOPOLOGY_TYPE_UNDEFINED _D3D12_PRIMITIVE_TOPOLOGY_TYPE = 0
	_D3D12_PRIMITIVE_TOPOLOGY_TYPE_POINT     _D3D12_PRIMITIVE_TOPOLOGY_TYPE = 1
	_D3D12_PRIMITIVE_TOPOLOGY_TYPE_LINE      _D3D12_PRIMITIVE_TOPOLOGY_TYPE = 2
	_D3D12_PRIMITIVE_TOPOLOGY_TYPE_TRIANGLE  _D3D12_PRIMITIVE_TOPOLOGY_TYPE = 3
	_D3D12_PRIMITIVE_TOPOLOGY_TYPE_PATCH     _D3D12_PRIMITIVE_TOPOLOGY_TYPE = 4
)

type _D3D12_RESOURCE_BARRIER_FLAGS int32

const (
	_D3D12_RESOURCE_BARRIER_FLAG_NONE _D3D12_RESOURCE_BARRIER_FLAGS = 0
)

type _D3D12_RESOURCE_BARRIER_TYPE int32

const (
	_D3D12_RESOURCE_BARRIER_TYPE_TRANSITION _D3D12_RESOURCE_BARRIER_TYPE = 0
)

type _D3D12_RESOURCE_DIMENSION int32

const (
	_D3D12_RESOURCE_DIMENSION_UNKNOWN   _D3D12_RESOURCE_DIMENSION = 0
	_D3D12_RESOURCE_DIMENSION_BUFFER    _D3D12_RESOURCE_DIMENSION = 1
	_D3D12_RESOURCE_DIMENSION_TEXTURE1D _D3D12_RESOURCE_DIMENSION = 2
	_D3D12_RESOURCE_DIMENSION_TEXTURE2D _D3D12_RESOURCE_DIMENSION = 3
	_D3D12_RESOURCE_DIMENSION_TEXTURE3D _D3D12_RESOURCE_DIMENSION = 4
)

type _D3D12_RESOURCE_FLAGS int32

const (
	_D3D12_RESOURCE_FLAG_NONE                        _D3D12_RESOURCE_FLAGS = 0
	_D3D12_RESOURCE_FLAG_ALLOW_RENDER_TARGET         _D3D12_RESOURCE_FLAGS = 0x1
	_D3D12_RESOURCE_FLAG_ALLOW_DEPTH_STENCIL         _D3D12_RESOURCE_FLAGS = 0x2
	_D3D12_RESOURCE_FLAG_ALLOW_UNORDERED_ACCESS      _D3D12_RESOURCE_FLAGS = 0x4
	_D3D12_RESOURCE_FLAG_DENY_SHADER_RESOURCE        _D3D12_RESOURCE_FLAGS = 0x8
	_D3D12_RESOURCE_FLAG_ALLOW_CROSS_ADAPTER         _D3D12_RESOURCE_FLAGS = 0x10
	_D3D12_RESOURCE_FLAG_ALLOW_SIMULTANEOUS_ACCESS   _D3D12_RESOURCE_FLAGS = 0x20
	_D3D12_RESOURCE_FLAG_VIDEO_DECODE_REFERENCE_ONLY _D3D12_RESOURCE_FLAGS = 0x40
)

type _D3D12_RESOURCE_STATES int32

const (
	_D3D12_RESOURCE_STATE_RENDER_TARGET         _D3D12_RESOURCE_STATES = 0x4
	_D3D12_RESOURCE_STATE_DEPTH_WRITE           _D3D12_RESOURCE_STATES = 0x10
	_D3D12_RESOURCE_STATE_PIXEL_SHADER_RESOURCE _D3D12_RESOURCE_STATES = 0x80
	_D3D12_RESOURCE_STATE_COPY_DEST             _D3D12_RESOURCE_STATES = 0x400
	_D3D12_RESOURCE_STATE_COPY_SOURCE           _D3D12_RESOURCE_STATES = 0x800
	_D3D12_RESOURCE_STATE_GENERIC_READ          _D3D12_RESOURCE_STATES = 0x1 | 0x2 | 0x40 | 0x80 | 0x200 | 0x800
	_D3D12_RESOURCE_STATE_PRESENT               _D3D12_RESOURCE_STATES = 0
)

type _D3D12_ROOT_PARAMETER_TYPE int32

const (
	_D3D12_ROOT_PARAMETER_TYPE_DESCRIPTOR_TABLE _D3D12_ROOT_PARAMETER_TYPE = iota
	_D3D12_ROOT_PARAMETER_TYPE_32BIT_CONSTANTS
	_D3D12_ROOT_PARAMETER_TYPE_CBV
	_D3D12_ROOT_PARAMETER_TYPE_SRV
	_D3D12_ROOT_PARAMETER_TYPE_UAV
)

type _D3D12_ROOT_SIGNATURE_FLAGS int32

const (
	_D3D12_ROOT_SIGNATURE_FLAG_ALLOW_INPUT_ASSEMBLER_INPUT_LAYOUT _D3D12_ROOT_SIGNATURE_FLAGS = 0x1
)

type _D3D12_RTV_DIMENSION int32

const (
	_D3D12_SHADER_COMPONENT_MAPPING_MASK                                     = 0x7
	_D3D12_SHADER_COMPONENT_MAPPING_SHIFT                                    = 3
	_D3D12_SHADER_COMPONENT_MAPPING_ALWAYS_SET_BIT_AVOIDING_ZEROMEM_MISTAKES = 1 << (_D3D12_SHADER_COMPONENT_MAPPING_SHIFT * 4)
	_D3D12_ENCODE_SHADER_4_COMPONENT_MAPPING_0_1_2_3                         = (0 & _D3D12_SHADER_COMPONENT_MAPPING_MASK) |
		((1 & _D3D12_SHADER_COMPONENT_MAPPING_MASK) << _D3D12_SHADER_COMPONENT_MAPPING_SHIFT) |
		((2 & _D3D12_SHADER_COMPONENT_MAPPING_MASK) << (_D3D12_SHADER_COMPONENT_MAPPING_SHIFT * 2)) |
		((3 & _D3D12_SHADER_COMPONENT_MAPPING_MASK) << (_D3D12_SHADER_COMPONENT_MAPPING_SHIFT * 3)) |
		_D3D12_SHADER_COMPONENT_MAPPING_ALWAYS_SET_BIT_AVOIDING_ZEROMEM_MISTAKES
	_D3D12_DEFAULT_SHADER_4_COMPONENT_MAPPING = _D3D12_ENCODE_SHADER_4_COMPONENT_MAPPING_0_1_2_3
)

type _D3D12_SHADER_VISIBILITY int32

const (
	_D3D12_SHADER_VISIBILITY_ALL           _D3D12_SHADER_VISIBILITY = 0
	_D3D12_SHADER_VISIBILITY_VERTEX        _D3D12_SHADER_VISIBILITY = 1
	_D3D12_SHADER_VISIBILITY_HULL          _D3D12_SHADER_VISIBILITY = 2
	_D3D12_SHADER_VISIBILITY_DOMAIN        _D3D12_SHADER_VISIBILITY = 3
	_D3D12_SHADER_VISIBILITY_GEOMETRY      _D3D12_SHADER_VISIBILITY = 4
	_D3D12_SHADER_VISIBILITY_PIXEL         _D3D12_SHADER_VISIBILITY = 5
	_D3D12_SHADER_VISIBILITY_AMPLIFICATION _D3D12_SHADER_VISIBILITY = 6
	_D3D12_SHADER_VISIBILITY_MESH          _D3D12_SHADER_VISIBILITY = 7
)

type _D3D12_SRV_DIMENSION int32

const (
	_D3D12_SRV_DIMENSION_UNKNOWN                           _D3D12_SRV_DIMENSION = 0
	_D3D12_SRV_DIMENSION_BUFFER                            _D3D12_SRV_DIMENSION = 1
	_D3D12_SRV_DIMENSION_TEXTURE1D                         _D3D12_SRV_DIMENSION = 2
	_D3D12_SRV_DIMENSION_TEXTURE1DARRAY                    _D3D12_SRV_DIMENSION = 3
	_D3D12_SRV_DIMENSION_TEXTURE2D                         _D3D12_SRV_DIMENSION = 4
	_D3D12_SRV_DIMENSION_TEXTURE2DARRAY                    _D3D12_SRV_DIMENSION = 5
	_D3D12_SRV_DIMENSION_TEXTURE2DMS                       _D3D12_SRV_DIMENSION = 6
	_D3D12_SRV_DIMENSION_TEXTURE2DMSARRAY                  _D3D12_SRV_DIMENSION = 7
	_D3D12_SRV_DIMENSION_TEXTURE3D                         _D3D12_SRV_DIMENSION = 8
	_D3D12_SRV_DIMENSION_TEXTURECUBE                       _D3D12_SRV_DIMENSION = 9
	_D3D12_SRV_DIMENSION_TEXTURECUBEARRAY                  _D3D12_SRV_DIMENSION = 10
	_D3D12_SRV_DIMENSION_RAYTRACING_ACCELERATION_STRUCTURE _D3D12_SRV_DIMENSION = 11
)

type _D3D12_STATIC_BORDER_COLOR int32

const (
	_D3D12_STATIC_BORDER_COLOR_TRANSPARENT_BLACK _D3D12_STATIC_BORDER_COLOR = 0
)

type _D3D12_STENCIL_OP int32

const (
	_D3D12_STENCIL_OP_KEEP     _D3D12_STENCIL_OP = 1
	_D3D12_STENCIL_OP_ZERO     _D3D12_STENCIL_OP = 2
	_D3D12_STENCIL_OP_REPLACE  _D3D12_STENCIL_OP = 3
	_D3D12_STENCIL_OP_INCR_SAT _D3D12_STENCIL_OP = 4
	_D3D12_STENCIL_OP_DECR_SAT _D3D12_STENCIL_OP = 5
	_D3D12_STENCIL_OP_INVERT   _D3D12_STENCIL_OP = 6
	_D3D12_STENCIL_OP_INCR     _D3D12_STENCIL_OP = 7
	_D3D12_STENCIL_OP_DECR     _D3D12_STENCIL_OP = 8
)

type _D3D12_TEXTURE_ADDRESS_MODE int32

const (
	_D3D12_TEXTURE_ADDRESS_MODE_WRAP        _D3D12_TEXTURE_ADDRESS_MODE = 1
	_D3D12_TEXTURE_ADDRESS_MODE_MIRROR      _D3D12_TEXTURE_ADDRESS_MODE = 2
	_D3D12_TEXTURE_ADDRESS_MODE_CLAMP       _D3D12_TEXTURE_ADDRESS_MODE = 3
	_D3D12_TEXTURE_ADDRESS_MODE_BORDER      _D3D12_TEXTURE_ADDRESS_MODE = 4
	_D3D12_TEXTURE_ADDRESS_MODE_MIRROR_ONCE _D3D12_TEXTURE_ADDRESS_MODE = 5
)

type _D3D12_TEXTURE_COPY_TYPE int32

const (
	_D3D12_TEXTURE_COPY_TYPE_SUBRESOURCE_INDEX _D3D12_TEXTURE_COPY_TYPE = 0
	_D3D12_TEXTURE_COPY_TYPE_PLACED_FOOTPRINT  _D3D12_TEXTURE_COPY_TYPE = 1
)

type _D3D12_TEXTURE_LAYOUT int32

const (
	_D3D12_TEXTURE_LAYOUT_UNKNOWN                _D3D12_TEXTURE_LAYOUT = 0
	_D3D12_TEXTURE_LAYOUT_ROW_MAJOR              _D3D12_TEXTURE_LAYOUT = 1
	_D3D12_TEXTURE_LAYOUT_64KB_UNDEFINED_SWIZZLE _D3D12_TEXTURE_LAYOUT = 2
	_D3D12_TEXTURE_LAYOUT_64KB_STANDARD_SWIZZLE  _D3D12_TEXTURE_LAYOUT = 3
)

type _D3D12XBOX_CREATE_DEVICE_FLAGS int32

type _D3D12XBOX_PROCESS_DEBUG_FLAGS int32

const (
	_D3D12XBOX_PROCESS_DEBUG_FLAG_NONE                          _D3D12XBOX_PROCESS_DEBUG_FLAGS = 0x0
	_D3D12XBOX_PROCESS_DEBUG_FLAG_DEBUG                         _D3D12XBOX_PROCESS_DEBUG_FLAGS = 0x00000001
	_D3D12_PROCESS_DEBUG_FLAG_DEBUG_LAYER_ENABLED               _D3D12XBOX_PROCESS_DEBUG_FLAGS = _D3D12XBOX_PROCESS_DEBUG_FLAG_DEBUG
	_D3D12XBOX_PROCESS_DEBUG_FLAG_INSTRUMENTED                  _D3D12XBOX_PROCESS_DEBUG_FLAGS = 0x00100000
	_D3D12XBOX_PROCESS_DEBUG_FLAG_VALIDATED                     _D3D12XBOX_PROCESS_DEBUG_FLAGS = _D3D12XBOX_PROCESS_DEBUG_FLAG_DEBUG
	_D3D12XBOX_PROCESS_DEBUG_FLAG_NULL_GPU                      _D3D12XBOX_PROCESS_DEBUG_FLAGS = 0x00400000
	_D3D12XBOX_PROCESS_DEBUG_FLAG_ENABLE_COMMON_STATE_PROMOTION _D3D12XBOX_PROCESS_DEBUG_FLAGS = 0x01000000
	_D3D12XBOX_PROCESS_DEBUG_FLAG_ENHANCED_VALIDATION           _D3D12XBOX_PROCESS_DEBUG_FLAGS = 0x02000000
)

type _DXGI_ALPHA_MODE uint32

const (
	_DXGI_ALPHA_MODE_UNSPECIFIED   _DXGI_ALPHA_MODE = 0
	_DXGI_ALPHA_MODE_PREMULTIPLIED _DXGI_ALPHA_MODE = 1
	_DXGI_ALPHA_MODE_STRAIGHT      _DXGI_ALPHA_MODE = 2
	_DXGI_ALPHA_MODE_IGNORE        _DXGI_ALPHA_MODE = 3
	_DXGI_ALPHA_MODE_FORCE_DWORD   _DXGI_ALPHA_MODE = 0xffffffff
)

type _DXGI_FORMAT int32

const (
	_DXGI_FORMAT_UNKNOWN            _DXGI_FORMAT = 0
	_DXGI_FORMAT_R32G32B32A32_FLOAT _DXGI_FORMAT = 2
	_DXGI_FORMAT_R32G32_FLOAT       _DXGI_FORMAT = 16
	_DXGI_FORMAT_R8G8B8A8_UNORM     _DXGI_FORMAT = 28
	_DXGI_FORMAT_D24_UNORM_S8_UINT  _DXGI_FORMAT = 45
	_DXGI_FORMAT_R16_UINT           _DXGI_FORMAT = 57
	_DXGI_FORMAT_B8G8R8A8_UNORM     _DXGI_FORMAT = 87
)

type _DXGI_MODE_SCANLINE_ORDER int32

type _DXGI_MODE_SCALING int32

type _DXGI_SCALING int32

type _DXGI_SWAP_EFFECT int32

const (
	_DXGI_SWAP_EFFECT_FLIP_DISCARD _DXGI_SWAP_EFFECT = 4
)

type _DXGI_USAGE uint32

const (
	_DXGI_USAGE_RENDER_TARGET_OUTPUT _DXGI_USAGE = 1 << (1 + 4)
)

const (
	_DXGI_ADAPTER_FLAG_SOFTWARE = 2

	_DXGI_CREATE_FACTORY_DEBUG = 0x01

	_DXGI_ERROR_NOT_FOUND = windows.Errno(0x887A0002)
)

var (
	_IID_ID3D12CommandAllocator    = windows.GUID{Data1: 0x6102dee4, Data2: 0xaf59, Data3: 0x4b09, Data4: [...]byte{0xb9, 0x99, 0xb4, 0x4d, 0x73, 0xf0, 0x9b, 0x24}}
	_IID_ID3D12CommandQueue        = windows.GUID{Data1: 0x0ec870a6, Data2: 0x5d7e, Data3: 0x4c22, Data4: [...]byte{0x8c, 0xfc, 0x5b, 0xaa, 0xe0, 0x76, 0x16, 0xed}}
	_IID_ID3D12Debug               = windows.GUID{Data1: 0x344488b7, Data2: 0x6846, Data3: 0x474b, Data4: [...]byte{0xb9, 0x89, 0xf0, 0x27, 0x44, 0x82, 0x45, 0xe0}}
	_IID_ID3D12DescriptorHeap      = windows.GUID{Data1: 0x8efb471d, Data2: 0x616c, Data3: 0x4f49, Data4: [...]byte{0x90, 0xf7, 0x12, 0x7b, 0xb7, 0x63, 0xfa, 0x51}}
	_IID_ID3D12DebugCommandList    = windows.GUID{Data1: 0x09e0bf36, Data2: 0x54ac, Data3: 0x484f, Data4: [...]byte{0x88, 0x47, 0x4b, 0xae, 0xea, 0xb6, 0x05, 0x3f}}
	_IID_ID3D12Device              = windows.GUID{Data1: 0x189819f1, Data2: 0x1db6, Data3: 0x4b57, Data4: [...]byte{0xbe, 0x54, 0x18, 0x21, 0x33, 0x9b, 0x85, 0xf7}}
	_IID_ID3D12Fence               = windows.GUID{Data1: 0x0a753dcf, Data2: 0xc4d8, Data3: 0x4b91, Data4: [...]byte{0xad, 0xf6, 0xbe, 0x5a, 0x60, 0xd9, 0x5a, 0x76}}
	_IID_ID3D12GraphicsCommandList = windows.GUID{Data1: 0x5b160d0f, Data2: 0xac1b, Data3: 0x4185, Data4: [...]byte{0x8b, 0xa8, 0xb3, 0xae, 0x42, 0xa5, 0xa4, 0x55}}
	_IID_ID3D12PipelineState       = windows.GUID{Data1: 0x765a30f3, Data2: 0xf624, Data3: 0x4c6f, Data4: [...]byte{0xa8, 0x28, 0xac, 0xe9, 0x48, 0x62, 0x24, 0x45}}
	_IID_ID3D12Resource            = windows.GUID{Data1: 0x696442be, Data2: 0xa72e, Data3: 0x4059, Data4: [...]byte{0xbc, 0x79, 0x5b, 0x5c, 0x98, 0x04, 0x0f, 0xad}}
	_IID_ID3D12RootSignature       = windows.GUID{Data1: 0xc54a6b66, Data2: 0x72df, Data3: 0x4ee8, Data4: [...]byte{0x8b, 0xe5, 0xa9, 0x46, 0xa1, 0x42, 0x92, 0x14}}

	_IID_IDXGIAdapter1 = windows.GUID{Data1: 0x29038f61, Data2: 0x3839, Data3: 0x4626, Data4: [...]byte{0x91, 0xfd, 0x08, 0x68, 0x79, 0x01, 0x1a, 0x05}}
	_IID_IDXGIFactory4 = windows.GUID{Data1: 0x1bc6ea02, Data2: 0xef36, Data3: 0x464f, Data4: [...]byte{0xbf, 0x0c, 0x21, 0xca, 0x39, 0xe5, 0x16, 0x8a}}
)

type _D3D12_BLEND_DESC struct {
	AlphaToCoverageEnable  _BOOL
	IndependentBlendEnable _BOOL
	RenderTarget           [8]_D3D12_RENDER_TARGET_BLEND_DESC
}

type _D3D12_BOX struct {
	left   uint32
	top    uint32
	front  uint32
	right  uint32
	bottom uint32
	back   uint32
}

type _D3D12_CACHED_PIPELINE_STATE struct {
	pCachedBlob           uintptr
	CachedBlobSizeInBytes uintptr
}

type _D3D12_CLEAR_VALUE struct {
	Format _DXGI_FORMAT
	Color  [4]float32 // Union
}

type _D3D12_CONSTANT_BUFFER_VIEW_DESC struct {
	BufferLocation _D3D12_GPU_VIRTUAL_ADDRESS
	SizeInBytes    uint32
}

type _D3D12_CPU_DESCRIPTOR_HANDLE struct {
	ptr uintptr
}

func (h *_D3D12_CPU_DESCRIPTOR_HANDLE) Offset(offsetInDescriptors int32, descriptorIncrementSize uint32) {
	h.ptr += uintptr(offsetInDescriptors) * uintptr(descriptorIncrementSize)
}

type _D3D12_DEPTH_STENCIL_VIEW_DESC struct {
	Format        _DXGI_FORMAT
	ViewDimension _D3D12_DSV_DIMENSION
	Flags         _D3D12_DSV_FLAGS
	_             [4]byte                                      // A padding (TODO: This can be different on 32bit)
	Texture2D     _D3D12_TEX2D_DSV                             // Union
	_             [12 - unsafe.Sizeof(_D3D12_TEX2D_DSV{})]byte // A padding for union
}

type _D3D12_DEPTH_STENCIL_DESC struct {
	DepthEnable      _BOOL
	DepthWriteMask   _D3D12_DEPTH_WRITE_MASK
	DepthFunc        _D3D12_COMPARISON_FUNC
	StencilEnable    _BOOL
	StencilReadMask  uint8
	StencilWriteMask uint8
	FrontFace        _D3D12_DEPTH_STENCILOP_DESC
	BackFace         _D3D12_DEPTH_STENCILOP_DESC
}

type _D3D12_DEPTH_STENCILOP_DESC struct {
	StencilFailOp      _D3D12_STENCIL_OP
	StencilDepthFailOp _D3D12_STENCIL_OP
	StencilPassOp      _D3D12_STENCIL_OP
	StencilFunc        _D3D12_COMPARISON_FUNC
}

type _D3D12_DESCRIPTOR_RANGE struct {
	RangeType                         _D3D12_DESCRIPTOR_RANGE_TYPE
	NumDescriptors                    uint32
	BaseShaderRegister                uint32
	RegisterSpace                     uint32
	OffsetInDescriptorsFromTableStart uint32
}

type _D3D12_GPU_DESCRIPTOR_HANDLE struct {
	ptr uint64
}

func (h *_D3D12_GPU_DESCRIPTOR_HANDLE) Offset(offsetInDescriptors int32, descriptorIncrementSize uint32) {
	h.ptr += uint64(offsetInDescriptors) * uint64(descriptorIncrementSize)
}

type _D3D12_GPU_VIRTUAL_ADDRESS uint64

type _D3D12_GRAPHICS_PIPELINE_STATE_DESC struct {
	pRootSignature        *_ID3D12RootSignature
	VS                    _D3D12_SHADER_BYTECODE
	PS                    _D3D12_SHADER_BYTECODE
	DS                    _D3D12_SHADER_BYTECODE
	HS                    _D3D12_SHADER_BYTECODE
	GS                    _D3D12_SHADER_BYTECODE
	StreamOutput          _D3D12_STREAM_OUTPUT_DESC
	BlendState            _D3D12_BLEND_DESC
	SampleMask            uint32
	RasterizerState       _D3D12_RASTERIZER_DESC
	DepthStencilState     _D3D12_DEPTH_STENCIL_DESC
	InputLayout           _D3D12_INPUT_LAYOUT_DESC
	IBStripCutValue       _D3D12_INDEX_BUFFER_STRIP_CUT_VALUE
	PrimitiveTopologyType _D3D12_PRIMITIVE_TOPOLOGY_TYPE
	NumRenderTargets      uint32
	RTVFormats            [8]_DXGI_FORMAT
	DSVFormat             _DXGI_FORMAT
	SampleDesc            _DXGI_SAMPLE_DESC
	NodeMask              uint32
	CachedPSO             _D3D12_CACHED_PIPELINE_STATE
	Flags                 _D3D12_PIPELINE_STATE_FLAGS
}

type _D3D12_HEAP_PROPERTIES struct {
	Type                 _D3D12_HEAP_TYPE
	CPUPageProperty      _D3D12_CPU_PAGE_PROPERTY
	MemoryPoolPreference _D3D12_MEMORY_POOL
	CreationNodeMask     uint32
	VisibleNodeMask      uint32
}

type _D3D12_INDEX_BUFFER_VIEW struct {
	BufferLocation _D3D12_GPU_VIRTUAL_ADDRESS
	SizeInBytes    uint32
	Format         _DXGI_FORMAT
}

type _D3D12_INPUT_ELEMENT_DESC struct {
	SemanticName         *byte
	SemanticIndex        uint32
	Format               _DXGI_FORMAT
	InputSlot            uint32
	AlignedByteOffset    uint32
	InputSlotClass       _D3D12_INPUT_CLASSIFICATION
	InstanceDataStepRate uint32
}

type _D3D12_INPUT_LAYOUT_DESC struct {
	pInputElementDescs *_D3D12_INPUT_ELEMENT_DESC
	NumElements        uint32
}

type _D3D12_RANGE struct {
	Begin uintptr
	End   uintptr
}

type _D3D12_RASTERIZER_DESC struct {
	FillMode              _D3D12_FILL_MODE
	CullMode              _D3D12_CULL_MODE
	FrontCounterClockwise _BOOL
	DepthBias             int32
	DepthBiasClamp        float32
	SlopeScaledDepthBias  float32
	DepthClipEnable       _BOOL
	MultisampleEnable     _BOOL
	AntialiasedLineEnable _BOOL
	ForcedSampleCount     uint32
	ConservativeRaster    _D3D12_CONSERVATIVE_RASTERIZATION_MODE
}

type _D3D12_RECT struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

type _D3D12_RESOURCE_BARRIER_Transition struct {
	Type       _D3D12_RESOURCE_BARRIER_TYPE
	Flags      _D3D12_RESOURCE_BARRIER_FLAGS
	Transition _D3D12_RESOURCE_TRANSITION_BARRIER
}

type _D3D12_RESOURCE_DESC struct {
	Dimension        _D3D12_RESOURCE_DIMENSION
	Alignment        uint64
	Width            uint64
	Height           uint32
	DepthOrArraySize uint16
	MipLevels        uint16
	Format           _DXGI_FORMAT
	SampleDesc       _DXGI_SAMPLE_DESC
	Layout           _D3D12_TEXTURE_LAYOUT
	Flags            _D3D12_RESOURCE_FLAGS
}

type _D3D12_RESOURCE_TRANSITION_BARRIER struct {
	pResource   *_ID3D12Resource
	Subresource uint32
	StateBefore _D3D12_RESOURCE_STATES
	StateAfter  _D3D12_RESOURCE_STATES
}

type _D3D12_ROOT_DESCRIPTOR_TABLE struct {
	NumDescriptorRanges uint32
	pDescriptorRanges   *_D3D12_DESCRIPTOR_RANGE
}

type _D3D12_ROOT_PARAMETER struct {
	ParameterType    _D3D12_ROOT_PARAMETER_TYPE
	DescriptorTable  _D3D12_ROOT_DESCRIPTOR_TABLE // Union
	ShaderVisibility _D3D12_SHADER_VISIBILITY
}

type _D3D12_ROOT_SIGNATURE_DESC struct {
	NumParameters     uint32
	pParameters       *_D3D12_ROOT_PARAMETER
	NumStaticSamplers uint32
	pStaticSamplers   *_D3D12_STATIC_SAMPLER_DESC
	Flags             _D3D12_ROOT_SIGNATURE_FLAGS
}

type _D3D12_SHADER_BYTECODE struct {
	pShaderBytecode uintptr
	BytecodeLength  uintptr
}

type _D3D12_SHADER_RESOURCE_VIEW_DESC struct {
	Format                  _DXGI_FORMAT
	ViewDimension           _D3D12_SRV_DIMENSION
	Shader4ComponentMapping uint32
	_                       [4]byte                                      // A padding (TODO: This can be different on 32bit)
	Texture2D               _D3D12_TEX2D_SRV                             // Union
	_                       [24 - unsafe.Sizeof(_D3D12_TEX2D_SRV{})]byte // A padding for union
}

type _D3D12_SO_DECLARATION_ENTRY struct {
	Stream         uint32
	SemanticName   *byte
	SemanticIndex  uint32
	StartComponent byte
	ComponentCount byte
	OutputSlot     byte
}

type _D3D12_STATIC_SAMPLER_DESC struct {
	Filter           _D3D12_FILTER
	AddressU         _D3D12_TEXTURE_ADDRESS_MODE
	AddressV         _D3D12_TEXTURE_ADDRESS_MODE
	AddressW         _D3D12_TEXTURE_ADDRESS_MODE
	MipLODBias       float32
	MaxAnisotropy    uint32
	ComparisonFunc   _D3D12_COMPARISON_FUNC
	BorderColor      _D3D12_STATIC_BORDER_COLOR
	MinLOD           float32
	MaxLOD           float32
	ShaderRegister   uint32
	RegisterSpace    uint32
	ShaderVisibility _D3D12_SHADER_VISIBILITY
}

type _D3D12_STREAM_OUTPUT_DESC struct {
	pSODeclaration   *_D3D12_SO_DECLARATION_ENTRY
	NumEntries       uint32
	pBufferStrides   *uint32
	NumStrides       uint32
	RasterizedStream uint32
}

type _D3D12_TEX2D_DSV struct {
	MipSlice uint32
}

type _D3D12_TEX2D_SRV struct {
	MostDetailedMip     uint32
	MipLevels           uint32
	PlaneSlice          uint32
	ResourceMinLODClamp float32
}

type _D3D12_TEXTURE_COPY_LOCATION_PlacedFootPrint struct {
	pResource       *_ID3D12Resource
	Type            _D3D12_TEXTURE_COPY_TYPE
	PlacedFootprint _D3D12_PLACED_SUBRESOURCE_FOOTPRINT
}

type _D3D12_TEXTURE_COPY_LOCATION_SubresourceIndex struct {
	pResource        *_ID3D12Resource
	Type             _D3D12_TEXTURE_COPY_TYPE
	SubresourceIndex uint32
	_                [unsafe.Sizeof(_D3D12_PLACED_SUBRESOURCE_FOOTPRINT{}) - unsafe.Sizeof(uint32(0))]byte // A padding for union
}

type _D3D12_VERTEX_BUFFER_VIEW struct {
	BufferLocation _D3D12_GPU_VIRTUAL_ADDRESS
	SizeInBytes    uint32
	StrideInBytes  uint32
}

type _D3D12_VIEWPORT struct {
	TopLeftX float32
	TopLeftY float32
	Width    float32
	Height   float32
	MinDepth float32
	MaxDepth float32
}

type _D3D12XBOX_CREATE_DEVICE_PARAMETERS struct {
	Version                              uint32
	ProcessDebugFlags                    _D3D12XBOX_PROCESS_DEBUG_FLAGS
	GraphicsCommandQueueRingSizeBytes    uint32
	pOffchipTessellationBuffer           _D3D12_GPU_VIRTUAL_ADDRESS
	GraphicsScratchMemorySizeBytes       uint32
	ComputeScratchMemorySizeBytes        uint32
	DisableGeometryShaderAllocations     _BOOL
	DisableTessellationShaderAllocations _BOOL

	// These members are for Xbox Series.
	DisableDXR                             _BOOL
	DisableAutomaticDPBBBreakBatchEvents   _BOOL
	pDXRStackBuffer                        _D3D12_GPU_VIRTUAL_ADDRESS
	DXRStackBufferOverrideSizeBytes        uint32
	CreateDeviceFlags                      _D3D12XBOX_CREATE_DEVICE_FLAGS
	AutoHDRPaperWhiteLevelNits             uint32
	DisableAutomaticCommandSegmentChaining _BOOL
}

var (
	d3d12       = windows.NewLazySystemDLL("d3d12.dll")
	d3d12x      = windows.NewLazySystemDLL(microsoftgdk.D3D12DLLName())
	d3dcompiler = windows.NewLazySystemDLL("d3dcompiler_47.dll")
	dxgi        = windows.NewLazySystemDLL("dxgi.dll")

	procD3D12CreateDevice           = d3d12.NewProc("D3D12CreateDevice")
	procD3D12GetDebugInterface      = d3d12.NewProc("D3D12GetDebugInterface")
	procD3D12SerializeRootSignature = d3d12.NewProc("D3D12SerializeRootSignature")

	procD3D12XboxCreateDevice = d3d12x.NewProc("D3D12XboxCreateDevice")

	procD3DCompile = d3dcompiler.NewProc("D3DCompile")

	procCreateDXGIFactory2 = dxgi.NewProc("CreateDXGIFactory2")
)

func _D3D12CreateDevice(pAdapter unsafe.Pointer, minimumFeatureLevel _D3D_FEATURE_LEVEL, riid *windows.GUID, ppDevice *unsafe.Pointer) error {
	r, _, _ := procD3D12CreateDevice.Call(uintptr(pAdapter), uintptr(minimumFeatureLevel), uintptr(unsafe.Pointer(riid)), uintptr(unsafe.Pointer(ppDevice)))
	if ppDevice == nil && uint32(r) != uint32(windows.S_FALSE) {
		return fmt.Errorf("directx: D3D12CreateDevice failed: HRESULT(%d)", uint32(r))
	}
	if ppDevice != nil && uint32(r) != uint32(windows.S_OK) {
		return fmt.Errorf("directx: D3D12CreateDevice failed: HRESULT(%d)", uint32(r))
	}
	return nil
}

func _D3D12GetDebugInterface() (*_ID3D12Debug, error) {
	var debug *_ID3D12Debug
	r, _, _ := procD3D12GetDebugInterface.Call(uintptr(unsafe.Pointer(&_IID_ID3D12Debug)), uintptr(unsafe.Pointer(&debug)))
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: D3D12GetDebugInterface failed: HRESULT(%d)", uint32(r))
	}
	return debug, nil
}

func _D3D12SerializeRootSignature(pRootSignature *_D3D12_ROOT_SIGNATURE_DESC, version _D3D_ROOT_SIGNATURE_VERSION) (*_ID3DBlob, error) {
	var blob *_ID3DBlob
	var errorBlob *_ID3DBlob
	r, _, _ := procD3D12SerializeRootSignature.Call(uintptr(unsafe.Pointer(pRootSignature)), uintptr(version), uintptr(unsafe.Pointer(&blob)), uintptr(unsafe.Pointer(&errorBlob)))
	if uint32(r) != uint32(windows.S_OK) {
		if errorBlob != nil {
			defer errorBlob.Release()
			return nil, fmt.Errorf("directx: D3D12SerializeRootSignature failed: %s: HRESULT(%d)", errorBlob.String(), uint32(r))
		}
		return nil, fmt.Errorf("directx: D3D12SerializeRootSignature failed: HRESULT(%d)", uint32(r))
	}
	return blob, nil
}

func _D3D12XboxCreateDevice(pAdapter unsafe.Pointer, pParameters *_D3D12XBOX_CREATE_DEVICE_PARAMETERS, riid *windows.GUID, ppDevice *unsafe.Pointer) error {
	r, _, _ := procD3D12XboxCreateDevice.Call(uintptr(pAdapter), uintptr(unsafe.Pointer(pParameters)), uintptr(unsafe.Pointer(riid)), uintptr(unsafe.Pointer(ppDevice)))
	if uint32(r) != uint32(windows.S_OK) {
		return fmt.Errorf("directx: D3D12XboxCreateDevice failed: HRESULT(%d)", uint32(r))
	}
	return nil
}

func _D3DCompile(srcData []byte, sourceName string, pDefines []_D3D_SHADER_MACRO, pInclude unsafe.Pointer, entryPoint string, target string, flags1 uint32, flags2 uint32) (*_ID3DBlob, error) {
	// TODO: Define _ID3DInclude for pInclude, but is it possible in Go?

	var defs unsafe.Pointer
	if len(pDefines) > 0 {
		defs = unsafe.Pointer(&pDefines[0])
	}
	sourceNameBytes := append([]byte(sourceName), 0)
	entryPointBytes := append([]byte(entryPoint), 0)
	targetBytes := append([]byte(target), 0)
	var code *_ID3DBlob
	var errorMsgs *_ID3DBlob
	r, _, _ := procD3DCompile.Call(
		uintptr(unsafe.Pointer(&srcData[0])), uintptr(len(srcData)), uintptr(unsafe.Pointer(&sourceNameBytes[0])),
		uintptr(defs), uintptr(unsafe.Pointer(pInclude)), uintptr(unsafe.Pointer(&entryPointBytes[0])),
		uintptr(unsafe.Pointer(&targetBytes[0])), uintptr(flags1), uintptr(flags2),
		uintptr(unsafe.Pointer(&code)), uintptr(unsafe.Pointer(&errorMsgs)))
	runtime.KeepAlive(pDefines)
	runtime.KeepAlive(pInclude)
	runtime.KeepAlive(sourceNameBytes)
	runtime.KeepAlive(entryPointBytes)
	runtime.KeepAlive(targetBytes)
	if uint32(r) != uint32(windows.S_OK) {
		if errorMsgs != nil {
			defer errorMsgs.Release()
			return nil, fmt.Errorf("directx: D3DCompile failed: %s: HRESULT(%d)", errorMsgs.String(), uint32(r))
		}
		return nil, fmt.Errorf("directx: D3DCompile failed: HRESULT(%d)", uint32(r))
	}
	return code, nil
}

func _CreateDXGIFactory2(flags uint32) (*_IDXGIFactory4, error) {
	var factory *_IDXGIFactory4
	r, _, _ := procCreateDXGIFactory2.Call(uintptr(flags), uintptr(unsafe.Pointer(&_IID_IDXGIFactory4)), uintptr(unsafe.Pointer(&factory)))
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: CreateDXGIFactory2 failed: HRESULT(%d)", uint32(r))
	}
	return factory, nil
}

type _D3D_SHADER_MACRO struct {
	Name       *byte
	Definition *byte
}

type _D3D12_COMMAND_QUEUE_DESC struct {
	Type     _D3D12_COMMAND_LIST_TYPE
	Priority int32
	Flags    _D3D12_COMMAND_QUEUE_FLAGS
	NodeMask uint32
}

type _D3D12_DESCRIPTOR_HEAP_DESC struct {
	Type           _D3D12_DESCRIPTOR_HEAP_TYPE
	NumDescriptors uint32
	Flags          _D3D12_DESCRIPTOR_HEAP_FLAGS
	NodeMask       uint32
}

type _D3D12_PLACED_SUBRESOURCE_FOOTPRINT struct {
	Offset    uint64
	Footprint _D3D12_SUBRESOURCE_FOOTPRINT
}

type _D3D12_RENDER_TARGET_BLEND_DESC struct {
	BlendEnable           _BOOL
	LogicOpEnable         _BOOL
	SrcBlend              _D3D12_BLEND
	DestBlend             _D3D12_BLEND
	BlendOp               _D3D12_BLEND_OP
	SrcBlendAlpha         _D3D12_BLEND
	DestBlendAlpha        _D3D12_BLEND
	BlendOpAlpha          _D3D12_BLEND_OP
	LogicOp               _D3D12_LOGIC_OP
	RenderTargetWriteMask uint8
}

type _D3D12_RENDER_TARGET_VIEW_DESC struct {
	Format        _DXGI_FORMAT
	ViewDimension _D3D12_RTV_DIMENSION
	_             [3]uint32 // Union: D3D12_BUFFER_RTV seems the biggest
}

type _D3D12_SAMPLER_DESC struct {
	Filter         _D3D12_FILTER
	AddressU       _D3D12_TEXTURE_ADDRESS_MODE
	AddressV       _D3D12_TEXTURE_ADDRESS_MODE
	AddressW       _D3D12_TEXTURE_ADDRESS_MODE
	MipLODBias     float32
	MaxAnisotropy  uint32
	ComparisonFunc _D3D12_COMPARISON_FUNC
	BorderColor    [4]float32
	MinLOD         float32
	MaxLOD         float32
}

type _D3D12_SUBRESOURCE_FOOTPRINT struct {
	Format   _DXGI_FORMAT
	Width    uint32
	Height   uint32
	Depth    uint32
	RowPitch uint32
}

type _DXGI_ADAPTER_DESC1 struct {
	Description           [128]uint16
	VendorId              uint32
	DeviceId              uint32
	SubSysId              uint32
	Revision              uint32
	DedicatedVideoMemory  uint
	DedicatedSystemMemory uint
	SharedSystemMemory    uint
	AdapterLuid           _LUID
	Flags                 uint32
}

type _DXGI_SWAP_CHAIN_FULLSCREEN_DESC struct {
	RefreshRate      _DXGI_RATIONAL
	ScanlineOrdering _DXGI_MODE_SCANLINE_ORDER
	Scaling          _DXGI_MODE_SCALING
	Windowed         _BOOL
}

type _DXGI_RATIONAL struct {
	Numerator   uint32
	Denominator uint32
}

type _DXGI_SAMPLE_DESC struct {
	Count   uint32
	Quality uint32
}

type _DXGI_SWAP_CHAIN_DESC1 struct {
	Width       uint32
	Height      uint32
	Format      _DXGI_FORMAT
	Stereo      _BOOL
	SampleDesc  _DXGI_SAMPLE_DESC
	BufferUsage _DXGI_USAGE
	BufferCount uint32
	Scaling     _DXGI_SCALING
	SwapEffect  _DXGI_SWAP_EFFECT
	AlphaMode   _DXGI_ALPHA_MODE
	Flags       uint32
}

type _LUID struct {
	LowPart  uint32
	HighPart int32
}

type _ID3D12CommandAllocator struct {
	vtbl *_ID3D12CommandAllocator_Vtbl
}

type _ID3D12CommandAllocator_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetPrivateData          uintptr
	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
	SetName                 uintptr
	GetDevice               uintptr
	Reset                   uintptr
}

func (i *_ID3D12CommandAllocator) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}

func (i *_ID3D12CommandAllocator) Reset() error {
	r, _, _ := syscallN(i.vtbl.Reset, uintptr(unsafe.Pointer(i)))
	if uint32(r) != uint32(windows.S_OK) {
		return fmt.Errorf("directx: ID3D12CommandAllocator::Reset failed: HRESULT(%d)", uint32(r))
	}
	return nil
}

type _ID3D12CommandQueue struct {
	vtbl *_ID3D12CommandQueue_Vtbl
}

type _ID3D12CommandQueue_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetPrivateData          uintptr
	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
	SetName                 uintptr
	GetDevice               uintptr
	UpdateTileMappings      uintptr
	CopyTileMappings        uintptr
	ExecuteCommandLists     uintptr
	SetMarker               uintptr
	BeginEvent              uintptr
	EndEvent                uintptr
	Signal                  uintptr
	Wait                    uintptr
	GetTimestampFrequency   uintptr
	GetClockCalibration     uintptr
	GetDesc                 uintptr
}

func (i *_ID3D12CommandQueue) ExecuteCommandLists(ppCommandLists []*_ID3D12GraphicsCommandList) {
	syscallN(i.vtbl.ExecuteCommandLists, uintptr(unsafe.Pointer(i)),
		uintptr(len(ppCommandLists)), uintptr(unsafe.Pointer(&ppCommandLists[0])))
	runtime.KeepAlive(ppCommandLists)
}

func (i *_ID3D12CommandQueue) Signal(signal *_ID3D12Fence, value uint64) error {
	r, _, _ := syscallN(i.vtbl.Signal, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(signal)), uintptr(value))
	runtime.KeepAlive(signal)
	if uint32(r) != uint32(windows.S_OK) {
		return fmt.Errorf("directx: ID3D12CommandQueue::Signal failed: HRESULT(%d)", uint32(r))
	}
	return nil
}

func (i *_ID3D12CommandQueue) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}

type _ID3D12Debug struct {
	vtbl *_ID3D12Debug_Vtbl
}

type _ID3D12Debug_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	EnableDebugLayer uintptr
}

func (i *_ID3D12Debug) As(debug **_ID3D12Debug3) {
	*debug = (*_ID3D12Debug3)(unsafe.Pointer(i))
}

func (i *_ID3D12Debug) EnableDebugLayer() {
	syscallN(i.vtbl.EnableDebugLayer, uintptr(unsafe.Pointer(i)))
}

func (i *_ID3D12Debug) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}

type _ID3D12Debug3 struct {
	vtbl *_ID3D12Debug3_Vtbl
}

type _ID3D12Debug3_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	EnableDebugLayer                            uintptr
	SetEnableGPUBasedValidation                 uintptr
	SetEnableSynchronizedCommandQueueValidation uintptr
	SetGPUBasedValidationFlags                  uintptr
}

func (i *_ID3D12Debug3) SetEnableGPUBasedValidation(enable bool) {
	syscallN(i.vtbl.SetEnableGPUBasedValidation, uintptr(unsafe.Pointer(i)), boolToUintptr(enable))
}

type _ID3D12DebugCommandList struct {
	vtbl *_ID3D12DebugCommandList_Vtbl
}

type _ID3D12DebugCommandList_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	AssertResourceState uintptr
	SetFeatureMask      uintptr
	GetFeatureMask      uintptr
}

func (i *_ID3D12DebugCommandList) SetFeatureMask(mask _D3D12_DEBUG_FEATURE) error {
	r, _, _ := syscallN(i.vtbl.SetFeatureMask, uintptr(unsafe.Pointer(i)), uintptr(mask))
	if uint32(r) != uint32(windows.S_OK) {
		return fmt.Errorf("directx: ID3D12DebugCommandList::SetFeatureMask failed: HRESULT(%d)", uint32(r))
	}
	return nil
}

type _ID3D12DescriptorHeap struct {
	vtbl *_ID3D12DescriptrHeap_Vtbl
}

type _ID3D12DescriptrHeap_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetPrivateData                     uintptr
	SetPrivateData                     uintptr
	SetPrivateDataInterface            uintptr
	SetName                            uintptr
	GetDevice                          uintptr
	GetDesc                            uintptr
	GetCPUDescriptorHandleForHeapStart uintptr
	GetGPUDescriptorHandleForHeapStart uintptr
}

func (i *_ID3D12DescriptorHeap) GetCPUDescriptorHandleForHeapStart() _D3D12_CPU_DESCRIPTOR_HANDLE {
	// There is a bug in the header file:
	// https://stackoverflow.com/questions/34118929/getcpudescriptorhandleforheapstart-stack-corruption
	var handle _D3D12_CPU_DESCRIPTOR_HANDLE
	syscallN(i.vtbl.GetCPUDescriptorHandleForHeapStart, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(&handle)))
	return handle
}

func (i *_ID3D12DescriptorHeap) GetGPUDescriptorHandleForHeapStart() _D3D12_GPU_DESCRIPTOR_HANDLE {
	// This has the same issue as GetCPUDescriptorHandleForHeapStart.
	var handle _D3D12_GPU_DESCRIPTOR_HANDLE
	syscallN(i.vtbl.GetGPUDescriptorHandleForHeapStart, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(&handle)))
	return handle
}

func (i *_ID3D12DescriptorHeap) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}

type _ID3D12Device struct {
	vtbl *_ID3D12Device_Vtbl
}

type _ID3D12Device_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetPrivateData                   uintptr
	SetPrivateData                   uintptr
	SetPrivateDataInterface          uintptr
	SetName                          uintptr
	GetNodeCount                     uintptr
	CreateCommandQueue               uintptr
	CreateCommandAllocator           uintptr
	CreateGraphicsPipelineState      uintptr
	CreateComputePipelineState       uintptr
	CreateCommandList                uintptr
	CheckFeatureSupport              uintptr
	CreateDescriptorHeap             uintptr
	GetDescriptorHandleIncrementSize uintptr
	CreateRootSignature              uintptr
	CreateConstantBufferView         uintptr
	CreateShaderResourceView         uintptr
	CreateUnorderedAccessView        uintptr
	CreateRenderTargetView           uintptr
	CreateDepthStencilView           uintptr
	CreateSampler                    uintptr
	CopyDescriptors                  uintptr
	CopyDescriptorsSimple            uintptr
	GetResourceAllocationInfo        uintptr
	GetCustomHeapProperties          uintptr
	CreateCommittedResource          uintptr
	CreateHeap                       uintptr
	CreatePlacedResource             uintptr
	CreateReservedResource           uintptr
	CreateSharedHandle               uintptr
	OpenSharedHandle                 uintptr
	OpenSharedHandleByName           uintptr
	MakeResident                     uintptr
	Evict                            uintptr
	CreateFence                      uintptr
	GetDeviceRemovedReason           uintptr
	GetCopyableFootprints            uintptr
	CreateQueryHeap                  uintptr
	SetStablePowerState              uintptr
	CreateCommandSignature           uintptr
	GetResourceTiling                uintptr
	GetAdapterLuid                   uintptr
}

func (i *_ID3D12Device) CreateCommandAllocator(typ _D3D12_COMMAND_LIST_TYPE) (*_ID3D12CommandAllocator, error) {
	var commandAllocator *_ID3D12CommandAllocator
	r, _, _ := syscallN(i.vtbl.CreateCommandAllocator, uintptr(unsafe.Pointer(i)),
		uintptr(typ), uintptr(unsafe.Pointer(&_IID_ID3D12CommandAllocator)), uintptr(unsafe.Pointer(&commandAllocator)))
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: ID3D12Device::CreateCommandAllocator failed: HRESULT(%d)", uint32(r))
	}
	return commandAllocator, nil
}

func (i *_ID3D12Device) CreateCommandList(nodeMask uint32, typ _D3D12_COMMAND_LIST_TYPE, pCommandAllocator *_ID3D12CommandAllocator, pInitialState *_ID3D12PipelineState) (*_ID3D12GraphicsCommandList, error) {
	var commandList *_ID3D12GraphicsCommandList
	r, _, _ := syscallN(i.vtbl.CreateCommandList,
		uintptr(unsafe.Pointer(i)), uintptr(nodeMask), uintptr(typ),
		uintptr(unsafe.Pointer(pCommandAllocator)), uintptr(unsafe.Pointer(pInitialState)), uintptr(unsafe.Pointer(&_IID_ID3D12GraphicsCommandList)),
		uintptr(unsafe.Pointer(&commandList)))
	runtime.KeepAlive(pCommandAllocator)
	runtime.KeepAlive(pInitialState)
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: ID3D12Device::CreateCommandList failed: HRESULT(%d)", uint32(r))
	}
	return commandList, nil
}

func (i *_ID3D12Device) CreateCommittedResource(pHeapProperties *_D3D12_HEAP_PROPERTIES, heapFlags _D3D12_HEAP_FLAGS, pDesc *_D3D12_RESOURCE_DESC, initialResourceState _D3D12_RESOURCE_STATES, pOptimizedClearValue *_D3D12_CLEAR_VALUE) (*_ID3D12Resource, error) {
	var resource *_ID3D12Resource
	r, _, _ := syscallN(i.vtbl.CreateCommittedResource,
		uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(pHeapProperties)), uintptr(heapFlags),
		uintptr(unsafe.Pointer(pDesc)), uintptr(initialResourceState), uintptr(unsafe.Pointer(pOptimizedClearValue)),
		uintptr(unsafe.Pointer(&_IID_ID3D12Resource)), uintptr(unsafe.Pointer(&resource)))
	runtime.KeepAlive(pHeapProperties)
	runtime.KeepAlive(pDesc)
	runtime.KeepAlive(pOptimizedClearValue)
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: ID3D12Device::CreateCommittedResource failed: HRESULT(%d)", uint32(r))
	}
	return resource, nil
}

func (i *_ID3D12Device) CreateCommandQueue(desc *_D3D12_COMMAND_QUEUE_DESC) (*_ID3D12CommandQueue, error) {
	var commandQueue *_ID3D12CommandQueue
	r, _, _ := syscallN(i.vtbl.CreateCommandQueue, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(desc)), uintptr(unsafe.Pointer(&_IID_ID3D12CommandQueue)), uintptr(unsafe.Pointer(&commandQueue)))
	runtime.KeepAlive(desc)
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: ID3D12Device::CreateCommandQueue failed: HRESULT(%d)", uint32(r))
	}
	return commandQueue, nil
}

func (i *_ID3D12Device) CreateConstantBufferView(pDesc *_D3D12_CONSTANT_BUFFER_VIEW_DESC, destDescriptor _D3D12_CPU_DESCRIPTOR_HANDLE) {
	syscallN(i.vtbl.CreateConstantBufferView, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(pDesc)), uintptr(destDescriptor.ptr))
	runtime.KeepAlive(pDesc)
}

func (i *_ID3D12Device) CreateDescriptorHeap(desc *_D3D12_DESCRIPTOR_HEAP_DESC) (*_ID3D12DescriptorHeap, error) {
	var descriptorHeap *_ID3D12DescriptorHeap
	r, _, _ := syscallN(i.vtbl.CreateDescriptorHeap, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(desc)), uintptr(unsafe.Pointer(&_IID_ID3D12DescriptorHeap)), uintptr(unsafe.Pointer(&descriptorHeap)))
	runtime.KeepAlive(desc)
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: ID3D12Device::CreateDescriptorHeap failed: HRESULT(%d)", uint32(r))
	}
	return descriptorHeap, nil
}

func (i *_ID3D12Device) CreateDepthStencilView(pResource *_ID3D12Resource, pDesc *_D3D12_DEPTH_STENCIL_VIEW_DESC, destDescriptor _D3D12_CPU_DESCRIPTOR_HANDLE) {
	syscallN(i.vtbl.CreateDepthStencilView, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(pResource)), uintptr(unsafe.Pointer(pDesc)), destDescriptor.ptr)
	runtime.KeepAlive(pResource)
	runtime.KeepAlive(pDesc)
}

func (i *_ID3D12Device) CreateFence(initialValue uint64, flags _D3D12_FENCE_FLAGS) (*_ID3D12Fence, error) {
	// TODO: Does this work on a 32bit machine?
	var fence *_ID3D12Fence
	r, _, _ := syscallN(i.vtbl.CreateFence, uintptr(unsafe.Pointer(i)),
		uintptr(initialValue), uintptr(flags), uintptr(unsafe.Pointer(&_IID_ID3D12Fence)), uintptr(unsafe.Pointer(&fence)))
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: ID3D12Device::CreateFence failed: HRESULT(%d)", uint32(r))
	}
	return fence, nil
}

func (i *_ID3D12Device) CreateGraphicsPipelineState(pDesc *_D3D12_GRAPHICS_PIPELINE_STATE_DESC) (*_ID3D12PipelineState, error) {
	var pipelineState *_ID3D12PipelineState
	r, _, _ := syscallN(i.vtbl.CreateGraphicsPipelineState, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(pDesc)), uintptr(unsafe.Pointer(&_IID_ID3D12PipelineState)), uintptr(unsafe.Pointer(&pipelineState)))
	runtime.KeepAlive(pDesc)
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: ID3D12Device::CreateGraphicsPipelineState failed: HRESULT(%d)", uint32(r))
	}
	return pipelineState, nil
}

func (i *_ID3D12Device) CreateRenderTargetView(pResource *_ID3D12Resource, pDesc *_D3D12_RENDER_TARGET_VIEW_DESC, destDescriptor _D3D12_CPU_DESCRIPTOR_HANDLE) {
	syscallN(i.vtbl.CreateRenderTargetView, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(pResource)), uintptr(unsafe.Pointer(pDesc)), destDescriptor.ptr)
	runtime.KeepAlive(pResource)
	runtime.KeepAlive(pDesc)
}

func (i *_ID3D12Device) CreateRootSignature(nodeMask uint32, pBlobWithRootSignature uintptr, blobLengthInBytes uintptr) (*_ID3D12RootSignature, error) {
	var signature *_ID3D12RootSignature
	r, _, _ := syscallN(i.vtbl.CreateRootSignature, uintptr(unsafe.Pointer(i)),
		uintptr(nodeMask), pBlobWithRootSignature, blobLengthInBytes,
		uintptr(unsafe.Pointer(&_IID_ID3D12RootSignature)), uintptr(unsafe.Pointer(&signature)))
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: ID3D12Device::CreateRootSignature failed: HRESULT(%d)", uint32(r))
	}
	return signature, nil
}

func (i *_ID3D12Device) CreateSampler(pDesc *_D3D12_SAMPLER_DESC, destDescriptor _D3D12_CPU_DESCRIPTOR_HANDLE) {
	syscallN(i.vtbl.CreateSampler, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(pDesc)), destDescriptor.ptr)
	runtime.KeepAlive(pDesc)
}

func (i *_ID3D12Device) CreateShaderResourceView(pResource *_ID3D12Resource, pDesc *_D3D12_SHADER_RESOURCE_VIEW_DESC, destDescriptor _D3D12_CPU_DESCRIPTOR_HANDLE) {
	syscallN(i.vtbl.CreateShaderResourceView, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(pResource)), uintptr(unsafe.Pointer(pDesc)), destDescriptor.ptr)
	runtime.KeepAlive(pResource)
	runtime.KeepAlive(pDesc)
}

func (i *_ID3D12Device) GetCopyableFootprints(pResourceDesc *_D3D12_RESOURCE_DESC, firstSubresource uint32, numSubresources uint32, baseOffset uint64) (layouts _D3D12_PLACED_SUBRESOURCE_FOOTPRINT, numRows uint, rowSizeInBytes uint64, totalBytes uint64) {
	syscallN(i.vtbl.GetCopyableFootprints, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(pResourceDesc)), uintptr(firstSubresource), uintptr(numSubresources),
		uintptr(baseOffset), uintptr(unsafe.Pointer(&layouts)), uintptr(unsafe.Pointer(&numRows)),
		uintptr(unsafe.Pointer(&rowSizeInBytes)), uintptr(unsafe.Pointer(&totalBytes)))
	runtime.KeepAlive(pResourceDesc)
	return
}

func (i *_ID3D12Device) GetDescriptorHandleIncrementSize(descriptorHeapType _D3D12_DESCRIPTOR_HEAP_TYPE) uint32 {
	r, _, _ := syscallN(i.vtbl.GetDescriptorHandleIncrementSize, uintptr(unsafe.Pointer(i)), uintptr(descriptorHeapType))
	return uint32(r)
}

func (i *_ID3D12Device) GetDeviceRemovedReason() error {
	r, _, _ := syscallN(i.vtbl.GetDeviceRemovedReason, uintptr(unsafe.Pointer(i)))
	if uint32(r) != uint32(windows.S_OK) {
		return fmt.Errorf("directx: ID3D12Device::GetDeviceRemovedReason failed: HRESULT(%d)", uint32(r))
	}
	return nil
}

type _ID3D12Fence struct {
	vtbl *_ID3D12Fence_Vtbl
}

type _ID3D12Fence_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetPrivateData          uintptr
	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
	SetName                 uintptr
	GetDevice               uintptr
	GetCompletedValue       uintptr
	SetEventOnCompletion    uintptr
	Signal                  uintptr
}

func (i *_ID3D12Fence) GetCompletedValue() uint64 {
	// TODO: Does this work on a 32bit machine?
	r, _, _ := syscallN(i.vtbl.GetCompletedValue, uintptr(unsafe.Pointer(i)))
	return uint64(r)
}

func (i *_ID3D12Fence) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}

func (i *_ID3D12Fence) SetEventOnCompletion(value uint64, hEvent windows.Handle) error {
	// TODO: Does this work on a 32bit machine?
	r, _, _ := syscallN(i.vtbl.SetEventOnCompletion, uintptr(unsafe.Pointer(i)), uintptr(value), uintptr(hEvent))
	if uint32(r) != uint32(windows.S_OK) {
		return fmt.Errorf("directx: ID3D12Fence::SetEventOnCompletion failed: HRESULT(%d)", uint32(r))
	}
	return nil
}

type _ID3D12GraphicsCommandList struct {
	vtbl *_ID3D12GraphicsCommandList_Vtbl
}

type _ID3D12GraphicsCommandList_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetPrivateData                     uintptr
	SetPrivateData                     uintptr
	SetPrivateDataInterface            uintptr
	SetName                            uintptr
	GetDevice                          uintptr
	GetType                            uintptr
	Close                              uintptr
	Reset                              uintptr
	ClearState                         uintptr
	DrawInstanced                      uintptr
	DrawIndexedInstanced               uintptr
	Dispatch                           uintptr
	CopyBufferRegion                   uintptr
	CopyTextureRegion                  uintptr
	CopyResource                       uintptr
	CopyTiles                          uintptr
	ResolveSubresource                 uintptr
	IASetPrimitiveTopology             uintptr
	RSSetViewports                     uintptr
	RSSetScissorRects                  uintptr
	OMSetBlendFactor                   uintptr
	OMSetStencilRef                    uintptr
	SetPipelineState                   uintptr
	ResourceBarrier                    uintptr
	ExecuteBundle                      uintptr
	SetDescriptorHeaps                 uintptr
	SetComputeRootSignature            uintptr
	SetGraphicsRootSignature           uintptr
	SetComputeRootDescriptorTable      uintptr
	SetGraphicsRootDescriptorTable     uintptr
	SetComputeRoot32BitConstant        uintptr
	SetGraphicsRoot32BitConstant       uintptr
	SetComputeRoot32BitConstants       uintptr
	SetGraphicsRoot32BitConstants      uintptr
	SetComputeRootConstantBufferView   uintptr
	SetGraphicsRootConstantBufferView  uintptr
	SetComputeRootShaderResourceView   uintptr
	SetGraphicsRootShaderResourceView  uintptr
	SetComputeRootUnorderedAccessView  uintptr
	SetGraphicsRootUnorderedAccessView uintptr
	IASetIndexBuffer                   uintptr
	IASetVertexBuffers                 uintptr
	SOSetTargets                       uintptr
	OMSetRenderTargets                 uintptr
	ClearDepthStencilView              uintptr
	ClearRenderTargetView              uintptr
	ClearUnorderedAccessViewUint       uintptr
	ClearUnorderedAccessViewFloat      uintptr
	DiscardResource                    uintptr
	BeginQuery                         uintptr
	EndQuery                           uintptr
	ResolveQueryData                   uintptr
	SetPredication                     uintptr
	SetMarker                          uintptr
	BeginEvent                         uintptr
	EndEvent                           uintptr
	ExecuteIndirect                    uintptr
}

func (i *_ID3D12GraphicsCommandList) ClearDepthStencilView(depthStencilView _D3D12_CPU_DESCRIPTOR_HANDLE, clearFlags _D3D12_CLEAR_FLAGS, depth float32, stencil uint8, numRects uint32, pRects *_D3D12_RECT) {
	syscallN(i.vtbl.ClearDepthStencilView, uintptr(unsafe.Pointer(i)),
		depthStencilView.ptr, uintptr(clearFlags), uintptr(math.Float32bits(depth)),
		uintptr(stencil), uintptr(numRects), uintptr(unsafe.Pointer(pRects)))
	runtime.KeepAlive(pRects)
}

func (i *_ID3D12GraphicsCommandList) ClearRenderTargetView(pRenderTargetView _D3D12_CPU_DESCRIPTOR_HANDLE, colorRGBA [4]float32, numRects uint32, pRects *_D3D12_RECT) {
	syscallN(i.vtbl.ClearRenderTargetView, uintptr(unsafe.Pointer(i)),
		pRenderTargetView.ptr, uintptr(unsafe.Pointer(&colorRGBA[0])), uintptr(numRects), uintptr(unsafe.Pointer(pRects)))
	runtime.KeepAlive(pRenderTargetView)
}

func (i *_ID3D12GraphicsCommandList) Close() error {
	r, _, _ := syscallN(i.vtbl.Close, uintptr(unsafe.Pointer(i)))
	if uint32(r) != uint32(windows.S_OK) {
		return fmt.Errorf("directx: ID3D12GraphicsCommandList::Close failed: HRESULT(%d)", uint32(r))
	}
	return nil
}

func (i *_ID3D12GraphicsCommandList) CopyTextureRegion_PlacedFootPrint_SubresourceIndex(pDst *_D3D12_TEXTURE_COPY_LOCATION_PlacedFootPrint, dstX uint32, dstY uint32, dstZ uint32, pSrc *_D3D12_TEXTURE_COPY_LOCATION_SubresourceIndex, pSrcBox *_D3D12_BOX) {
	syscallN(i.vtbl.CopyTextureRegion, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(pDst)), uintptr(dstX), uintptr(dstY),
		uintptr(dstZ), uintptr(unsafe.Pointer(pSrc)), uintptr(unsafe.Pointer(pSrcBox)))
	runtime.KeepAlive(pDst)
	runtime.KeepAlive(pSrc)
	runtime.KeepAlive(pSrcBox)
}

func (i *_ID3D12GraphicsCommandList) CopyTextureRegion_SubresourceIndex_PlacedFootPrint(pDst *_D3D12_TEXTURE_COPY_LOCATION_SubresourceIndex, dstX uint32, dstY uint32, dstZ uint32, pSrc *_D3D12_TEXTURE_COPY_LOCATION_PlacedFootPrint, pSrcBox *_D3D12_BOX) {
	syscallN(i.vtbl.CopyTextureRegion, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(pDst)), uintptr(dstX), uintptr(dstY),
		uintptr(dstZ), uintptr(unsafe.Pointer(pSrc)), uintptr(unsafe.Pointer(pSrcBox)))
	runtime.KeepAlive(pDst)
	runtime.KeepAlive(pSrc)
	runtime.KeepAlive(pSrcBox)
}

func (i *_ID3D12GraphicsCommandList) DrawIndexedInstanced(indexCountPerInstance uint32, instanceCount uint32, startIndexLocation uint32, baseVertexLocation int32, startInstanceLocation uint32) {
	syscallN(i.vtbl.DrawIndexedInstanced, uintptr(unsafe.Pointer(i)),
		uintptr(indexCountPerInstance), uintptr(instanceCount), uintptr(startIndexLocation), uintptr(baseVertexLocation), uintptr(startInstanceLocation))
}

func (i *_ID3D12GraphicsCommandList) IASetIndexBuffer(pView *_D3D12_INDEX_BUFFER_VIEW) {
	syscallN(i.vtbl.IASetIndexBuffer, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(pView)))
	runtime.KeepAlive(pView)
}

func (i *_ID3D12GraphicsCommandList) IASetPrimitiveTopology(primitiveTopology _D3D_PRIMITIVE_TOPOLOGY) {
	syscallN(i.vtbl.IASetPrimitiveTopology, uintptr(unsafe.Pointer(i)), uintptr(primitiveTopology))
}

func (i *_ID3D12GraphicsCommandList) IASetVertexBuffers(startSlot uint32, numViews uint32, pViews *_D3D12_VERTEX_BUFFER_VIEW) {
	syscallN(i.vtbl.IASetVertexBuffers, uintptr(unsafe.Pointer(i)),
		uintptr(startSlot), uintptr(numViews), uintptr(unsafe.Pointer(pViews)))
	runtime.KeepAlive(pViews)
}

func (i *_ID3D12GraphicsCommandList) OMSetRenderTargets(numRenderTargetDescriptors uint32, pRenderTargetDescriptors *_D3D12_CPU_DESCRIPTOR_HANDLE, rtsSingleHandleToDescriptorRange bool, pDepthStencilDescriptor *_D3D12_CPU_DESCRIPTOR_HANDLE) {
	syscallN(i.vtbl.OMSetRenderTargets, uintptr(unsafe.Pointer(i)),
		uintptr(numRenderTargetDescriptors), uintptr(unsafe.Pointer(pRenderTargetDescriptors)), boolToUintptr(rtsSingleHandleToDescriptorRange), uintptr(unsafe.Pointer(pDepthStencilDescriptor)))
	runtime.KeepAlive(pRenderTargetDescriptors)
	runtime.KeepAlive(pDepthStencilDescriptor)
}

func (i *_ID3D12GraphicsCommandList) OMSetStencilRef(stencilRef uint32) {
	syscallN(i.vtbl.OMSetStencilRef, uintptr(unsafe.Pointer(i)), uintptr(stencilRef))
}

func (i *_ID3D12GraphicsCommandList) QueryInterface(riid *windows.GUID, ppvObject *unsafe.Pointer) error {
	r, _, _ := syscallN(i.vtbl.QueryInterface, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(riid)), uintptr(unsafe.Pointer(ppvObject)))
	runtime.KeepAlive(riid)
	if uint32(r) != uint32(windows.S_OK) {
		return fmt.Errorf("directx: ID3D12GraphicsCommandList::QueryInterface failed: HRESULT(%d)", uint32(r))
	}
	return nil
}

func (i *_ID3D12GraphicsCommandList) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}

func (i *_ID3D12GraphicsCommandList) Reset(pAllocator *_ID3D12CommandAllocator, pInitialState *_ID3D12PipelineState) error {
	r, _, _ := syscallN(i.vtbl.Reset, uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(pAllocator)), uintptr(unsafe.Pointer(pInitialState)))
	runtime.KeepAlive(pAllocator)
	runtime.KeepAlive(pInitialState)
	if uint32(r) != uint32(windows.S_OK) {
		return fmt.Errorf("directx: ID3D12GraphicsCommandList::Reset failed: HRESULT(%d)", uint32(r))
	}
	return nil
}

func (i *_ID3D12GraphicsCommandList) ResourceBarrier(numBarriers uint32, pBarriers *_D3D12_RESOURCE_BARRIER_Transition) {
	syscallN(i.vtbl.ResourceBarrier, uintptr(unsafe.Pointer(i)),
		uintptr(numBarriers), uintptr(unsafe.Pointer(pBarriers)))
	runtime.KeepAlive(pBarriers)
}

func (i *_ID3D12GraphicsCommandList) RSSetViewports(numViewports uint32, pViewports *_D3D12_VIEWPORT) {
	syscallN(i.vtbl.RSSetViewports, uintptr(unsafe.Pointer(i)),
		uintptr(numViewports), uintptr(unsafe.Pointer(pViewports)))
	runtime.KeepAlive(pViewports)
}

func (i *_ID3D12GraphicsCommandList) RSSetScissorRects(numRects uint32, pRects *_D3D12_RECT) {
	syscallN(i.vtbl.RSSetScissorRects, uintptr(unsafe.Pointer(i)),
		uintptr(numRects), uintptr(unsafe.Pointer(pRects)))
	runtime.KeepAlive(pRects)
}

func (i *_ID3D12GraphicsCommandList) SetDescriptorHeaps(ppDescriptorHeaps []*_ID3D12DescriptorHeap) {
	syscallN(i.vtbl.SetDescriptorHeaps, uintptr(unsafe.Pointer(i)),
		uintptr(len(ppDescriptorHeaps)), uintptr(unsafe.Pointer(&ppDescriptorHeaps[0])))
	runtime.KeepAlive(ppDescriptorHeaps)
}

func (i *_ID3D12GraphicsCommandList) SetGraphicsRootDescriptorTable(rootParameterIndex uint32, baseDescriptor _D3D12_GPU_DESCRIPTOR_HANDLE) {
	syscallN(i.vtbl.SetGraphicsRootDescriptorTable, uintptr(unsafe.Pointer(i)),
		uintptr(rootParameterIndex), uintptr(baseDescriptor.ptr))
}

func (i *_ID3D12GraphicsCommandList) SetGraphicsRootSignature(pRootSignature *_ID3D12RootSignature) {
	syscallN(i.vtbl.SetGraphicsRootSignature, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(pRootSignature)))
	runtime.KeepAlive(pRootSignature)
}

func (i *_ID3D12GraphicsCommandList) SetPipelineState(pPipelineState *_ID3D12PipelineState) {
	syscallN(i.vtbl.SetPipelineState, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(pPipelineState)))
	runtime.KeepAlive(pPipelineState)
}

type _ID3D12PipelineState struct {
	vtbl *_ID3D12PipelineState_Vtbl
}

type _ID3D12PipelineState_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetPrivateData          uintptr
	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
	SetName                 uintptr
	GetDevice               uintptr
	GetCachedBlob           uintptr
}

func (i *_ID3D12PipelineState) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}

type _ID3D12Resource struct {
	vtbl *_ID3D12Resource_Vtbl
}

type _ID3D12Resource_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetPrivateData          uintptr
	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
	SetName                 uintptr
	GetDevice               uintptr
	Map                     uintptr
	Unmap                   uintptr
	GetDesc                 uintptr
	GetGPUVirtualAddress    uintptr
	WriteToSubresource      uintptr
}

func (i *_ID3D12Resource) GetDesc() _D3D12_RESOURCE_DESC {
	var resourceDesc _D3D12_RESOURCE_DESC
	syscallN(i.vtbl.GetDesc, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(&resourceDesc)))
	return resourceDesc
}

func (i *_ID3D12Resource) GetGPUVirtualAddress() _D3D12_GPU_VIRTUAL_ADDRESS {
	r, _, _ := syscallN(i.vtbl.GetGPUVirtualAddress, uintptr(unsafe.Pointer(i)))
	return _D3D12_GPU_VIRTUAL_ADDRESS(r)
}

func (i *_ID3D12Resource) Map(subresource uint32, pReadRange *_D3D12_RANGE) (uintptr, error) {
	var retryCount int
retry:
	var data uintptr
	r, _, _ := syscallN(i.vtbl.Map, uintptr(unsafe.Pointer(i)),
		uintptr(subresource), uintptr(unsafe.Pointer(pReadRange)), uintptr(unsafe.Pointer(&data)))
	runtime.KeepAlive(pReadRange)
	if uint32(r) != uint32(windows.S_OK) {
		return 0, fmt.Errorf("directx: ID3D12Resource::Map failed: HRESULT(%d)", uint32(r))
	}
	if data == 0 {
		// This is very mysterious, but sometimes Map fails especially on tests with Warp and/or Proton (Steam Deck) (#2113).
		if retryCount >= 5 {
			return 0, fmt.Errorf("directx: ID3D12Resource::Map failed: nothing is mapped")
		}
		retryCount++
		goto retry
	}
	return data, nil
}

func (i *_ID3D12Resource) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}

func (i *_ID3D12Resource) Unmap(subresource uint32, pWrittenRange *_D3D12_RANGE) {
	syscallN(i.vtbl.Unmap, uintptr(unsafe.Pointer(i)),
		uintptr(subresource), uintptr(unsafe.Pointer(pWrittenRange)))
	runtime.KeepAlive(pWrittenRange)
}

type _ID3DBlob struct {
	vtbl *_ID3DBlob_Vtbl
}

type _ID3DBlob_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetBufferPointer uintptr
	GetBufferSize    uintptr
}

func (i *_ID3DBlob) GetBufferPointer() uintptr {
	r, _, _ := syscallN(i.vtbl.GetBufferPointer, uintptr(unsafe.Pointer(i)))
	return r
}

func (i *_ID3DBlob) GetBufferSize() uintptr {
	r, _, _ := syscallN(i.vtbl.GetBufferSize, uintptr(unsafe.Pointer(i)))
	return r
}

func (i *_ID3DBlob) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}

func (i *_ID3DBlob) String() string {
	var str string
	h := (*reflect.StringHeader)(unsafe.Pointer(&str))
	h.Data = i.GetBufferPointer()
	h.Len = int(i.GetBufferSize())
	return str
}

type _IDXGIAdapter1 struct {
	vtbl *_IDXGIAdapter1_Vtbl
}

type _IDXGIAdapter1_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
	GetPrivateData          uintptr
	GetParent               uintptr
	EnumOutputs             uintptr
	GetDesc                 uintptr
	CheckInterfaceSupport   uintptr
	GetDesc1                uintptr
}

func (i *_IDXGIAdapter1) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}

func (i *_IDXGIAdapter1) GetDesc1() (*_DXGI_ADAPTER_DESC1, error) {
	var desc _DXGI_ADAPTER_DESC1
	r, _, _ := syscallN(i.vtbl.GetDesc1, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(&desc)))
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: IDXGIAdapter1::GetDesc1 failed: HRESULT(%d)", uint32(r))
	}
	return &desc, nil
}

type _IDXGIFactory4 struct {
	vtbl *_IDXGIFactory4_Vtbl
}

type _IDXGIFactory4_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	SetPrivateData                uintptr
	SetPrivateDataInterface       uintptr
	GetPrivateData                uintptr
	GetParent                     uintptr
	EnumAdapters                  uintptr
	MakeWindowAssociation         uintptr
	GetWindowAssociation          uintptr
	CreateSwapChain               uintptr
	CreateSoftwareAdapter         uintptr
	EnumAdapters1                 uintptr
	IsCurrent                     uintptr
	IsWindowedStereoEnabled       uintptr
	CreateSwapChainForHwnd        uintptr
	CreateSwapChainForCoreWindow  uintptr
	GetSharedResourceAdapterLuid  uintptr
	RegisterStereoStatusWindow    uintptr
	RegisterStereoStatusEvent     uintptr
	UnregisterStereoStatus        uintptr
	RegisterOcclusionStatusWindow uintptr
	RegisterOcclusionStatusEvent  uintptr
	UnregisterOcclusionStatus     uintptr
	CreateSwapChainForComposition uintptr
	GetCreationFlags              uintptr
	EnumAdapterByLuid             uintptr
	EnumWarpAdapter               uintptr
}

func (i *_IDXGIFactory4) CreateSwapChainForComposition(pDevice unsafe.Pointer, pDesc *_DXGI_SWAP_CHAIN_DESC1, pRestrictToOutput *_IDXGIOutput) (*_IDXGISwapChain1, error) {
	var swapChain *_IDXGISwapChain1
	r, _, _ := syscallN(i.vtbl.CreateSwapChainForComposition,
		uintptr(unsafe.Pointer(i)), uintptr(pDevice), uintptr(unsafe.Pointer(pDesc)),
		uintptr(unsafe.Pointer(pRestrictToOutput)), uintptr(unsafe.Pointer(&swapChain)))
	runtime.KeepAlive(pDesc)
	runtime.KeepAlive(pRestrictToOutput)
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: IDXGIFactory4::CreateSwapChainForComposition failed: HRESULT(%d)", uint32(r))
	}
	return swapChain, nil
}

func (i *_IDXGIFactory4) CreateSwapChainForHwnd(pDevice unsafe.Pointer, hWnd windows.HWND, pDesc *_DXGI_SWAP_CHAIN_DESC1, pFullscreenDesc *_DXGI_SWAP_CHAIN_FULLSCREEN_DESC, pRestrictToOutput *_IDXGIOutput) (*_IDXGISwapChain1, error) {
	var swapChain *_IDXGISwapChain1
	r, _, _ := syscallN(i.vtbl.CreateSwapChainForHwnd,
		uintptr(unsafe.Pointer(i)), uintptr(pDevice), uintptr(hWnd),
		uintptr(unsafe.Pointer(pDesc)), uintptr(unsafe.Pointer(pFullscreenDesc)), uintptr(unsafe.Pointer(pRestrictToOutput)),
		uintptr(unsafe.Pointer(&swapChain)))
	runtime.KeepAlive(pDesc)
	runtime.KeepAlive(pFullscreenDesc)
	runtime.KeepAlive(pRestrictToOutput)
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: IDXGIFactory4::CreateSwapChainForHwnd failed: HRESULT(%d)", uint32(r))
	}
	return swapChain, nil
}

func (i *_IDXGIFactory4) EnumAdapters1(adapter uint32) (*_IDXGIAdapter1, error) {
	var ptr *_IDXGIAdapter1
	r, _, _ := syscallN(i.vtbl.EnumAdapters1, uintptr(unsafe.Pointer(i)), uintptr(adapter), uintptr(unsafe.Pointer(&ptr)))
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: IDXGIFactory4::EnumAdapters1 failed: HRESULT(%d)", uint32(r))
	}
	return ptr, nil
}

func (i *_IDXGIFactory4) EnumWarpAdapter() (*_IDXGIAdapter1, error) {
	var ptr *_IDXGIAdapter1
	r, _, _ := syscallN(i.vtbl.EnumWarpAdapter, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(&_IID_IDXGIAdapter1)), uintptr(unsafe.Pointer(&ptr)))
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: IDXGIFactory4::EnumWarpAdapter failed: HRESULT(%d)", uint32(r))
	}
	return ptr, nil
}

func (i *_IDXGIFactory4) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}

type _IDXGIOutput struct {
	vtbl *_IDXGIOutput_Vtbl
}

type _IDXGIOutput_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	SetPrivateData              uintptr
	SetPrivateDataInterface     uintptr
	GetPrivateData              uintptr
	GetParent                   uintptr
	GetDesc                     uintptr
	GetDisplayModeList          uintptr
	FindClosestMatchingMode     uintptr
	WaitForVBlank               uintptr
	TakeOwnership               uintptr
	ReleaseOwnership            uintptr
	GetGammaControlCapabilities uintptr
	SetGammaControl             uintptr
	GetGammaControl             uintptr
	SetDisplaySurface           uintptr
	GetDisplaySurfaceData       uintptr
	GetFrameStatistics          uintptr
}

type _ID3D12RootSignature struct {
	vtbl *_ID3D12RootSignature_Vtbl
}

type _ID3D12RootSignature_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetPrivateData          uintptr
	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
	SetName                 uintptr
	GetDevice               uintptr
}

func (i *_ID3D12RootSignature) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}

type _IDXGISwapChain1 struct {
	vtbl *_IDXGISwapChain1_Vtbl
}

type _IDXGISwapChain1_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	SetPrivateData           uintptr
	SetPrivateDataInterface  uintptr
	GetPrivateData           uintptr
	GetParent                uintptr
	GetDevice                uintptr
	Present                  uintptr
	GetBuffer                uintptr
	SetFullscreenState       uintptr
	GetFullscreenState       uintptr
	GetDesc                  uintptr
	ResizeBuffers            uintptr
	ResizeTarget             uintptr
	GetContainingOutput      uintptr
	GetFrameStatistics       uintptr
	GetLastPresentCount      uintptr
	GetDesc1                 uintptr
	GetFullscreenDesc        uintptr
	GetHwnd                  uintptr
	GetCoreWindow            uintptr
	Present1                 uintptr
	IsTemporaryMonoSupported uintptr
	GetRestrictToOutput      uintptr
	SetBackgroundColor       uintptr
	GetBackgroundColor       uintptr
	SetRotation              uintptr
	GetRotation              uintptr
}

func (i *_IDXGISwapChain1) As(swapChain **_IDXGISwapChain4) {
	*swapChain = (*_IDXGISwapChain4)(unsafe.Pointer(i))
}

type _IDXGISwapChain4 struct {
	vtbl *_IDXGISwapChain4_Vtbl
}

type _IDXGISwapChain4_Vtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	SetPrivateData           uintptr
	SetPrivateDataInterface  uintptr
	GetPrivateData           uintptr
	GetParent                uintptr
	GetDevice                uintptr
	Present                  uintptr
	GetBuffer                uintptr
	SetFullscreenState       uintptr
	GetFullscreenState       uintptr
	GetDesc                  uintptr
	ResizeBuffers            uintptr
	ResizeTarget             uintptr
	GetContainingOutput      uintptr
	GetFrameStatistics       uintptr
	GetLastPresentCount      uintptr
	GetDesc1                 uintptr
	GetFullscreenDesc        uintptr
	GetHwnd                  uintptr
	GetCoreWindow            uintptr
	Present1                 uintptr
	IsTemporaryMonoSupported uintptr
	GetRestrictToOutput      uintptr
	SetBackgroundColor       uintptr
	GetBackgroundColor       uintptr
	SetRotation              uintptr
	GetRotation              uintptr

	SetSourceSize                 uintptr
	GetSourceSize                 uintptr
	SetMaximumFrameLatency        uintptr
	GetMaximumFrameLatency        uintptr
	GetFrameLatencyWaitableObject uintptr
	SetMatrixTransform            uintptr
	GetMatrixTransform            uintptr
	GetCurrentBackBufferIndex     uintptr
	CheckColorSpaceSupport        uintptr
	SetColorSpace1                uintptr
	ResizeBuffers1                uintptr
	SetHDRMetaData                uintptr
}

func (i *_IDXGISwapChain4) GetBuffer(buffer uint32) (*_ID3D12Resource, error) {
	var resource *_ID3D12Resource
	r, _, _ := syscallN(i.vtbl.GetBuffer, uintptr(unsafe.Pointer(i)),
		uintptr(buffer), uintptr(unsafe.Pointer(&_IID_ID3D12Resource)), uintptr(unsafe.Pointer(&resource)))
	if uint32(r) != uint32(windows.S_OK) {
		return nil, fmt.Errorf("directx: IDXGISwapChain4::GetBuffer failed: HRESULT(%d)", uint32(r))
	}
	return resource, nil
}

func (i *_IDXGISwapChain4) GetCurrentBackBufferIndex() uint32 {
	r, _, _ := syscallN(i.vtbl.GetCurrentBackBufferIndex, uintptr(unsafe.Pointer(i)))
	return uint32(r)
}

func (i *_IDXGISwapChain4) Present(syncInterval uint32, flags uint32) error {
	r, _, _ := syscallN(i.vtbl.Present, uintptr(unsafe.Pointer(i)), uintptr(syncInterval), uintptr(flags))
	if uint32(r) != uint32(windows.S_OK) {
		return fmt.Errorf("directx: IDXGISwapChain4::Present failed: HRESULT(%d)", uint32(r))
	}
	return nil
}

func (i *_IDXGISwapChain4) ResizeBuffers(bufferCount uint32, width uint32, height uint32, newFormat _DXGI_FORMAT, swapChainFlags uint32) error {
	r, _, _ := syscallN(i.vtbl.ResizeBuffers,
		uintptr(unsafe.Pointer(i)), uintptr(bufferCount), uintptr(width),
		uintptr(height), uintptr(newFormat), uintptr(swapChainFlags))
	if uint32(r) != uint32(windows.S_OK) {
		return fmt.Errorf("directx: IDXGISwapChain4::ResizeBuffers failed: HRESULT(%d)", uint32(r))
	}
	return nil
}

func (i *_IDXGISwapChain4) Release() {
	syscallN(i.vtbl.Release, uintptr(unsafe.Pointer(i)))
}
