
#include <AL/al.h>

typedef int ALCenum;
typedef int ALCint;
typedef char ALCchar;
typedef int ALCboolean;

typedef struct{} ALCdevice;
typedef struct{} ALCcontext;

ALCdevice *alcOpenDevice(const ALCchar *devicename);
ALCcontext * alcCreateContext(ALCdevice *device, ALCint* attrlist);
ALCboolean alcMakeContextCurrent(ALCcontext *context);
