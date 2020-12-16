package encrypt

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestGenRsaKey(t *testing.T) {
	private, public, err := GenRsaKey(1024)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println(private)
	fmt.Println(public)
}

func TestRsaEncrypt(t *testing.T) {
	private, public, err := GenRsaKey(4096)
	println(private)
	println(public)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	src := "{\"phone\":\"86-12345678900\",\"amount\":\"500\",\"msgCode\":\"236745\",\"timeStamp\":\"1551095160653\"}"
	enc, err := RsaEncrypt(public, src)
	str := base64.StdEncoding.EncodeToString(enc[:])
	ecd, _ := base64.StdEncoding.DecodeString(str)
	fmt.Println("enc: ", str)
	dec, _ := RsaDecrypt(private, string(ecd[:]))
	fmt.Println("dec: ", string(dec[:]))
	if string(dec[:]) != src {
		t.Errorf(err.Error())
		return
	}
}

func TestAesDecrypt(t *testing.T) {
	iv, err := base64.StdEncoding.DecodeString("JqCvx/OxR/MN4REmBGDJxQ==")
	key := []byte("69u92Jg5SBWOgH41oB0tKY5rzTIrsjhu")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	str, err := AesDecryptIv("ns6Mi2GyaFJJ8Z3HiUwLs3kjShAqQepSyKRFzNv1FViXgOpccwFH6Gab1MSkyZ25OSoCWuQV1a5c+pHG/ZnJnCcoGqfPbkWeMcLpaLYhh0M=", key, iv)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if string(str) == "a8971729fbc199fb3459529cebcd8704791fc699d88ac89284f23ff8e7fca7d6" {
		fmt.Println("ok")
	}
	fmt.Println(string(str))
}
