package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/almaz-uno/searching-helper/pkg/runt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ryboe/q"
)

var (
	cloudID  = "vgpc:dXMtY2VudHJhbDEuZ2NwLmNsb3VkLmVzLmlvOjQ0MyRmZTkzMWJhNDlkODQ0NWQ2YTVmM2FjMGYzOWRkMTNkMiQyYWUwOGNmZjc5NmM0N2YzOTQ4OTU0M2IxZDM3ZWQzNQ=="
	cloudKey = "Y0VpN21KTUJ0VHZNS2xNSVVqb2o6T2hQZHVJeGtTdC1UU0I3M0c2QTRnZw=="

	cfgListenOn = runt.CfgEnv("LISTEN_ON", ":32080")
	cfgCertFile = runt.CfgEnv("CERT_FILE", "")
	cfgKeyFile  = runt.CfgEnv("KEY_FILE", "")

	indexName = "videogames"
)

type (
	gravitsappa struct {
		client *elasticsearch.Client
	}
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	runt.Main(func(ctx context.Context, cancel context.CancelFunc) error {
		defer cancel()

		client, err := elasticsearch.NewClient(elasticsearch.Config{
			CloudID: cloudID,
			APIKey:  cloudKey,
		})
		if err != nil {
			log.Error().Err(err).Msg("Unable to create Elasticsearch client")
			return err
		}

		echoServer := echo.New()
		echoServer.Use(middleware.CORS())
		echoServer.Use(middleware.Logger())
		echoServer.Debug = true

		gravitsappa := &gravitsappa{
			client: client,
		}
		echoServer.GET("/debug", gravitsappa.getDebug, logErrors)
		echoServer.POST("/debug", gravitsappa.postDebug, logErrors)

		go func() {
			defer cancel()

			if cfgCertFile != "" && cfgKeyFile != "" {
				log.Info().Msg("Starting in HTTPS mode")
				if e := echoServer.StartTLS(cfgListenOn, cfgCertFile, cfgKeyFile); e != nil && !errors.Is(e, http.ErrServerClosed) {
					log.Error().Err(e).Msg("Unable to start echo server")
				} else {
					log.Info().Msg("Exiting echo due to end of work")
				}
			} else {
				log.Info().Msg("Starting in HTTP mode")
				if e := echoServer.Start(cfgListenOn); e != nil && !errors.Is(e, http.ErrServerClosed) {
					log.Error().Err(e).Msg("Unable to start echo server")
				} else {
					log.Info().Msg("Exiting echo due to end of work")
				}
			}
		}()

		<-ctx.Done()

		return nil
	})
}

func logErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			log.Error().Err(err).Stack().
				Str("requestURI", c.Request().RequestURI).
				Msg("Error while processing request")
		}
		return err
	}
}

//go:embed assets/debug.html
var debugHtmlTemplate string

//go:embed assets/defaultTemplate.json
var defaultTemplate string

const (
	placeHolder = "$$query$$"
)

func (gravitsappa *gravitsappa) getDebug(c echo.Context) error {
	data := map[string]any{
		"q": c.QueryParams().Get("q"),
		"t": c.QueryParams().Get("t"),
	}

	if data["t"] == "" {
		data["t"] = defaultTemplate
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	tmpl, err := template.New("debug").Parse(debugHtmlTemplate)
	if err != nil {
		log.Fatal().Err(err).Send()
		panic(err)
	}

	return tmpl.Execute(c.Response(), data)
}

func (gravitsappa *gravitsappa) postDebug(c echo.Context) error {
	data := map[string]any{
		"q": c.FormValue("q"),
		"t": c.FormValue("t"),
	}

	if data["t"] == "" {
		data["t"] = defaultTemplate
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	tmpl, err := template.New("debug").Parse(debugHtmlTemplate)
	if err != nil {
		log.Fatal().Err(err).Send()
		panic(err)
	}

	sp := strconv.Quote(data["q"].(string))
	sp = sp[1 : len(sp)-1]
	sr := strings.Replace(data["t"].(string), placeHolder, sp, -1)
	data["sr"] = sr
	data["indexName"] = indexName

	res, err := gravitsappa.client.Search(
		gravitsappa.client.Search.WithContext(c.Request().Context()),
		gravitsappa.client.Search.WithIndex(indexName),
		gravitsappa.client.Search.WithHuman(),
		gravitsappa.client.Search.WithBody(strings.NewReader(sr)),
	)
	if err != nil {
		return err
	}

	bb, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.IsError() {
		return errors.New(fmt.Sprintf("search operation returned status %s: %s", res.Status(), string(bb)))
	}

	r := new(SearchResponse)
	err = json.Unmarshal(bb, r)
	if err != nil {
		return err
	}

	for i, h := range r.Hits.Hits {
		zz, _ := json.MarshalIndent(h.Source, "", "  ")
		r.Hits.Hits[i].SourceStr = string(zz)
		r.Hits.Hits[i].Number = i + 1
	}

	buff := &bytes.Buffer{}
	json.Indent(buff, bb, "", "  ")
	data["rawResponse"] = string(buff.Bytes())

	q.Q(r)
	data["r"] = r

	return tmpl.Execute(c.Response(), data)
}

// func testAccess(ctx context.Context, client *elasticsearch.Client) error {
// 	request := `{
//   "query": {
//     "bool": {
//       "should": [
//         {
//           "fuzzy": {
//             "Name": {
//               "boost": 1,
//               "value": "Indiana Jones RPG"
//             }
//           }
//         },
//         {
//           "match": {
//             "Genre": {
//               "boost": 2,
//               "query": "Indiana Jones RPG"
//             }
//           }
//         },
//         {
//           "match": {
//             "ConsoleName": {
//               "boost": 2,
//               "query": "Indiana Jones RPG"
//             }
//           }
//         }
//       ]
//     }
//   }
// }`

// 	res, err := client.Search(
// 		client.Search.WithContext(ctx),
// 		client.Search.WithIndex("videogames"),
// 		client.Search.WithBody(strings.NewReader(request)),
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	bb, _ := io.ReadAll(res.Body)
// 	res.Body.Close()
// 	if res.IsError() {
// 		return errors.New(fmt.Sprintf("search operation returned status %s: %s", res.Status(), string(bb)))
// 	}

// 	r := new(SearchResponse)
// 	err = json.Unmarshal(bb, r)
// 	if err != nil {
// 		return err
// 	}

// 	q.Q(r)
// 	return nil
// }
