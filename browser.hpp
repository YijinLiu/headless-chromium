#ifndef BROWSER_HPP_
#define BROWSER_HPP_

#include <functional>
#include <string>

#include <base/macros.h>
#include <headless/public/domains/page.h>
#include <headless/public/domains/runtime.h>
#include <headless/public/headless_browser.h>
#include <headless/public/headless_devtools_client.h>
#include <headless/public/headless_devtools_target.h>
#include <headless/public/headless_web_contents.h>

class Browser : public headless::HeadlessWebContents::Observer, headless::page::Observer {
  public:
    Browser();
    ~Browser() override;

    // Only returns after Shutdown is called or errors happen.
    int Run(std::function<void()> readyCb);

    // Open URL in new tab.
    bool OpenUrl(const std::string& url, int width, int height, std::function<void()> readyCb);

    void Evaluate(const std::string& script, std::function<void(bool, const std::string&)> resultCb);

    void Shutdown();

  private:
    void onStart(std::function<void()> readyCb, headless::HeadlessBrowser* browser);

    // headless::HeadlessWebContents::Observer method.
    void DevToolsTargetReady() override;

    // headless::page::Observer implementation:
    void OnLoadEventFired(const headless::page::LoadEventFiredParams& params) override;

    // Not owned.
    headless::HeadlessBrowser* browser_;
    // Not owned.
    headless::HeadlessWebContents* web_contents_; 
    std::unique_ptr<headless::HeadlessDevToolsClient> devtools_client_;
    std::function<void()> pageLoadedCb_;

    DISALLOW_COPY_AND_ASSIGN(Browser);
};

#endif  // BROWSER_HPP_
