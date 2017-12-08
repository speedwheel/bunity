package session

import (
	"github.com/kataras/iris/sessions"
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/sessions/sessiondb/file"
	"app/config"
	"time"
)

var (
	cookieName = "mycustomsessionid"
	hashKey = []byte("the-big-and-secret-fash-key-here")
	blockKey = []byte("lot-secret-of-characters-big-too")
	secureCookie = securecookie.New(hashKey, blockKey)
	Sessions = sessions.New(sessions.Config{
		Cookie: cookieName,
		Encode: secureCookie.Encode,
		Expires: time.Hour * 24,
		Decode: secureCookie.Decode,
		DisableSubdomainPersistence: true,
	})

)

func init() {
	db, _ := file.New(config.GetAppPath()+"/sessions/", 0666)
	db.Async(true)
	Sessions.UseDatabase(db)
}


/*package db

import (
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
	"github.com/gorilla/securecookie"
)

func InitRedisSesstion() sessions.Sessions {
	cookieName := "kazelisession"
	hashKey := []byte("the-big-and-secret-fash-key-here")
	blockKey := []byte("lot-secret-of-characters-big-too")
	secureCookie := securecookie.New(hashKey, blockKey)
	db := redis.New(service.Config{Network: service.DefaultRedisNetwork,
		Addr:          service.DefaultRedisAddr,
		Password:      "",
		Database:      "",
		MaxIdle:       0,
		MaxActive:     0,
		IdleTimeout:   service.DefaultRedisIdleTimeout,
		Prefix:        "",
		MaxAgeSeconds: service.DefaultRedisMaxAgeSeconds})
	
	mySessions := sessions.New(sessions.Config{
		Cookie: cookieName,
		Encode: secureCookie.Encode,
		Decode: secureCookie.Decode,
	})
	mySessions.UseDatabase(db)
	
	return mySessions
	
}*/