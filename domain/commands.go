package domain

import (
	"time"

	"golang.org/x/oauth2"
)

type (
	TokenSource struct {
		AccessToken string
	}

	CommandSet map[Order]Command

	Command interface {
		Execute(*AutoScalingGroup) error
	}

	CMDError struct {
		Code    int
		Message string
	}

	BaseCommand struct {
		Provider Provider
		State    CommandState
		Error    *CMDError
		Timeout  time.Duration
	}

	BaseCommands []BaseCommand
)

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
