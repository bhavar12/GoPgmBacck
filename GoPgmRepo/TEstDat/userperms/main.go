package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func userPermission(ctx context.Context) error {
	req, err := http.NewRequest("GET", "http://localhost:10000/", nil)
	if err != nil {
		fmt.Println("Error in creating req", err.Error())
		return err
	}
	// Propagate the incoming context.
	req = req.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in req", err.Error())
		return err
	}
	//defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in resp", err.Error())
		return err
	}
	fmt.Println("Response from api", string(bodyBytes))
	fmt.Println("userPermission task completed")
	return nil
}

func getGuest(num int) error {
	time.Sleep(2 * time.Second)
	if num/2 == 0 {
		return errors.New("test")
	}
	fmt.Println("getGuest task completed")
	return nil
}

func main() {
	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error, 1)
	userChan := make(chan bool, 1)
	guestChan := make(chan bool, 1)
	defer func() {
		close(errChan)
		close(userChan)
		close(guestChan)
	}()
	userId := "11"
	var err error
	if userId != "" {
		go func() {
			defer func() {
				fmt.Println("go routine completed userPermission")
			}()
			err = userPermission(ctx)
			if err != nil {

				fmt.Println("userPermission err", err)
				fmt.Println("started write to errChan userPermission")
				errChan <- err
				fmt.Println("completed writing to errChan userPermission")
				return
			}
			fmt.Println("writing to userChan")
			userChan <- true

		}()
	}

	go func() {
		defer func() {
			fmt.Println("go routine completed getGuest")
		}()
		err = getGuest(2)
		if err != nil {
			cancel()
			fmt.Println("getGuest err", err)
			fmt.Println("started write to errChan getGuest")
			errChan <- err
			fmt.Println("completed writing to errChan getGuest")
			return
		}
		fmt.Println("writing to guestChan")
		guestChan <- true

	}()

	var upflag bool
	var guestflag bool
	var errFlag bool
	for {
		select {
		case <-errChan:
			fmt.Println("error occured")
			errFlag = true
		case <-userChan:
			fmt.Println("userPermission success")
			upflag = true
		case <-guestChan:
			fmt.Println("getGuest success")
			guestflag = true
		default:
		}
		if (userId != "" && upflag && guestflag) || (userId == "" && guestflag) || errFlag {
			break
		}
	}

	if errFlag {
		fmt.Println("return error")
	}
	if userId == "" && guestflag {
		fmt.Println("return guest data without filter")
	}
	if userId != "" && upflag && guestflag {
		fmt.Println("return guest data with filter")
	}

	time.Sleep(5 * time.Second)
	fmt.Println("completed in ", time.Since(start))
}
