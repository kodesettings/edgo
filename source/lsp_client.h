#ifndef _LSP_CLIENT_H
#define _LSP_CLIENT_H

#include <string>
#include <glog/logging.h>
#include <map>
#include <sstream>

typedef struct {
	std::string lang;
	Cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stop   context.CancelFunc
	reader *bufio.Reader

	bool isready; // flag for lsp initialization

	std::stringstream usermessages;
	std::stringstream diagnosticschannel;

	int id;
	std::map<std::string, diagnosticparams_t> file2diagnostic;
} lspclient_t;

bool Start(const std::string cmd, std::string args...);
void Send(const std::string message);

diagnosticresponse_t receiveDiagnostics(void);
void receiveLoop(void);
diagnosticparams_t GetDiagnostic(const std::string filename);
void InitLspClient(const std::string dir);

void DidOpen(const std::string file, std::string *text);
void DidChange(const std::string file, std::string *text, int version);
void DidClose(const std::string file);

#endif // _LSP_CLIENT_H
