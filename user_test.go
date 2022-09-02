package sockshope2e_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/thanhpk/randstr"
)

const USER_URL = "http://34.69.149.34"

type UserReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserRespBody struct {
	Id string `json:"id"`
}

// taken from https://apitest.dev/docs/config#networking - configure as needed
var cookieJar, _ = cookiejar.New(nil)
var cli = &http.Client{
	Timeout: time.Second * 15,
	Jar:     cookieJar,
}

func TestUserCRUD(t *testing.T) {
	reqBodyStruct := UserReqBody{
		Username: randstr.String(8), // This is a random 3rd party library for the sake of the interview, I'd use something different in the real world
		Password: randstr.String(8),
		Email:    randstr.String(8),
	}
	reqBody, _ := json.Marshal(reqBodyStruct)
	respBody := UserRespBody{}

	apitest.New().
		EnableNetworking(cli).
		Post(fmt.Sprintf("%s/register", USER_URL)).
		JSON(string(reqBody)).
		Expect(t).
		Assert(jsonpath.Present("id")).
		Status(200). // getting 200, but create should probably be 201
		End().
		JSON(&respBody)
	fmt.Println(respBody.Id) // Useful for debugging, would not be for production
	// Defer delete as cleanup
	defer apitest.New().
		EnableNetworking(cli).
		Delete(fmt.Sprintf("%s/customers/%s", USER_URL, respBody.Id)).
		Expect(t).
		Body("{\"status\": true}").
		Status(200).
		End()

	apitest.New().
		EnableNetworking(cli).
		Get(fmt.Sprintf("%s/customers/%s", USER_URL, respBody.Id)).
		Expect(t).
		Assert(jsonpath.Equal("username", reqBodyStruct.Username)).
		Assert(jsonpath.Equal("firstName", "")).
		Assert(jsonpath.Equal("lastName", "")). // email is apparently not returned
		Status(200).
		End()

}
