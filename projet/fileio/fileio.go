// Package fileio lit et écrit les fichiers texte.
package fileio

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
)

// Lire renvoie le contenu du fichier sous forme de texte.
func Lire(chemin string) (string, error) {
	f, err := os.Open(chemin)
	if errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("fichier introuvable : %s", chemin)
	}
	if err != nil {
		return "", fmt.Errorf("lecture de %s : %w", chemin, err)
	}
	defer f.Close()
	contenu, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("lecture de %s : %w", chemin, err)
	}
	return string(contenu), nil
}

// Ecrire écrit le texte dans le fichier (créé ou remplacé).
func Ecrire(chemin string, texte string) error {
	f, err := os.Create(chemin)
	if err != nil {
		return fmt.Errorf("écriture de %s : %w", chemin, err)
	}
	defer f.Close()
	_, err = f.WriteString(texte)
	if err != nil {
		return fmt.Errorf("écriture de %s : %w", chemin, err)
	}
	return nil
}
