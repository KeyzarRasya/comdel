package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type MidtransTime time.Time

func (mt *MidtransTime) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	if str == "" {
		*mt = MidtransTime(time.Time{})
		return nil
	}

	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return err
	}

	*mt = MidtransTime(t)
	return nil
}

func (mt MidtransTime) Time() time.Time {
	return time.Time(mt)
}

func Status(orderId string) (*Subscription, error) {
	var statusResponse Subscription;
	var client http.Client;
	var endpoint = fmt.Sprintf("https://api.sandbox.midtrans.com/v2/%s/status", orderId);

	var serverKey = fmt.Sprintf("%s:", os.Getenv("MIDTRANS_SERVER_KEY"))
	var encodedServerKey = base64.StdEncoding.EncodeToString([]byte(serverKey))

	req, err := http.NewRequest("GET", endpoint, nil);

	if err != nil {
		return nil, err;
	}

	req.Header.Add("Accept", "application/json");
	req.Header.Add("Content-Type", "application/json");
	req.Header.Add("Authorization", encodedServerKey);

	resp, err := client.Do(req)

	if err != nil {
		return nil, err;
	}

	byteBody, err := io.ReadAll(resp.Body);

	if err != nil {
		return nil, err;
	}

	if err := json.Unmarshal(byteBody, &statusResponse); err != nil {
		return nil, err;
	}

	return &statusResponse, nil;
}