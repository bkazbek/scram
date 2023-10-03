package internal

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/xdg-go/pbkdf2"
	"github.com/xdg-go/scram"
	mathRand "math/rand"
	"net"
	"os"
	"os/signal"
	"scram_tcp/utils"
	"syscall"
	"time"
)

const (
	HOST         = "0.0.0.0:"
	TYPE         = "tcp"
	AUTH_TIMEOUT = 15
)

type Server struct {
	*scram.Server
	ctx         context.Context
	listener    net.Listener
	logger      *log.Logger
	authTimeout time.Duration
}

func NewServer(username, password, port string) (*Server, error) {
	listen, err := net.Listen(TYPE, HOST+port)
	if err != nil {
		return nil, err
	}
	credentialsLookupFunc := configureUserCreds(username, password)

	scramServer, err := scram.SHA256.NewServer(credentialsLookupFunc)
	if err != nil {
		return nil, err
	}
	return &Server{
		Server:      scramServer,
		listener:    listen,
		ctx:         context.Background(),
		logger:      log.New(),
		authTimeout: time.Second * AUTH_TIMEOUT,
	}, nil
}

func (s *Server) Start() {
	// close listener
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				s.logger.Error("listen accept err:", err)
				break
			}
			go s.handleIncomingRequest(conn)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	s.logger.Println("starting to shutdown")
	if err := s.Close(); err != nil {
		s.logger.Println("got error on shutdown", err)
	}

	s.logger.Println("good bye!")
}

func (s *Server) Close() error {
	s.ctx.Done()
	return s.listener.Close()
}

func (s *Server) authorize(ctx context.Context, conn net.Conn) error {
	conv := s.NewConversation()

	var err error

	for {
		if conv.Done() {
			s.logger.Info("conv done")
			break
		}
		var buffer []byte
		buffer, err = utils.ReadConnWithCtx(ctx, conn)
		if err != nil {
			s.logger.Error("conn read err:", err)
			break
		}

		var got string
		got, err = conv.Step(string(buffer))
		if err != nil {
			s.logger.Error("got step err:", err, got, string(buffer))
			break
		}

		err = utils.WriteConn(conn, []byte(got))
		if err != nil {
			s.logger.Error("write err:", err)
			break
		}
	}

	if err != nil {
		s.logger.Error("got err:", err)
		return err
	}

	if !conv.Valid() {
		s.logger.Error("SHA-256: Conversation is not valid")
		return errors.New("conv is not valid")
	}

	if !conv.Done() {
		return errors.New("conv is not finished")
	}

	return nil
}

func (s *Server) handleIncomingRequest(conn net.Conn) {
	defer conn.Close()
	ctx, cancelFunc := context.WithTimeout(s.ctx, s.authTimeout)
	defer cancelFunc()
	err := s.authorize(ctx, conn)
	if err != nil {
		s.logger.Error("scram auth failed:", err)
		return
	}

	randomQuote := Quotes[mathRand.Intn(len(Quotes))]
	data, _ := json.Marshal(randomQuote)
	err = utils.WriteConn(conn, data)
	if err != nil {
		s.logger.Error("write conn err:", err)
		return
	}

}

func configureUserCreds(username, password string) scram.CredentialLookup {
	salt := utils.GenerateRandomSalt(16)
	saltedPassword := pbkdf2.Key([]byte(password), salt, 4096, sha256.New().Size(), sha256.New)

	clientKeyHmac := utils.ComputeHMAC(sha256.New, saltedPassword, []byte("Client Key"))
	serverKeyHmac := utils.ComputeHMAC(sha256.New, saltedPassword, []byte("Server Key"))
	storedKeyHmac := utils.ComputeHash(sha256.New, clientKeyHmac)

	return func(iUser string) (scram.StoredCredentials, error) {
		if iUser != username {
			return scram.StoredCredentials{}, errors.New("access denied")
		}
		return scram.StoredCredentials{
			KeyFactors: scram.KeyFactors{
				Salt:  string(salt),
				Iters: 4096,
			},
			ServerKey: serverKeyHmac,
			StoredKey: storedKeyHmac,
		}, nil
	}
}
