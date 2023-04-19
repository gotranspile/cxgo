
#include <stddef.h>
#include <sys/socket.h>
#include <unistd.h>

#define __cdecl
#define __fastcall
#define __stdcall
#define WINAPI
#define WINAPIV

#define VOID void
#define FALSE 0
#define TRUE 1

#define MAX_PATH 260
#define INFINITE ((DWORD)-1)

#define BYTE _cxgo_uint8
#define WORD _cxgo_uint16
#define DWORD _cxgo_uint32
#define QWORD _cxgo_uint64

#define _BYTE _cxgo_uint8
#define _WORD _cxgo_uint16
#define _DWORD _cxgo_uint32
#define _QWORD _cxgo_uint64

#define BOOL int

#define CHAR _cxgo_int8
#define WCHAR wchar_t

#define INT int
#define UINT unsigned int
#define LONG long
#define ULONG unsigned long
#define LONGLONG long long
#define ULONGLONG unsigned long long

#define PLONG LONG*

#define LPVOID VOID*
#define LPBOOL BOOL*
#define LPBYTE BYTE*
#define LPWORD WORD*
#define LPDWORD DWORD*
#define LPQWORD QWORD*

#define WORD_PTR WORD*
#define DWORD_PTR DWORD*
#define QWORD_PTR QWORD*

#define LPCCH const char*
#define LPCSTR const char*
#define LPSTR char*
#define LPCWSTR const wchar_t*
#define LPWSTR wchar_t*

#define INT_PTR intptr_t
#define UINT_PTR uintptr_t

typedef void* HINSTANCE;
typedef void* HMODULE;
typedef int HANDLE;
typedef int HWND;
typedef int HDC;
typedef int HIMC;

typedef INT_PTR LRESULT;
typedef INT_PTR LSTATUS;
typedef INT_PTR LPARAM;
typedef UINT_PTR WPARAM;

typedef struct _REGKEY* HKEY;
typedef HKEY* PHKEY;
#define HKEY_LOCAL_MACHINE ((HKEY)1)

typedef int SOCKET;

typedef int LCID;
typedef int REGSAM;

typedef struct {
	WORD wFormatTag;
	WORD nChannels;
	DWORD nSamplesPerSec;
	DWORD nAvgBytesPerSec;
	WORD nBlockAlign;
} WAVEFORMAT;

typedef struct _WAVEOUT {};

#define HWAVEOUT struct _WAVEOUT*

#define LPWAVEFORMAT WAVEFORMAT*
#define LPHWAVEOUT HWAVEOUT*

#pragma pack(push, 1)
typedef struct {
	WORD wFormatTag;
	WORD nChannels;
	DWORD nSamplesPerSec;
	DWORD nAvgBytesPerSec;
	WORD nBlockAlign;
	WORD wBitsPerSample;
	WORD cbSize;
} WAVEFORMATEX;

typedef struct {
	WAVEFORMATEX wfx;
	WORD wID;
	DWORD fdwFlags;
	WORD nBlockSize;
	WORD nFramesPerBlock;
	WORD nCodecDelay;
} MPEGLAYER3WAVEFORMAT;
#pragma pack(pop)

typedef struct _GUID {
	DWORD Data1;
	WORD Data2;
	WORD Data3;
	BYTE Data4[8];
} GUID;
typedef GUID IID;

typedef struct _SYSTEMTIME {
	WORD wYear;
	WORD wMonth;
	WORD wDayOfWeek;
	WORD wDay;
	WORD wHour;
	WORD wMinute;
	WORD wSecond;
	WORD wMilliseconds;
} SYSTEMTIME, *LPSYSTEMTIME;

typedef struct _MEMORYSTATUS {
	DWORD dwLength;
} MEMORYSTATUS, *LPMEMORYSTATUS;

typedef union _LARGE_INTEGER {
	QWORD QuadPart;
	DWORD LowPart;
} LARGE_INTEGER, *PLARGE_INTEGER;

typedef struct _FILETIME {
	DWORD dwLowDateTime;
	DWORD dwHighDateTime;
} FILETIME, *LPFILETIME;

typedef struct _SECURITY_ATTRIBUTES {
} SECURITY_ATTRIBUTES, *LPSECURITY_ATTRIBUTES;

typedef struct _CRITICAL_SECTION {
	void* opaque;
} CRITICAL_SECTION, *LPCRITICAL_SECTION;

typedef struct _OSVERSIONINFOA {
	DWORD dwOSVersionInfoSize;
	DWORD dwMajorVersion;
	DWORD dwMinorVersion;
	DWORD dwBuildNumber;
	DWORD dwPlatformId;
	BYTE szCSDVersion[128];
} OSVERSIONINFOA, *LPOSVERSIONINFOA;

