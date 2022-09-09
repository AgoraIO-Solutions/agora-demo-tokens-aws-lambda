package main

import (
	"encoding/json"
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/RtcTokenBuilder"
	rtmtokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/RtmTokenBuilder"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Token struct {
	UID      uint64            `json:"uid"`
	RTMUID   string            `json:"rtmuid"`
	CHANNELS map[string]string `json:"channels"`
	RTM      string            `json:"rtm"`
}

func generateRTMToken(uid string) string {
	appID := os.Getenv("APP_ID")
	appCertificate := os.Getenv("CERTIFICATE")
	expireTimeInSeconds := uint32(24 * 60 * 60)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	result, err := rtmtokenbuilder.BuildToken(appID, appCertificate, uid, rtmtokenbuilder.RoleRtmUser, expireTimestamp)

	if err != nil {
		log.Printf("Error %+v", err)
	}
	return result
}

func generateRtcToken(uid uint32, channelName string, role rtctokenbuilder.Role) string {
	appID := os.Getenv("APP_ID")
	appCertificate := os.Getenv("CERTIFICATE")
	tokenExpireTimeInSeconds := uint32(60 * 60 * 24)
	result, err := rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, uid, role, tokenExpireTimeInSeconds)
	if err != nil {
		log.Printf("Error %+v", err)
	}
	return result
}

func generateARandomUID() uint32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Uint32()
}

func getToken(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	uidStr := request.QueryStringParameters["uid"]
	uid, _ := strconv.ParseUint(uidStr, 0, 64)
	channels := request.MultiValueQueryStringParameters["channels"]
	rtmUid := strconv.FormatUint(uint64(uid), 10)
	rtmToken := generateRTMToken(rtmUid)

	tokensMap := make(map[string]string)
	for _, channel := range channels {
		rtcToken := generateRtcToken(uint32(uid), channel, rtctokenbuilder.RolePublisher)
		tokensMap[channel] = rtcToken
	}

	token := Token{
		UID: uid, CHANNELS: tokensMap, RTM: rtmToken, RTMUID: rtmUid,
	}
	tokenText, _ := json.Marshal(token)

	return events.APIGatewayProxyResponse{
		Body:       string(tokenText),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(getToken)
}
