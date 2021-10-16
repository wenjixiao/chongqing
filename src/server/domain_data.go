package main

/* basic type of the system */
type Pid string //player id
type Gid int //game id

type Color bool
const Black Color = true
const White Color = false

type WhoFirst bool
const Random WhoFirst = true
const Earmark WhoFirst = false

type Proto struct {
	Handicap byte //让子
	WhoFirst WhoFirst //谁先下
	Komi float32 //贴目
	FirstPlayer Pid 
	Time Time
}

type EndType byte 
const Counting EndType = 1
const NotCounting EndType = 2
const LineBroken EndType = 3
const Timeout EndType = 4

type Result struct {
	EndType EndType
	Winner Color
	Mount float32
}

type Player struct {
	Pid Pid 
	Level string
}

type Point struct {
    X,Y byte
}

type Time struct {
    BaoLiu int32
    DuMiao int32
    TimesYingJi int32
    EveryYingJi int32
}

type Hand struct {
    Seq int32
    Color Color
    Point Point
    Time Time
}