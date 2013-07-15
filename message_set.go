package kafka

type MessageBlock struct {
	Offset int64
	Msg    *Message
}

func (msb *MessageBlock) encode(pe packetEncoder) {
	pe.putInt64(msb.Offset)
	pe.pushLength32()
	msb.Msg.encode(pe)
	pe.pop()
}

func (msb *MessageBlock) decode(pd packetDecoder) (err error) {
	msb.Offset, err = pd.getInt64()
	if err != nil {
		return err
	}

	err = pd.pushLength32()
	if err != nil {
		return err
	}

	msb.Msg = new(Message)
	err = msb.Msg.decode(pd)
	if err != nil {
		return err
	}

	err = pd.pop()
	if err != nil {
		return err
	}

	return nil
}

type MessageSet struct {
	Messages []*MessageBlock
}

func (ms *MessageSet) encode(pe packetEncoder) {
	for i := range ms.Messages {
		ms.Messages[i].encode(pe)
	}
}

func (ms *MessageSet) decode(pd packetDecoder) (err error) {
	ms.Messages = nil

	for pd.remaining() > 0 {
		msb := new(MessageBlock)
		err = msb.decode(pd)
		if err != nil {
			return err
		}
		ms.Messages = append(ms.Messages, msb)
	}

	return nil
}

func (ms *MessageSet) addMessage(msg *Message) {
	block := new(MessageBlock)
	block.Msg = msg
	ms.Messages = append(ms.Messages, block)
}

func newMessageSet() *MessageSet {
	set := new(MessageSet)
	set.Messages = make([]*MessageBlock, 0)
	return set
}
