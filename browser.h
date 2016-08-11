#ifndef BROWSER_H_
#define BROWSER_H_

typedef void* BrowserPtr;

BrowserPtr create_browser();

void destroy_browser(BrowserPtr browser);

int run_browser(BrowserPtr browser);

void shutdown_browser(BrowserPtr browser);

int open_url(BrowserPtr browser, const char* url, int width, int height);

void evaluate_script(BrowserPtr browser, const char* script);

#endif  // BROWSER_H_
