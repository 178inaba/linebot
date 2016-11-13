package hello

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"

	"conf"
)

var cf *conf.Conf

func init() {
	var err error
	cf, err = conf.LoadConf("etc/conf.toml")
	if err != nil {
		panic(err)
		//log.Criticalf(appengine.NewContext(nil), "Load config error: %s", err)
		//return
	}

	http.Handle("/", routes())
}

func routes() *gin.Engine {
	r := gin.Default()

	r.POST("/", index)

	return r
}

func index(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	lc, err := linebot.New(cf.Secret, cf.Token,
		linebot.WithHTTPClient(urlfetch.Client(ctx)))
	if err != nil {
		log.Errorf(ctx, "%v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	events, err := lc.ParseRequest(c.Request)
	if err != nil {
		log.Errorf(ctx, "%v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, e := range events {
		if e.Type != linebot.EventTypeMessage {
			continue
		}

		switch m := e.Message.(type) {
		case *linebot.TextMessage:
			_, err := lc.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(m.Text)).Do()
			if err != nil {
				log.Errorf(ctx, "%v", err)
			}
		}
	}
}
