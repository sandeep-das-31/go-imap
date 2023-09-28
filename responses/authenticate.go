package responses

import (
	"encoding/base64"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-sasl"
)

// An AUTHENTICATE response.
type Authenticate struct {
	Mechanism       sasl.Client
	InitialResponse []byte
	RepliesCh       chan []byte
}

// Implements
func (r *Authenticate) Replies() <-chan []byte {
	return r.RepliesCh
}

func (r *Authenticate) writeLine(l string) error {
	r.RepliesCh <- []byte(l + "\r\n")
	return nil
}

func (r *Authenticate) cancel() error {
	return r.writeLine("*")
}

func (r *Authenticate) Handle(resp imap.Resp) error {
	log.Println("running handle from Authenticate")
	cont, ok := resp.(*imap.ContinuationReq)
	if !ok {
		log.Println("running handle from Authenticate ErrUnhandled ", ErrUnhandled)
		return ErrUnhandled
	}

	// Empty challenge, send initial response as stated in RFC 2222 section 5.1
	if cont.Info == "" && r.InitialResponse != nil {
		log.Println("running handle from Authenticate Empty challenge ")
		encoded := base64.StdEncoding.EncodeToString(r.InitialResponse)
		if err := r.writeLine(encoded); err != nil {
			log.Println("running handle from Authenticate error in writing ", err)
			return err
		}
		r.InitialResponse = nil
		log.Println("running handle from Authenticate in writing return nil ")
		return nil
	}
	log.Println("running handle from Authenticate challenge start ")
	challenge, err := base64.StdEncoding.DecodeString(cont.Info)
	if err != nil {
		log.Println("running handle from Authenticate challenge err ", err)
		r.cancel()
		return err
	}
	log.Println("running handle from Authenticate challenge end ")

	reply, err := r.Mechanism.Next(challenge)
	if err != nil {
		log.Println("running handle from Authenticate reply error ", err)
		r.cancel()
		return err
	}
	log.Println("running handle from Authenticate reply end ")
	encoded := base64.StdEncoding.EncodeToString(reply)
	log.Println("running handle from Authenticate encode and return ")
	return r.writeLine(encoded)
}
