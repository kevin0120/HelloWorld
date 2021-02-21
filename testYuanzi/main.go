package main

import "time"

type Human struct {
	sex  string
	live int64
	tag  string
}
func main() {
	var a Human
	var b *Human
	for i := 0; i < 100; i++ {
		go func() {
			a = Human{
				sex:  "man",
				live: 80,
				tag:  "i am a man",
			}
			b = &Human{
				sex:  "girl",
				live: 89,
				tag:  "i am a girl",
			}
		}()
		go func() {
			a = Human{
				sex:  "girl",
				live: 89,
				tag:  "i am a girl",
			}
			b = &Human{
				sex:  "man",
				live: 80,
				tag:  "i am a man",
			}
		}()
	}
	for {
		time.Sleep(1 * time.Second)
	}
}
