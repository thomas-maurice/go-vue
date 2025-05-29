package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/spf13/cobra"
)

var genKeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "Generates a public/private api key pair for token signing",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		privKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return err
		}

		pubKey := privKey.PublicKey
		encodedPriv, err := x509.MarshalECPrivateKey(privKey)
		if err != nil {
			return err
		}
		encodedPub, err := x509.MarshalPKIXPublicKey(&pubKey)
		if err != nil {
			return err
		}

		pemBlock := pem.EncodeToMemory(&pem.Block{Type: "ECDSA PRIVATE KEY", Bytes: encodedPriv})
		fmt.Println(string(pemBlock))
		pemBlock = pem.EncodeToMemory(&pem.Block{Type: "ECDSA PUBLIC KEY", Bytes: encodedPub})
		fmt.Println(string(pemBlock))
		return nil
	},
}

func initGenKeyCmd() {

}
