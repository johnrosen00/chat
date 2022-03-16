package rabbitpayload

//Payload is what a message from rabbit should look like
type Payload struct {
	UserIDs []int64 `json:"userIDs"`
	Body    []byte  `json:"body"`
}

// //MessageBody is marshaled into body when sending messages over q
// type MessageBody struct {
// 	Type string            `json:"type"`
// 	Body *messages.Message `json:"body"`
// }

// //ChannelBody is marshaled into body when sending channels over q
// type ChannelBody struct {
// 	Type string            `json:"type"`
// 	Body *channels.Channel `json:"body"`
// }

// //Int64Body is marshaled into body when sending channels over q
// type Int64Body struct {
// 	Type string `json:"type"`
// 	Body int64  `json:"body"`
// }
