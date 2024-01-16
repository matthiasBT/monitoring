package adapters

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"net"
)

var ErrResponseNotOK = errors.New("response not OK")

func getIPFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	return ip
}

func encryptData(payload []byte, rsaKey *rsa.PublicKey) ([]byte, error) {
	key, encryptedPayload, err := encryptAES(payload)
	if err != nil {
		return nil, err
	}
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaKey, key, nil)
	if err != nil {
		return nil, err
	}
	return append(encryptedKey, encryptedPayload...), nil
}

func encryptAES(plaintext []byte) ([]byte, []byte, error) {
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return key, ciphertext, nil
}

func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue // interface down or loopback
		}

		// Get associated unicast interface addresses.
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		// Iterate over the addresses looking for a non-loopback IPv4 address.
		for _, addr := range addrs {
			ip := getIPFromAddr(addr)
			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				return ip.String(), nil // return the first non-loopback IPv4 address
			}
		}
	}
	return "", fmt.Errorf("no active network interface found")
}
