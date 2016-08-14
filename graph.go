package main

import (
	"image/color"
	"github.com/scalingdata/gosigar"
	"image"
)

const GRAPH_SIZE = 100

var background = color.RGBA{50, 50, 50, 255} //grey
var memCol = color.RGBA{3, 177, 88, 255} //dark green
var totalMemCol = color.RGBA{159, 252, 206, 255} //green
var swapCol = color.RGBA{255, 0, 255, 255} //purple
var loadColor = color.RGBA{255, 0, 0, 255} //red
var userCpuCol = color.RGBA{51, 153, 255, 255} //blue
var systemCpuCol = color.RGBA{0, 0, 153, 255} //dark blue
var downloadCol = color.RGBA{238, 207, 25, 255} //dark yellow
var uploadCol = color.RGBA{242, 235, 113, 255} //yellow
var readCol = color.RGBA{178, 96, 1, 255} //dark orange
var writeCol = color.RGBA{245, 102, 17, 255} //orange

type Graph struct {
	startX        int
	endX          int
}

type MonitorData interface {
	collectData()
	drawGraph(icon *image.RGBA)
}

type MemoryGraph struct {
	graph Graph
	mem sigar.Mem
}

func (memoryGraph *MemoryGraph) collectData() {
	//update values
	memoryGraph.mem.Get()
}

func (memoryGraph *MemoryGraph) drawGraph(icon *image.RGBA) {
	//including cache
	totalMemPct := calculatePct(convertToKilo(memoryGraph.mem.Used), convertToKilo(memoryGraph.mem.Total))
	vLine(icon, memoryGraph.graph.endX, 100 - totalMemPct, GRAPH_SIZE, totalMemCol)

	usedPct := calculatePct(convertToKilo(memoryGraph.mem.ActualUsed), convertToKilo(memoryGraph.mem.Total))
	vLine(icon, memoryGraph.graph.endX, 100 - usedPct, GRAPH_SIZE, memCol)
}

type SwapGraph struct {
	graph Graph
	swap sigar.Swap
}

func (swapGraph *SwapGraph) collectData() {
	//update values
	swapGraph.swap.Get()
}

func (swapGraph *SwapGraph) drawGraph(icon *image.RGBA) {
	freeSwapPct := calculatePct(convertToKilo(swapGraph.swap.Used), convertToKilo(swapGraph.swap.Total))
	vLine(icon, swapGraph.graph.endX, 100 - freeSwapPct, GRAPH_SIZE, swapCol)
}

type LoadGraph struct {
	graph Graph
	load sigar.LoadAverage
}

func (loadGraph *LoadGraph) collectData() {
	//update values
	loadGraph.load.Get()
}

func (loadGraph *LoadGraph) drawGraph(icon *image.RGBA) {
	height := int(loadGraph.load.One * 100)
	vLine(icon, loadGraph.graph.endX, 100 - height, GRAPH_SIZE, loadColor)
}

type CpuGraph struct {
	graph Graph
	cpu sigar.Cpu

	diffUser int
	diffSystem int
	diffTotal int
}

func (cpuGraph *CpuGraph) collectData() {
	oldCpu := cpuGraph.cpu

	newCpu := sigar.Cpu{}
	newCpu.Get()

	cpuGraph.diffUser = int(newCpu.User - oldCpu.User)
	cpuGraph.diffSystem = int(newCpu.Sys - oldCpu.Sys)
	cpuGraph.diffTotal = int(newCpu.Total() - oldCpu.Total())

	cpuGraph.cpu = newCpu
}

func (cpuGraph *CpuGraph) drawGraph(icon *image.RGBA) {
	systemPct := cpuGraph.diffSystem
	vLine(icon, cpuGraph.graph.endX, 100 - systemPct, GRAPH_SIZE, systemCpuCol)

	userPct := cpuGraph.diffUser
	vLine(icon, cpuGraph.graph.endX, 100 - userPct, GRAPH_SIZE, userCpuCol)
}

type DiskGraph struct {
	graph Graph
	disk sigar.DiskIo

	diffRead int
	diffWrite int
}

func (diskGraph *DiskGraph) collectData() {
	oldDisk := diskGraph.disk

	diskList := sigar.DiskList{}
	diskList.Get()

	newDisk := sigar.DiskIo{}
	for _, disk := range diskList.List {
		newDisk = disk
		break
	}

	diskGraph.diffRead = int(convertToKilo(newDisk.ReadBytes - oldDisk.ReadBytes))
	diskGraph.diffWrite = int(convertToKilo(newDisk.WriteBytes - oldDisk.WriteBytes))

	diskGraph.disk = newDisk
}

func (diskGraph *DiskGraph) drawGraph(icon *image.RGBA) {
	readPct := calculatePct(diskGraph.diffRead, 50 * 1024)
	vLine(icon, diskGraph.graph.endX, 100 - readPct, GRAPH_SIZE, readCol)

	writePct := calculatePct(diskGraph.diffWrite, 50 * 1024)
	vLine(icon, diskGraph.graph.endX, 100 - writePct, GRAPH_SIZE, writeCol)
}

type NetworkGraph struct {
	graph Graph
	netIFace sigar.NetIface

	diffDownload int
	diffUpload int
}

func (netGraph *NetworkGraph) collectData() {
	oldNet := netGraph.netIFace

	netList := sigar.NetIfaceList{}
	netList.Get()

	newNet := sigar.NetIface{}
	for _, net := range netList.List {
		newNet = net
		break
	}

	netGraph.diffDownload = int(convertToKilo(newNet.RecvBytes - oldNet.RecvBytes))
	netGraph.diffUpload = int(convertToKilo(newNet.SendBytes - oldNet.SendBytes))

	netGraph.netIFace = newNet
}

func (netGraph *NetworkGraph) drawGraph(icon *image.RGBA) {
	readPct := calculatePct(netGraph.diffDownload, 50 * 1024)
	vLine(icon, netGraph.graph.endX, 100 - readPct, GRAPH_SIZE, downloadCol)

	writePct := calculatePct(netGraph.diffUpload, 50 * 1024)
	vLine(icon, netGraph.graph.endX, 100 - writePct, GRAPH_SIZE, uploadCol)
}
