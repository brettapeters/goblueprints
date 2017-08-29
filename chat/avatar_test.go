package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := &client{}

	// no avatar value, should throw an error
	url, err := authAvatar.GetAvatarURL(client)
	if g, w := err, ErrNoAvatarURL; g != w {
		t.Errorf("\nAuthAvatar.GetAvatarURL returned: %v\nWant %v", g, w)
	}

	// set a value
	testURL := "http://url-to-gravatar/"
	client.userData = map[string]interface{}{"avatar_url": testURL}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	} else {
		if g, w := url, testURL; g != w {
			t.Errorf("\nAuthAvatar.GetAvatarURL returned: %v\nWant %v", g, w)
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := &client{}
	client.userData = map[string]interface{}{
		"userid": "0bc83cb571cd1c50ba6f3e8a78ef1346",
	}

	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}

	if g, w := url, "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346"; g != w {
		t.Errorf("\nGravatarAvatar.GetAvatarURL returned: %v\nWant %v", g, w)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	// make a test avatar file
	userID := "abc123"
	filepath, close := tmpAvatarFile(userID)
	defer close()

	var fileSystemAvatar FileSystemAvatar
	client := &client{}
	client.userData = map[string]interface{}{"userid": userID}

	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
	}

	if g, w := url, filepath; g != w {
		t.Errorf("\nFileSystemAvatar.GetAvatarURL returned %v\nWant %v", g, w)
	}
}

func tmpAvatarFile(userID string) (string, func()) {
	filepath := "/avatars/" + userID + ".jpg"
	ioutil.WriteFile(filepath, []byte{}, 0777)
	return filepath, func() { os.Remove(filepath) }
}
