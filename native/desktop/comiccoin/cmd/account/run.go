package account

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	account_http "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/interface/http"
	account_httphandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/interface/http/handler"
	httpmiddle "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/interface/http/middleware"
	account_repo "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/repo"
	account_s "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/service"
	account_usecase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the account",
	Run: func(cmd *cobra.Command, args []string) {
		//
		// STEP 1
		// Load up our dependencies and configuration
		//

		// Load up our operating system interaction handlers, more specifically
		// signals. The OS sends our application various signals based on the
		// OS's state, we want to listen into the termination signals.
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGUSR1)

		// bootstrapPeers, err := StringToAddres(flagBootstrapPeers)
		// if err != nil {
		// 	log.Fatalf("Failed converting string to multi-addresses: %v\n", err)
		// }

		cfg := &config.Config{
			// Blockchain: config.BlockchainConfig{
			// 	ChainID:       constants.ChainIDMainNet,
			// 	TransPerBlock: 1,
			// 	Difficulty:    2,
			// },
			// App: config.AppConfig{
			// 	HTTPPort: flagListenHTTPPort,
			// 	HTTPIP:   flagListenHTTPIP,
			// 	DirPath:  flagDataDir,
			// },
			// Peer: config.PeerConfig{
			// 	ListenPort:       flagListenPeerToPeerPort,
			// 	KeyName:          flagKeypairName,
			// 	RendezvousString: flagRendezvousString,
			// 	BootstrapPeers:   bootstrapPeers,
			// },
			DB: config.DBConfig{
				DataDir: flagDataDir,
			},
		}
		logger := logger.NewLogger()
		db := dbase.NewDatabase(cfg, logger)
		accountRepo := account_repo.NewAccountRepo(cfg, logger, db)
		createAccountUseCase := account_usecase.NewCreateAccountUseCase(cfg, logger, accountRepo)
		createAccountService := account_s.NewCreateAccountService(cfg, logger, createAccountUseCase)
		createAccountHTTPHandler := account_httphandler.NewCreateAccountHTTPHandler(cfg, logger, createAccountService)
		httpMiddleware := httpmiddle.NewMiddleware(cfg, logger)
		httpServ := account_http.NewHTTPServer(cfg, logger, httpMiddleware, createAccountHTTPHandler)

		//
		// STEP 3
		// Run the main loop blocking code while other input ports run in
		// background.
		//

		// Run in background the peer to peer node which will synchronize our
		// blockchain with the network.
		// go peerNode.Run()
		go httpServ.Run()
		defer httpServ.Shutdown()

		<-done
	},
}
