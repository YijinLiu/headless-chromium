#include <signal.h>
#include <unistd.h>

#include <base/logging.h>

#include "browser.hpp"

Browser* browser;

void signal_handler(int signo) {
    LOG(INFO) << "Received signal " << signo;
    browser->Shutdown();
}

int main(int argc, const char** argv) {
    headless::RunChildProcessIfNeeded(argc, argv);
    if (signal(SIGINT, signal_handler) == SIG_ERR) {
        LOG(ERROR) << "Failed to catch SIGINT";
    }
    if (signal(SIGHUP, signal_handler) == SIG_ERR) {
        LOG(ERROR) << "Failed to catch SIGHUP";
    }
    browser = new Browser;
    browser->Run(argc, argv, []{
        LOG(INFO) << "Headless Chromium is ready!";
    });
    delete browser;
    return 0;
}
