package crane

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"sync"
)

type Crane struct {
	echo          *echo.Echo
	melody        *melody.Melody
	redis         *redis.Client
	playerSession PlayerSession
	Handlers      Handler
	Logger        *log.Logger
}

var instance = &Crane{}

func GetInstance() *Crane {
	return instance
}

func Initialize() error {
	return GetInstance().initialize()
}

func Start() error {
	return GetInstance().start()
}

func HandleAPI(url string, f func(*Crane, echo.Context) error) {
	GetInstance().handleAPI(url, f)
}

func HandleOperation(ope string, handler func(*melody.Session, *Crane, []byte)) {
	GetInstance().handleOperation(ope, handler)
}

func (c *Crane) initialize() error {
	file, err := os.OpenFile("/app/log/crane.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	c.echo = echo.New()
	c.melody = melody.New()
	c.redis = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
		PoolSize: 1000,
	})
	c.playerSession = PlayerSession{
		sessions: sync.Map{},
	}
	c.Logger = log.New(file, "", log.LstdFlags)
	return nil
}

func (c *Crane) start() error {
	c.Logger.Print("exec start")
	if c.echo == nil {
		return fmt.Errorf("web server not create")
	}
	if c.melody == nil {
		return fmt.Errorf("websocket server not available")
	}
	c.echo.Use(middleware.Logger())
	c.echo.Use(middleware.Recover())
	c.echo.GET("/ws", func(con echo.Context) error {
		c.melody.HandleRequest(con.Response().Writer, con.Request())
		return nil
	})

	c.melody.HandleConnect(func(s *melody.Session) {
		c.Logger.Print("exec HandleConnect")
		var uuid string
		if s.Request.Header.Get("x-uuid") != "" {
			uuid = s.Request.Header.Get("x-uuid")
		} else if s.Request.FormValue("uuid") != "" {
			uuid = s.Request.FormValue("uuid")
		} else {
			return
		}
		c.playerSession.Set(uuid, s)
		s.Set("uuid", uuid)
	})

	c.melody.HandleDisconnect(func(s *melody.Session) {
		uuid, _ := s.Get("uuid")
		c.playerSession.Delete(uuid.(string))
		s.Set("uuid", nil)
	})

	c.melody.HandleMessage(func(s *melody.Session, msg []byte) {
		c.Logger.Print("exec HandleMessage")
		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err != nil {
			c.Logger.Print(err.Error())
		}
		val, flg := data["ope"]
		if flg {
			ope := val.(string)
			var handler func(*melody.Session, *Crane, []byte)
			c.Logger.Print(ope)
			if ope == "get" {
				handler = c.Handlers.Get(ope)
			} else {
				handler = c.Handlers.Get("start")
			}
			handler(s, c, msg)
		}
	})

	go c.CreateSender()
	errChan := make(chan error, 1)
	go func(addr string) {
		errChan <- c.echo.Start(addr)
	}(":1323")
	err := <-errChan
	return err
}

func (c *Crane) handleAPI(uri string, f func(*Crane, echo.Context) error) {
	c.echo.POST(uri, func(ctx echo.Context) error {
		return f(c, ctx)
	})
}

func (c *Crane) handleOperation(ope string, handler func(*melody.Session, *Crane, []byte)) {
	c.Logger.Print("exec handleOperation")
	c.Handlers.Set(ope, handler)
}

func (c *Crane) Broadcast(msg []byte) {
	c.melody.Broadcast(msg)
}

func (c *Crane) Send(msg []byte) error {
	ctx := context.Background()
	c.redis.Publish(ctx, "crane_subscribe_channel", msg)
	return nil
}

func (c *Crane) CreateSender() {
	ctx := context.Background()
	pubsub := c.redis.Subscribe(ctx, "crane_subscribe_channel")
	for {
		for reply := range pubsub.Channel() {
			jsonMessage := map[string]string{}
			json.Unmarshal([]byte(reply.Payload), &jsonMessage)
			target, ok := jsonMessage["target"]
			if target == "broadcast" {
				c.Broadcast([]byte(reply.Payload))
			} else if ok {
				session := c.playerSession.Get(jsonMessage["target"])
				if session != nil {
					session.Write([]byte(reply.Payload))
				}
			}
		}
	}
}
