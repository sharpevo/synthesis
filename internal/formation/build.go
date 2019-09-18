package formation

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
	"synthesis/internal/geometry"
	"synthesis/internal/log"
	"synthesis/internal/reagent"
)

const (
	BINFILE = "test.bin"
)

func (b *Bin) Build( // {{{
	step int,
	lotamountc chan int,
	img *image.RGBA,
	paintedc chan struct{},
	preview bool,
) (countc chan int) {
	fmt.Println("preview", preview)
	//lotamount := 1
	//stepwise := false
	countc = make(chan int)
	go func() {
		stepCount := 0
		lot := 0
		for cycleIndex := 0; cycleIndex < b.CycleCount; cycleIndex++ {
			<-paintedc
			img.Pix = make([]uint8, 4*b.Substrate().Width*(b.Substrate().Height+1))
			log.V("cycle index", cycleIndex).Debug()
			stripSum := b.Substrate().Strip()
			for stripCount := 0; stripCount < stripSum; stripCount++ {
				log.V("strip count", stripCount).Debug()
				posx := stripCount * 1280
				posy := b.Substrate().Top()
				rowIndex := 3
				b.PrintheadArray.MoveBottomRow(rowIndex, posx, posy)
				log.Df("move row %d to %d %d\n", rowIndex, posx, posy)
				for b.PrintheadArray.Top() >= b.Substrate().Bottom() {
					dataMap, count := b.genData(
						cycleIndex, img, &stepCount, lotamountc, &lot, paintedc, countc, preview)
					if count > 0 {
						log.Vs(log.M{
							"count":    count,
							"data map": dataMap,
						}).Debug("data downward #1")
						x, y := b.RawPos()
						b.AddFormation(cycleIndex, x, y, dataMap)
					}
					posy -= step
					b.PrintheadArray.MoveBottomRow(rowIndex, posx, posy)
				}
				if step == 1 {
					continue
				}
				posy = b.Substrate().Bottom()
				rowIndex = 2
				b.PrintheadArray.MoveTopRow(rowIndex, posx, posy)
				log.Df("move row %d to %d %d\n", rowIndex, posx, posy)
				for b.PrintheadArray.Bottom() <= b.Substrate().Top() {
					dataMap, count := b.genData(
						cycleIndex, img, &stepCount, lotamountc, &lot, paintedc, countc, preview)
					if count > 0 {
						log.Vs(log.M{
							"count":    count,
							"data map": dataMap,
						}).Debug("data upward #2")
						x, y := b.RawPos()
						b.AddFormation(cycleIndex, x, y, dataMap)
					}
					posy += step
					b.PrintheadArray.MoveTopRow(rowIndex, posx, posy)
				}
				posy = b.Substrate().Top()
				rowIndex = 1
				b.PrintheadArray.MoveBottomRow(rowIndex, posx, posy)
				log.Df("move row %d to %d %d\n", rowIndex, posx, posy)
				for b.PrintheadArray.Top() >= b.Substrate().Bottom() {
					dataMap, count := b.genData(
						cycleIndex, img, &stepCount, lotamountc, &lot, paintedc, countc, preview)
					if count > 0 {
						log.Vs(log.M{
							"count":    count,
							"data map": dataMap,
						}).Debug("data downward #3")
						x, y := b.RawPos()
						b.AddFormation(cycleIndex, x, y, dataMap)
					}
					posy -= step
					b.PrintheadArray.MoveBottomRow(rowIndex, posx, posy)
				}
				posy = b.Substrate().Bottom()
				rowIndex = 0
				b.PrintheadArray.MoveTopRow(rowIndex, posx, posy)
				log.Df("move row %d to %d %d\n", rowIndex, posx, posy)
				for b.PrintheadArray.Bottom() <= b.Substrate().Top() {
					dataMap, count := b.genData(
						cycleIndex, img, &stepCount, lotamountc, &lot, paintedc, countc, preview)
					if count > 0 {
						log.Vs(log.M{
							"count":    count,
							"data map": dataMap,
						}).Debug("data upward #4")
						x, y := b.RawPos()
						b.AddFormation(cycleIndex, x, y, dataMap)
					}
					posy += step
					b.PrintheadArray.MoveTopRow(rowIndex, posx, posy)
				}
			}
			if preview && b.Mode == MODE_CIJ {
				outputFile, _ := os.Create(
					fmt.Sprintf("output/CIJ.%03d.png", cycleIndex))
				png.Encode(outputFile, img)
				outputFile.Close()
				//if lotamount, stepwise = <-lotamountc; stepwise {
				//if lotamount < 0 {
				//cycleIndex += lotamount - 2
				//if cycleIndex < -1 {
				//cycleIndex = 0 - 1
				//}
				//go func() {
				//paintedc <- struct{}{}
				//lotamountc <- 1
				//lot = 0
				//}()
				//continue
				//} else {
				//lot++
				//}
				//if lot == lotamount ||
				//cycleIndex+1 >= b.CycleCount { // even step may not 100%
				//countc <- cycleIndex + 1
				//lot = 0
				//} else {
				//go func() {
				//paintedc <- struct{}{}
				//lotamountc <- lotamount
				//}()
				//}
				//} else {
				//countc <- cycleIndex + 1
				//}
			}
			if b.check(lotamountc, &cycleIndex, &lot, paintedc, countc) {
				continue
			}
		}
		err := b.SaveToFile(BINFILE)
		if err != nil {
			log.E(err.Error())
		}
		close(countc)
	}()
	paintedc <- struct{}{}
	return countc
} // }}}

