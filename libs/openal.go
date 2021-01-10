package libs

const (
	alH  = "AL/al.h"
	alcH = "AL/alc.h"
)

func init() {
	RegisterLibrary(alH, func(c *Env) *Library {
		return &Library{
			// TODO
			Header: `
#include <` + BuiltinH + `>

#define ALvoid void

typedef unsigned int ALuint;
typedef int ALint;
typedef int ALsizei;
typedef int ALenum;
typedef float ALfloat;

const int INT16_MIN = 0;
const int INT16_MAX = 0;

const int AL_NO_ERROR = 0;
const int AL_BUFFERS_PROCESSED = 0;
const int AL_PITCH = 0;
const int AL_GAIN = 0;
const int AL_POSITION = 0;
const int AL_SOURCE_STATE = 0;
const int AL_PLAYING = 0;
const int AL_FORMAT_STEREO16 = 0;
const int AL_FORMAT_MONO16 = 0;

ALenum alGetError(ALvoid);

void alGetSourcei(ALuint source, ALenum pname, ALint* value);
void alSourcef(ALuint source, ALenum param, ALfloat value);
void alSourcefv(ALuint source, ALenum param, ALfloat* values);
void alListenerf(ALenum param, ALfloat value);
void alListener3f(ALenum param, ALfloat v1, ALfloat v2, ALfloat v3);
void alGenSources(ALsizei n, ALuint* sources);
void alGenBuffers(ALsizei n, ALuint* buffers);
void alBufferData(ALuint buffer, ALenum format, const ALvoid *data, ALsizei size, ALsizei freq);
void alSourceQueueBuffers(ALuint source, ALsizei n, ALuint* buffers);
void alSourceUnqueueBuffers(ALuint source, ALsizei n, ALuint* buffers);
void alSourcePlay(ALuint source);
void alSourceStop(ALuint source);
`,
		}
	})
	RegisterLibrary(alcH, func(c *Env) *Library {
		return &Library{
			// TODO
			Header: `
#include <` + alH + `>

typedef int ALCenum;
typedef int ALCint;
typedef char ALCchar;
typedef int ALCboolean;

typedef struct{} ALCdevice;
typedef struct{} ALCcontext;

ALCdevice *alcOpenDevice(const ALCchar *devicename);
ALCcontext * alcCreateContext(ALCdevice *device, ALCint* attrlist);
ALCboolean alcMakeContextCurrent(ALCcontext *context);
`,
		}
	})
}
