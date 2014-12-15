//TODO: Check for API errors returned at API level
package gotrade

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mrjones/oauth"
)

var (
	baseURL    = "https://etws.etrade.com/%s/rest/%s"
	sandboxURL = "https://etwssandbox.etrade.com/%s/sandbox/rest/%s"
	jsonURL    = ".json"
)

type ETradeClient struct {
	consumer    *oauth.Consumer
	accessToken *oauth.AccessToken
	url         string
}

type strURLParameter []string

// Creates a comma separated string
func (q strURLParameter) String() string {
	return strings.Join(q, ",")
}

type IntDollar int64

func (i IntDollar) String() (s string) {
	dollars := i / 100
	cents := i % 100
	return fmt.Sprintf("%d.%02d", dollars, cents)
}

func New(consumerID, consumerSecret string) (client ETradeClient, err error) {

	c := oauth.NewConsumer(
		consumerID,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://etws.etrade.com/oauth/request_token",
			AuthorizeTokenUrl: "https://us.etrade.com/e/t/etws/authorize",
			AccessTokenUrl:    "https://etws.etrade.com/oauth/access_token",
		})
	requestToken, url, err := c.GetRequestTokenAndUrl("oob")
	if err != nil {
		return client, err
	}
	url = fmt.Sprintf("https://us.etrade.com/e/t/etws/authorize?key=%s&token=%s", consumerID, requestToken.Token)
	fmt.Println("(1) Go to: " + url + "")
	fmt.Println("(2) Grant access, you should get back a verification code.")
	fmt.Println("(3) Enter that verification code here: ")

	verificationCode := ""
	fmt.Scanln(&verificationCode)

	accessToken, err := c.AuthorizeToken(requestToken, verificationCode)
	if err != nil {
		return client, err
	}

	client.consumer = c
	client.accessToken = accessToken
	client.url = baseURL
	return
}

func (client *ETradeClient) SetToSandBox() {
	client.url = sandboxURL
}

func (client *ETradeClient) SetToProduction() {
	client.url = baseURL
}

func (client *ETradeClient) requestAndUnmarshal(requestURL string, v interface{}) (raw string, err error) {

	r, err := client.consumer.Get(requestURL, nil, client.accessToken)
	if err != nil {
		return
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNotModified {
		return raw, nil
	} else if r.StatusCode != http.StatusOK {
		return raw, fmt.Errorf("Status code not valid for request: %s", r.Status)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &v)
	if err != nil {
		return
	}
	return string(body), err
}
