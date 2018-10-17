package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	_CircleCommandIndex = 0
	_CircleCommand      = " CIRCLE "      //Текст с которого начинается описание круга
	_CirclePrefixCoord  = "center point," //Текст с которого начинаются координаты у круга

	_PolylineCommandIndex = 1
	_PolylineCommand      = "POLYLINE" //Текст с которого начинается описание полилинии ВНИМАНИЕ: "LWPOLYLINE" тоже содержит "POLYLINE"
	_PolylinePrefixCoord  = "at "      //Текст с которого начинаются координаты полилинии

)

//TCommands - define log record
type TCommands struct {
	i         int
	keyString string
	pntString string
}

var cmd [2]TCommands

//init - initialize function
//fill array cmd
//at index 0 - keep parameters of Circle
//at index 1 - keep parameters of Polyline
//at index 2 - keep parameters of LwPolyline
func init() {
	//Circle log record
	cmd[0].i = _CircleCommandIndex
	cmd[0].keyString = _CircleCommand
	cmd[0].pntString = _CirclePrefixCoord
	//Polyline log record
	cmd[1].i = _PolylineCommandIndex
	cmd[1].keyString = _PolylineCommand
	cmd[1].pntString = _PolylinePrefixCoord
}

func main() {
	var (
		inputType     int //0 - обрабатываем круги, 1 - обрабатываем полилинии
		index         int //index of current command in cmd[]
		countCommands int //подсчёт найденных объектов: полилиний
		//нужно для полилиний, чтобы писать "-999 -999 -999" в конце полилинии
		outputFileName string
		inputFileName  string
		iFile          *os.File
		oFile          *os.File
		err            error
	)

	//exit if params not true
	if !TestInputParams() {
		os.Exit(1)
	}
	fmt.Fprint(os.Stderr, "params preview ok \n")

	//input type: circle -> 0, line -> 1
	if strings.Contains(os.Args[1], "/c") {
		fmt.Fprint(os.Stderr, "process CIRCLE \n")
		inputType = 0
		inputFileName = os.Args[2]
		if len(os.Args) == 3 {
			//if output file name not present in parameters
			//make output file name from input file name, change ext to "xyz"
			outputFileName = strings.TrimSuffix(os.Args[2], filepath.Ext(os.Args[2])) + ".xyz"
		} else {
			outputFileName = os.Args[3]
		}
	} else {
		fmt.Fprint(os.Stderr, "process POLYLINE \n")
		inputType = 1
		inputFileName = os.Args[1]
		if len(os.Args) == 2 {
			//if output file name not present in parameters
			//make output file name from input file name, change ext to "xyz"
			outputFileName = strings.TrimSuffix(os.Args[1], filepath.Ext(os.Args[1])) + ".xyz"
		} else {
			outputFileName = os.Args[2]
		}
	}

	fmt.Println("input  file: ", inputFileName)

	iFile, err = os.Open(inputFileName) //Open file to READ
	if err != nil {
		fmt.Fprint(os.Stderr, "file: "+os.Args[1]+" can't open to read")
		os.Exit(2)
	}
	fmt.Fprint(os.Stderr, "opened \n")

	fmt.Println("output file: ", outputFileName)
	oFile, err = os.Create(outputFileName) //Open file to WRITE
	if err != nil {
		fmt.Fprint(os.Stderr, "file: "+outputFileName+" can't open to write")
		os.Exit(3)
	}
	fmt.Fprint(os.Stderr, "opened \n")

	iScanner := bufio.NewScanner(iFile)
	for i := 0; iScanner.Scan(); i++ {
		s := iScanner.Text()
		//Поиск начала объекта
		if inputType == 1 {
			//при разборе линии надо запомнить что линия началась, когда начнётся
			//новая линия запишем -999
			if strings.Contains(s, _PolylineCommand) {
				index = _PolylineCommandIndex
				//запишем конец полилинии (если это не первая)
				if countCommands > 0 {
					fmt.Fprintf(oFile, "%s\n", "-999 -999 -999")
				}
				countCommands++
			} else {
				//пишем в консоль, чтобы видно было, что программа работает
				fmt.Fprint(os.Stderr, ".")
				//Поиск строк координат и их разбор
				if strings.Contains(s, cmd[index].pntString) {
					//строка с координатами найдена
					ProcessLine(oFile, s)
				}
			}
		} else { //обработка кругов
			//пишем в консоль, чтобы видно было, что программа работает
			fmt.Fprint(os.Stderr, ".")
			//при обработке кругов координата начинается с "center point,"
			if strings.Contains(s, _CirclePrefixCoord) {
				//строка с координатами найдена
				ProcessLine(oFile, s)
			}
		}
	}
	if countCommands > 0 {
		fmt.Fprintf(oFile, "%s\n", "-999 -999 -999")
	}

	fmt.Fprint(os.Stderr, "\n")
	fmt.Fprint(os.Stderr, "done: ", countCommands, " objects\n")
}

//ProcessLine Обработка строки "s" с координатами x= 111 y=111 z =111
//результат записываем в файл "f"
func ProcessLine(f *os.File, s string) {
	i := strings.LastIndex(s, "=") + 1
	//находим "=" c конца, от этой позиции "i" до конца это Z
	//выдёргиваем Z
	Z, errZ := strconv.ParseFloat(strings.TrimSpace(s[i:]), 64)
	if errZ != nil {
		fmt.Printf("strvconv not pass: '%s'\n", s[i:])
	}

	//выделяем в новую строку всё от начала до позиции "i"
	//теперь здесь только X и Y
	s = s[1:(i - 2)] //удаляем найденное из строки, на 2 символа влево - там стоит 'Z'
	//получим: 'at point  X=   7.7218  Y=  11.7846  '

	i = strings.LastIndex(s, "=") + 1 //находим позицию Y=
	Y, errY := strconv.ParseFloat(strings.TrimSpace(s[i:]), 64)
	if errY != nil {
		fmt.Printf("strvconv not pass: '%s'\n", s[i:])
	}

	s = s[1:(i - 2)] //удаляем найденное из строки,  на 2 символа влево - там стоит 'Y'
	//получим: 'at point  X=   7.7218  '

	i = strings.LastIndex(s, "=") + 1
	X, errX := strconv.ParseFloat(strings.TrimSpace(s[i:]), 64)
	if errX != nil {
		fmt.Printf("strvconv not pass: '%s'\n", s[i:])
	}
	fmt.Fprintf(f, "%f   %f   %f\n", X, Y, Z)

}

//TestInputParams - test input params
//exit if params not true
func TestInputParams() bool {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "using:> alg c:\\log\\i.log v:\\dat\\o.xyz \n")
		fmt.Fprint(os.Stderr, "for polyline \n\n")
		fmt.Fprint(os.Stderr, "using:> alg /c c:\\log\\i.log v:\\dat\\o.xyz \n")
		fmt.Fprint(os.Stderr, "for circle \n\n")
		return false
	}
	return true
}
