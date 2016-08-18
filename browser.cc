#include "browser.hpp"

#include <base/bind.h>
#include <base/callback.h>
#include <ui/gfx/geometry/size.h>
#include <url/gurl.h>

Browser::Browser() : browser_(nullptr), browser_context_(nullptr), web_contents_(nullptr),
                     devtools_client_(headless::HeadlessDevToolsClient::Create()) {}

Browser::~Browser() {}

int Browser::Run(std::function<void()> readyCb) {
    headless::HeadlessBrowser::Options::Builder opts_builder;
    return headless::HeadlessBrowserMain(
        opts_builder.Build(), base::Bind(&Browser::onStart, base::Unretained(this), readyCb));
}

bool Browser::OpenUrl(const std::string& url, int width, int height, std::function<void()> readyCb) {
    GURL gurl(url);
    if (web_contents_ != nullptr) web_contents_->Close();
    browser_context_ = browser_->CreateBrowserContextBuilder().Build();
    headless::HeadlessWebContents::Builder builder(browser_context_->CreateWebContentsBuilder());
	builder.SetInitialURL(gurl);
    web_contents_ = builder.Build();
    if (web_contents_ == nullptr) return false;
    web_contents_->AddObserver(this);
    return true;
}

namespace {

void EvaluateCallback(std::function<void(bool, const std::string&)> cb,
                      std::unique_ptr<headless::runtime::EvaluateResult> result) {
    if (result->HasExceptionDetails()) {
        std::string exception;
        if (result->HasExceptionDetails()) {
            exception = result->GetExceptionDetails()->GetText();
        }
        cb(false, exception);
    } else {
        std::string value;
        result->GetResult()->Serialize()->GetAsString(&value);
        cb(true, value);
    }
}

}  // namespace

void Browser::Evaluate(
    const std::string& script, std::function<void(bool, const std::string&)> resultCb) {
    devtools_client_->GetRuntime()->Evaluate(script, base::Bind(&EvaluateCallback, resultCb));
}

void Browser::Shutdown() {
    if (!browser_) return;
    if (web_contents_) {
        web_contents_->RemoveObserver(this);
        web_contents_ = nullptr;
    }
    browser_context_->Close();
    browser_context_ = nullptr;
    browser_->Shutdown();
    browser_ = nullptr;
}

void Browser::onStart(std::function<void()> readyCb, headless::HeadlessBrowser* browser) {
    browser_ = browser;
    readyCb();
}

void Browser::DevToolsTargetReady() {
    web_contents_->GetDevToolsTarget()->AttachClient(devtools_client_.get());
    devtools_client_->GetPage()->AddObserver(this);
    devtools_client_->GetPage()->Enable();
}

void Browser::OnLoadEventFired(const headless::page::LoadEventFiredParams& params) {
    pageLoadedCb_();
}

extern "C" {

Browser* create_browser() {
    return new Browser();
}

void destroy_browser(Browser* browser) {
    delete browser;
}

extern void signalReady(Browser*);

int run_browser(Browser* browser) {
    return browser->Run([browser]() { signalReady(browser); });
}

void shutdown_browser(Browser* browser) {
    browser->Shutdown();
}

int open_url(Browser* browser, const char* cstr_url, int width, int height) {
    std::string url(cstr_url);
    free((void*)cstr_url);
    return browser->OpenUrl(url, width, height, [browser]() { signalReady(browser); }) ? 1: 0;
}

extern void signalEvaluateResult(Browser* browser, int success, const char* result);

void evaluate_script(Browser* browser, const char* cstr_script) {
    std::string script(cstr_script);
    free((void*)cstr_script);
    browser->Evaluate(script, [browser](bool success, const std::string& result) {
       signalEvaluateResult(browser, (success ? 1 : 0 ), result.c_str());
    });
}

}  // extern "C"
