package main

const (
	_CircleCommandIndex = 0
	_CircleCommand      = " CIRCLE "      //Текст с которого начинается описание круга
	_CirclePrefixCoord  = "center point," //Текст с которого начинаются координаты у круга

	_PolylineCommandIndex = 1
	_PolylineCommand      = "POLYLINE" //Текст с которого начинается описание полилинии ВНИМАНИЕ: "LWPOLYLINE" тоже содержит "POLYLINE"
	_PolylinePrefixCoord  = "at "      //Текст с которого начинаются координаты полилинии

)
