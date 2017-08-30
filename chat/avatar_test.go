package main

import (
	"io/ioutil"
	"os"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}

	// no avatar value, should throw an error
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}

	// set a value
	testURL := "http://url-to-avatar"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testURL, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}

	if g, w := url, testURL; g != w {
		t.Errorf("\nAuthAvatar.GetAvatarURL returned: %v\nWant %v", g, w)
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc123"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}

	if g, w := url, "//www.gravatar.com/avatar/abc123"; g != w {
		t.Errorf("\nGravatarAvitar.GetAvatarURL returned %s\nWant %s", g, w)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	// make a test avatar file
	userID := "abc123"
	filepath, close := tmpAvatarFile(userID)
	defer close()

	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: userID}

	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
	}

	if g, w := url, filepath; g != w {
		t.Errorf("\nFileSystemAvatar.GetAvatarURL returned %v\nWant %v", g, w)
	}
}

func tmpAvatarFile(userID string) (string, func()) {
	filepath := "avatars/" + userID + ".jpg"
	ioutil.WriteFile(filepath, []byte{}, 0777)
	return filepath, func() { os.Remove(filepath) }
}
