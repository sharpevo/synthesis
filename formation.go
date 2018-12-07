package formation

import (
	"fmt"
	"sync"
)

type MotorConf struct {
	Path  string
	Speed string
	Accel string
}

func NewMotorConf(
	path string,
	speed string,
	accel string,
) MotorConf {
	return MotorConf{
		Path:  fmt.Sprintf("%s/CONN/DEVICE_NAME", path),
		Speed: speed,
		Accel: accel,
	}
}

type MoveInstruction struct {
	MotorConf MotorConf
	PositionX string
	PositionY string
}

type PrintheadConf struct {
	Path         string
	BitsPerPixel string
	Width        string
	BufferSize   string
}

func NewPrintheadConf(
	path string,
	bitsPerPixel string,
	width string,
	bufferSize string,
) PrintheadConf {
	return PrintheadConf{
		Path:         fmt.Sprintf("%s/CONN/ADDRESS", path),
		BitsPerPixel: bitsPerPixel,
		Width:        width,
		BufferSize:   bufferSize,
	}
}

type PrintInstruction struct {
	PrintheadConf PrintheadConf
	LineBuffers   []string
}

type Formation struct {
	CycleIndex int
	Move       MoveInstruction
	Prints     []PrintInstruction
	Print      PrintInstruction
}

type Cycle struct {
	CycleIndex     int
	ReagentCycle   Formation
	ActivatorCycle Formation
}

type Bin struct {
	//sync.RWMutex
	CycleCount    int
	MotorConf     MotorConf
	PrintheadConf PrintheadConf
	Formations    map[int][]Formation
	CycleMap      map[int]Cycle
}

func NewBin(
	cycleCount int,
	motorConf MotorConf,
	printheadConf PrintheadConf,
) *Bin {
	return &Bin{
		CycleCount:    cycleCount,
		MotorConf:     motorConf,
		PrintheadConf: printheadConf,
	}
}

func (b *Bin) AddFormation(
	cycleIndex int,
	posx string,
	posy string,
	dataSlice []string,
) {
	//b.Lock()
	//defer b.Unlock()
	formation := Formation{
		CycleIndex: cycleIndex,
	}
	x, y := posx, posy
	//fmt.Println("move to", x, y)
	moveIns := MoveInstruction{
		MotorConf: b.MotorConf,
		PositionX: x,
		PositionY: y,
	}
	formation.Move = moveIns

	formation.Print = PrintInstruction{
		PrintheadConf: b.PrintheadConf,
		LineBuffers:   dataSlice,
	}

	if b.Formations == nil {
		b.Formations = make(map[int][]Formation)
	}
	if _, ok := b.Formations[cycleIndex]; !ok {
		b.Formations[cycleIndex] = []Formation{}
	}
	b.Formations[cycleIndex] = append(b.Formations[cycleIndex], formation)
}
