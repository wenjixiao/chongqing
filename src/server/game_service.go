/* 注意：一个Game创建时，game id 和 2 pid 是确定的！不管断线还是任何情况，这个数据都永远不变。*/

package main

import (
    "fmt"
    "math/rand"
    "os"
)

type PlayerResultAgree struct {
    player Pid
    result Result
    agree bool
}
/* collect two pid's agree or not of the game's result */
func (gs GameService) answerResult(player Pid,result Result,agree bool) {
    gs.ch_player_result_agree <- PlayerResultAgree{player,result,agree}
}

//add hand service
type HandParam struct {
    player Pid
    hand Hand
}
                                                                 
/* hand */
func (gs *GameService) hand(player Pid,hand Hand) {
    gs.ch_hand <- HandParam{player,hand}
}

//time tick
type TickParam struct {
    player Pid
    time Time
}

/* tick */
func (gs *GameService) tick(player Pid,time Time) {
    gs.ch_tick <- TickParam{player,time}
}

/* linebroken */
func (gs *GameService) linebroken(player Pid){
    gs.ch_linebroken <- player
}

/* comeback */
func (gs *GameService) comeback(player Pid){
    gs.ch_comeback <- player
}

/* join */
func (gs *GameService) join(player Pid){
    gs.ch_join <- player
}

