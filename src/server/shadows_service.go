package main

import (
    "fmt"
    )

//---------------------------------------------------------
//put shadow service
type param1 struct {
    pid Pid
    shadow *Shadow
}

func (ss ShadowsService) Put(pid Pid,shadow *Shadow) {
    ss.ch_service1 <- param1{pid,shadow}
}

//---------------------------------------------------------
//get shadow service
type result2 struct {
    shadow *Shadow
    ok bool
}

type param2 struct {
    pid Pid
    ch_result chan result2
}

func (ss ShadowsService) Get(pid Pid) (*Shadow,bool) {
    ch_result := make(chan result2)
    ss.ch_service2 <- param2{pid,ch_result}
    result := <- ch_result
    return result.shadow,result.ok
}

//---------------------------------------------------------
//list shadows service
func (ss ShadowsService) List() []*Shadow {
    ch_result := make(chan []*Shadow)
    ss.ch_service3 <- ch_result
    return <- ch_result
}

//---------------------------------------------------------
func (ss *ShadowsService) start() {
    fmt.Println("<shadows service> running...")
    ss.pid2shadow = make(map[Pid]*Shadow)
    
    ss.ch_service1 = make(chan param1,3)
    ss.ch_service2 = make(chan param2,3)
    ss.ch_service3 = make(chan chan []*Shadow,3)
    
    go func(){
        for {
            select {
            case param := <- ss.ch_service1:
                //Put
                ss.pid2shadow[param.pid] = param.shadow
                
            case param := <- ss.ch_service2:
                //Get
                shadow,ok := ss.pid2shadow[param.pid]
                param.ch_result <- result2{shadow,ok} 
                
            case ch_result := <- ss.ch_service3:
                //List
                var shadows []*Shadow
                for _,v := range ss.pid2shadow {
                    shadows = append(shadows,v)
                }
                ch_result <- shadows
            }
        }
    }()
}

type ShadowsService struct {
    pid2shadow map[Pid]*Shadow
    ch_service1 chan param1 //Put
    ch_service2 chan param2 //Get
    ch_service3 chan chan []*Shadow //List
}

//---------------------------------------------------------