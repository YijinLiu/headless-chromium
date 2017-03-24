# Headless Chromium
==========
Run headless Chromium as server and talk to it in client using devtools protocol.
We don't require installation of Chrome, which makes it suitable for environments that don't have
many resources.

It compiles headless Chromium into binary "hc_server" and generate libraries to start / talk to it.
Currently, only Golang library is supported. See go/demos/render for how to use it.

## Manual
<pre>
$ go get -d github.com/yijinliu/headless-chromium
$ src/github.com/yijinliu/headless-chromium/scripts/install_chromium.sh
$ src/github.com/yijinliu/headless-chromium/scripts/setup.sh
$ go install github.com/yijinliu/headless-chromium/go/demos/render
$ ./bin/render
</pre>
