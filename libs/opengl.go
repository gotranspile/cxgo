package libs

import "github.com/gotranspile/cxgo/types"

const (
	glH = "GL/gl.h"
)

func init() {
	RegisterLibrary(glH, func(env *Env) *Library {
		int32T := types.IntT(4)
		uint32T := types.UintT(4)
		intT := env.Go().Int()
		uint32PtrT := env.PtrT(types.UintT(4))
		uPtrT := env.Go().UnsafePtr()
		float32T := env.PtrT(types.FloatT(4))
		return &Library{
			Imports: map[string]string{
				"gl": "github.com/go-gl/gl/v3.3-core/gl", // last supported version on macos
			},
			Idents: map[string]*types.Ident{
				//functions
				"glGenBuffers":              types.NewIdentGo("glGenBuffers", "gl.GenBuffers", env.FuncTT(nil, int32T, uint32PtrT)),
				"glBindBuffer":              types.NewIdentGo("glBindBuffer", "gl.BindBuffer", env.FuncTT(nil, uint32T, uint32T)),
				"glBufferData":              types.NewIdentGo("glBufferData", "gl.BufferData", env.FuncTT(nil, uint32T, intT, uPtrT, uint32T)),
				"glCreateShader":            types.NewIdentGo("glCreateShader", "gl.CreateShader", env.FuncTT(uint32T, uint32T)),
				"glShaderSource":            types.NewIdentGo("glShaderSource", "gl.ShaderSource", env.FuncTT(nil, uint32T, int32T, env.PtrT(env.C().String()), env.PtrT(int32T))),
				"glCompileShader":           types.NewIdentGo("glCompileShader", "gl.CompileShader", env.FuncTT(nil, uint32T)),
				"glCreateProgram":           types.NewIdentGo("glCreateProgram", "gl.CreateProgram", env.FuncTT(uint32T)),
				"glAttachShader":            types.NewIdentGo("glAttachShader", "gl.AttachShader", env.FuncTT(nil, uint32T, uint32T)),
				"glLinkProgram":             types.NewIdentGo("glLinkProgram", "gl.LinkProgram", env.FuncTT(nil, uint32T)),
				"glGetUniformLocation":      types.NewIdentGo("glGetUniformLocation", "gl.GetUniformLocation", env.FuncTT(int32T, uint32T, env.C().String())),
				"glGetAttribLocation":       types.NewIdentGo("glGetAttribLocation", "gl.GetAttribLocation", env.FuncTT(int32T, uint32T, env.C().String())),
				"glEnableVertexAttribArray": types.NewIdentGo("glEnableVertexAttribArray", "gl.EnableVertexAttribArray", env.FuncTT(nil, uint32T)),
				"glVertexAttribPointer":     types.NewIdentGo("glVertexAttribPointer", "gl.VertexAttribPointerWithOffset", env.FuncTT(nil, uint32T, int32T, uint32T, env.Go().Bool(), int32T, env.Go().Uintptr())),
				"glViewport":                types.NewIdentGo("glViewport", "gl.Viewport", env.FuncTT(nil, int32T, int32T, int32T, int32T)),
				"glClear":                   types.NewIdentGo("glClear", "gl.Clear", env.FuncTT(nil, uint32T)),
				"glUseProgram":              types.NewIdentGo("glUseProgram", "gl.UseProgram", env.FuncTT(nil, uint32T)),
				"glUniformMatrix4fv":        types.NewIdentGo("glUniformMatrix4fv", "gl.UniformMatrix4fv", env.FuncTT(nil, int32T, int32T, env.Go().Bool(), float32T)),
				"glDrawArrays":              types.NewIdentGo("glDrawArrays", "gl.DrawArrays", env.FuncTT(nil, uint32T, int32T, int32T)),
				// constants
				"GL_FALSE":                      types.NewIdentGo("GL_FALSE", "gl.FALSE", env.Go().Int()),
				"GL_TRUE":                       types.NewIdentGo("GL_TRUE", "gl.TRUE", env.Go().Int()),
				"GL_TEXTURE_2D":                 types.NewIdentGo("GL_TEXTURE_2D", "gl.TEXTURE_2D", env.Go().Int()),
				"GL_UNPACK_ROW_LENGTH":          types.NewIdentGo("GL_UNPACK_ROW_LENGTH", "gl.UNPACK_ROW_LENGTH", env.Go().Int()),
				"GL_BGRA":                       types.NewIdentGo("GL_BGRA", "gl.BGRA", env.Go().Int()),
				"GL_RGBA":                       types.NewIdentGo("GL_RGBA", "gl.RGBA", env.Go().Int()),
				"GL_TEXTURE_WRAP_S":             types.NewIdentGo("GL_TEXTURE_WRAP_S", "gl.TEXTURE_WRAP_S", env.Go().Int()),
				"GL_TEXTURE_WRAP_T":             types.NewIdentGo("GL_TEXTURE_WRAP_T", "gl.TEXTURE_WRAP_T", env.Go().Int()),
				"GL_TEXTURE_MIN_FILTER":         types.NewIdentGo("GL_TEXTURE_MIN_FILTER", "gl.TEXTURE_MIN_FILTER", env.Go().Int()),
				"GL_FRAMEBUFFER":                types.NewIdentGo("GL_FRAMEBUFFER", "gl.FRAMEBUFFER", env.Go().Int()),
				"GL_COLOR_BUFFER_BIT":           types.NewIdentGo("GL_COLOR_BUFFER_BIT", "gl.COLOR_BUFFER_BIT", env.Go().Int()),
				"GL_UNSIGNED_SHORT_1_5_5_5_REV": types.NewIdentGo("GL_UNSIGNED_SHORT_1_5_5_5_REV", "gl.UNSIGNED_SHORT_1_5_5_5_REV", env.Go().Int()),
				"GL_UNSIGNED_SHORT_5_5_5_1":     types.NewIdentGo("GL_UNSIGNED_SHORT_5_5_5_1", "gl.UNSIGNED_SHORT_5_5_5_1", env.Go().Int()),
				"GL_LINEAR":                     types.NewIdentGo("GL_LINEAR", "gl.LINEAR", env.Go().Int()),
				"GL_ARRAY_BUFFER":               types.NewIdentGo("GL_ARRAY_BUFFER", "gl.ARRAY_BUFFER", env.Go().Int()),
				"GL_FLOAT":                      types.NewIdentGo("GL_FLOAT", "gl.FLOAT", env.Go().Int()),
				"GL_TRIANGLE_STRIP":             types.NewIdentGo("GL_TRIANGLE_STRIP", "gl.TRIANGLE_STRIP", env.Go().Int()),
				"GL_CLAMP_TO_EDGE":              types.NewIdentGo("GL_CLAMP_TO_EDGE", "gl.CLAMP_TO_EDGE", env.Go().Int()),
				"GL_TEXTURE0":                   types.NewIdentGo("GL_TEXTURE0", "gl.TEXTURE0", env.Go().Int()),
				"GL_VERTEX_SHADER":              types.NewIdentGo("GL_VERTEX_SHADER", "gl.VERTEX_SHADER", env.Go().Int()),
				"GL_COMPILE_STATUS":             types.NewIdentGo("GL_COMPILE_STATUS", "gl.COMPILE_STATUS", env.Go().Int()),
				"GL_FRAGMENT_SHADER":            types.NewIdentGo("GL_FRAGMENT_SHADER", "gl.FRAGMENT_SHADER", env.Go().Int()),
				"GL_LINK_STATUS":                types.NewIdentGo("GL_LINK_STATUS", "gl.LINK_STATUS", env.Go().Int()),
				"GL_STATIC_DRAW":                types.NewIdentGo("GL_STATIC_DRAW", "gl.STATIC_DRAW", env.Go().Int()),
				"GL_TRIANGLES":                  types.NewIdentGo("GL_TRIANGLES", "gl.TRIANGLES", env.Go().Int()),
			},
		}
	})
}