func (b *Bin) genData( // {{{
	cycleIndex int,
	img *image.RGBA,
	stepCount *int,

	lotamountc chan int,
	lot *int,
	paintedc chan struct{},
	countc chan int,
	preview bool,
) ([]string, int) {
	count := 0
	dataSlice := make([][]string, b.PrintheadArray.PrintheadCount)
	for _, nozzle := range b.PrintheadArray.Nozzles {
		if dataSlice[nozzle.Printhead.Index] == nil {
			dataSlice[nozzle.Printhead.Index] = make([]string, 1280)
		}
		dataSlice[nozzle.Printhead.Index][nozzle.Index] = "0"
		if nozzle.Reagent.Equal(reagent.Nil) {
			continue
		}
		//log.D(nozzle.Pos.X, nozzle.Pos.Y, b.Substrate.Width, b.Substrate.Height)
		if nozzle.Pos.Y >= b.Substrate().Height ||
			nozzle.Pos.Y < 0 ||
			nozzle.Pos.X >= b.Substrate().Width ||
			nozzle.Pos.X < 0 {
			continue
		}
		spot := b.Substrate().Spots[nozzle.Pos.Y][nozzle.Pos.X]
		if spot == nil || cycleIndex > len(spot.Reagents)-1 {
			//log.Vs(log.M{"spot": spot, "cycle index": cycleIndex}).
			//Error("not enough reagents")
			continue
		}
		if spot != nil &&
			nozzle.Reagent.Equal(spot.Reagents[cycleIndex]) {
			count += 1
			dataSlice[nozzle.Printhead.Index][nozzle.Index] = "1"

			//log.D("printing ", nozzle.Reagent.Name, nozzle.Pos.X, nozzle.Pos.Y)
			log.D("printing ", nozzle.Reagent.Name, nozzle.Pos.X, b.Substrate().Height-nozzle.Pos.Y)
			if preview {
				img.Set(nozzle.Pos.X, b.Substrate().Height-nozzle.Pos.Y, nozzle.Reagent.Color)
			}
		}
	}

	output := make([]string, b.PrintheadArray.PrintheadCount)
	if count > 0 {
		for deviceIndex, dataBinSlice := range dataSlice {
			dataHexSlice := make([]string, 160)
			for i := 0; i < len(dataBinSlice); i += 8 {
				value, _ := strconv.ParseInt(
					strings.Join(dataBinSlice[i:i+8], ""), 2, 64)
				dataHexSlice = append(dataHexSlice, fmt.Sprintf("%02x", value))
			}
			output[deviceIndex] = strings.Join(dataHexSlice, "")
			log.V("device index", deviceIndex).Debug()
			log.V("dataBinSlice", dataBinSlice[:16]).Debug()
			log.V("linebuffer", output[deviceIndex][:8]).Debug()
		}
		if preview && b.Mode == MODE_DOD {
			outputFile, _ := os.Create(
				fmt.Sprintf("output/DoD.%06d.%03d.png", *stepCount, cycleIndex))
			png.Encode(outputFile, img)
			outputFile.Close()
			*stepCount = *stepCount + 1
		}
		b.check(lotamountc, &cycleIndex, lot, paintedc, countc)
	}
	return output, count
} // }}}

func (b *Bin) RawPos() (string, string) { // {{{
	return fmt.Sprintf("%.6f", geometry.Raw(b.PrintheadArray.SightBottom.Pos.X, b.PrintheadArray.OffsetX())),
		fmt.Sprintf("%.6f", geometry.Raw(b.PrintheadArray.SightBottom.Pos.Y, b.PrintheadArray.OffsetY()))
} // }}}

