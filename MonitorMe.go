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

func main() {
	//initial paint the complete background
	for posX := 0; posX <= GRAPH_SIZE * GRAPH_AMOUNT; posX++ {
		vLine(icon, posX, 0, GRAPH_SIZE, background)
	}

	flushDataToIcon()
	systray.Run(onTrayReady)
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

	quitClickChan := systray.AddMenuItem("Quit", "").ClickedCh
	go func() {
		<-quitClickChan
		systray.Quit()
	}()

	go worker()
}

//draws a vertically line
func vLine(icon *image.RGBA, posX int, topY int, bottomY int, col color.RGBA) {
	for ; topY <= bottomY; topY++ {
		icon.Set(posX, topY, col)
	}
}

func worker() {
	graphs := [GRAPH_AMOUNT]MonitorData{}

	cpu := sigar.Cpu{}
	cpu.Get()

	//minus 1 because indexes starts with 0
	memGraph := &MemoryGraph{graph: Graph{0, GRAPH_SIZE - 1, systray.AddMenuItem("", "")}, mem: sigar.Mem{}}
	swapGraph := &SwapGraph{graph: Graph{GRAPH_SIZE, GRAPH_SIZE * 2 - 1, systray.AddMenuItem("", "")}, swap: sigar.Swap{}}
	loadGraph := &LoadGraph{graph: Graph{GRAPH_SIZE * 2, GRAPH_SIZE * 3 - 1, systray.AddMenuItem("", "")}, load: sigar.LoadAverage{}}
	cpuGraph := &CpuGraph{graph: Graph{GRAPH_SIZE * 3, GRAPH_SIZE * 4 - 1, systray.AddMenuItem("", "")}, cpu: cpu}
	diskGraph := &DiskGraph{graph: Graph{GRAPH_SIZE * 4, GRAPH_SIZE * 5 - 1, systray.AddMenuItem("", "")}, disk: sigar.DiskIo{}}
	netGraph := &NetworkGraph{graph: Graph{GRAPH_SIZE * 5, GRAPH_SIZE * 6 - 1, systray.AddMenuItem("", "")}, netIFace: sigar.NetIface{}}
	graphs[0] = memGraph
	graphs[1] = swapGraph
	graphs[2] = loadGraph
	graphs[3] = cpuGraph
	graphs[4] = diskGraph
	graphs[5] = netGraph

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

	//set the background of the new created row at the end
	vLine(icon, GRAPH_SIZE * GRAPH_AMOUNT, 0, GRAPH_SIZE, background)
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