package formation

func build(
	step int,
	cycleCount int,
	printheadArray *printhead.Array,
	subs *substrate.Substrate,
	motorPath string,
	motorSpeed string,
	motorAccel string,
	printhead0Path string,
	buildProgressbar *widgets.QProgressBar,
) {
	bin := formation.NewBin(
		cycleCount,
		formation.NewMotorConf(
			motorPath,
			motorSpeed,
			motorAccel,
		),
		formation.NewPrintheadConf(
			printhead0Path,
			"1",
			"2560",
			"320",
		),
	)
	fmt.Println("create bin", bin)
	imageIndex := 0

	go func() {
		for cycleIndex := 0; cycleIndex < cycleCount; cycleIndex++ {
			img := image.NewRGBA(image.Rect(0, 0, subs.Width, subs.Height+1))
			fmt.Println("cycle ", cycleIndex)
			stripSum := subs.Strip()
			for stripCount := 0; stripCount < stripSum; stripCount++ {

				fmt.Println("strip ", stripCount)
				posx := stripCount * 1280
				posy := subs.Top()

				rowIndex := 3

				printheadArray.MoveBottomRow(rowIndex, posx, posy)
				fmt.Printf("move row %v to %v %v\n", rowIndex, posx, posy)

				for printheadArray.Top() >= subs.Bottom() {
					dataMap, count := genData(cycleIndex, printheadArray, subs, img, &imageIndex)
					if count > 0 {
						if DEBUGABLE {
							fmt.Println("data downward #1", count, dataMap)
						}
						x, y := RawPos(
							printheadArray.SightBottom.Pos.X,
							printheadArray.SightBottom.Pos.Y,
						)
						bin.AddFormation(
							cycleIndex, x, y, dataMap,
						)
					}
					posy -= step
					printheadArray.MoveBottomRow(rowIndex, posx, posy)
				}

				if step == 1 {
					continue
				}

				posy = subs.Bottom()
				// distance is integer multiple of 4
				// so that move one row will match the others
				rowIndex = 2
				printheadArray.MoveTopRow(rowIndex, posx, posy)
				fmt.Printf("move row %v to %v %v\n", rowIndex, posx, posy)

				for printheadArray.Bottom() <= subs.Top() {
					dataMap, count := genData(cycleIndex, printheadArray, subs, img, &imageIndex)
					if count > 0 {
						if DEBUGABLE {
							fmt.Println("data upward #2", count, dataMap)
						}
						// use the bottom position
						// sinc the offset is bottomed
						x, y := RawPos(
							printheadArray.SightBottom.Pos.X,
							printheadArray.SightBottom.Pos.Y,
						)
						bin.AddFormation(
							cycleIndex, x, y, dataMap,
						)
					}
					posy += step
					printheadArray.MoveTopRow(rowIndex, posx, posy)
				}

				posy = subs.Top()
				rowIndex = 1
				printheadArray.MoveBottomRow(rowIndex, posx, posy)
				fmt.Printf("move row %v to %v %v\n", rowIndex, posx, posy)

				for printheadArray.Top() >= subs.Bottom() {
					dataMap, count := genData(cycleIndex, printheadArray, subs, img, &imageIndex)
					if count > 0 {
						if DEBUGABLE {
							fmt.Println("data downward #3", count, dataMap)
						}
						x, y := RawPos(
							printheadArray.SightBottom.Pos.X,
							printheadArray.SightBottom.Pos.Y,
						)
						bin.AddFormation(
							cycleIndex, x, y, dataMap,
						)
					}
					posy -= step
					printheadArray.MoveBottomRow(rowIndex, posx, posy)
				}

				posy = subs.Bottom()
				rowIndex = 0
				printheadArray.MoveTopRow(rowIndex, posx, posy)
				fmt.Printf("move row %v to %v %v\n", rowIndex, posx, posy)

				for printheadArray.Bottom() <= subs.Top() {
					dataMap, count := genData(cycleIndex, printheadArray, subs, img, &imageIndex)
					if count > 0 {
						if DEBUGABLE {
							fmt.Println("data upward #4", count, dataMap)
						}
						// use the bottom position
						// sinc the offset is bottomed
						x, y := RawPos(
							printheadArray.SightBottom.Pos.X,
							printheadArray.SightBottom.Pos.Y,
						)
						bin.AddFormation(
							cycleIndex, x, y, dataMap,
						)
					}
					posy += step
					printheadArray.MoveTopRow(rowIndex, posx, posy)
				}
			}
			buildProgressbar.SetValue((cycleIndex + 1) * buildProgressbar.Maximum() / cycleCount)
		}
		err := bin.SaveToFile("test.bin")
		if err != nil {
			fmt.Println(err)
		}
	}()

}

