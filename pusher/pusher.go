package pusher

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nats-io/nats.go"
)

type Pusher struct {
	natsUrl   string
	stream    string
	interval  int
	userCreds string
	work      chan struct{}
	nc        *nats.Conn
}

type nscCreds struct {
	UserCreds string `json:"user_creds"`
	Operator  struct {
		Service []string `json:"service"`
	} `json:"operator"`
}

type Person struct {
	Name  string `fake:"{firstname}"`
	Email string `fake:"{email}"`
	Phone string `fake:"{phone}"`
}

type Order struct {
	ID   string      `fake:"{uuid}"`
	Date time.Time   `fake:"{date}"`
	CC   interface{} `fake:"{creditcard}"`
}

type FakeOrder struct {
	Person *Person `json:"person"`
	Order  *Order  `json:"order"`
}

func New(stream, user, service string, interval int) (*Pusher, error) {
	var cresdsFile string
	if strings.HasPrefix(user, "nsc://") {
		path, err := exec.LookPath("nsc")
		if err != nil {
			return nil, fmt.Errorf("nsc required: %w", err)
		}

		cmd := exec.Command(path, "generate", "profile", user)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("nsc invoke failed: %s", string(out))
		}

		var creds nscCreds
		err = json.Unmarshal(out, &creds)
		if err != nil {
			return nil, fmt.Errorf("nsc parse error: %s", err)
		}

		if len(creds.Operator.Service) == 0 {
			return nil, errors.New("no services defined for operator")
		}

		cresdsFile = creds.UserCreds

		// prefer the service definition in the operator
		service = strings.Join(creds.Operator.Service, ",")
	} else {
		if _, err := os.Stat(user); err != nil {
			return nil, fmt.Errorf("unable to open creds file: %s", err.Error())
		}

		cresdsFile = user
	}

	if interval == 0 {
		interval = 2
	}

	return &Pusher{
		natsUrl:   service,
		stream:    stream,
		interval:  interval,
		userCreds: cresdsFile,
		work:      make(chan struct{}),
	}, nil
}

func (p *Pusher) Start() error {
	fmt.Println("starting pusher")
	nc, err := nats.Connect(p.natsUrl, nats.UserCredentials(p.userCreds))
	if err != nil {
		return err
	}

	p.nc = nc

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-p.work:
				fmt.Println("pusher exit")
				return
			case <-time.After(1 * time.Second):
				fmt.Println("sending message")
				js.Publish(p.stream, newRandomMessage())
			}
		}
	}()

	return nil
}

func (p *Pusher) Stop() {
	fmt.Println("Stopping pusher")
	close(p.work)
	p.nc.Close()
}

func newRandomMessage() []byte {
	var fo FakeOrder
	gofakeit.Struct(&fo)

	bytes, _ := json.Marshal(fo)
	return bytes
}
