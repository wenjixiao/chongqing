package main

import (
    "fmt"
    "time"
    )

const Countdown int = 5*60 //seconds,5 minutes

type LinebrokenService struct {
    tick <-chan time.Time
    ch_linebroken chan Pid //player linebroken
    ch_check chan Pid //check player if linebroken
    pid2countdown map[Pid]int
}

func (ls *LinebrokenService) linebroken(pid Pid) {
    ls.ch_linebroken <- pid
}

func (ls *LinebrokenService) check(pid Pid) {
    ls.ch_check <- pid
}

func (ls *LinebrokenService) start() {
    fmt.Println("<linebroken manager> running...")
    ls.tick = time.Tick(1 * time.Second)
    
    ls.ch_linebroken = make(chan Pid,3)
    ls.ch_check = make(chan Pid,3)
    
    go func(){
        for {
            select {
            case <- ls.tick:
                for pid,countdown := range ls.pid2countdown {
                    if countdown > 0 {
                        ls.pid2countdown[pid] = countdown - 1
                    }else{
                        //@todo player linebroken countdown timeout!
                        for _,gameService := range games_service.ListGameServices(pid) {
                            result := Result{} //make a result,because player timeout
                            result.EndType = LineBroken
                            result.Winner = gameService.getColor(gameService.otherPid(pid))
                            /* timeout,has result now,game should over */
                            gameService.setResult(result)
                        }
                    }
                }
                
            case pid := <- ls.ch_linebroken:
                ls.pid2countdown[pid] = Countdown
                //tell games of the player playing
                for _,gameService := range games_service.ListGameServices(pid) {
                    gameService.linebroken(pid)
                }
                
            case pid := <- ls.ch_check:
                if _,ok := ls.pid2countdown[pid]; ok {
                    //comeback
                    delete(ls.pid2countdown,pid)
                    for _,gameService := range games_service.ListGameServices(pid) {
                        /* trigger game_service player comeback */
                        gameService.comeback(pid)
                    }
                }
            }//select end
        }
    }()
}