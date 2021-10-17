/* 注意：一个Game创建时，game id 和 2 pid 是确定的！不管断线还是任何情况，这个数据都永远不变。*/

package main

import (
    "fmt"
    "math/rand"
    "os"
)

type PidResultAgree struct {
    pid Pid
    result Result
    agree bool
}
/* collect two pid's agree or not of the game's result */
func (gs GameService) answerResult(pid Pid,result Result,agree bool) {
    gs.ch_pid_result_agree <- PidResultAgree{pid,result,agree}
}

//add hand service
type HandParam struct {
    pid Pid
    hand Hand
}
                                                                 
/* hand */
func (gs *GameService) hand(pid Pid,hand Hand) {
    gs.ch_hand <- HandParam{pid,hand}
}

//time tick
type TickParam struct {
    pid Pid
    time Time
}

/* tick */
func (gs *GameService) tick(pid Pid,time Time) {
    gs.ch_tick <- TickParam{pid,time}
}

/* linebroken */
func (gs *GameService) linebroken(pid Pid){
    gs.ch_linebroken <- pid
}

/* comeback */
func (gs *GameService) comeback(pid Pid){
    gs.ch_comeback <- pid
}

/* join */
func (gs *GameService) join(player Player){
    gs.ch_join <- player
}

/* unjoin */
func (gs *GameService) unjoin(pid Pid){
    gs.ch_unjoin <- pid
}

/* set proto */
func (gs *GameService) setProto(proto Proto) {
    gs.ch_proto <- proto
}

/* set result */
func (gs *GameService) setResult(result Result) {
    gs.ch_result <- result
}

/* all things about game. when someone join the game or reload,he should know the game data */
func (gs *GameService) getGameData() *MsgGameData {
    //get game data for pid
    ch_result := make(chan *MsgGameData)
    gs.ch_game_data <- ch_result
    return <- ch_result
}

type OtherPidsParam struct {
    pid Pid
    ch_result chan []Pid
}

/* one pid + all watchers */
func (gs *GameService) getOtherPids(pid Pid) []Pid {
    ch_result := make(chan []Pid,3)
    gs.ch_other_pids <- OtherPidsParam{pid,ch_result}
    return <- ch_result
}

/* when we start a game,we must have a game id and two pids */
func (gs *GameService) start(gid Gid,players [2]Player) {
    gs.gid = gid
    gs.players = players
    gs.status = Game_Preparing
    gs.hands = []Hand{}
    gs.watchers = []Player{}
    
    gs.ch_linebroken = make(chan Pid,3)
    gs.ch_comeback = make(chan Pid,3)
    gs.ch_hand = make(chan HandParam,3)
    gs.ch_tick = make(chan TickParam,3)
    gs.ch_join = make(chan Player,3)
    gs.ch_unjoin = make(chan Pid,3)
    gs.ch_game_data = make(chan chan *MsgGameData,3)
    gs.ch_other_pids = make(chan OtherPidsParam,3)
    gs.ch_proto = make(chan Proto,3)
    gs.ch_result = make(chan Result,3)
    gs.ch_pid_result_agree = make(chan PidResultAgree,3)
    
    go func(){
        
        var my_pras []PidResultAgree
        
        for {
            select {
                
            case pra := <- gs.ch_pid_result_agree:
                //pid agree the result?
                my_pras = append(my_pras,pra)
                if len(my_pras) == 2 {
                    if my_pras[0].agree && my_pras[1].agree {
                        /* you two agree the result,now game over */
                        gs.set_result(my_pras[0].result)
                    }else{
                        my_pras = []PidResultAgree{}
                        gs.game_status_changed(Game_Running)   
                    }
                }
                
            case param := <- gs.ch_hand:
                //add hand
                gs.hands = append(gs.hands,param.hand)
                //
                for _,pid := range gs.otherPids(param.pid) {
                    if myshadow,ok := shadows_service.Get(pid); ok {
                        myshadow.outMsgHand(gs.gid,param.hand)
                    }
                }
                
            case param := <- gs.ch_tick:
                //time tick
                //update time
                if i,ok := gs.indexOf(param.pid); ok {
                    gs.times[i] = param.time
                }
                //
                for _,pid := range gs.otherPids(param.pid) {
                    if myshadow,ok := shadows_service.Get(pid); ok {
                        myshadow.outMsgTick(gs.gid,param.time)
                    }
                }
                
            case pid := <- gs.ch_linebroken:
                //linebroken
                if i,ok := gs.indexOf(pid); ok {
                    gs.linebrokens[i] = true
                }
                //tell everybody ,game paused!
                gs.game_status_changed(Game_Paused)
                
            case pid := <- gs.ch_comeback:
                //comeback
                if i,ok := gs.indexOf(pid); ok {
                    gs.linebrokens[i] = false
                    //have a look,if can restart
                    if gs.linebrokens[0]==false && gs.linebrokens[1]==false {
                        //when every things ok,just tell all players 'we started again'
                        gs.game_status_changed(Game_Running)
                    }
                }
                
            case player := <- gs.ch_join:
                //first tell every client,we should add one pid in watchers
                for _,pid := range gs.allPids() {
                    if myshadow,ok := shadows_service.Get(pid); ok {
                        myshadow.playerJoin(gs.gid,player)
                    }
                }
                //then,i change too
                gs.watchers = append(gs.watchers,player)
                
            case pid := <- gs.ch_unjoin:
                //first,remove the pid from watchers
                var index int = -1
                for i,player := range gs.watchers {
                    if player.Pid == pid {
                        index = i
                        break
                    }
                }
                if index != -1 {
                    gs.watchers = append(gs.watchers[:index],gs.watchers[index+1:]...)
                }
                //tell the left players , someone gone
                for _,mypid := range gs.allPids() {
                    if myshadow,ok := shadows_service.Get(mypid); ok { 
                        myshadow.playerUnjoin(gs.gid,pid)
                    }
                }
                
            case ch_result := <- gs.ch_game_data:
                //get game data
                ch_result <- gs.get_game_data()
                
            case param:= <- gs.ch_other_pids:
                param.ch_result <- gs.otherPids(param.pid)
                
            case proto := <- gs.ch_proto:
                //set proto,means game started
                gs.proto = proto
                gs.implement_proto()
                /* tell every pid the proto */
                for _,pid := range gs.allPids() {
                    if myshadow,ok := shadows_service.Get(pid); ok { 
                        myshadow.setProto(gs.gid,proto)
                    }
                }
                /* tell game status change too! */
                gs.game_status_changed(Game_Running)
                
            case result := <- gs.ch_result:
                //set result,means game over!
                gs.set_result(result)
            } //select end
        }
    }()
}

