package main

import (
	"net"
	"os"
	"fmt"
	"encoding/json"
)

func (shadow *Shadow) Connected(conn net.Conn){
	fmt.Println("server connected")
	shadow.conn = conn
}
                                
func (shadow Shadow) Closed(conn net.Conn){
	fmt.Println("server closed")
}

func (shadow Shadow) Exception(conn net.Conn,err error){
	fmt.Println("server exception: ",err)
	//@todo struct zero compare?
	if shadow.player != nil {
	    //tell linebroken_server
	    linebroken_service.linebroken(shadow.player.Pid)
	}
}

func (shadow Shadow) Process(conn net.Conn,msgBytes []byte) {
	fmt.Println("---server msg process---")
	var err error
	jsonMsg := &JsonMsg{}
	JsonDecode(msgBytes,jsonMsg)
	
	fmt.Println("got json msg: ",jsonMsg)
	
    switch jsonMsg.Head {
    case msg_ans_result_head:
        msgAnsResult := &MsgAnsResult{}
        JsonDecode(jsonMsg.Body,msgAnsResult)
        gameService := games_service.Get(msgAnsResult.Gid)
        gameService.answerResult(shadow.player.Pid,msgAnsResult.Result,msgAnsResult.Agree)
        
    case msg_ask_end_head:
        msgAskEnd := &MsgAskEnd{}
        JsonDecode(jsonMsg.Body,msgAskEnd)
        gameService := games_service.Get(msgAskEnd.Gid)
        _shadow := shadows_service.Get(gameService.otherPlayer(shadow.player.Pid))
        WriteMsg(_shadow.conn,msgBytes)
        
    case msg_ans_end_head:
        msgAnsEnd := &MsgAnsEnd{}
        JsonDecode(jsonMsg.Body,msgAnsEnd)
        gameService := games_service.Get(msgAnsEnd.Gid)
        if msgAnsEnd.Agree {
            /* we should compute stones now! */
            gameService.game_status_changed(Game_Paused)
            gameService.ask_result(gameService.compute_result())
        }else{
            _shadow := shadows_service.Get(gameService.otherPlayer(shadow.player.Pid))
            WriteMsg(_shadow.conn,msgBytes)
        }
        
    case msg_invite_head:
        msgInvite := &MsgInvitePlayer{}
        JsonDecode(jsonMsg.Body,msgInvite)
        _shadow := shadows_service.Get(msgInvite.Pid)
        msgInvite1 := &MsgInvitePlayer{_shadow.player.Pid}
        msg := &JsonMsg{msg_invite_head,JsonEncode(msgInvite1)}
        WriteJsonMsg(shadow.conn,msg)
        
    case msg_ask_proto_head:
        //msg ask proto
        p := &MsgAskProto{}
        JsonDecode(jsonMsg.Body,p)
        shadow.sendOtherPlayer(p.Gid,msgBytes)
        
    case msg_ans_proto_head:
        //msg answer proto
        msgAnsProto := &MsgAnsProto{}
        JsonDecode(jsonMsg.Body,msgAnsProto)
        //decode ok
        if msgAnsProto.Agree {
            //yes
            /* here,triger the action: Game_Preparing -> Game_Running */
            games_service.Get(msgAnsProto.Gid).setProto(msgAnsProto.Proto)
        }else {
            //no
            shadow.sendOtherPlayer(msgAnsProto.Gid,msgBytes)
        }
	}
}

//---------------------------------------------------------
//msg_hand
func (shadow *Shadow) inMsgHand(ms MsgHand) {
    games_service.Get(ms.Gid).hand(shadow.player.Pid,ms.Hand)
}

func (shadow *Shadow) outMsgHand(gid Gid,hand Hand) {
    msgHand := &MsgHand{gid,hand}
    msg := &JsonMsg{msg_hand_head,JsonEncode(msgHand)}
    WriteJsonMsg(shadow.conn,msg)
}
//---------------------------------------------------------
//msg_tick
func (shadow *Shadow) inMsgTick(msgTick MsgTick) {
    games_service.Get(msgTick.Gid).tick(shadow.player.Pid,msgTick.Time)
}

