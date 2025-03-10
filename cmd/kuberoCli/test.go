package kuberoCli

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/ssh"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests and generate RSA certificates",
	Long:  `This command runs tests and generates encrypted RSA certificates with random passwords stored in a keyring.`,
	Run: func(cmd *cobra.Command, args []string) {
		runTests()
		generateRSACertificates()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func runTests() {
	fmt.Println("Running tests...")
	// Add your test logic here
}

func generateRSACertificates() {
	fmt.Println("Generating RSA certificates...")

	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Encode private key to PEM format
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Generate random password
	password := generateRandomPassword(16)

	// Store password in keyring
	err = keyring.Set("kubero-cli", "rsa-password", password)
	if err != nil {
		log.Fatalf("Failed to store password in keyring: %v", err)
	}

	// Encrypt private key with password
	encryptedPrivateKeyPEM, err := x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", privateKeyPEM, []byte(password), x509.PEMCipherAES256)
	if err != nil {
		log.Fatalf("Failed to encrypt private key: %v", err)
	}

	// Save encrypted private key to file
	err = os.WriteFile("id_rsa", pem.EncodeToMemory(encryptedPrivateKeyPEM), 0600)
	if err != nil {
		log.Fatalf("Failed to save private key: %v", err)
	}

	// Generate public key
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatalf("Failed to generate public key: %v", err)
	}

	// Save public key to file
	err = os.WriteFile("id_rsa.pub", ssh.MarshalAuthorizedKey(publicKey), 0644)
	if err != nil {
		log.Fatalf("Failed to save public key: %v", err)
	}

	fmt.Println("RSA certificates generated successfully.")
}

func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, length)
	for i := range password {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Fatalf("Failed to generate random password: %v", err)
		}
		password[i] = charset[randomIndex.Int64()]
	}
	return string(password)
}
