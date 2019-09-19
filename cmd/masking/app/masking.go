package app

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"synthesis/internal/formation"
	"synthesis/internal/geometry"
	"synthesis/internal/printhead"
	"synthesis/internal/reagent"
	"synthesis/internal/substrate"
	"synthesis/pkg/config"
)

const (
	CONF_SEQ_FILE   = "seqfile"
	CONF_ACT_ENABLE = "activator.enabled"

	CONF_MOTOR_STROKE_X = "motor.stroke.x"
	CONF_MOTOR_STROKE_Y = "motor.stroke.y"
	CONF_MOTOR_PATH     = "motor.path"
	CONF_MOTOR_SPEED    = "motor.speed"
	CONF_MOTOR_ACCEL    = "motor.accel"

	CONF_PRINTHEAD_PATH_0          = "printhead.path.0"
	CONF_PRINTHEAD_PATH_1          = "printhead.path.1"
	CONF_PRINTHEAD_SPACE           = "printhead.space"
	CONF_PRINTHEAD_OFFSET_0_X      = "printhead.offset.0.x"
	CONF_PRINTHEAD_OFFSET_1_X      = "printhead.offset.1.x"
	CONF_PRINTHEAD_OFFSET_0_Y      = "printhead.offset.0.y"
	CONF_PRINTHEAD_OFFSET_1_Y      = "printhead.offset.1.y"
	CONF_PRINTHEAD_REAGENT_0_LINE0 = "printhead.reagent.0.0"
	CONF_PRINTHEAD_REAGENT_0_LINE1 = "printhead.reagent.0.1"
	CONF_PRINTHEAD_REAGENT_0_LINE2 = "printhead.reagent.0.2"
	CONF_PRINTHEAD_REAGENT_0_LINE3 = "printhead.reagent.0.3"
	CONF_PRINTHEAD_REAGENT_1_LINE0 = "printhead.reagent.1.0"
	CONF_PRINTHEAD_REAGENT_1_LINE1 = "printhead.reagent.1.1"
	CONF_PRINTHEAD_REAGENT_1_LINE2 = "printhead.reagent.1.2"
	CONF_PRINTHEAD_REAGENT_1_LINE3 = "printhead.reagent.1.3"
	CONF_PRINTHEAD_ROW_RESOLUTION  = "printhead.row.resoultion"

	CONF_PRINT_MODE = "print.mode"

	CONF_SLIDE_COUNT_H      = "slide.count.horizon"
	CONF_SLIDE_COUNT_V      = "slide.count.vertical"
	CONF_SLIDE_WIDTH        = "slide.width"
	CONF_SLIDE_HEIGHT       = "slide.height"
	CONF_SLIDE_SPACE_H      = "slide.space.horizon"
	CONF_SLIDE_SPACE_V      = "slide.space.vertical"
	CONF_SLIDE_RESOLUTION_X = "slide.resolution.horizon"
	CONF_SLIDE_RESOLUTION_Y = "slide.resolution.vertical"
)

var offsetX, offsetY float64

type MaskCommand struct {
}

func NewMaskCommand() *MaskCommand {
	setDefaultConfig()
	return &MaskCommand{}
}

func (c *MaskCommand) Validate() (err error) {
	return nil
}

func (c *MaskCommand) Execute() error {
	array, err := buildPrintheadArray()
	if err != nil {
		return err
	}
	substrate, cycleCount, err := buildSubstrate()
	if err != nil {
		return err
	}
	build(
		cycleCount,
		array,
		substrate,
		config.GetString(CONF_MOTOR_PATH),
		config.GetString(CONF_MOTOR_SPEED),
		config.GetString(CONF_MOTOR_ACCEL),
		config.GetString(CONF_PRINTHEAD_PATH_0),
		config.GetInt(CONF_PRINT_MODE),
	)
	return nil
}

