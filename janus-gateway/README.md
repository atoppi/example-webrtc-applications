## video-room
This example demonstrates how to subscribe to a stream in a Janus video-room using Pion WebRTC.
It can be used to test RFC8888 feedbacks.

### Installing
OSX
```sh
brew install pkg-config
https://gstreamer.freedesktop.org/data/pkg/osx/

export PKG_CONFIG_PATH=/Library/Frameworks/GStreamer.framework/Versions/Current/lib/pkgconfig
```
Ubuntu
```sh
apt install pkg-config
apt install libgstreamer*
```

Build
```sh
cd video-room
go build
```

### Running
```
./video-room [--ws=ws://localhost:8188/janus] --room=1234 --feed=1000 [--enable-rfc8888] [--enable-stun]

```