/* unjoin */
func (gs *GameService) unjoin(player Pid){
    gs.ch_unjoin <- player
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

type OtherPlayersParam struct {
    player Pid
    ch_result chan []Pid
}

/* one pid + all watchers */
func (gs *GameService) getOtherPlayers(player Pid) []Pid {
    ch_result := make(chan []Pid,3)
    gs.ch_other_players <- OtherPlayersParam{player,ch_result}
    return <- ch_result
}

/* when we start a game,we must have a game id and two pids */
func (gs *GameService) start(gid Gid,players [2]Pid) {
    gs.gid = gid
    gs.players = players
    gs.status = Game_Preparing
    gs.hands = []Hand{}
    gs.watchers = []Pid{}
    
    gs.ch_linebroken = make(chan Pid,3)
    gs.ch_comeback = make(chan Pid,3)
    gs.ch_hand = make(chan HandParam,3)
    gs.ch_tick = make(chan TickParam,3)
    gs.ch_join = make(chan Pid,3)
    gs.ch_unjoin = make(chan Pid,3)
    gs.ch_game_data = make(chan chan *MsgGameData,3)
    gs.ch_other_players = make(chan OtherPlayersParam,3)
    gs.ch_proto = make(chan Proto,3)
    gs.ch_result = make(chan Result,3)
    gs.ch_player_result_agree = make(chan PlayerResultAgree,3)
    
    go func(){
        
        var my_pras []PlayerResultAgree
        
        for {
            select {
                
            case pra := <- gs.ch_player_result_agree:
                //pid agree the result?
                my_pras = append(my_pras,pra)
                if len(my_pras) == 2 {
                    if my_pras[0].agree && my_pras[1].agree {
                        /* you two agree the result,now game over */
                        gs.set_result(my_pras[0].result)
                    }else{
                        my_pras = []PlayerResultAgree{}
                        gs.game_status_changed(Game_Running)   
                    }
                }
                
            case param := <- gs.ch_hand:
                //add hand
                gs.hands = append(gs.hands,param.hand)
                //
                for _,player := range gs.otherPlayers(param.player) {
                    shadows_service.Get(player).outMsgHand(gs.gid,param.hand)
                }
                
            case param := <- gs.ch_tick:
                //time tick
                //update time
                if i,ok := gs.indexOf(param.player); ok {
                    gs.times[i] = param.time
                }
                //
                for _,player := range gs.otherPlayers(param.player) {
                    shadows_service.Get(player).outMsgTick(gs.gid,param.time)
                }
                
            case player := <- gs.ch_linebroken:
                //linebroken
                if i,ok := gs.indexOf(player); ok {
                    gs.linebrokens[i] = true
                }
                //tell everybody ,game paused!
                gs.game_status_changed(Game_Paused)
                
            case player := <- gs.ch_comeback:
                //comeback
                if i,ok := gs.indexOf(player); ok {
                    gs.linebrokens[i] = false
                    //have a look,if can restart
                    if gs.linebrokens[0]==false && gs.linebrokens[1]==false {
                        //when every things ok,just tell all players 'we started again'
                        gs.game_status_changed(Game_Running)
                    }
                }
                
            case player := <- gs.ch_join:
                //first tell every client,we should add one pid in watchers
                for _,p := range gs.allPlayers() {
                    shadows_service.Get(p).playerJoin(gs.gid,p)
                }
                //then,i change too
                gs.watchers = append(gs.watchers,player)
                
            case player := <- gs.ch_unjoin:
                //first,remove the pid from watchers
                var index int = -1
                for i,p := range gs.watchers {
                    if p == player {
                        index = i
                        break
                    }
                }
                if index != -1 {
                    gs.watchers = append(gs.watchers[:index],gs.watchers[index+1:]...)
                }
                //tell the left players , someone gone
                for _,p := range gs.allPlayers() {
                    shadows_service.Get(p).playerUnjoin(gs.gid,p)
                }
                
            case ch_result := <- gs.ch_game_data:
                //get game data
                ch_result <- gs.get_game_data()
                
            case param:= <- gs.ch_other_players:
                param.ch_result <- gs.otherPlayers(param.player)
                
            case proto := <- gs.ch_proto:
                //set proto,means game started
                gs.proto = proto
                gs.implement_proto()
                /* tell every pid the proto */
                for _,player := range gs.allPlayers() {
                    shadows_service.Get(player).setProto(gs.gid,proto)
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
    players [2]Pid
    times [2]Time
    linebrokens [2]bool
    hands []Hand
    watchers []Pid
    proto Proto
    result Result
    status byte
    firstIndex int // 0 or 1
    
    //chans
    ch_linebroken chan Pid
    ch_comeback chan Pid
    ch_join chan Pid
    ch_unjoin chan Pid
    ch_hand chan HandParam
    ch_tick chan TickParam
    ch_game_data chan chan *MsgGameData
    ch_other_players chan OtherPlayersParam
    ch_proto chan Proto
    ch_result chan Result
    ch_player_result_agree chan PlayerResultAgree
}

//---------------------------------------------------------

func (gs GameService) isLinebroken() bool {
    return gs.linebrokens[0] || gs.linebrokens[1]
}

func (gs GameService) hasPlayer(player Pid) bool {
    return gs.players[0] == player || gs.players[1] == player
}

/* does the pid in game's players array? */
func (gs GameService) indexOf(player Pid) (index int,ok bool) {
    for i,p := range gs.players {
        if p == player {
            index = i
            ok = true
        }
    }
    return
}

func (gs GameService) otherPlayer(player Pid) (other Pid) {
    if gs.players[0] == player {
        other = gs.players[1]
    }else{
        other = gs.players[0]
    }
    return
}

/* 1 pid + all watchers */
func (gs GameService) otherPlayers(player Pid) []Pid {
    other := gs.otherPlayer(player)
    return append([]Pid{other},gs.watchers...)
}

func (gs GameService) allPlayers() (all []Pid) {
    all = append(all,gs.players[0],gs.players[1])
    all = append(all,gs.watchers...)
    return                                                           
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
        shadows_service.Get(player).askResult(gs.gid,result)
    }
}

/* here,game over */
func (gs GameService) set_result(result Result) {
    gs.result = result
    /* tell every pid the result */
    for _,player := range gs.allPlayers() {
        shadows_service.Get(player).setResult(gs.gid,result)
    }
    /* tell game status change too! */
    gs.game_status_changed(Game_Stopped)
}

func (gs GameService) getColor(player Pid) (color Color) {
    if gs.proto.Handicap == 0 {
        /* fristIndex is White*/
        if gs.players[gs.firstIndex] == player {
            color = White
        }else{
            color = Black
        }
    }else{
        /* fristIndex is Black*/
        if gs.players[gs.firstIndex] == player {
            color = Black
        }else{
            color = White
        }
    }
    return
}

func (gs GameService) get_game_data() *MsgGameData {
    pid2player := func(pid Pid) *Player {
        return shadows_service.Get(pid).player
    }

    game_data := &MsgGameData{}
    game_data.Gid = gs.gid
    /* game_data.Players */
    for index,pid := range gs.players {
        game_data.Players[index] = *pid2player(pid)
    }
    game_data.Times = gs.times
    game_data.Hands = gs.hands
    /* game_data.Watchers */
    for index,pid := range gs.watchers {
        game_data.Watchers[index] = *pid2player(pid)
    }
    game_data.Proto = gs.proto
    game_data.FirstIndex = gs.firstIndex
    return game_data
}

func pid2player(pid Pid) *Player {
    return shadows_service.Get(pid).player
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
    for _,player := range gs.allPlayers() {
        shadows_service.Get(player).gameStatusChanged(gs.gid,gs.status)
    }
}
//---------------------------------------------------------