func (b *Bin) check( // {{{
	lotamountc chan int,
	cycleIndex *int,
	lot *int,
	paintedc chan struct{},
	countc chan int,
) (continued bool) {
	if lotamount, stepwise := <-lotamountc; stepwise {
		if lotamount < 0 {
			*cycleIndex += lotamount - 2
			if *cycleIndex < -1 {
				*cycleIndex = 0 - 1
			}
			go func() {
				paintedc <- struct{}{}
				lotamountc <- 1
				*lot = 0
			}()
			return true
		} else {
			*lot++
		}
		if *lot == lotamount ||
			*cycleIndex+1 >= b.CycleCount { // even step may not 100%
			countc <- *cycleIndex + 1
			*lot = 0
		} else {
			go func() {
				paintedc <- struct{}{}
				lotamountc <- lotamount
			}()
		}
	} else {
		countc <- *cycleIndex + 1
	}
	return false
} // }}}

type Output [][]string

var SEPARATOR = []string{strings.Repeat("-", 320), strings.Repeat("-", 320)}

func (b *Bin) BuildWithoutMotor(
	step int,
	lotamountc chan int,
	img *image.RGBA,
	paintedc chan struct{},
	preview bool,
) (countc chan int, outputc chan Output) {
	output := Output{}
	fmt.Println("preview", preview)
	countc = make(chan int)
	outputc = make(chan Output)
	go func() {
		stepCount := 0
		lot := 0
		for cycleIndex := 0; cycleIndex < b.CycleCount; cycleIndex++ {
			<-paintedc
			img.Pix = make([]uint8, 4*b.Substrate().Width*(b.Substrate().Height+1))
			log.V("cycle index", cycleIndex).Debug()
			stripSum := b.Substrate().Strip()
			for stripCount := 0; stripCount < stripSum; stripCount++ {
				log.V("strip count", stripCount).Debug()
				posx := stripCount * 1280
				posy := b.Substrate().Top()

				output = append(output, SEPARATOR)

				rowIndex := 3 // {{{
				b.PrintheadArray.MoveBottomRow(rowIndex, posx, posy)
				log.Df("move row %d to %d %d\n", rowIndex, posx, posy)
				for b.PrintheadArray.Top() >= b.Substrate().Bottom() {
					dataMap, count := b.genDataWithoutMotor(
						cycleIndex, img, &stepCount, lotamountc, &lot, paintedc, countc, preview)
					if count > 0 {
						log.Vs(log.M{
							"count":    count,
							"data map": dataMap,
						}).Debug("data downward #1")
						x, y := b.RawPos()
						b.AddFormation(cycleIndex, x, y, dataMap)
						output = append(output, dataMap)
					}
					posy -= step
					b.PrintheadArray.MoveBottomRow(rowIndex, posx, posy)
				} // }}}

				if step == 1 {
					continue
				}

				output = append(output, SEPARATOR)
				posy = b.Substrate().Bottom()

				rowIndex = 2 // {{{
				b.PrintheadArray.MoveTopRow(rowIndex, posx, posy)
				log.Df("move row %d to %d %d\n", rowIndex, posx, posy)
				for b.PrintheadArray.Bottom() <= b.Substrate().Top() {
					dataMap, count := b.genDataWithoutMotor(
						cycleIndex, img, &stepCount, lotamountc, &lot, paintedc, countc, preview)
					if count > 0 {
						log.Vs(log.M{
							"count":    count,
							"data map": dataMap,
						}).Debug("data upward #2")
						x, y := b.RawPos()
						b.AddFormation(cycleIndex, x, y, dataMap)
						output = append(output, dataMap)
					}
					posy += step
					b.PrintheadArray.MoveTopRow(rowIndex, posx, posy)
				} // }}}

				output = append(output, SEPARATOR)
				posy = b.Substrate().Top()

				rowIndex = 1 // {{{
				b.PrintheadArray.MoveBottomRow(rowIndex, posx, posy)
				log.Df("move row %d to %d %d\n", rowIndex, posx, posy)
				for b.PrintheadArray.Top() >= b.Substrate().Bottom() {
					dataMap, count := b.genDataWithoutMotor(
						cycleIndex, img, &stepCount, lotamountc, &lot, paintedc, countc, preview)
					if count > 0 {
						log.Vs(log.M{
							"count":    count,
							"data map": dataMap,
						}).Debug("data downward #3")
						x, y := b.RawPos()
						b.AddFormation(cycleIndex, x, y, dataMap)
						output = append(output, dataMap)
					}
					posy -= step
					b.PrintheadArray.MoveBottomRow(rowIndex, posx, posy)
				} // }}}

				output = append(output, SEPARATOR)
				posy = b.Substrate().Bottom()

				rowIndex = 0 // {{{
				b.PrintheadArray.MoveTopRow(rowIndex, posx, posy)
				log.Df("move row %d to %d %d\n", rowIndex, posx, posy)
				for b.PrintheadArray.Bottom() <= b.Substrate().Top() {
					dataMap, count := b.genDataWithoutMotor(
						cycleIndex, img, &stepCount, lotamountc, &lot, paintedc, countc, preview)
					if count > 0 {
						log.Vs(log.M{
							"count":    count,
							"data map": dataMap,
						}).Debug("data upward #4")
						x, y := b.RawPos()
						b.AddFormation(cycleIndex, x, y, dataMap)
						output = append(output, dataMap)
					}
					posy += step
					b.PrintheadArray.MoveTopRow(rowIndex, posx, posy)
				} // }}}

			}
			//if preview && b.Mode == MODE_CIJ {
			//outputFile, _ := os.Create(
			//fmt.Sprintf("output/CIJ.%03d.png", cycleIndex))
			//png.Encode(outputFile, img)
			//outputFile.Close()
			//}
			if b.check(lotamountc, &cycleIndex, &lot, paintedc, countc) {
				continue
			}
		}
		err := b.SaveToFile(BINFILE)
		if err != nil {
			log.E(err.Error())
		}
		close(countc)
		outputc <- output
	}()
	paintedc <- struct{}{}
	return countc, outputc
}

