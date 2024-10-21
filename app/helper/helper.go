package helper

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/umutteroll07/ShortenURL/app/model"
	"github.com/umutteroll07/ShortenURL/internal/database"
)

var ctx = context.Background()

func RedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})

	fmt.Println("redis connection is success")
	return rdb
	
}

func GetTTL() time.Duration {
    ttl, err := time.ParseDuration(os.Getenv("TTL") + "s")
    if err != nil {
        panic("Invalid TTL value")
    }
    return ttl
}


func CheckExpiration(shortUrl string) (bool, error) {
    rdb := RedisClient()
    ttl, err := rdb.TTL(ctx, shortUrl).Result()
    if err != nil {
        return false, err
    }
    if ttl <= 0 {
        return false, nil
    }
    return true, nil
}

var seededRand *rand.Rand = rand.New(
    rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int) string {


	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[seededRand.Intn(len(charset))]
    }

	UniqueShortUrlControl(string(b))

    return string(b)
}


func UniqueShortUrlControl(shrtUrl string){
	 db := database.Connection()
	 var urls []model.Url
	 db.Select("short_url").Find(&urls)
 
	 var shortUrlList []string
	 for _, url := range urls {
		 shortUrlList = append(shortUrlList, url.Short_url)
	 }

	for _, value := range shortUrlList{
		if value == os.Getenv("BASE_URL") + shrtUrl{
			fmt.Println("same url")
			StringWithCharset(6)
		}
	}
}

func UniqueOriginalUrlControl(orginalUrl string) bool{
	db := database.Connection()
	var urls []model.Url
	db.Select("original_url").Find(&urls)

	var originalUrlList []string
	for _, url := range urls{
		originalUrlList = append(originalUrlList, url.Original_url)
	}

	var isContain bool
	for _, value := range originalUrlList{
		if value == orginalUrl{
			isContain = true
			break
		}
	}
	if isContain{
		return true
	}
	return false
	
}
