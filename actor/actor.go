package actor

import (
	"context"
	"github.com/kunal-sanghvi/git2consul-go/fs"
	"github.com/kunal-sanghvi/git2consul-go/git"
	"github.com/kunal-sanghvi/git2consul-go/parser"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	rootDirectory = "/tmp/go-git2consul/"
)

var (
	wg sync.WaitGroup
)

func Git2Consul(configFile, pemFile string) {
	// Setup background context for app
	ctx := context.Background()

	// Create root directory for all operations
	log.Print("Creating root directory")
	fileSystem := fs.NewLocalFS(rootDirectory)
	if err := fileSystem.CreateStorage(); err != nil {
		log.Printf("Storage already exists, cleaning up")
		if err = fileSystem.CleanUp(); err != nil {
			log.Fatalf("Failed to create storage. ERROR:: %v", err)
		}
		log.Printf("Clean up complete, now creating new storage")
		if err = fileSystem.CreateStorage(); err != nil {
			log.Fatalf("Failed to create storage. ERROR:: %v", err)
		}
	}
	log.Print("Root directory created")

	confParser := parser.NewJSONParser(rootDirectory, configFile)
	config, err := confParser.ParseFromFile()
	if err != nil {
		log.Fatalf("Failed to parse config from file %s. ERROR:: %v", configFile, err)
	}

	defer func() {
		log.Print("Cleaning up root directory")
		if err = fileSystem.CleanUp(); err != nil {
			log.Fatalf("Failed to clean up. ERROR: %v", err)
		}
		log.Print("Cleanup complete")
	}()

	gitController := git.NewGitOps(fileSystem, pemFile)

	for _, repo := range config.Repos {
		log.Printf("Checking out %s consul repo", repo.Name)
		gitRepo, err := gitController.PlainCloneCtx(ctx, repo.URL, repo.Name)
		if err != nil {
			log.Fatalf("Failed to clone %s consul repo. ERROR :%v", repo.Name, err)
		}
		if err = gitController.Fetch(ctx, gitRepo); err != nil {
			log.Fatalf("Failed to fetch %s consul repo. ERROR :%v", repo.Name, err)
		}
		for _, branch := range repo.Branches {
			if err = gitController.Checkout(gitRepo, branch); err != nil {
				log.Fatalf("Failed to checkout %s branch for %s consul repo. ERROR: %v", branch, repo.Name, err)
			}
			SweepRepository(ctx, branch, repo, confParser)
		}
	}
}

func SweepRepository(ctx context.Context, branch string, repo *parser.Repo, confParser parser.Parser) {
	if err := os.Chdir(rootDirectory + repo.Name); err != nil {
		log.Fatalf("Failed to parse %s consul repo. ERROR: %v", repo.Name, err)
	}
	files := make([]string, 0)
	if err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Failed to parse %s path in %s consul repo. ERROR: %s", path, repo.Name, err)
		}
		if info.IsDir() && strings.HasPrefix(info.Name(), ".git") {
			log.Printf("Ignoring dot directory %s", info.Name())
			return filepath.SkipDir
		}

		if info.IsDir() && strings.HasPrefix(info.Name(), ".idea") {
			log.Printf("Ignoring dot directory %s", info.Name())
			return filepath.SkipDir
		}

		if info.IsDir() && strings.HasPrefix(info.Name(), "message") {
			log.Printf("Ignoring dot directory %s", info.Name())
			return filepath.SkipDir
		}
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".md") || strings.HasPrefix(info.Name(), ".git")) {
			log.Printf("Ignoring %s file", info.Name())
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		log.Fatalf("Failed to recursively parse %s consul repo. ERROR: %v", repo.Name, err)
	}
	for _, file := range files {
		wg.Add(1)
		log.Printf("Populating from %s", file)
		go populateFromFile(ctx, branch, repo.Name, file, confParser)
	}
	wg.Wait()
}

func populateFromFile(ctx context.Context, branch, repo, filePath string, confParser parser.Parser) {
	defer wg.Done()
	if err := confParser.ParseConfigFile(ctx, branch, repo, filePath); err != nil {
		log.Fatalf("Failed to parse %s. ERROR: %v", filePath, err)
	}
	log.Printf("Done with %s", filePath)
}