func (b *Bin) genDataWithoutMotor(
	cycleIndex int,
	img *image.RGBA,
	stepCount *int,

	lotamountc chan int,
	lot *int,
	paintedc chan struct{},
	countc chan int,
	preview bool,
) ([]string, int) {
	count := 0
	valid := false
	dataSlice := make([][]string, b.PrintheadArray.PrintheadCount)
	for _, nozzle := range b.PrintheadArray.Nozzles {
		if dataSlice[nozzle.Printhead.Index] == nil {
			dataSlice[nozzle.Printhead.Index] = make([]string, 1280)
		}
		dataSlice[nozzle.Printhead.Index][nozzle.Index] = "0"
		if nozzle.Reagent.Equal(reagent.Nil) {
			count += 1
			continue
		}
		//log.D(nozzle.Pos.X, nozzle.Pos.Y, b.Substrate.Width, b.Substrate.Height)
		if nozzle.Pos.Y >= b.Substrate().Height ||
			nozzle.Pos.Y < 0 ||
			nozzle.Pos.X >= b.Substrate().Width ||
			nozzle.Pos.X < 0 {
			continue
		}
		spot := b.Substrate().Spots[nozzle.Pos.Y][nozzle.Pos.X]
		if spot == nil || cycleIndex > len(spot.Reagents)-1 {
			//log.Vs(log.M{"spot": spot, "cycle index": cycleIndex}).
			//Error("not enough reagents")
			continue
		}
		if spot != nil &&
			nozzle.Reagent.Equal(spot.Reagents[cycleIndex]) {
			count += 1
			valid = true
			dataSlice[nozzle.Printhead.Index][nozzle.Index] = "1"

			//log.D("printing ", nozzle.Reagent.Name, nozzle.Pos.X, nozzle.Pos.Y)
			log.D("printing ", nozzle.Reagent.Name, nozzle.Pos.X, b.Substrate().Height-nozzle.Pos.Y)
			if preview {
				img.Set(nozzle.Pos.X, b.Substrate().Height-nozzle.Pos.Y, nozzle.Reagent.Color)
			}
		}
	}

	output := make([]string, b.PrintheadArray.PrintheadCount)
	if count > 0 {
		for deviceIndex, dataBinSlice := range dataSlice {
			dataHexSlice := make([]string, 160)
			for i := 0; i < len(dataBinSlice); i += 8 {
				value, _ := strconv.ParseInt(
					strings.Join(dataBinSlice[i:i+8], ""), 2, 64)
				dataHexSlice = append(dataHexSlice, fmt.Sprintf("%02x", value))
				//dataHexSlice = append(dataHexSlice, byte(value))
			}
			output[deviceIndex] = strings.Join(dataHexSlice, "")
			//output[deviceIndex] = dataHexSlice
			log.V("device index", deviceIndex).Debug()
			log.V("dataBinSlice", dataBinSlice[:16]).Debug()
			log.V("linebuffer", output[deviceIndex][:8]).Debug()
		}
		//if preview && b.Mode == MODE_DOD {
		if preview && valid {
			outputFile, _ := os.Create(
				//fmt.Sprintf("output/DoD.%06d.%03d.png", *stepCount, cycleIndex))
				fmt.Sprintf("output/%06d.%03d.png", *stepCount, cycleIndex))
			png.Encode(outputFile, img)
			outputFile.Close()
			*stepCount = *stepCount + 1
		}
		b.check(lotamountc, &cycleIndex, lot, paintedc, countc)
	}
	return output, count
}
