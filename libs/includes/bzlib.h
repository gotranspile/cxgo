const int BZ_OK = 0;

int BZ2_bzBuffToBuffCompress(
      char*         dest,
      unsigned int* destLen,
      char*         source,
      unsigned int  sourceLen,
      int           blockSize100k,
      int           verbosity,
      int           workFactor
   );

int BZ2_bzBuffToBuffDecompress(
      char*         dest,
      unsigned int* destLen,
      char*         source,
      unsigned int  sourceLen,
      int           small,
      int           verbosity
   );