package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	_CircleCommandIndex = 0        //internal index
	_CircleCommand      = "CIRCLE" //Текст с которого начинается описание круга
	_CirclePrefixCoord  = "center" //Текст с которого начинаются координаты у круга

	_PolylineCommandIndex = 1          //internal index
	_PolylineCommand      = "POLYLINE" //Текст с которого начинается описание полилинии ВНИМАНИЕ: "LWPOLYLINE" тоже содержит "POLYLINE"
	_PolylinePrefixCoord  = "at "      //Текст с которого начинаются координаты полилинии

)

//TCommands - define log record
type TCommands struct {
	i         int
	keyString string
	pntString string
}

var _Cmd [2]TCommands

//init - initialize function
//fill array _Cmd
//at index 0 - keep parameters of Circle
//at index 1 - keep parameters of Polyline
//at index 2 - keep parameters of LwPolyline
func init() {
	//Circle log record
	_Cmd[0].i = _CircleCommandIndex
	_Cmd[0].keyString = _CircleCommand
	_Cmd[0].pntString = _CirclePrefixCoord
	//Polyline log record
	_Cmd[1].i = _PolylineCommandIndex
	_Cmd[1].keyString = _PolylineCommand
	_Cmd[1].pntString = _PolylinePrefixCoord
}

func main() {
	var (
		index          int  //index of current command in _Cmd[]
		isNewCommand   bool //началась новая команда
		countPolylines int  //подсчёт полилиний, нужно чтобы писать "-999 -999 -999" в конце полилинии, а не в начале
		iFile          *os.File
		oFile          *os.File
		err            error
	)

	TestInputParams()

	fmt.Println("input  file: ", os.Args[1])
	fmt.Println("output file: ", os.Args[2])

	iFile, err = os.Open(os.Args[1]) //Open file to READ
	if err != nil {
		fmt.Fprint(os.Stderr, "file: "+os.Args[1]+" can't open to read")
		os.Exit(2)
	}

	oFile, err = os.Create(os.Args[2]) //Open file to WRITE
	if err != nil {
		fmt.Fprint(os.Stderr, "file: "+os.Args[2]+" can't open to write")
		os.Exit(3)
	}

	fmt.Println("input  file: ", iFile)
	fmt.Println("output file: ", oFile)

	iScanner := bufio.NewScanner(iFile)
	for i := 0; iScanner.Scan(); i++ {
		s := iScanner.Text()
		//Поиск начала объекта
		if strings.Contains(s, _CircleCommand) {
			isNewCommand = true //найден круг
			index = _CircleCommandIndex
		}
		if strings.Contains(s, _PolylineCommand) {
			isNewCommand = true //найдена новая полилиния
			index = _PolylineCommandIndex
			//запишем конец предыдущей полилинии (если это не первая)
			//если countPolylines == 0 то это самая первая полилиния в файле и писать "концовку" не надо
			if countPolylines > 0 {
				fmt.Fprintf(oFile, "%s\n", "-999 -999 -999")
			}
			countPolylines++
		}

		//Поиск строк координат и их разбор
		if strings.Contains(s, _Cmd[index].pntString) { //находим строку с координатами
			////			if (index == _PolylineCommandIndex) || (index == _LwPolylineCommandIndex) {
			//находим "=" c конца, от этой позиции до конца это Z
			//выделяем в новую строку всё от начала до этой позиции
			//теперь здесь только X и Y
			i := strings.LastIndex(s, "=") + 1
			//выдёргиваем Z
			Z, errZ := strconv.ParseFloat(strings.TrimSpace(s[i:]), 64)
			if errZ != nil {
				fmt.Printf("strvconv not pass: '%s'\n", s[i:])
			}

			s = s[1:(i - 2)] //удаляем найденное из строки, на 2 символа влево - там стоит 'Z'
			//получим: 'at point  X=   7.7218  Y=  11.7846  '

			i = strings.LastIndex(s, "=") + 1
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
			fmt.Fprintf(oFile, "%f   %f   %f\n", X, Y, Z)
			////			}
		}
	}
	fmt.Println(isNewCommand, index)
}

//TestInputParams - test input params
func TestInputParams() {
	if len(os.Args) < 3 {
		fmt.Fprint(os.Stderr, "using:> alg c:\\log\\inputfile.log v:\\dat\\outputfile.xyz \n")
		os.Exit(1)
	}

}
