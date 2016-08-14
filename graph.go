package main

import (
	"image/color"
	"github.com/cloudfoundry/gosigar"
	"image"
	"log"
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
	log.Println(cpuGraph.graph.endX)
	vLine(icon, cpuGraph.graph.endX, 100 - systemPct, GRAPH_SIZE, systemCpuCol)

	userPct := cpuGraph.diffUser
	vLine(icon, cpuGraph.graph.endX, 100 - userPct, GRAPH_SIZE, userCpuCol)
}
