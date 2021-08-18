package libs

import "github.com/gotranspile/cxgo/types"

const (
	glH = "GL/gl.h"
)

func init() {
	RegisterLibrary(glH, func(env *Env) *Library {
		return &Library{
			// TODO
			Imports: map[string]string{
				"gl": "github.com/go-gl/gl/v3.3-core/gl", // last supported version on macos
			},
			Idents: map[string]*types.Ident{
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
			Header: `
#include <` + BuiltinH + `>

typedef int GLint;
typedef float GLfloat;
typedef unsigned int GLuint;
typedef unsigned int GLsizei;
typedef unsigned int GLbitfield;
typedef unsigned int GLsizeiptr;
typedef int GLenum;
typedef char GLchar;
#define GLvoid void
typedef int GLboolean;

// fake enum values
const int GL_FALSE = 0;
const int GL_TRUE = 1;
const int GL_TEXTURE_2D = 0;
const int GL_UNPACK_ROW_LENGTH = 0;
const int GL_BGRA = 0;
const int GL_RGBA = 0;
const int GL_TEXTURE_WRAP_S = 0;
const int GL_TEXTURE_WRAP_T = 0;
const int GL_TEXTURE_MIN_FILTER = 0;
const int GL_FRAMEBUFFER = 0;
const int GL_COLOR_BUFFER_BIT = 0;
const int GL_UNSIGNED_SHORT_1_5_5_5_REV = 0;
const int GL_UNSIGNED_SHORT_5_5_5_1 = 0;
const int GL_LINEAR = 0;
const int GL_ARRAY_BUFFER = 0;
const int GL_FLOAT = 0;
const int GL_TRIANGLE_STRIP = 0;
const int GL_CLAMP_TO_EDGE = 0;
const int GL_TEXTURE0 = 0;
const int GL_VERTEX_SHADER = 0;
const int GL_COMPILE_STATUS = 0;
const int GL_FRAGMENT_SHADER = 0;
const int GL_LINK_STATUS = 0;
const int GL_STATIC_DRAW = 0;
const int GL_TRIANGLES = 0x0004;

void glViewport(GLint x,  GLint y,  GLsizei width,  GLsizei height);
void glBindTexture(GLenum target,  GLuint texture);
GLenum glGetError( void);
void glPixelStorei(GLenum pname,  GLint param);
void glTexSubImage2D( 	GLenum target,
  	GLint level,
  	GLint xoffset,
  	GLint yoffset,
  	GLsizei width,
  	GLsizei height,
  	GLenum format,
  	GLenum type,
  	const GLvoid * pixels);
void glTexParameteri( 	GLenum target,
  	GLenum pname,
  	GLint param);
void glBindFramebuffer( 	GLenum target,
  	GLuint framebuffer);
void glClear(GLbitfield mask);
void glUseProgram( 	GLuint program);
void glUniform1i( 	GLint location, GLint v0);
void glUniform1f(GLint location,  GLfloat v0);
void glUniformMatrix2fv(GLint location,  GLsizei count,  GLboolean transpose,  const GLfloat *value);
void glBindBuffer( 	GLenum target,
  	GLuint buffer);

void glVertexAttribPointer( 	GLuint index,
  	GLint size,
  	GLenum type,
  	GLboolean normalized,
  	GLsizei stride,
  	const GLvoid * pointer);

void glDrawArrays( 	GLenum mode,
  	GLint first,
  	GLsizei count);

void glGenTextures( 	GLsizei n,
  	GLuint * textures);
void glActiveTexture( 	GLenum texture);
void glTexImage2D( 	GLenum target,
  	GLint level,
  	GLint internalformat,
  	GLsizei width,
  	GLsizei height,
  	GLint border,
  	GLenum format,
  	GLenum type,
  	const GLvoid * data);
GLuint glCreateShader( 	GLenum shaderType);
void glShaderSource( 	GLuint shader,
  	GLsizei count,
  	const GLchar **string,
  	const GLint *length);
void glCompileShader( 	GLuint shader);
void glGetShaderiv(GLuint shader,  GLenum pname,  GLint *params);
void glGetShaderInfoLog( 	GLuint shader,
  	GLsizei maxLength,
  	GLsizei *length,
  	GLchar *infoLog);
GLuint glCreateProgram( 	void);
void glAttachShader( 	GLuint program,
  	GLuint shader);
void glLinkProgram( 	GLuint program);
void glGetProgramiv(GLuint program,  GLenum pname,  GLint *params);
GLint glGetAttribLocation( 	GLuint program,
  	const GLchar *name);
void glEnableVertexAttribArray( 	GLuint index);
GLint glGetUniformLocation( 	GLuint program,
  	const GLchar *name);
void glGenBuffers( 	GLsizei n,
  	GLuint * buffers);
void glBufferData( 	GLenum target,
  	GLsizeiptr size,
  	const GLvoid * data,
  	GLenum usage);
void glClearColor(GLfloat r, GLfloat g, GLfloat b, GLfloat a);
void glGenVertexArrays(GLsizei n, GLuint *arrays);
void glBindVertexArray(GLuint array);
void glUniformMatrix4fv(GLint, GLsizei count, GLboolean transpose, const GLfloat *value);
`,
		}
	})
}
