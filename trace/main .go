package main

import (
	"fmt"
	"os"
)

func main() {
	contenu := []byte("Premiere trace\n")

	err := os.WriteFile("sortie.txt", contenu, 0644)

	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}

	fmt.Println("Fichier crée")
}
