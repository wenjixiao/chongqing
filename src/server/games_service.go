package main

import (
    "fmt"
    "os"
    )

//---------------------------------------------------------
//put game serivce
func (gss *GamesService) Put(gameService *GameService) {
    gss.ch_put_gameservice <- gameService
}
//---------------------------------------------------------
//get gameService
type gid_param struct {
    gid Gid
    ch_result chan *GameService
}

func (gss *GamesService) Get(gid Gid) *GameService {
    ch_result := make(chan *GameService)
    gss.ch_get_gameservice <- gid_param{gid,ch_result}
    return <- ch_result
}
//---------------------------------------------------------
func (gss *GamesService) List() []*GameService {
    ch_result := make(chan []*GameService)
    gss.ch_list_gameservice <- ch_result
    return <- ch_result
}
//---------------------------------------------------------
type player_param struct {
    player Pid
    ch_result chan []*GameService
}

func (gss *GamesService) ListGameServices(pid Pid) []*GameService {
    ch_result := make(chan []*GameService)
    gss.ch_player_gameservice <- player_param{pid,ch_result}
    return <- ch_result
}
//---------------------------------------------------------

func (gss *GamesService) start() {
    gss.gid2gameservice = make(map[Gid]*GameService)
    
    gss.ch_put_gameservice = make(chan *GameService,3)
    gss.ch_get_gameservice = make(chan gid_param,3)
    gss.ch_list_gameservice = make(chan chan []*GameService,3)
    gss.ch_player_gameservice = make(chan player_param,3) //get player serivces
    
    go func(){
        for {
            select {
            case gameService := <- gss.ch_put_gameservice:
                //put
                gss.gid2gameservice[gameService.gid] = gameService
                
            case param := <- gss.ch_get_gameservice:
                //get
                if gameService,ok := gss.gid2gameservice[param.gid]; ok {
                    param.ch_result <- gameService
                }else{
                    fmt.Println("error: no GameService of gid=",param.gid)
                    os.Exit(1)
                }
                
            case ch_result := <- gss.ch_list_gameservice:
                //list
                var gameServices []*GameService
                for _,gameService := range gss.gid2gameservice {
                    gameServices = append(gameServices,gameService)
                }
                ch_result <- gameServices
                
            case param := <- gss.ch_player_gameservice:
                //player in game services
                var gameServices []*GameService
                for _,gameService := range gss.gid2gameservice {
                    if gameService.hasPlayer(param.player) {
                        gameServices = append(gameServices,gameService)
                    }
                }
                param.ch_result <- gameServices
            }
        }
    }()
}

type GamesService struct {
    gid2gameservice map[Gid]*GameService
    ch_put_gameservice chan *GameService //put game service
    ch_get_gameservice chan gid_param //get game service
    ch_list_gameservice chan chan []*GameService //list game service
    ch_player_gameservice chan player_param //game serivces who has player
}
//---------------------------------------------------------
