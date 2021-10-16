package main
//=========================================================
const msg_test_head byte = 0
const msg_hand_head byte = 1
const msg_get_game_data_head byte = 2
const msg_game_data_head byte = 3
const msg_ask_proto_head byte = 4
const msg_ans_proto_head byte = 5
const msg_game_status_changed_head byte = 6
const msg_tick_head byte = 7
const msg_invite_head byte = 8
const msg_ask_end_head byte = 9
const msg_ask_result_head byte = 10
const msg_ans_end_head byte = 11
const msg_ans_result_head byte = 12
//=========================================================
type MsgHand struct {
    Gid Gid  `json:"gid"`
    Hand Hand `json:"hand"`
}

type MsgTick struct {
    Gid Gid `json:"gid"`
    Time Time `json:"time"`
}

type MsgAskProto struct {
    Gid Gid `json:"gid"`
    Proto Proto `json:"proto"`
}

type MsgAnsProto struct {
    Gid Gid `json:"gid"`
    Agree bool `json:"agree"`
    Proto Proto `json:"proto"`
}

type MsgGetGameData struct {
    Gid Gid `json:"gid"`
}

type MsgGameData struct {
    Gid Gid `json:"gid"`
    Players [2]Player `json:"players"`
    Times [2]Time `json:"times"`
    Hands []Hand `json:"hands"`
    Watchers []Player `json:"watchers"`
    Proto Proto `json:"proto"`
    FirstIndex int `json:"firstIndex"`
}

type MsgGameStatusChanged struct {
    Gid Gid `json:"gid"`
    Status byte `json:"status"`
}

type MsgInvitePlayer struct {
    Pid Pid `json:"pid"`
}

type MsgAskEnd struct {
    Gid Gid `json:"gid"`
}

type MsgAnsEnd struct {
    Gid Gid `json:"gid"`
    Agree bool `json:"agree"`
}

type MsgAskResult struct {
    Gid Gid `json:"gid"`
    Result Result
}

type MsgAnsResult struct {
    Gid Gid `json:"gid"`
    Agree bool `json:"agree"`
    Result Result `json:"result"`
}

type MsgSetProto struct {
    Gid Gid `json:"gid"`
    Proto Proto `json:"proto"`
}

type MsgSetResult struct {
    Gid Gid `json:"gid"`
    Result Result `json:"result"`
}
//=========================================================