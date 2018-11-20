package formation

import ()

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
		Path:  path,
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
		Path:         path,
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
	Print0     PrintInstruction
	Print1     PrintInstruction
}

type Bin struct {
	CycleCount     int
	MotorConf      MotorConf
	PrintheadConfs map[int]PrintheadConf
	Formations     map[int][]Formation
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
		PrintheadConfs: map[int]PrintheadConf{
			0: printhead0Conf,
			1: printhead1Conf,
		},
	}
}

func (b *Bin) AddFormation(
	cycleIndex int,
	posx string,
	posy string,
	print0Data string,
	print1Data string,
) {
	formation := Formation{
		CycleIndex: cycleIndex,
	}
	moveIns := MoveInstruction{
		MotorConf: b.MotorConf,
		PositionX: posx,
		PositionY: posy,
	}
	formation.Move = moveIns
	if print0Data != "" {
		print0Ins := PrintInstruction{
			PrintheadConf: b.PrintheadConfs[0],
			LineBuffer:    print0Data,
		}
		formation.Print0 = print0Ins
	}
	if print1Data != "" {
		print1Ins := PrintInstruction{
			PrintheadConf: b.PrintheadConfs[1],
			LineBuffer:    print1Data,
		}
		formation.Print1 = print1Ins
	}
	if b.Formations == nil {
		b.Formations = make(map[int][]Formation)
	}
	if _, ok := b.Formations[cycleIndex]; !ok {
		b.Formations[cycleIndex] = []Formation{}
	}
	b.Formations[cycleIndex] = append(b.Formations[cycleIndex], formation)
}
