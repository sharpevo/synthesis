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
	LineBuffer    string
}

type Formation struct {
	CycleIndex int
	Move       MoveInstruction
	Prints     []PrintInstruction
}

type Cycle struct {
	CycleIndex     int
	ReagentCycle   Formation
	ActivatorCycle Formation
}

type Bin struct {
	sync.RWMutex
	CycleCount     int
	MotorConf      MotorConf
	PrintheadConfs map[string]PrintheadConf
	Formations     map[int][]Formation
	CycleMap       map[int]Cycle
}

func NewBin(
	cycleCount int,
	motorConf MotorConf,
	printhead0Conf PrintheadConf,
	printhead1Conf PrintheadConf,
) *Bin {
	return &Bin{
		CycleCount: cycleCount,
		MotorConf:  motorConf,
		PrintheadConfs: map[string]PrintheadConf{
			printhead0Conf.Path: printhead0Conf,
			printhead1Conf.Path: printhead1Conf,
		},
	}
}

func (b *Bin) AddFormation(
	cycleIndex int,
	posx string,
	posy string,
	dataMap map[string]string,
) {
	b.Lock()
	defer b.Unlock()
	formation := Formation{
		CycleIndex: cycleIndex,
	}
	//x := fmt.Sprintf("%.6f", float64(posx)/float64(geometry.MM))
	//y := fmt.Sprintf("%.6f", float64(posy)/float64(geometry.MM))
	x, y := posx, posy
	fmt.Println("motion in aoztech", x, y)
	moveIns := MoveInstruction{
		MotorConf: b.MotorConf,
		PositionX: x,
		PositionY: y,
	}
	formation.Move = moveIns

	for devicePath, data := range dataMap {
		ins := PrintInstruction{
			PrintheadConf: b.PrintheadConfs[devicePath],
			LineBuffer:    data,
		}
		formation.Prints = append(formation.Prints, ins)
	}

	if b.Formations == nil {
		b.Formations = make(map[int][]Formation)
	}
	if _, ok := b.Formations[cycleIndex]; !ok {
		b.Formations[cycleIndex] = []Formation{}
	}
	b.Formations[cycleIndex] = append(b.Formations[cycleIndex], formation)
}