func genData(
	cycleIndex int,
	printheadArray *printhead.Array,
	subs *substrate.Substrate,
	img *image.RGBA,
	imageIndex *int,
) ([]string, int) {
	count := 0
	dataSlice := make([][]string, printheadArray.PrintheadCount)
	for _, nozzle := range printheadArray.Nozzles {
		if dataSlice[nozzle.Printhead.Index] == nil {
			dataSlice[nozzle.Printhead.Index] = make([]string, 1280)
		}
		dataSlice[nozzle.Printhead.Index][nozzle.Index] = "0"
		if nozzle.Reagent.Equal(reagent.Nil) {
			continue
		}
		//fmt.Println(nozzle.Pos.X, nozzle.Pos.Y, subs.Width, subs.Height)
		if nozzle.Pos.Y >= subs.Height ||
			nozzle.Pos.Y < 0 ||
			nozzle.Pos.X >= subs.Width ||
			nozzle.Pos.X < 0 {
			continue
		}
		spot := subs.Spots[nozzle.Pos.Y][nozzle.Pos.X]
		if spot == nil || cycleIndex > len(spot.Reagents)-1 {
			//fmt.Println("not enough reagents")
			continue
		}
		if spot != nil &&
			nozzle.Reagent.Equal(spot.Reagents[cycleIndex]) {
			count += 1
			dataSlice[nozzle.Printhead.Index][nozzle.Index] = "1"

			if DEBUGABLE {
				//fmt.Printf(" | printing ", nozzle.Reagent.Name, nozzle.Pos.X, nozzle.Pos.Y)
			}
			if IMAGABLE {
				img.Set(nozzle.Pos.X, subs.Height-nozzle.Pos.Y, nozzle.Reagent.Color)
			}
		}
	}

	output := make([]string, printheadArray.PrintheadCount)
	if count > 0 {
		for deviceIndex, dataBinSlice := range dataSlice {
			dataHexSlice := make([]string, 160)
			for i := 0; i < len(dataBinSlice); i += 8 {
				value, _ := strconv.ParseInt(strings.Join(dataBinSlice[i:i+8], ""), 2, 64)
				dataHexSlice = append(dataHexSlice, fmt.Sprintf("%02x", value))
			}
			output[deviceIndex] = strings.Join(dataHexSlice, "")
			if DEBUGABLE {
				fmt.Println("print device", deviceIndex)
				fmt.Printf("data: %#v\n", dataBinSlice[:16])
				fmt.Printf("linebuffer: %#v\n", output[deviceIndex][:8])
			}
		}

		if IMAGABLE {
			outputFile, _ := os.Create(fmt.Sprintf("output/%06d.%03d.png", *imageIndex, cycleIndex))
			png.Encode(outputFile, img)
			outputFile.Close()
			*imageIndex = *imageIndex + 1
		}
	}
	return output, count
}

func RawPos(
	posx int,
	posy int,
) (string, string) {
	x := offsetX - geometry.Mm(posx)
	y := offsetY - geometry.Mm(posy)
	if DEBUGABLE {
		fmt.Println("move to", x, y)
	}
	return fmt.Sprintf("%.6f", x), fmt.Sprintf("%.6f", y)
}

func MostLeftSpot(cycleIndex int, spots []*slide.Spot) *slide.Spot {
	var target *slide.Spot
	for _, spot := range spots {
		if spot.Reagents[cycleIndex].Printed {
			continue
		}
		if target == nil {
			target = spot
		} else {
			if spot.Pos.AtLeft(target.Pos) {
				target = spot
			}
		}
	}
	return target
}

func getStep(stepInput *widgets.QLineEdit) int {
	result, err := strconv.Atoi(stepInput.Text())
	if err != nil {
		stepInput.SetText("1")
		return 1
	}
	return result
}
