## MonitorMe

## Description

Cross-platform application to display performance charts in the taskbar as a trayicon.

## Supported data

* Memory usage (including cache)
* Swap
* Load average (*nix only)

### Planned 

* CPU usage (system and user) - have to compare it with a every second - not implemented yet
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
* Mac (untested)
* Windows (supports only one icon and load average isn't available)

## Credits/Dependencies

* Golang 1.6+
* https://github.com/getlantern/systray
* https://github.com/cloudfoundry/gosigar
* https://github.com/scalingdata/gosigar (assuming Linux only - cross compile to windows failed)