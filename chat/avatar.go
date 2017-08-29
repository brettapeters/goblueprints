package main

import (
	"errors"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: unable to get an avatar url")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar implements Avatar using the avatar url
// provided by the auth service
type AuthAvatar struct{}

// UseAuthAvatar is an instance of AuthAvatar
var UseAuthAvatar AuthAvatar

// GetAvatarURL returns the avatar_url from the
// client userData map, or an error if none is found
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar implements Avatar using the Gravatar service
type GravatarAvatar struct{}

// UseGravatar is an instance of GravatarAvatar
var UseGravatar GravatarAvatar

// GetAvatarURL returns a Gravatar url for the passed client
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// FileSystemAvatar implements Avatar using image files
// from the local filesystem
type FileSystemAvatar struct{}

// UseFileSystemAvatar is an instance of FileSystemAvatar
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL returns the path to an image file for a passed client
func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "/avatars/" + useridStr + ".jpg", nil
		}
	}
	return "", ErrNoAvatarURL
}
