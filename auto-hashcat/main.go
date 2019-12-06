package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	log.Println()
	fmt.Println(`                 (__) 
                 (oo) 
           /------\/ 
          / |    ||   
         *  /\---/\ 
            ~~   ~~   
..."Mooing out NTLM to clear-text!"...`)
	var hashcatPath, passwordDir, hashes string
	var rules bool
	flag.StringVar(&hashcatPath, "hashcat", "", "Pass the non default hashcat path")
	flag.StringVar(&passwordDir, "dir", "", "Pass the value to the directory with .txt files")
	flag.StringVar(&hashes, "hashes", "", "File containing NTLM hashes")
	flag.BoolVar(&rules, "rules", false, "Dictionary with rules attack")
	flag.Parse()
	// Get hashcat Path
	if hashcatPath == "" {
		defaultHashcatPath, err := exec.LookPath("hashcat")
		if err != nil {
			fmt.Println("It seems as though hashcat is not installed")
			fmt.Println("Please specify the path to hashcat or # apt install hashcat")
			log.Panicln(err)
		}
		hashcatPath = defaultHashcatPath
	}
	// Chech for hash path
	if hashes == "" {
		flag.Usage()
		log.Fatalln("hashes not specified")
	} else {
		if _, err := os.Lstat(hashes); err != nil {
			flag.Usage()
			log.Fatalln("hashes file does not exist")
		}
	}
	// Check dictionary path
	if passwordDir == "" {
		flag.Usage()
		log.Fatalln("password directory not supplied")
	} else {
		if _, err := os.Lstat(passwordDir); err != nil {
			flag.Usage()
			log.Fatalln("password directory does not exist")
		}
	}
	start := time.Now()
	dictionary(hashcatPath, hashes, passwordDir)
	if rules {
		dictionaryAndRules(hashcatPath, hashes, passwordDir)
	}
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println(elapsed)
}

func dictionary(hashcatPath string, hashes string, passwordDir string) {
	log.Println("Running hashes againts dictionaries only")
	args := []string{"--force", "-o", "cracked.txt", "-a", "0", "-m", "1000", hashes, passwordDir}
	// get sub directories that contain password files
	files, err := ioutil.ReadDir(passwordDir)
	if err != nil {
		log.Fatalln(err)
	}
	// append other password directories
	for _, file := range files {
		if file.IsDir() {
			args = append(args, passwordDir+file.Name()+"/")
		}
	}
	cracker := exec.Command(hashcatPath, args...)
	cracker.Stdout = os.Stdout
	cracker.Stderr = os.Stderr
	cracker.Stdin = os.Stdin
	cracker.Run()
	log.Println("Dictionary check done")
}

func dictionaryAndRules(hashcatPath string, hashes string, passwordDir string) {
	log.Println("Running hashes against rulesets")
	ruleDir := "/usr/share/hashcat/rules/"
	rulesets := []string{"best64.rule", "combinator.rule", "leetspeak.rule", "rockyou-30000.rule"}
	fmt.Println(ruleDir, rulesets)
	// add for loop which appends rules to args
	for _, rule := range rulesets {
		args := []string{"--force", "-a", "0", "-m", "1000", hashes, passwordDir, "-r", ruleDir + rule}
		cracker := exec.Command(hashcatPath, args...)
		cracker.Stdout = os.Stdout
		cracker.Stderr = os.Stderr
		cracker.Stdin = os.Stdin
		cracker.Run()
	}
	log.Println("Dictionary + Rules check done")
}
