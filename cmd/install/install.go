package install

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/ssh"
)

const (
	keySize       = 2048
	keyringService = "kubero-cli"
	keyringUser    = "kubero"
)

// InstallCmd represents the install command
type InstallCmd struct {
	// Add fields as needed
}

// NewInstallCmd creates a new InstallCmd
func NewInstallCmd() *InstallCmd {
	return &InstallCmd{}
}

// Execute runs the install command
func (cmd *InstallCmd) Execute() error {
	privateKey, publicKey, err := generateRSAKeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	password, err := generateRandomPassword()
	if err != nil {
		return fmt.Errorf("failed to generate random password: %w", err)
	}

	err = storePasswordInKeyring(password)
	if err != nil {
		return fmt.Errorf("failed to store password in keyring: %w", err)
	}

	err = savePrivateKeyToFile(privateKey, password)
	if err != nil {
		return fmt.Errorf("failed to save private key to file: %w", err)
	}

	err = savePublicKeyToFile(publicKey)
	if err != nil {
		return fmt.Errorf("failed to save public key to file: %w", err)
	}

	fmt.Println("Installation completed successfully.")
	return nil
}

func generateRSAKeyPair() (*rsa.PrivateKey, ssh.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

func generateRandomPassword() (string, error) {
	const passwordLength = 16
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	password := make([]byte, passwordLength)
	_, err := rand.Read(password)
	if err != nil {
		return "", err
	}

	for i := range password {
		password[i] = charset[int(password[i])%len(charset)]
	}

	return string(password), nil
}

func storePasswordInKeyring(password string) error {
	return keyring.Set(keyringService, keyringUser, password)
}

func savePrivateKeyToFile(privateKey *rsa.PrivateKey, password string) error {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	privateKeyFile, err := os.Create("private_key.pem")
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	return pem.Encode(privateKeyFile, privateKeyBlock)
}

func savePublicKeyToFile(publicKey ssh.PublicKey) error {
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicKey)

	publicKeyFile, err := os.Create("public_key.pub")
	if err != nil {
		return err
	}
	defer publicKeyFile.Close()

	_, err = publicKeyFile.Write(publicKeyBytes)
	return err
}
