package outputdto

type ErrMsgDTO struct {
	Err         error  `json:"error"`
	ListenerTag string `json:"listener_tag"`
	Msg         []byte `json:"msg"`
}