//---------------------------------------------------------

type GameService struct {
    gid Gid
    players [2]Player
    times [2]Time
    linebrokens [2]bool
    hands []Hand
    watchers []Player
    proto Proto
    result Result
    status byte
    firstIndex int // 0 or 1
    
    //chans
    ch_linebroken chan Pid
    ch_comeback chan Pid
    ch_join chan Player
    ch_unjoin chan Pid
    ch_hand chan HandParam
    ch_tick chan TickParam
    ch_game_data chan chan *MsgGameData
    ch_other_pids chan OtherPidsParam
    ch_proto chan Proto
    ch_result chan Result
    ch_pid_result_agree chan PidResultAgree
}

//---------------------------------------------------------

func (gs GameService) isLinebroken() bool {
    return gs.linebrokens[0] || gs.linebrokens[1]
}

func (gs GameService) hasPid(pid Pid) bool {
    return gs.players[0].Pid == pid || gs.players[1].Pid == pid
}

/* does the pid in game's players array? */
func (gs GameService) indexOf(pid Pid) (index int,ok bool) {
    for i,player := range gs.players {
        if player.Pid == pid {
            index = i
            ok = true
        }
    }
    return
}

func (gs GameService) otherPid(pid Pid) (other Pid) {
    if gs.players[0].Pid == pid {
        other = gs.players[1].Pid
    }else{
        other = gs.players[0].Pid
    }
    return
}

/* 1 pid + all watchers */
func (gs GameService) otherPids(pid Pid) []Pid {
    mypids := []Pid{gs.otherPid(pid)}
    for _,player := range gs.watchers {
        mypids = append(mypids,player.Pid)
    }
    return mypids
}

func (gs GameService) allPids() []Pid {
    mypids := []Pid{gs.players[0].Pid,gs.players[1].Pid}
    for _,player := range gs.watchers {
        mypids = append(mypids,player.Pid)
    }
    return mypids                                  
}

/* proto like a class,when we use the proto,we should create objects */
func (gs GameService) implement_proto() {
    /* who first should implement,we should know firstIndex */
    if gs.proto.WhoFirst == Random {
        gs.firstIndex = rand.Intn(2)
    }else{
        if index,ok := gs.indexOf(gs.proto.FirstPlayer); ok {
            gs.firstIndex = index
        }else{
            fmt.Println("error: proto's FirstPlayer is not in game's players")
            os.Exit(1)
        }
    }
    /* time should implement */
    gs.times[0] = gs.proto.Time
    gs.times[1] = gs.proto.Time
}

/* players agree to end,we should compute the stones,make a result */
func (gs GameService) compute_result() (result Result) {
    //@todo
    return
}

/* when we compute a result of Counting(EndType),we should ask players if agree */
func (gs GameService) ask_result(result Result) {
    for _,player := range gs.players {
        if myshadow,ok := shadows_service.Get(player.Pid); ok { 
            myshadow.askResult(gs.gid,result)
        }
    }
}

/* here,game over */
func (gs GameService) set_result(result Result) {
    gs.result = result
    /* tell every pid the result */
    for _,pid := range gs.allPids() {
        if myshadow,ok := shadows_service.Get(pid); ok { 
            myshadow.setResult(gs.gid,result)
        }
    }
    /* tell game status change too! */
    gs.game_status_changed(Game_Stopped)
}

func (gs GameService) getColor(pid Pid) (color Color) {
    if gs.proto.Handicap == 0 {
        /* fristIndex is White*/
        if gs.players[gs.firstIndex].Pid == pid {
            color = White
        }else{
            color = Black
        }
    }else{
        /* fristIndex is Black*/
        if gs.players[gs.firstIndex].Pid == pid {
            color = Black
        }else{
            color = White
        }
    }
    return
}

func (gs GameService) get_game_data() *MsgGameData {
    game_data := &MsgGameData{}
    game_data.Gid = gs.gid
    game_data.Players = gs.players
    game_data.Times = gs.times
    game_data.Hands = gs.hands
    game_data.Watchers = gs.watchers
    game_data.Proto = gs.proto
    game_data.FirstIndex = gs.firstIndex
    return game_data
}

//---------------------------------------------------------
//game statuses
const Game_Preparing byte = 0 //preparing
const Game_Running byte = 1 //running
const Game_Stopped byte = 2 //stopped
const Game_Paused byte = 3 //paused

/* tell every pid in the game,the change status changed */
func (gs *GameService) game_status_changed(status byte) {
    gs.status = status
    for _,pid := range gs.allPids() {
        if myshadow,ok := shadows_service.Get(pid); ok {
            myshadow.gameStatusChanged(gs.gid,gs.status)
        }
    }
}
//---------------------------------------------------------