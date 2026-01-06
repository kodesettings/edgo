#ifndef _LOGGING_H_
#define _LOGGING_H_

#include <glog/logging.h>

// Initialization macros for google's logging system.
// By default we have to log into files on disk.
#define EDGO_LOGGING_INIT google::InitGoogleLogging("edgo");
#define EDGO_LOGGING_TO_FILE FLAGS_logtostderr = 0;

#endif // _LOGGING_H_
