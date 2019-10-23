package main
//
//import (
//	"fmt"
//	"os"
//)
//
//func LoadFromFile(fp string) error {
//	fmt.Printf("Opening file: %s\n", fp)
//	f, err := os.OpenFile(fp, os.O_RDONLY, 0)
//	if err != nil {
//		return fmt.Errorf("could not open file: %v", err)
//	}
//	defer f.Close()
//
//	err = r.load(f)
//	if err != nil {
//		return fmt.Errorf("error decoding csv: %v\n", err)
//	}
//	return nil
//}
