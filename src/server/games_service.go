package main

//---------------------------------------------------------
//put game serivce
func (gss *GamesService) Put(gameService *GameService) {
    gss.ch_put_gameservice <- gameService
}
//---------------------------------------------------------
//get gameService
type get_result struct {
    gameService *GameService
    ok bool
}

type get_param struct {
    gid Gid
    ch_result chan get_result
}

func (gss *GamesService) Get(gid Gid) (*GameService,bool) {
    ch_result := make(chan get_result)
    gss.ch_get_gameservice <- get_param{gid,ch_result}
    result := <- ch_result
    return result.gameService,result.ok
}
//---------------------------------------------------------
func (gss *GamesService) List() []*GameService {
    ch_result := make(chan []*GameService)
    gss.ch_list_gameservice <- ch_result
    return <- ch_result
}
//---------------------------------------------------------
type pid_param struct {
    pid Pid
    ch_result chan []*GameService
}

func (gss *GamesService) ListGameServices(pid Pid) []*GameService {
    ch_result := make(chan []*GameService)
    gss.ch_pid_gameservice <- pid_param{pid,ch_result}
    return <- ch_result
}
//---------------------------------------------------------

func (gss *GamesService) start() {
    gss.gid2gameservice = make(map[Gid]*GameService)
    
    gss.ch_put_gameservice = make(chan *GameService,3)
    gss.ch_get_gameservice = make(chan get_param,3)
    gss.ch_list_gameservice = make(chan chan []*GameService,3)
    gss.ch_pid_gameservice = make(chan pid_param,3) //get player serivces
    
    go func(){
        for {
            select {
            case gameService := <- gss.ch_put_gameservice:
                //put
                gss.gid2gameservice[gameService.gid] = gameService
                
            case param := <- gss.ch_get_gameservice:
                //get
                gameService,ok := gss.gid2gameservice[param.gid]
                param.ch_result <- get_result{gameService,ok}
                
            case ch_result := <- gss.ch_list_gameservice:
                //list
                var gameServices []*GameService
                for _,gameService := range gss.gid2gameservice {
                    gameServices = append(gameServices,gameService)
                }
                ch_result <- gameServices
                
            case param := <- gss.ch_pid_gameservice:
                //player in game services
                var gameServices []*GameService
                for _,gameService := range gss.gid2gameservice {
                    if gameService.hasPid(param.pid) {
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
    ch_get_gameservice chan get_param //get game service
    ch_list_gameservice chan chan []*GameService //list game service
    ch_pid_gameservice chan pid_param //game serivces who has player
}
//---------------------------------------------------------