typedef struct _WIN32_FIND_DATAA {
	DWORD dwFileAttributes;
	FILETIME ftCreationTime;
	FILETIME ftLastAccessTime;
	FILETIME ftLastWriteTime;
	DWORD nFileSizeHigh;
	DWORD nFileSizeLow;
	DWORD dwReserved0;
	DWORD dwReserved1;
	CHAR cFileName[MAX_PATH];
	CHAR cAlternateFileName[14];
} WIN32_FIND_DATAA, *LPWIN32_FIND_DATAA;

typedef struct _OVERLAPPED {
} OVERLAPPED, *LPOVERLAPPED;

typedef struct _RECT {
	LONG left;
	LONG top;
	LONG right;
	LONG bottom;
} RECT, *LPRECT;

typedef struct _POINT {
	LONG x;
	LONG y;
} POINT;

typedef struct WSAData {
	WORD wVersion;
	WORD wHighVersion;
	CHAR szDescription[257];
	CHAR szSystemStatus[129];
	WORD iMaxSockets;
	WORD iMaxUdpDg;
	LPVOID lpVendorInfo;
} WSADATA, *LPWSADATA;

struct _stat {
	DWORD st_dev;
	WORD st_ino;
	WORD st_mode;
	WORD st_nlink;
	WORD st_uid;
	WORD st_gid;
	DWORD st_rdev;
	DWORD st_size;
	DWORD st_mtime;
	DWORD st_atime;
	DWORD st_ctime;
};

struct in_addr {
	union {
		DWORD S_addr;
	} S_un;
};

#define _strdup strdup
#define _strcmpi strcasecmp
#define _strnicmp strncasecmp
#define _fileno fileno
#define _read read
#define _write write
#define _close close
#define _onexit atexit

unsigned int _control87(unsigned int new_, unsigned int mask);
unsigned int _controlfp(unsigned int new_, unsigned int mask);
int _open(const char* filename, int oflag, ...);
int _chmod(const char* filename, int mode);
int _access(const char* filename, int mode);
int _stat(const char* path, struct _stat* buffer);
int _mkdir(const char* path);
int _unlink(const char* filename);
char* _getcwd(char* buffer, int maxlen);
uintptr_t _beginthread(void(__cdecl* start_address)(void*), unsigned int stack_size, void* arglist);
char* _strrev(char* str);
char* _itoa(int val, char* s, int radix);
wchar_t* _itow(int val, wchar_t* s, int radix);
void _makepath(char* path, const char* drive, const char* dir, const char* fname, const char* ext);
void _splitpath(const char* path, char* drive, char* dir, char* fname, char* ext);

VOID WINAPI DebugBreak();
BOOL WINAPI CloseHandle(HANDLE hObject);
DWORD WINAPI GetLastError();
VOID WINAPI GetLocalTime(LPSYSTEMTIME lpSystemTime);
HANDLE WINAPI FindFirstFileA(LPCSTR lpFileName, LPWIN32_FIND_DATAA lpFindFileData);
BOOL WINAPI FindNextFileA(HANDLE hFindFile, LPWIN32_FIND_DATAA lpFindFileData);
BOOL WINAPI FindClose(HANDLE hFindFile);
HANDLE WINAPI CreateFileA(LPCSTR lpFileName, DWORD dwDesiredAccess, DWORD dwShareMode,
			  LPSECURITY_ATTRIBUTES lpSecurityAttributes, DWORD dwCreationDisposition,
			  DWORD dwFlagsAndAttributes, HANDLE hTemplateFile);
BOOL WINAPI ReadFile(HANDLE hFile, LPVOID lpBuffer, DWORD nNumberOfBytesToRead, LPDWORD lpNumberOfBytesRead,
		     LPOVERLAPPED lpOverlapped);
DWORD WINAPI SetFilePointer(HANDLE hFile, LONG lDistanceToMove, PLONG lpDistanceToMoveHigh, DWORD dwMoveMethod);
BOOL WINAPI CopyFileA(LPCSTR lpExistingFileName, LPCSTR lpNewFileName, BOOL bFailIfExists);
BOOL WINAPI DeleteFileA(LPCSTR lpFileName);
BOOL WINAPI MoveFileA(LPCSTR lpExistingFileName, LPCSTR lpNewFileName);
BOOL WINAPI CreateDirectoryA(LPCSTR lpPathName, LPSECURITY_ATTRIBUTES lpSecurityAttributes);
BOOL WINAPI RemoveDirectoryA(LPCSTR lpPathName);
DWORD WINAPI GetCurrentDirectoryA(DWORD nBufferLength, LPSTR lpBuffer);
BOOL WINAPI SetCurrentDirectoryA(LPCSTR lpPathName);
int WINAPI GetDateFormatA(LCID Locale, DWORD dwFlags, const SYSTEMTIME* lpDate, LPCSTR lpFormat, LPSTR lpDateStr,
			  int cchDate);
