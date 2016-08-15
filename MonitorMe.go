package main

import (
	"image"
	"image/color"
	"bytes"
	"image/jpeg"
	"log"
	"github.com/getlantern/systray"
	"time"
	"github.com/scalingdata/gosigar"
)

const GRAPH_AMOUNT = 6

var icon = image.NewRGBA(image.Rect(0, 0, GRAPH_SIZE * GRAPH_AMOUNT, GRAPH_SIZE))
var iconData []byte
var done = make(chan bool)

func main() {
	//initial paint the complete background
	for posX := 0; posX <= GRAPH_SIZE * GRAPH_AMOUNT; posX++ {
		vLine(icon, posX, 0, GRAPH_SIZE, background)
	}

	flushDataToIcon()
	go systray.Run(onTrayReady)
	<-done
}

func calculatePct(value int, max int) int {
	return value * 100 / max
}

func convertToKilo(val uint64) int {
	return int(val / 1024)
}

func onTrayReady() {
	//set a initial blank icon
	systray.SetIcon(iconData)

	loadAvg := sigar.LoadAverage{}
	netFace := sigar.NetIface{}
	cpu := sigar.Cpu{}
	cpu.Get()

	graphs := [GRAPH_AMOUNT]MonitorData{}
	//minus 1 because indexes starts with 0
	memGraph := &MemoryGraph{graph: Graph{0, GRAPH_SIZE - 1, blankItem()}, mem: sigar.Mem{}}
	swapGraph := &SwapGraph{graph: Graph{GRAPH_SIZE, GRAPH_SIZE * 2 - 1, blankItem()}, swap: sigar.Swap{}}

	loadGraph := &LoadGraph{graph: Graph{GRAPH_SIZE * 2, GRAPH_SIZE * 3 - 1, blankItem()}, load: loadAvg}
	cpuGraph := &CpuGraph{graph: Graph{GRAPH_SIZE * 3, GRAPH_SIZE * 4 - 1, blankItem()}, cpu: cpu}
	diskGraph := &DiskGraph{graph: Graph{GRAPH_SIZE * 4, GRAPH_SIZE * 5 - 1, blankItem()}, disk: sigar.DiskIo{}}
	netGraph := &NetworkGraph{graph: Graph{GRAPH_SIZE * 5, GRAPH_SIZE * 6 - 1, blankItem()}, netIFace: netFace}
	graphs[0] = memGraph
	graphs[1] = swapGraph
	graphs[2] = loadGraph
	graphs[3] = cpuGraph
	graphs[4] = diskGraph
	graphs[5] = netGraph

	quitClickChan := systray.AddMenuItem("Quit", "").ClickedCh
	go func() {
		<-quitClickChan
		systray.Quit()
		done <- true
	}()

	go worker(graphs)
}

//draws a vertically line
func vLine(icon *image.RGBA, posX int, topY int, bottomY int, col color.RGBA) {
	for ; topY <= bottomY; topY++ {
		icon.Set(posX, topY, col)
	}
}

func worker(graphs [GRAPH_AMOUNT]MonitorData) {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C

		scrollForward()
		for _, graph := range graphs {
			graph.collectData()
			graph.drawGraph(icon)
		}

		flushDataToIcon()
	}
}

//move the old data one row back
func scrollForward() {
	//start with the second row and just override the first row
	for posX := 1; posX < GRAPH_SIZE * GRAPH_AMOUNT; posX++ {
		for posY := 0; posY <= GRAPH_SIZE; posY++ {
			col := icon.At(posX, posY)
			icon.Set(posX - 1, posY, col)
		}
	}
}

func flushDataToIcon() {
	//encode image into slice array
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, icon, nil)
	if err != nil {
		log.Fatal(err)
	}

	iconData = buf.Bytes()
	systray.SetIcon(iconData)
}

func blankItem() *systray.MenuItem {
	return systray.AddMenuItem("", "")
}