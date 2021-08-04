package bootstrap

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	HttpServ  *http.Server
	GinEngibe *gin.Engine
}

func NewApp(debug bool) *App {
	rand.Seed(time.Now().UnixNano())
	app := new(App)
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	app.GinEngibe = gin.Default()
	return app
}

func (app *App) Use(fc ...func()) {
	for _, f := range fc {
		f()
	}
}

func (app *App) RegisterRoutes(registeRoutes func(*gin.Engine)) {
	registeRoutes(app.GinEngibe)
}

func (app *App) Run(port string) {
	app.HttpServ = &http.Server{
		Addr:    ":" + port,
		Handler: app.GinEngibe,
	}
	// 优雅终止
	quit := make(chan os.Signal, 4)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go handleSignal(quit, app)

	log.Printf("http server expect run in port %s", port)
	err := app.HttpServ.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Panic(err)
	}
}

func handleSignal(c <-chan os.Signal, app *App) {
	switch <-c {
	case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
		log.Printf("Shutdown quickly, bye...")
	case syscall.SIGHUP:
		log.Printf("Shutdown gracefully, bye...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := app.HttpServ.Shutdown(ctx); err != nil {
			log.Printf("http Server Shutdown err:%v", err)
		}
	}
	os.Exit(0)
}

var (
	Param struct {
		C string
		H bool
	}
)

func InitFlag() {
	flag.StringVar(&Param.C, "config", "conf/app.json", "配置文件地址")
	flag.BoolVar(&Param.H, "help", false, "帮助")
}

func Flag() bool {
	flag.Parse()
	if Param.H {
		flag.PrintDefaults()
		return false
	}
	// 存到viper中
	return true
}
