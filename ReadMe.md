## MonitorMe

## Description

Cross-platform application to display performance charts in the taskbar as a trayicon.

## Supported data

* Memory usage (including cache)
* Swap
* Load average (*nix only)
* CPU usage (system and user)
* Network usage (only available in cloudfoundry/gosigar)
* Disk usage (only available in cloudfoundry/gosigar)

## Downloads

https://github.com/games647/MonitorMe/releases

## How to use

Using Commandline:

go run *.go

or:

./MonitorMe

Windows

MonitorMe.exe

## Image

![ubuntu showcase](http://i.imgur.com/9s8vXIz.png)

## Supported OS

* Linux
* Mac
* Windows (supports only one icon and load average isn't available)
Requires to compile a dll file: https://github.com/getlantern/systray#windows

Cross compiles doesn't seem to work, but it looks like if you compile it on the target machine:
example:
for linux compile it on a linux
for windows compile it on a windows
for mac compile it on a mac

## Credits/Dependencies

* Golang 1.6+
* https://github.com/getlantern/systray
* https://github.com/cloudfoundry/gosigar
* https://github.com/scalingdata/gosigar