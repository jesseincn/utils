package encrypt

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	aesKey  = []byte("69u92Jg5SBWOgH41oB0tKY5rzTIrsjhu")
	aesData = []byte("hello ryan")
)

func TestGenRsaKey(t *testing.T) {
	private, public, err := GenRsaKey(1024)
	require.NoError(t, err)
	fmt.Println(private)
	fmt.Println(public)
}

func TestRsaEncrypt(t *testing.T) {
	private, public, err := GenRsaKey(4096)
	require.NoError(t, err)
	t.Log(private)
	t.Log(public)

	src := "{\"phone\":\"86-12345678900\",\"amount\":\"500\",\"msgCode\":\"236745\",\"timeStamp\":\"1551095160653\"}"
	enc, err := RsaEncrypt(public, src)
	require.NoError(t, err)
	str := base64.StdEncoding.EncodeToString(enc[:])
	ecd, _ := base64.StdEncoding.DecodeString(str)
	t.Log("enc: ", str)
	dec, _ := RsaDecrypt(private, string(ecd[:]))
	t.Log("dec: ", string(dec[:]))

	assert.Equal(t, string(dec[:]), src)
}

func TestAesEncryptCBC(t *testing.T) {
	encrypted, hexStr, base64Str := AesEncryptCBC(aesData, aesKey)
	t.Log("encrypted: ", string(encrypted))
	t.Log("hexStr: ", hexStr)
	t.Log("base64Str: ", base64Str)
}

func TestAesDecryptCBC(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString("Z80Me/eOk3lSGm94V4bSwg==")
	require.NoError(t, err)
	str := AesDecryptCBC(data, aesKey)

	t.Log("str: ", string(str))
	t.Log("origin data: ", string(aesData))
	assert.Equal(t, string(str), string(aesData))
}

func TestAesEncryptECB(t *testing.T) {
	encrypted, hexStr, base64Str := AesEncryptECB(aesData, aesKey)
	t.Log("encrypted: ", string(encrypted))
	t.Log("hexStr: ", hexStr)
	t.Log("base64Str: ", base64Str)
}

func TestAesDecryptECB(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString("D4f3s+AOimtIbUgdQ4Dx+g==")
	require.NoError(t, err)
	str := AesDecryptECB(data, aesKey)

	t.Log("str: ", string(str))
	t.Log("origin data: ", string(aesData))
	assert.Equal(t, string(str), string(aesData))
}

func TestAesEncryptCFB(t *testing.T) {
	encrypted, hexStr, base64Str := AesEncryptCFB(aesData, aesKey)
	t.Log("encrypted: ", string(encrypted))
	t.Log("hexStr: ", hexStr)
	t.Log("base64Str: ", base64Str)
}

func TestAesDecryptCFB(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString("GmfndeMTcs8LYIkc5jzMVD5W3HdqbDjRTMI=")
	require.NoError(t, err)
	str := AesDecryptCFB(data, aesKey)

	t.Log("str: ", string(str))
	t.Log("origin data: ", string(aesData))
	assert.Equal(t, string(str), string(aesData))
}
