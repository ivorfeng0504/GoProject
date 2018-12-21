package yqq


type YqqEntire struct{
	TypeGroup TypeGroup
	YqqHomeData HomeData
	YqqStat Stats
}

type TypeGroup struct{
	Operate []YqqRoomSimple `json:Operate,[]yqqRoomSimple`
	ShortTerm []YqqRoomSimple `json:ShortTerm,[]yqqRoomSimple`
	Technology []YqqRoomSimple `json:Technology,[]yqqRoomSimple`
	Value []YqqRoomSimple `json:Value,[]yqqRoomSimple`
}



