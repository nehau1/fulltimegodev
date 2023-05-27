package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type UserProfile struct {
	ID       int
	Comments []string
	Likes    int
	Friends  []int
}

func main() {
	start := time.Now()
	userProfile, err := handleGetUserProfile(10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(userProfile)
	fmt.Println("fetching the user profile took ", time.Since(start))
}

type Response struct {
	data any
	err  error
}

func handleGetUserProfile(id int) (*UserProfile, error) {
	var (
		respch = make(chan Response, 3)
		wg     = &sync.WaitGroup{}
	)
	//handle 3 request inside their goroutines
	go getComments(id, respch, wg)
	go getLikes(id, respch, wg)
	go getFriends(id, respch, wg)

	wg.Add(3)
	wg.Wait() // block until the wg counter == 0 we unblock

	close(respch)

	userProfile := &UserProfile{}

	for response := range respch {
		if response.err != nil {
			return nil, response.err
		}
		switch msg := response.data.(type) {
		case int:
			userProfile.Likes = msg
		case []string:
			userProfile.Comments = msg
		case []int:
			userProfile.Friends = msg
		}

	}
	return userProfile, nil
}

func getComments(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 200)
	comments := []string{
		"Hey, that was great!",
		"Yeah buddy",
		"Ow, I didn't know that",
	}

	respch <- Response{
		data: comments,
		err:  nil,
	}

	wg.Done()

}

func getLikes(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 200)
	respch <- Response{
		data: 20,
		err:  nil,
	}
	//done
	wg.Done()
}

func getFriends(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 100)
	respch <- Response{
		data: []int{11, 22, 33, 44, 55},
		err:  nil,
	}
	wg.Done()
}
