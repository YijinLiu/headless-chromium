# Headless Chromium
Run headless Chromium as server and talk to it in client using devtools protocol.
We don't require installation of Chrome, which makes it suitable for resource constrained
environments.

It compiles headless Chromium into binary "hc_server" and generate libraries to start / talk to it.
Currently, only Golang library is supported. See go/demos/render for how to use it.

## Manual
<pre>
$ docker run -it --cap-add=SYS_ADMIN --name=hc.${USER} yijinliu/hc:57.0.2987.110
</pre>
Inside docker:
<pre>
$ go get github.com/yijinliu/headless-chromium/go/demos/render
$ ./bin/render
</pre>

## Build Docker container
### Dev container
<pre>
$ cd docker/dev
$ make VERSION=57.0.2987.129
</pre>
### HC container
<pre>
$ docker run -it --cap-add=SYS_ADMIN --name=hcdev.${USER} \
    -v /home/${USER}/projects/headless-chromium:/home/hcdev/external \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v $(which docker):/usr/bin/docker yijinliu/hc:dev
</pre>
Inside docker:
<pre>
$ git config --global user.name "XXX"
$ git config --global user.email "YYY@ZZZ"
$ cd external/docker/hc
$ make VERSION=57.0.2987.110
</pre>
