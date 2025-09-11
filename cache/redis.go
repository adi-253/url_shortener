package cache

import(
	"context"
	"github.com/redis/go-redis/v9"
	"time"
	"log"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	Password: "",
	DB: 0,
})

func CachePut(short_url string, long_url string) error{
	err := rdb.Set(ctx, short_url, long_url, time.Hour).Err()
	log.Printf("Putting to Cache")
	return err

}

func CacheGet(short_url string) (string,error){
	long_url, err := rdb.Get(ctx, short_url).Result()
	log.Printf("Fetching from Cache")
	return long_url, err

}
