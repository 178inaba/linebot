package hello

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"github.com/PuerkitoBio/goquery"
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

	client := urlfetch.Client(ctx)
	lc, err := linebot.New(cf.Secret, cf.Token, linebot.WithHTTPClient(client))
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

		switch e.Message.(type) {
		case *linebot.TextMessage:
			resp, err := client.Get("http://b.hatena.ne.jp/ranking/daily")
			if err != nil {
				log.Errorf(ctx, "%v", err)
				return
			}
			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromResponse(resp)
			if err != nil {
				log.Errorf(ctx, "%v", err)
				return
			}

			dom := doc.Find(".entry-list-l .entry-link").First()
			url, _ := dom.Attr("href")
			title := dom.Text()

			msg := fmt.Sprintf("%s\n%s", title, url)
			_, err = lc.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(msg)).Do()
			if err != nil {
				log.Errorf(ctx, "%v", err)
			}
		}
	}
}
