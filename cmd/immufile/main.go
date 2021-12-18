package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"immu/immufile/pkg/hashtree"
)

const (
	Verified     = "verified"
	VerifiedLong  = "verified: total equality"
	Tampered              = "tampered"
	TamperedLine       = "tampered: file hash is valid, but the provided line content are not equal"
	TamperedLineInTree = "tampered: exactly this line was tampered in the file"
	TamperedTree = "tampered: line content is valid, but the file content was potentially tampered"
	TamperedBoth = "tampered: file was tampered as well as line content"
)

var (
	filePath string
	lineContent string
	lineNumber int
	fileHash string
	short bool
)

func init() {
	flag.StringVar(&filePath, "f", "", "Path to a file")
	flag.StringVar(&lineContent, "c", "", "Content of a line to verify for tempering")
	flag.IntVar(&lineNumber, "l", -1, "Line number")
	flag.StringVar(&fileHash,"d", "", "File hash")
	flag.BoolVar(&short,"short", false, "Short answer tampered|verified")
}

func main() {
	flag.Parse()

	if len(strings.TrimSpace(filePath)) == 0 {
		printHelpAndExit("file was not provided")
	}

	switch flag.NFlag() {
	case 1:
		hash, err := calculateFileHash(filePath)
		if err != nil {
			printHelpAndExit(err.Error())
		}
		fmt.Println(hash)
	case 4, 5:
		result, err := checkLineContent()
		if err != nil {
			printHelpAndExit(err.Error())
		}
		fmt.Println(result)
	default:
		printHelpAndExit(fmt.Sprintf("some parameters were not provided: %d flags previded, while accepting %d, %d or %d", flag.NFlag(), 1, 4, 5))
	}
}

func printHelpAndExit(msg string) {
	fmt.Println(msg)
	flag.PrintDefaults()
	os.Exit(1)
}

func calculateFileHash(filePath string) (string, error) {
	hashTree, err := buildHashTree(filePath)
	if err != nil {
		return "", err
	}
	return hashtree.ToHashString(hashTree.Hash()), nil
}

func buildHashTree(filePath string) (hashtree.Tree, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return hashtree.Tree{}, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error during file closing: ", err)
		}
	}()

	reader := bufio.NewReader(file)
	return hashtree.ByLinesFromReader(reader)
}

func checkLineContent() (string, error) {
	if lineNumber == 0 {
		return "", fmt.Errorf("line number %d is negative", lineNumber)
	}

	providedHash, err := hashtree.HashFromString(fileHash)
	if err != nil {
		return "", err
	}

	tree, err := buildHashTree(filePath)
	if err != nil {
		return "", err
	}

	fileLineContentHash, err := tree.LeafHashAt(lineNumber)
	if err != nil {
		return "", err
	}

	lineContentHash := hashtree.CalculateHash([]byte(lineContent))

	if tree.Hash() == providedHash {
		if fileLineContentHash == lineContentHash {
			if short {
				return Verified, nil
			}
			return VerifiedLong, nil
		} else {
			if short {
				return Tampered, nil
			}
			return TamperedLine, nil
		}
	} else {
		if fileLineContentHash == lineContentHash {
			if short {
				return Tampered, nil
			}
			return TamperedTree, nil
		} else {
			newTree, err := tree.UpdateLeafHashAt(lineNumber, lineContentHash)
			if err != nil {
				return "", err
			}

			if newTree.Hash() == providedHash {
				if short {
					return Tampered, nil
				}
				return TamperedLineInTree, nil
			} else {
				if short {
					return Tampered, nil
				}
				return TamperedBoth, nil
			}
		}
	}
}
