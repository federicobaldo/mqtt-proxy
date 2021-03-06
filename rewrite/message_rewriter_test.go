package rewrite

import (
	"testing"

	"github.com/wolfeidau/mqtt"
	. "launchpad.net/gocheck"
)

func Test2(t *testing.T) { TestingT(t) }

type MessagRewriterSuite struct {
	msgRewriter *MsgRewriter
}

var _ = Suite(&MessagRewriterSuite{})

func (s *MessagRewriterSuite) SetUpTest(c *C) {
	s.msgRewriter = &MsgRewriter{
		CredentialsRewriter: &CredentialsReplaceRewriter{
			User:   "guest",
			Pass:   "123",
			UserId: "1",
			MqttId: "1",
		},
		IngressRewriter: &TopicPartRewriter{
			Token:     "123",
			Direction: INGRESS,
		},
		EgressRewriter: &TopicPartRewriter{
			Token:     "123",
			Direction: EGRESS,
		},
	}
}

func (s *MessagRewriterSuite) TestIngressMsgPublish(c *C) {

	// client publish a message to a topic
	pub := createPublish("$cloud/456/789")
	expectedPub := createPublish("123/$cloud/456/789")

	modPub := s.msgRewriter.RewriteIngress(pub)
	c.Assert(modPub, DeepEquals, expectedPub)

}

func (s *MessagRewriterSuite) TestIngressMsgSubscribe(c *C) {

	// client subscribe a message to a topic
	sub := createSubscribe("$cloud/456/789")
	expectedSub := createSubscribe("123/$cloud/456/789")

	modSub := s.msgRewriter.RewriteIngress(sub)
	c.Assert(modSub, DeepEquals, expectedSub)

}

func (s *MessagRewriterSuite) TestConnect(c *C) {

	// connection request message
	connect := createConnect("bob", "11223344", "abc")
	expectedConnect := createConnect("guest", "123", "1")

	modConnect := s.msgRewriter.RewriteIngress(connect)
	c.Assert(modConnect, DeepEquals, expectedConnect)
}

func (s *MessagRewriterSuite) TestIngressMsgUnsubscribe(c *C) {

	// client unsubscribe to a topic
	unsub := createUnsubscribe("$cloud/456/789")
	expectedUnsub := createUnsubscribe("123/$cloud/456/789")

	modUnsub := s.msgRewriter.RewriteIngress(unsub)
	c.Assert(modUnsub, DeepEquals, expectedUnsub)
}

func (s *MessagRewriterSuite) TestEgressMsgPublish(c *C) {

	// client publish a message to a topic
	pub := createPublish("123/$block/456/789")
	expectedPub := createPublish("$block/456/789")

	modPub := s.msgRewriter.RewriteEgress(pub)
	c.Assert(modPub, DeepEquals, expectedPub)

}

func createConnect(user string, pass string, clientId string) mqtt.Message {
	return &mqtt.Connect{
		ProtocolName:    "MQIsdp",
		ProtocolVersion: 3,
		UsernameFlag:    true,
		PasswordFlag:    true,
		WillRetain:      false,
		WillQos:         1,
		WillFlag:        true,
		CleanSession:    true,
		KeepAliveTimer:  10,
		ClientId:        clientId,
		WillTopic:       "topic",
		WillMessage:     "message",
		Username:        user,
		Password:        pass,
	}
}

func createPublish(topic string) mqtt.Message {
	return &mqtt.Publish{
		Header: mqtt.Header{
			DupFlag:  false,
			QosLevel: mqtt.QosAtMostOnce,
			Retain:   false,
		},
		TopicName: topic,
		Payload:   mqtt.BytesPayload{1, 2, 3},
	}
}

func createSubscribe(topic string) mqtt.Message {
	return &mqtt.Subscribe{
		Header: mqtt.Header{
			DupFlag:  false,
			QosLevel: mqtt.QosAtLeastOnce,
		},
		MessageId: 0x4321,
		Topics: []mqtt.TopicQos{
			{topic, mqtt.QosExactlyOnce},
		},
	}
}

func createUnsubscribe(topic string) mqtt.Message {
	return &mqtt.Unsubscribe{
		Header: mqtt.Header{
			DupFlag:  false,
			QosLevel: mqtt.QosAtLeastOnce,
		},
		MessageId: 0x4321,
		Topics:    []string{topic},
	}
}
