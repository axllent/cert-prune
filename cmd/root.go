// Package cmd handles the CLI frontend
package cmd

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	certPath       = "/etc/letsencrypt"
	expireDays     int
	log            *logrus.Logger
	verboseLogging bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cert-prune [path]",
	Short: "A utility to delete expired Let's Encrypt certficates",
	Long: `A utility to delete expired Let's Encrypt certficates.

All unused certificates, and (by default) all csrs & keys older than 60 days are deleted.

If no path is provided then /etc/letsencrypt is assumed.

Support:
  https://github.com/axllent/cert-prune`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var nrCerts, nrCSRs, nrKeys int
		if verboseLogging {
			log.SetLevel(logrus.DebugLevel)
		}

		if len(args) == 1 {
			certPath = args[0]
		}

		if !verifyCertPath(certPath) {
			log.Errorf("Path \"%s\" does not look like a Let's Encrypt folder\n", certPath)
			os.Exit(1)
		}

		keepCerts := make(map[string]bool)

		pems, err := filepath.Glob(filepath.Join(certPath, "live/*/*.pem"))
		if err != nil {
			panic(err)
		}

		for _, p := range pems {
			l, err := filepath.EvalSymlinks(p)
			if err != nil {
				log.Error(err)
			}

			keepCerts[l] = true
		}

		archives, err := filepath.Glob(filepath.Join(certPath, "archive/*/*.pem"))
		if err != nil {
			log.Error(err)
		}

		for _, a := range archives {
			_, keep := keepCerts[a]
			if !keep {
				log.Debugf("deleting %s\n", a)
				if err := os.Remove(a); err != nil {
					log.Error(err)
					continue
				}

				nrCerts++
			}
		}

		csrs, err := filepath.Glob(filepath.Join(certPath, "csr/*.pem"))
		if err != nil {
			panic(err)
		}

		for _, csr := range csrs {
			info, err := os.Stat(csr)
			if err != nil {
				log.Error(err)
				continue
			}

			if time.Now().Sub(info.ModTime()) > time.Duration(expireDays)*24*time.Hour {
				log.Debugf("deleting %s\n", csr)
				if err := os.Remove(csr); err != nil {
					log.Error(err)
					continue
				}

				nrCSRs++
			}
		}

		keys, err := filepath.Glob(filepath.Join(certPath, "keys/*.pem"))
		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			info, err := os.Stat(key)
			if err != nil {
				log.Error(err)
				continue
			}

			if time.Now().Sub(info.ModTime()) > time.Duration(expireDays)*24*time.Hour {
				log.Debugf("deleting %s\n", key)
				if err := os.Remove(key); err != nil {
					log.Error(err)
					continue
				}

				nrKeys++
			}
		}

		log.Infof("Certs deleted:   %d", nrCerts)
		log.Infof("CSRs  deleted:   %d", nrCSRs)
		log.Infof("Keys  deleted:   %d", nrKeys)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// hide autocompletion
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	// hide help command
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	// hide help flag
	rootCmd.PersistentFlags().BoolP("help", "h", false, "This help")
	rootCmd.PersistentFlags().Lookup("help").Hidden = true

	rootCmd.PersistentFlags().IntVarP(&expireDays, "nr-days", "n", 60, "Delete generation CSRs and Keys older than X days")
	rootCmd.PersistentFlags().BoolVarP(&verboseLogging, "verbose", "v", false, "Verbose logging")

	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)

	log.Out = os.Stdout

	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})
}

func verifyCertPath(p string) bool {
	if !isDir(p) {
		return false
	}

	for _, i := range []string{"live", "archive", "csr", "keys"} {
		if !isDir(path.Join(p, i)) {
			return false
		}
	}

	return true
}

// IsDir returns whether a path is a directory
func isDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.IsDir() {
		return false
	}

	return true
}
