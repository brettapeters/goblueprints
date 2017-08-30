package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
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
	GetAvatarURL(ChatUser) (string, error)
}

// TryAvatars combines all Avatar implementations
// and returns a value when one is found
type TryAvatars []Avatar

// GetAvatarURL loops over a slice of Avatar objects and
// returns an avatar url when one is found
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

// AuthAvatar implements Avatar using the avatar url
// provided by the auth service
type AuthAvatar struct{}

// UseAuthAvatar is an instance of AuthAvatar
var UseAuthAvatar AuthAvatar

// GetAvatarURL returns the avatar_url from the
// client userData map, or an error if none is found
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if len(u.AvatarURL()) == 0 {
		return "", ErrNoAvatarURL
	}
	return u.AvatarURL(), nil
}

// GravatarAvatar implements Avatar using the Gravatar service
type GravatarAvatar struct{}

// UseGravatar is an instance of GravatarAvatar
var UseGravatar GravatarAvatar

// GetAvatarURL returns a Gravatar url for the passed client
func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if len(u.UniqueID()) == 0 {
		return "", ErrNoAvatarURL
	}
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// FileSystemAvatar implements Avatar using image files
// from the local filesystem
type FileSystemAvatar struct{}

// UseFileSystemAvatar is an instance of FileSystemAvatar
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL returns the path to an image file for a passed client
func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if len(u.UniqueID()) == 0 {
		return "", ErrNoAvatarURL
	}
	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		return "", ErrNoAvatarURL
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fname := file.Name()
		if u.UniqueID() == strings.TrimSuffix(fname, filepath.Ext(fname)) {
			return "avatars/" + fname, nil
		}
	}
	return "", ErrNoAvatarURL
}
