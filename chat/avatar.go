package main

import "errors"

//ErrNoAvatarURL is error that is returned when Avatar instance is unabled to provide an avatar URL
var ErrNoAvatarURL = errors.New("chat: unable to get an avatar URL")

// Avatar represents types capable of represengint user profile pictures
type Avatar interface {
	// gets the avatar URL for the specified client,
	// or returns ErrNoAvatarURL if object is unable to get a URL for the specified client
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}
