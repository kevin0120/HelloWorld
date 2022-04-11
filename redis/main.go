package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

func ExampleNewClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func ExampleClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Set("储婷玉1", 126666, 0).Err()
	if err != nil {
		panic(err)
	}
	err = client.Set("杨敬伟", 12656545, 0).Err()
	if err != nil {
		panic(err)
	}
	val, err := client.Get("杨敬伟").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("杨敬伟", val)

	val2, err := client.Get("储婷玉").Result()
	if err == redis.Nil {
		fmt.Println("储婷玉 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("储婷玉", val2)
	}
	// Output: key value
	// key2 does not exist
}
func main() {
	ExampleNewClient()
	ExampleClient()

}
