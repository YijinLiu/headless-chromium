#include "browser.hpp"

#include <base/bind.h>
#include <base/callback.h>
#include <base/command_line.h>
#include <base/logging.h>
#include <base/strings/string_number_conversions.h>
#include <net/base/ip_address.h>
#include <ui/gfx/geometry/size.h>
#include <url/gurl.h>

Browser::Browser() : browser_(nullptr), browser_context_(nullptr), web_contents_(nullptr),
                     devtools_client_(headless::HeadlessDevToolsClient::Create()) {}

Browser::~Browser() {}

namespace {

const char kPort[] = "port";
const char kAddr[] = "addr";
const char kProxy[] = "proxy";

const int kDefaultPort = 9222;
const char kDefaultAddr[] = "127.0.0.1";

}  // namespace

int Browser::Run(int argc, const char** argv, std::function<void()> readyCb) {
	base::CommandLine command_line(argc, argv);
    command_line.AppendSwitch("disable-extensions");
    command_line.AppendSwitch("disable-gpu");
    base::CommandLine::SwitchMap switche_map = command_line.GetSwitches();

    // Parse port.
    int port = kDefaultPort;
    if (command_line.HasSwitch(kPort)) {
        const std::string port_str = command_line.GetSwitchValueASCII(kPort);
        switche_map.erase(kPort);
        if (!base::StringToInt(port_str, &port) ||
            !base::IsValueInRangeForNumericType<uint16_t>(port)) {
            LOG(FATAL) << "Invalid devtools server port: " << port_str;
        }
    }
    
    // Parse addr.
    std::string addr = kDefaultAddr;
    if (command_line.HasSwitch(kAddr)) {
        addr = command_line.GetSwitchValueASCII(kAddr);
        switche_map.erase(kAddr);
    }
    net::IPAddress parsed_addr;
    if (!net::ParseURLHostnameToAddress(addr, &parsed_addr)) {
        LOG(FATAL) << "Invalid devtools server address: " << addr;
    }

    // Has proxy server?
    std::string proxy_server;
    if (command_line.HasSwitch(kProxy)) {
        proxy_server = command_line.GetSwitchValueASCII(kProxy);
        switche_map.erase(kProxy);
    }

    argc = switche_map.size();
    argv = new const char*[argc];
    auto* args = new std::string[argc];
    int i = 0;
    for (const auto& pair : switche_map) {
        args[i] = "--" + pair.first + "=" + pair.second;
        argv[i] = args[i].c_str();
        i++;
    }
    headless::HeadlessBrowser::Options::Builder builder(argc, argv);
    builder.EnableDevToolsServer(net::IPEndPoint(parsed_addr, base::checked_cast<uint16_t>(port)));
    LOG(INFO) << "Opening devtools port on " << addr << ":" << port << " ...";
    if (!proxy_server.empty()) {
        net::HostPortPair parsed_proxy_server = net::HostPortPair::FromString(proxy_server);
        if (parsed_proxy_server.host().empty() || !parsed_proxy_server.port()) {
            LOG(FATAL) << "Malformed proxy server: " << proxy_server;
        }
        builder.SetProxyServer(parsed_proxy_server);
    }
    return headless::HeadlessBrowserMain(
        builder.Build(), base::Bind(&Browser::onStart, base::Unretained(this), readyCb));
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
