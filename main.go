package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xLib"
)

const (
	_CircleCommandIndex = 0
	_CircleCommand      = " CIRCLE "      //Текст с которого начинается описание круга
	_CirclePrefixCoord  = "center point," //Текст с которого начинаются координаты у круга

	_PolylineCommandIndex = 1
	_PolylineCommand      = "POLYLINE" //Текст с которого начинается описание полилинии ВНИМАНИЕ: "LWPOLYLINE" тоже содержит "POLYLINE"
	_PolylinePrefixCoord  = "at "      //Текст с которого начинаются координаты полилинии

)

//init - initialize function
//fill array cmd
//at index 0 - keep parameters of Circle
//at index 1 - keep parameters of Polyline
func init() {
}

//parseParameters() - обработка параметров, возврат имён файлов и режима обработки
func parseParameters() (int, string, string) {
	inputType := 0
	inputFileName := ""
	outputFileName := ""
	if strings.Contains(os.Args[1], "/c") {
		inputType = 0
		inputFileName = os.Args[2]
		if len(os.Args) == 3 {
			//if output file name not present in parameters
			//make output file name from input file name, change ext to "xyz"
			outputFileName = xLib.ChangeFileExt(os.Args[2], ".xyz")
		} else {
			outputFileName = os.Args[3]
		}
	} else {
		inputType = 1
		inputFileName = os.Args[1]
		if len(os.Args) == 2 {
			//if output file name not present in parameters
			//make output file name from input file name, change ext to "xyz"
			outputFileName = xLib.ChangeFileExt(os.Args[1], ".xyz")
		} else {
			outputFileName = os.Args[2]
		}
	}
	return inputType, inputFileName, outputFileName
}

func main() {
	var (
		inputType int //0 - обрабатываем круги, 1 - обрабатываем полилинии
		//index         int //index of current command in cmd[]
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
	fmt.Print("params preview ok \n")

	inputType, inputFileName, outputFileName = parseParameters()
	if inputType == 0 {
		fmt.Print("process CIRCLE \n")
	} else {
		fmt.Print("process POLYLINE \n")
	}

	fmt.Print("input  file: '", inputFileName)
	iFile, err = os.Open(inputFileName) //Open file to READ
	if err != nil {
		fmt.Println("file: " + os.Args[1] + " can't open to read")
		os.Exit(2)
	}
	fmt.Println("' opened successfully")

	fmt.Print("output file: '", outputFileName)
	oFile, err = os.Create(outputFileName) //Open file to WRITE
	if err != nil {
		fmt.Println("file: " + outputFileName + " can't open to write")
		os.Exit(3)
	}
	fmt.Println("' opened successfully")

	iScanner := bufio.NewScanner(iFile)
	for i := 0; iScanner.Scan(); i++ {
		s := iScanner.Text()
		//Поиск начала объекта
		if inputType == 1 {
			//ищем начало полилинии, строка с новой полилинией содержит "POLYLINE"
			if strings.Contains(s, _PolylineCommand) {
				//запишем конец полилинии (если это не первая)
				if countCommands > 0 {
					fmt.Fprintf(oFile, "%s\n", "-999 -999 -999")
				}
				//при разборе линии надо запомнить что линия началась, когда начнётся
				//новая линия запишем в поток "-999 -999 -999"
				countCommands++
			} else {
				fmt.Print(".") //пишем в консоль, чтобы видно было, что программа работает
				//Поиск строк координат и их разбор
				if strings.Contains(s, _PolylinePrefixCoord) {
					//строка с координатами найдена
					ProcessLine(oFile, s)
				}
			}
		} else { //обработка кругов
			fmt.Print(".") //пишем в консоль, чтобы видно было, что программа работает
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

	fmt.Print("\n")
	fmt.Print("done: ", countCommands, " objects\n")
}

//ProcessLine Обработка строки "s" с координатами x= 111 y=111 z =111
//результат записываем в файл "f"
//--функцию необходимо переделать. ей не надо писать в консоль, надо возвращять ошибку
//--и большой вопрос надо ли ей писать в файл, может вернуть строку которую надо писать
//  вдруг захочу писать не в файл а в поток или в базу
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
		fmt.Print("using:> alg c:\\log\\i.log v:\\dat\\o.xyz \n")
		fmt.Print("for polyline \n\n")
		fmt.Print("using:> alg /c c:\\log\\i.log v:\\dat\\o.xyz \n")
		fmt.Print("for circle \n\n")
		return false
	}
	return true
}
