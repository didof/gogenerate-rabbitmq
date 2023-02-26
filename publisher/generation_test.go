package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

const generatorName = "publisher.go"

type GenerationSuite struct {
	suite.Suite
	tmpDir string
}

func (suite *GenerationSuite) createTmp() string {
	suite.T().Helper()

	var err error
	e, err := os.Executable()
	if err != nil {
		suite.T().Errorf("Error retrieving executable: %v", err)
	}

	suite.tmpDir, err = ioutil.TempDir("", filepath.Base(e))
	if err != nil {
		suite.T().Errorf("Error creating tmp dir: %v", err)
	}
	return suite.tmpDir
}

func (suite *GenerationSuite) removeTmp() {
	suite.T().Helper()

	err := os.RemoveAll(suite.tmpDir)
	if err != nil {
		suite.T().Errorf("Error removing temporary directory: %v", err)
	}
}

func (suite *GenerationSuite) copyGenerator() string {
	suite.T().Helper()

	name, err := CopyIntoDir(generatorName, suite.tmpDir)
	if err != nil {
		suite.T().Fatal(err)
	}

	return name
}

func TestGenerationSuite(t *testing.T) {
	suite.Run(t, new(GenerationSuite))
}

func (suite *GenerationSuite) SetupTest() {
	suite.createTmp()
	suite.copyGenerator()

	// TODO setup: copy into tmp dir the generator (publisher.go) and each test data file

	// TODO exec:  run in the test the cmd go generate ./<testdata-name>
	// TODO exec:  assert creation of expected file
	// TODO cleanup: remove all
}

// NOTA BENE > Since there is only one test (TestGeneration), this hook is run once. Thus, it is used to remove the tmp dir.
func (suite *GenerationSuite) AfterTest(suiteName, testName string) {
	suite.removeTmp()
}

func (suite *GenerationSuite) TestGeneration() {
	entries, err := os.ReadDir("testdata")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		name := e.Name()
		suite.Run(name, func() {
			// Copy test file into tmp dir
			path, err := CopyIntoDir(filepath.Join("testdata", name), suite.tmpDir)
			if err != nil {
				suite.T().Fatal(err)
			}

			// Schedule removal of it
			suite.T().Cleanup(func() {
				if _, err := os.ReadFile(path); err != nil {
					suite.T().Logf("Error removing %s: %v", path, err)
				}
			})

			cmd := exec.Command("go", "generate", name)
			cmd.Dir = suite.tmpDir

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				suite.T().Errorf("Error running 'go generate': %v", err)
				return
			}

			outputName := AddSuffix(name, "rabbitmq_publisher")
			outPath := filepath.Join(suite.tmpDir, outputName)

			_, err = os.Stat(outPath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					suite.T().Errorf("File has not been generated: %v", err)
					return
				}

				suite.T().Errorf("Error checking existence of %s: %v", outPath, err)
			}
		})
	}
}

func BaseWithoutExt(name string) string {
	return strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
}

func AddSuffix(base, suffix string) string {
	base = BaseWithoutExt(base)
	return fmt.Sprintf("%s_%s.go", base, suffix)
}

func CopyIntoDir(src, dst string) (string, error) {
	// Create destination file
	d, err := os.Create(filepath.Join(dst, filepath.Base(src)))
	if err != nil {
		return "", fmt.Errorf("error creating destination file: %v", err)
	}
	defer d.Close()

	// Open the source file
	s, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("error opening file %s: %v", src, err)
	}
	defer s.Close()

	// Copy the source file contents into the destination file
	_, err = io.Copy(d, s)
	if err != nil {
		return "", fmt.Errorf("error copying file contents from %s to %s: %v", s.Name(), d.Name(), err)
	}

	return d.Name(), nil
}
