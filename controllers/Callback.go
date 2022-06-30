package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	data,err:=getUserData(state,code)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Fprintf(w,"Success!  Data: %s",data)
}

func getUserData(state string, code string) ([]byte,error) {
	if state != RandomString {
		return nil,fmt.Errorf("invalid state")
	}
	
	token,err := ssogolang.Exchange(context.Background(),code)
	if err!=nil{
		return nil,err
	}
	
	response,err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token="+token.AccessToken)
	if err!=nil{
		return nil,err
	}
	defer response.Body.Close()

	data,err:=io.ReadAll(response.Body)
	if err!=nil{
		return nil,err
	}

	return data,nil

}