func build(
	cycleCount int,
	array *printhead.Array,
	substrate *substrate.Substrate,
	motorPath string,
	motorSpeed string,
	motorAccel string,
	p0p string,
	mode int,
) {
	bin := formation.NewBin(
		cycleCount,
		formation.NewMotorConf(
			motorPath,
			motorSpeed,
			motorAccel,
		),
		formation.NewPrintheadConf(
			p0p,
			"1",
			"2560",
			"320",
		),
		byte(mode),
		substrate,
		array,
	)
	fmt.Println("create bin", bin)
	img := image.NewRGBA(image.Rect(0, 0, substrate.Width, substrate.Height+1))
	lotc := make(chan int)
	close(lotc)
	paintedc := make(chan struct{})
	//countc := bin.Build(step, lotc, img, paintedc, true)
	countc, outputc := bin.BuildWithoutMotor(lotc, img, paintedc, true)
	//go func() {
	for _ = range countc {
		//for count := range countc {
		//fmt.Println(count)
		go func() {
			paintedc <- struct{}{}
		}()
	}
	//}()
	//for i := 0; i < len(bin.Formations); i++ {
	//fmt.Println("## cycle", i)
	//for _, f := range bin.Formations[i] {
	//fmt.Printf(">> %#v || (%v, %v)\n", f.Print.LineBuffers, dpi(f.Move.PositionX), dpi(f.Move.PositionY))
	//fmt.Println()
	//}
	//}
	output := <-outputc
	//fmt.Printf("%#v\n", output)
	fmt.Println(len(output))

	zero := strings.Repeat("0", 320)
	dash := strings.Repeat("-", 320)
	var f *os.File
	var err error
	count := 0
	for _, o := range output {
		if o[0] == dash {
			if f != nil {
				f.Close()
			}
			f, err = os.OpenFile(
				fmt.Sprintf("output-%d.txt", count),
				os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
			if err != nil {
				fmt.Println(err)
			}
			count++
		} else {
			fmt.Fprintln(f, o[0])
			if o[0] != zero {
				//fmt.Println(i, o[0])
			}
		}
	}
	f.Close()
}

func dpi(input string) int {
	x, _ := strconv.ParseFloat(input, 64)
	return int(x*600/25.4 + 0.5)
}

func setDefaultConfig() {
	config.SetDefault(CONF_SEQ_FILE, "input.txt")
	config.SetDefault(CONF_ACT_ENABLE, false)
	config.SetDefault(CONF_MOTOR_STROKE_X, 100.00)
	config.SetDefault(CONF_MOTOR_STROKE_Y, 100.00)
	config.SetDefault(CONF_MOTOR_PATH, "/AOZTECH/Motor")
	config.SetDefault(CONF_MOTOR_SPEED, 10)
	config.SetDefault(CONF_MOTOR_ACCEL, 100)
	config.SetDefault(CONF_PRINTHEAD_PATH_0, "/Ricoh-G5/Printer#2")
	config.SetDefault(CONF_PRINTHEAD_PATH_1, "/Ricoh-G5/Printer#2")
	config.SetDefault(CONF_PRINTHEAD_SPACE, 0)
	config.SetDefault(CONF_PRINTHEAD_REAGENT_0_LINE0, "A")
	config.SetDefault(CONF_PRINTHEAD_REAGENT_0_LINE1, "C")
	config.SetDefault(CONF_PRINTHEAD_REAGENT_0_LINE2, "G")
	config.SetDefault(CONF_PRINTHEAD_REAGENT_0_LINE3, "T")
	config.SetDefault(CONF_PRINTHEAD_REAGENT_1_LINE0, "-")
	config.SetDefault(CONF_PRINTHEAD_REAGENT_1_LINE1, "-")
	config.SetDefault(CONF_PRINTHEAD_REAGENT_1_LINE2, "-")
	config.SetDefault(CONF_PRINTHEAD_REAGENT_1_LINE3, "-")
	config.SetDefault(CONF_PRINTHEAD_ROW_RESOLUTION, 150)

	config.SetDefault(CONF_PRINT_MODE, formation.MODE_CIJ)
	config.SetDefault(CONF_SLIDE_COUNT_H, 3)
	config.SetDefault(CONF_SLIDE_COUNT_V, 1)
	config.SetDefault(CONF_SLIDE_WIDTH, 20)
	config.SetDefault(CONF_SLIDE_HEIGHT, 29)
	config.SetDefault(CONF_SLIDE_SPACE_H, 5)
	config.SetDefault(CONF_SLIDE_SPACE_V, 25)
	config.SetDefault(CONF_SLIDE_RESOLUTION_X, 600)
	config.SetDefault(CONF_SLIDE_RESOLUTION_Y, 600)
	config.SetDefault(CONF_PRINTHEAD_OFFSET_0_X, 35)
	config.SetDefault(CONF_PRINTHEAD_OFFSET_0_Y, 20)
	config.SetDefault(CONF_PRINTHEAD_OFFSET_1_X, 35)
	config.SetDefault(CONF_PRINTHEAD_OFFSET_1_Y, 65)

	//config.SafeWriteConfig()
}

func toFloat(inputString string) (float64, error) {
	if inputString == "" {
		return 0.0, fmt.Errorf("invalid float")
	}
	inputFloat, err := strconv.ParseFloat(inputString, 64)
	if err != nil {
		return 0.0, fmt.Errorf(
			"failed to convert %q to float: %v",
			inputString,
			err.Error(),
		)
	}
	return inputFloat, nil
}

func buildPrintheadArray() (*printhead.Array, error) {
	p0x := config.GetFloat(CONF_PRINTHEAD_OFFSET_0_X)
	p0y := config.GetFloat(CONF_PRINTHEAD_OFFSET_0_Y)
	p1x := config.GetFloat(CONF_PRINTHEAD_OFFSET_1_X)
	p1y := config.GetFloat(CONF_PRINTHEAD_OFFSET_1_Y)
	if err := validateOffsetX(p0x, config.GetFloat(CONF_MOTOR_STROKE_X)); err != nil {
		return nil, err
	}
	if err := validateOffsetY(p0y, config.GetFloat(CONF_SLIDE_HEIGHT)); err != nil {
		return nil, err
	}
	if err := validateOffsetX(p1x, config.GetFloat(CONF_MOTOR_STROKE_X)); err != nil {
		return nil, err
	}
	if err := validateOffsetY(p1y, config.GetFloat(CONF_SLIDE_HEIGHT)); err != nil {
		return nil, err
	}
	offsetX = p0x
	offsetY = p0y
	p0 := printhead.NewPrinthead(
		0,
		[]*reagent.Reagent{
			reagent.NewReagent(config.GetString(CONF_PRINTHEAD_REAGENT_0_LINE0)),
			reagent.NewReagent(config.GetString(CONF_PRINTHEAD_REAGENT_0_LINE1)),
			reagent.NewReagent(config.GetString(CONF_PRINTHEAD_REAGENT_0_LINE2)),
			reagent.NewReagent(config.GetString(CONF_PRINTHEAD_REAGENT_0_LINE3)),
		},
		config.GetString(CONF_PRINTHEAD_PATH_0),
		false,
		p0x,
		p0y,
		config.GetInt(CONF_PRINTHEAD_ROW_RESOLUTION),
	)
	p0xu := geometry.Millimeter2Dot(p0x)
	p0yu := geometry.Millimeter2Dot(p0y)
	n0 := p0.MakeNozzles(p0xu, p0yu)

	p1 := printhead.NewPrinthead(
		1,
		[]*reagent.Reagent{
			reagent.NewReagent(config.GetString(CONF_PRINTHEAD_REAGENT_1_LINE0)),
			reagent.NewReagent(config.GetString(CONF_PRINTHEAD_REAGENT_1_LINE1)),
			reagent.NewReagent(config.GetString(CONF_PRINTHEAD_REAGENT_1_LINE2)),
			reagent.NewReagent(config.GetString(CONF_PRINTHEAD_REAGENT_1_LINE3)),
		},
		config.GetString(CONF_PRINTHEAD_PATH_1),
		false,
		p1x,
		p1y,
		config.GetInt(CONF_PRINTHEAD_ROW_RESOLUTION),
	)
	step := geometry.DPI / config.GetInt(CONF_PRINTHEAD_ROW_RESOLUTION)
	deltay := geometry.Millimeter2Dot(p1y - p0y)
	yrem := deltay % step
	if yrem > step/2 {
		deltay += step - yrem
	} else {
		deltay -= yrem
	}
	deltax := getDeltax()
	xrem := deltax % step
	if xrem > step/2 {
		deltax += step - xrem
	} else {
		deltax -= xrem
	}
	n1 := p1.MakeNozzles(p0xu-deltax, p0yu+deltay)

	return printhead.NewArray(
		append(n0, n1...),
		2,
		[]*printhead.Printhead{p0, p1},
	), nil
}

func validateOffsetX(offset float64, stroke float64) error {
	max := stroke - 4*25.4/600
	if offset > max || offset < -max {
		return fmt.Errorf(
			"invalid printhead offset: %v > %v", offset, max)
	}
	return nil
}

func getDeltax() int {
	p0x := config.GetFloat(CONF_PRINTHEAD_OFFSET_0_X)
	p1x := config.GetFloat(CONF_PRINTHEAD_OFFSET_1_X)
	return geometry.Millimeter2Dot(p1x - p0x)
}

func buildSubstrate() (*substrate.Substrate, int, error) {
	sch := config.GetInt(CONF_SLIDE_COUNT_H)
	scv := config.GetInt(CONF_SLIDE_COUNT_V)
	ssh := config.GetFloat(CONF_SLIDE_SPACE_H)
	ssv := config.GetFloat(CONF_SLIDE_SPACE_V)
	sw := config.GetFloat(CONF_SLIDE_WIDTH)
	sh := config.GetFloat(CONF_SLIDE_HEIGHT)
	rx := config.GetInt(CONF_SLIDE_RESOLUTION_X)
	ry := config.GetInt(CONF_SLIDE_RESOLUTION_Y)
	seqText, err := loadSeqText()
	if err != nil {
		return nil, 0, err
	}
	spots, cycleCount := substrate.ParseSpots(
		seqText,
		config.GetBool(CONF_ACT_ENABLE),
	)
	sub, err := substrate.NewSubstrate(
		sch,
		scv,
		sw,
		sh,
		ssh,
		ssv,
		spots,
		getDeltax(),
		rx,
		ry,
	)
	if err != nil {
		return nil, 0, err
	}
	return sub, cycleCount, nil
}

func validateOffsetY(offset float64, height float64) error {
	max := offset - height
	if offset < 0 {
		max = height - offset
	}
	if max <= -config.GetFloat(CONF_MOTOR_STROKE_Y) {
		return fmt.Errorf(
			"invalid config: slide height '%v', position y '%v'",
			offset,
			height,
		)
	}
	return nil
}

func loadSeqText() (string, error) {
	seqBytes, err := ioutil.ReadFile(config.GetString(CONF_SEQ_FILE))
	if err != nil {
		return "", err
	}
	return string(seqBytes), nil
}
