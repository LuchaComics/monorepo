package blockchain

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config/constants"
	account_http "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http"
	account_httphandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http/handler"
	httpmiddle "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http/middleware"
	account_rpc "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/rpc"
	account_repo "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/repo"
	ik_repo "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/repo"
	account_s "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
	ik_s "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
	account_usecase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	ik_use "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
)

func runCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "run",
		Short: "Runs the blockchain application",
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

			bootstrapPeers, err := StringToAddres(flagBootstrapPeers)
			if err != nil {
				log.Fatalf("Failed converting string to multi-addresses: %v\n", err)
			}

			cfg := &config.Config{
				Blockchain: config.BlockchainConfig{
					ChainID:       constants.ChainIDMainNet,
					TransPerBlock: 1,
					Difficulty:    2,
				},
				App: config.AppConfig{
					DirPath:     flagDataDir,
					HTTPAddress: flagListenHTTPAddress,
					RPCAddress:  flagListenRPCAddress,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
				Peer: config.PeerConfig{
					ListenPort:       flagListenPeerToPeerPort,
					KeyName:          flagKeypairName,
					RendezvousString: flagRendezvousString,
					BootstrapPeers:   bootstrapPeers,
				},
			}
			logger := logger.NewLogger()
			db := dbase.NewDatabase(cfg.DB.DataDir, logger)

			// Repo
			accountRepo := account_repo.NewAccountRepo(cfg, logger, db)
			ikRepo := ik_repo.NewIdentityKeyRepo(cfg, logger, db)

			// Use-case
			createAccountUseCase := account_usecase.NewCreateAccountUseCase(cfg, logger, accountRepo)
			getAccountUseCase := account_usecase.NewGetAccountUseCase(cfg, logger, accountRepo)
			accountDecryptKeyUseCase := account_usecase.NewAccountDecryptKeyUseCase(cfg, logger, accountRepo)
			accountEncryptKeyUseCase := account_usecase.NewAccountEncryptKeyUseCase(cfg, logger, accountRepo)
			ikCreateUseCase := ik_use.NewCreateIdentityKeyUseCase(cfg, logger, ikRepo)
			ikGetUseCase := ik_use.NewGetIdentityKeyUseCase(cfg, logger, ikRepo)

			// Service
			createAccountService := account_s.NewCreateAccountService(cfg, logger, createAccountUseCase, getAccountUseCase, accountEncryptKeyUseCase)
			getAccountService := account_s.NewGetAccountService(cfg, logger, getAccountUseCase)
			getKeyService := account_s.NewGetKeyService(cfg, logger, getAccountUseCase, accountDecryptKeyUseCase)
			ikCreateService := ik_s.NewCreateIdentityKeyService(cfg, logger, ikCreateUseCase, ikGetUseCase)
			ikGetService := ik_s.NewGetIdentityKeyService(cfg, logger, ikGetUseCase)

			// Get our identity key.
			ik, err := ikGetService.Execute(flagIdentityKeyID)
			if err != nil {
				log.Fatalf("Failed getting identity key: %v", err)
			}
			if ik == nil {
				log.Fatal("Failed getting identity key: d.n.e.")
			}
			logger.Debug("Identity key found")

			privateKey, _ := ik.GetPrivateKey()
			publicKey, _ := ik.GetPublicKey()
			libp2pnet := p2p.NewLibP2PNetwork(cfg, logger, privateKey, publicKey)

			//TODO
			_ = ikCreateService
			_ = ikGetService
			_ = libp2pnet

			//TODO
			// USE CASES - NETWORK
			// - Share Signed Pending Transaction (Publisher)
			// - Receive Signed Pending Transaction (Subscriber)
			// - Share Purpose Block Data (Publisher)
			// - Receive Purpose Block Data (Subscriber)
			// - Share Block Data (Publisher)
			// - Receive Block Data (Subscriber)
			// - Ask Latest Block Hash (Req-Res)
			// - Receive Latest Block Hash (Req-Res)
			// - Ask Block Data (Req-Res)
			// - Receive Block Data (Req-Res)

			// USE CASES - HTTP
			// Send/Receive Signed Pending Transaction
			// Send/Receive Purpose Block Data
			// Send/Receive Block Data
			// Send/Receive Latest Block Hash
			// Send/Receive Block Data

			// HTTP
			createAccountHTTPHandler := account_httphandler.NewCreateAccountHTTPHandler(cfg, logger, createAccountService)
			getAccountHTTPHandler := account_httphandler.NewGetAccountHTTPHandler(cfg, logger, getAccountService)
			createTransactionHTTPHandler := account_httphandler.NewCreateTransactionHTTPHandler(cfg, logger, getKeyService)
			httpMiddleware := httpmiddle.NewMiddleware(cfg, logger)
			httpServ := account_http.NewHTTPServer(
				cfg, logger, httpMiddleware,
				createAccountHTTPHandler,
				getAccountHTTPHandler,
				createTransactionHTTPHandler,
			)

			// RPC
			rpcServ := account_rpc.NewRPCServer(
				cfg, logger,
				getKeyService,
			)

			//
			// STEP 3
			// Run the main loop blocking code while other input ports run in
			// background.
			//

			// Run in background the peer to peer node which will synchronize our
			// blockchain with the network.
			// go peerNode.Run()
			go httpServ.Run()
			go rpcServ.Run()
			defer httpServ.Shutdown()
			defer rpcServ.Shutdown()

			<-done
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagListenHTTPAddress, "listen-http-address", "127.0.0.1:8000", "The IP and port to attach for our HTTP JSON API server")
	cmd.Flags().StringVar(&flagListenRPCAddress, "listen-rpc-address", "localhost:8001", "The ip and port to listen to for the TCP RPC server")
	cmd.Flags().StringVar(&flagIdentityKeyID, "identitykey-id", "", "The unique identifier  to use to lookup the identity key and assign to this peer")
	cmd.MarkFlagRequired("identitykey-id")

	return cmd
}
