package blockchain

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	blockdata_repo "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/repo"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
)

func runCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "run",
		Short: "Runs the blockchain",
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
				App: config.AppConfig{
					DirPath:     flagDataDir,
					HTTPAddress: flagListenHTTPAddress,
					RPCAddress:  flagListenRPCAddress,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}
			logger := logger.NewLogger()
			db := dbase.NewDatabase(cfg.DB.DataDir, logger)

			// Repo
			blockdataRepo := blockdata_repo.NewBlockDataRepo(cfg, logger, db)

			// // Use-case
			// createBlockchainUseCase := blockchain_usecase.NewCreateBlockchainUseCase(cfg, logger, blockchainRepo)
			// getBlockchainUseCase := blockchain_usecase.NewGetBlockchainUseCase(cfg, logger, blockchainRepo)
			// blockchainDecryptKeyUseCase := blockchain_usecase.NewBlockchainDecryptKeyUseCase(cfg, logger, blockchainRepo)
			// blockchainEncryptKeyUseCase := blockchain_usecase.NewBlockchainEncryptKeyUseCase(cfg, logger, blockchainRepo)
			//
			// // Service
			// createBlockchainService := blockchain_s.NewCreateBlockchainService(cfg, logger, createBlockchainUseCase, getBlockchainUseCase, blockchainEncryptKeyUseCase)
			// getBlockchainService := blockchain_s.NewGetBlockchainService(cfg, logger, getBlockchainUseCase)
			// getKeyService := blockchain_s.NewGetKeyService(cfg, logger, getBlockchainUseCase, blockchainDecryptKeyUseCase)
			//
			// // HTTP
			// createBlockchainHTTPHandler := blockchain_httphandler.NewCreateBlockchainHTTPHandler(cfg, logger, createBlockchainService)
			// getBlockchainHTTPHandler := blockchain_httphandler.NewGetBlockchainHTTPHandler(cfg, logger, getBlockchainService)
			// httpMiddleware := httpmiddle.NewMiddleware(cfg, logger)
			// httpServ := blockchain_http.NewHTTPServer(
			// 	cfg, logger, httpMiddleware,
			// 	createBlockchainHTTPHandler,
			// 	getBlockchainHTTPHandler,
			// )
			//
			// // RPC
			// rpcServ := blockchain_rpc.NewRPCServer(
			// 	cfg, logger,
			// 	getKeyService,
			// )
			//
			// //
			// // STEP 3
			// // Run the main loop blocking code while other input ports run in
			// // background.
			// //
			//
			// // Run in background the peer to peer node which will synchronize our
			// // blockchain with the network.
			// // go peerNode.Run()
			// go httpServ.Run()
			// go rpcServ.Run()
			// defer httpServ.Shutdown()
			// defer rpcServ.Shutdown()

			<-done
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagListenHTTPAddress, "listen-http-address", "127.0.0.1:8010", "The IP and port to attach for our HTTP JSON API server")
	cmd.Flags().StringVar(&flagListenRPCAddress, "listen-rpc-address", "localhost:8011", "The ip and port to listen to for the TCP RPC server")

	return cmd
}
