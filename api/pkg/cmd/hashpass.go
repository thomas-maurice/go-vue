package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

var hashPassCmd = &cobra.Command{
	Use:   "hashpass",
	Short: "Hashes a password for a user",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		hash, err := bcrypt.GenerateFromPassword(password, 12)
		if err != nil {
			panic(err)
		}
		fmt.Println()
		fmt.Println(string(hash))

		return nil
	},
}

func initHashPassCmd() {

}
