#include "browser.hpp"

int main(int argc, const char** argv) {
    headless::RunChildProcessIfNeeded(argc, argv);
    Browser browser;
    browser.Run(argc, argv, []{
        LOG(INFO) << "Headless Chromium is ready!";
    });
    return 0;
}