LSTATUS WINAPI RegOpenKeyExA(HKEY hKey, LPCSTR lpSubKey, DWORD ulOptions, REGSAM samDesired, PHKEY phkResult);
LSTATUS WINAPI RegQueryValueExA(HKEY hKey, LPCSTR lpValueName, LPDWORD lpReserved, LPDWORD lpType, LPBYTE lpData,
				LPDWORD lpcbData);
LSTATUS WINAPI RegSetValueExA(HKEY, LPCSTR lpValueName, DWORD Reserved, DWORD dwType, const BYTE* lpData, DWORD cbData);
LSTATUS WINAPI RegCloseKey(HKEY hKey);
int WINAPI MulDiv(int nNumber, int nNumerator, int nDenominator);
LSTATUS WINAPI RegCreateKeyExA(HKEY hKey, LPCSTR lpSubKey, DWORD Reserved, LPSTR lpClass, DWORD dwOptions,
			       REGSAM samDesired, const LPSECURITY_ATTRIBUTES lpSecurityAttributes, PHKEY phkResult,
			       LPDWORD lpdwDisposition);
VOID WINAPI GlobalMemoryStatus(LPMEMORYSTATUS lpBuffer);
DWORD WINAPI GetModuleFileNameA(HMODULE hModule, LPSTR lpFileName, DWORD nSize);
BOOL WINAPI QueryPerformanceCounter(LARGE_INTEGER* lpPerformanceCount);
BOOL WINAPI QueryPerformanceFrequency(LARGE_INTEGER* lpFrequency);
VOID WINAPI InitializeCriticalSection(LPCRITICAL_SECTION lpCriticalSection);
VOID WINAPI DeleteCriticalSection(LPCRITICAL_SECTION lpCriticalSection);
VOID WINAPI EnterCriticalSection(LPCRITICAL_SECTION lpCriticalSection);
VOID WINAPI LeaveCriticalSection(LPCRITICAL_SECTION lpCriticalSection);
BOOL WINAPI HeapDestroy(HANDLE hHeap);
BOOL WINAPI GetVersionExA(LPOSVERSIONINFOA lpVersionInformation);
VOID WINAPI OutputDebugStringA(LPCSTR lpOutputString);
HINSTANCE WINAPI ShellExecuteA(HWND hwnd, LPCSTR lpOperation, LPCSTR lpFile, LPCSTR lpParameters, LPCSTR lpDirectory,
			       INT nShowCmd);
int WINAPI GetTimeFormatA(LCID Locale, DWORD dwFlags, const SYSTEMTIME* lpTime, LPCSTR lpFormat, LPSTR lpTimeStr,
			  int cchTime);
BOOL WINAPI SystemTimeToFileTime(const SYSTEMTIME* lpSystemTime, LPFILETIME lpFileTime);
LONG WINAPI CompareFileTime(const FILETIME* lpFileTime1, const FILETIME* lpFileTime2);
HANDLE WINAPI CreateMutexA(LPSECURITY_ATTRIBUTES lpSecurityAttributes, BOOL bInitialOwner, LPCSTR lpName);
BOOL WINAPI ReleaseMutex(HANDLE hMutex);
BOOL WINAPI SetEvent(HANDLE hEvent);
DWORD WINAPI WaitForSingleObject(HANDLE hHandle, DWORD dwMilliseconds);
char* WINAPI inet_ntoa(struct in_addr in);
int WINAPI WideCharToMultiByte(UINT CodePage, DWORD dwFlags, LPCWSTR lpWideCharStr, int cchWideChar,
			       LPSTR lpMultiByteStr, int cbMultiByte, LPCCH lpDefaultChar, LPBOOL lpUsedDefaultChar);

LONG InterlockedExchange(volatile LONG* Target, LONG Value);
LONG InterlockedDecrement(volatile LONG* Addend);
LONG InterlockedIncrement(volatile LONG* Addend);
int WINAPI MessageBoxA(HWND hWnd, LPCSTR lpText, LPCSTR lpCaption, UINT uType);
int WINAPI WSAStartup(WORD wVersionRequested, struct WSAData* lpWSAData);
int WINAPI WSACleanup();
int WINAPI closesocket(SOCKET s);
int WINAPI ioctlsocket(SOCKET s, long cmd, unsigned long* argp);
int WINAPI WSAGetLastError();
SOCKET WINAPI socket(int domain, int type, int protocol);
int WINAPI setsockopt(SOCKET s, int level, int opt, const void* value, unsigned int len);
int WINAPI bind(int sockfd, const struct sockaddr* addr, unsigned int addrlen);
int WINAPI recvfrom(int sockfd, void* buffer, unsigned int length, int flags, struct sockaddr* addr,
		    unsigned int* addrlen);
int WINAPI sendto(int sockfd, void* buffer, unsigned int length, int flags, const struct sockaddr* addr,
		  unsigned int addrlen);
