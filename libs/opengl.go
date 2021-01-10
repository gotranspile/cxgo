package libs

const (
	glH = "GL/gl.h"
)

func init() {
	RegisterLibrary(glH, func(c *Env) *Library {
		return &Library{
			// TODO
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
`,
		}
	})
}