func (shadow *Shadow) outMsgTick(gid Gid,time Time) {
    msgTick := &MsgTick{gid,time}
    msg := &JsonMsg{msg_tick_head,JsonEncode(msgTick)}
    WriteJsonMsg(shadow.conn,msg)
}
//---------------------------------------------------------
func (shadow Shadow) inInvite(pid Pid) {
    shadows_service.Get(pid).outInvite(shadow.player.Pid)
}

func (shadow Shadow) outInvite(pid Pid) {
    msgInvite := &MsgInvitePlayer{pid}
    msg := &JsonMsg{msg_invite_head,JsonEncode(msgInvite)}
    WriteJsonMsg(shadow.conn,msg)
}
//the game result ,are you agree?
func (shadow Shadow) askResult(gid Gid,result Result){
    msgAskResult := &MsgAskResult{gid,result}
    msg := &JsonMsg{msg_ask_result_head,JsonEncode(msgAskResult)}
    WriteJsonMsg(shadow.conn,msg)
}

func (shadow Shadow) playerJoin(gid Gid,pid Pid) {
}

func (shadow Shadow) playerUnjoin(gid Gid,pid Pid) {
}
//---------------------------------------------------------
//game status changed
func (shadow *Shadow) gameStatusChanged(gid Gid,status byte) {
    msgBytes := JsonEncode(&MsgGameStatusChanged{gid,status})
    msg := &JsonMsg{msg_game_status_changed_head,msgBytes}
    WriteJsonMsg(shadow.conn,msg)
}
//---------------------------------------------------------
//msg_game_data
func (shadow *Shadow) getGameData(msgGetGameData MsgGetGameData) {
    gameService := games_service.Get(msgGetGameData.Gid)
    game_data := gameService.getGameData()
    msg := &JsonMsg{msg_game_data_head,JsonEncode(game_data)}
    WriteJsonMsg(shadow.conn,msg)
}
//---------------------------------------------------------
/* send msg to other player of game */
func (shadow *Shadow) sendOtherPlayer(gid Gid,msgBytes []byte) {
    gameService := games_service.Get(gid)
    other := gameService.otherPlayer(shadow.player.Pid)
    _shadow := shadows_service.Get(other)
    WriteMsg(_shadow.conn,msgBytes)
}
//---------------------------------------------------------
/* when we agree a proto,we should tell every player the proto by the game_service */
func (shadow Shadow) setProto(gid Gid,proto Proto) {
    msgSetProto := &MsgSetProto{gid,proto}
    msg := &JsonMsg{msg_game_data_head,JsonEncode(msgSetProto)}
    WriteJsonMsg(shadow.conn,msg)
}

func (shadow Shadow) setResult(gid Gid,result Result) {
    msgSetResult := &MsgSetResult{gid,result}
    msg := &JsonMsg{msg_game_data_head,JsonEncode(msgSetResult)}
    WriteJsonMsg(shadow.conn,msg)
}
//---------------------------------------------------------
type JsonMsg struct {
	Head byte `json:"head"`
	Body []byte `json:"body"`
}

type Shadow struct {
    conn net.Conn
    player *Player
}
//---------------------------------------------------------
func WriteJsonMsg(conn net.Conn,msg *JsonMsg){
	WriteMsg(conn,JsonEncode(msg))
}

func WriteMsg(conn net.Conn,msgBytes []byte) {
	_,err := conn.Write(AddHeader(msgBytes))
	if err != nil {
		fmt.Println("conn write error: ",err)
		os.Exit(1)
	}
}
//---------------------------------------------------------
func JsonEncode(v interface{}) []byte {
	msgBytes,err := json.Marshal(v)
	CheckJsonError(err)
	return msgBytes
}

func JsonDecode(data []byte,v interface{}) {
    err := json.Unmarshal(data,v)
    CheckJsonError(err)
}

func CheckJsonError(err error) {
	if err != nil {
		fmt.Println("json marshal or unmarshal error: ",err)
		os.Exit(1)
	}
}
//---------------------------------------------------------
//file end