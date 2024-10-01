package daemon

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config/constants"
	http "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http"
	httphandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http/handler"
	httpmiddle "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http/middleware"
	rpc "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/rpc"
	taskmng "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/task"
	taskmnghandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/task/handler"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/repo"
	service "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
	usecase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
)

func DaemonCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Run a blockchain node",
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
			kmutex := kmutexutil.NewKMutexProvider()

			// ------ Peer-to-Peer (P2P) ------
			ikRepo := repo.NewIdentityKeyRepo(cfg, logger, db)
			ikCreateUseCase := usecase.NewCreateIdentityKeyUseCase(cfg, logger, ikRepo)
			ikGetUseCase := usecase.NewGetIdentityKeyUseCase(cfg, logger, ikRepo)
			ikCreateService := service.NewCreateIdentityKeyService(cfg, logger, ikCreateUseCase, ikGetUseCase)
			ikGetService := service.NewGetIdentityKeyService(cfg, logger, ikGetUseCase)

			// If nothing was set then we use a default value. We do this to
			// simplify the user's experience.
			if flagIdentityKeyID == "" {
				flagIdentityKeyID = constants.DefaultIdentityKeyID
			}

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
			libP2PNetwork := p2p.NewLibP2PNetwork(cfg, logger, privateKey, publicKey)
			h := libP2PNetwork.GetHost()

			// Build host multiaddress
			hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", h.ID()))

			// Now we can build a full multiaddress to reach this host
			// by encapsulating both addresses:
			addr := h.Addrs()[0]
			fullAddr := addr.Encapsulate(hostAddr)

			logger.Info("Blockchain node ready",
				slog.Any("peer identity", h.ID()),
				slog.Any("full address", fullAddr),
			)

			//TODO
			_ = ikCreateService
			_ = ikGetService

			//TODO
			// USE CASES - NETWORK
			// - ✅ Share Signed Pending Transaction (Publisher)
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

			// ------ Repo ------
			accountRepo := repo.NewAccountRepo(cfg, logger, db)
			signedTxRepo := repo.NewSignedTransactionRepo(cfg, logger, db)
			signedTxDTORepo := repo.NewSignedTransactionDTORepo(cfg, logger, libP2PNetwork)

			// ------ Use-case ------
			createAccountUseCase := usecase.NewCreateAccountUseCase(cfg, logger, accountRepo)
			getAccountUseCase := usecase.NewGetAccountUseCase(cfg, logger, accountRepo)
			accountDecryptKeyUseCase := usecase.NewAccountDecryptKeyUseCase(cfg, logger, accountRepo)
			accountEncryptKeyUseCase := usecase.NewAccountEncryptKeyUseCase(cfg, logger, accountRepo)
			broadcastSignedTxDTOUseCase := usecase.NewBroadcastSignedTransactionDTOUseCase(cfg, logger, signedTxDTORepo)
			receiveSignedTxDTOUseCase := usecase.NewReceiveSignedTransactionDTOUseCase(cfg, logger, signedTxDTORepo)
			createSignedTransactionUseCase := usecase.NewCreateSignedTransactionUseCase(cfg, logger, signedTxRepo)

			// ------ Service ------
			createAccountService := service.NewCreateAccountService(cfg, logger, createAccountUseCase, getAccountUseCase, accountEncryptKeyUseCase)
			getAccountService := service.NewGetAccountService(cfg, logger, getAccountUseCase)
			getKeyService := service.NewGetKeyService(cfg, logger, getAccountUseCase, accountDecryptKeyUseCase)
			_ = getKeyService
			createTxService := service.NewCreateTransactionService(cfg, logger, getAccountUseCase, accountDecryptKeyUseCase, broadcastSignedTxDTOUseCase)
			mempoolReceiveSignedTxService := service.NewMempoolReceiveService(cfg, logger, kmutex, receiveSignedTxDTOUseCase, createSignedTransactionUseCase)

			// ------ Interface ------
			// HTTP
			createAccountHTTPHandler := httphandler.NewCreateAccountHTTPHandler(cfg, logger, createAccountService)
			getAccountHTTPHandler := httphandler.NewGetAccountHTTPHandler(cfg, logger, getAccountService)
			createTransactionHTTPHandler := httphandler.NewCreateTransactionHTTPHandler(cfg, logger, createTxService)
			httpMiddleware := httpmiddle.NewMiddleware(cfg, logger)
			httpServ := http.NewHTTPServer(
				cfg, logger, httpMiddleware,
				createAccountHTTPHandler,
				getAccountHTTPHandler,
				createTransactionHTTPHandler,
			)

			// RPC
			rpcServ := rpc.NewRPCServer(
				cfg, logger,
				getKeyService,
			)

			// TASK MANAGER
			th1 := taskmnghandler.NewMempoolReceiveTaskHandler(cfg, logger, mempoolReceiveSignedTxService)
			taskManager := taskmng.NewTaskManager(
				cfg,
				logger,
				th1,
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
			go taskManager.Run()
			defer httpServ.Shutdown()
			defer rpcServ.Shutdown()
			defer taskManager.Shutdown()

			<-done
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagListenHTTPAddress, "listen-http-address", "127.0.0.1:8000", "The IP and port to attach for our HTTP JSON API server")
	cmd.Flags().StringVar(&flagListenRPCAddress, "listen-rpc-address", "localhost:8001", "The ip and port to listen to for the TCP RPC server")
	cmd.Flags().StringVar(&flagIdentityKeyID, "identitykey-id", "", "If you would like to use a custom identity then this is the identifier used to lookup a custom identity profile to assign for this blockchain node.")
	cmd.Flags().IntVar(&flagListenPeerToPeerPort, "listen-p2p-port", 26642, "The port to listen to for other peers")
	cmd.Flags().StringVar(&flagBootstrapPeers, "bootstrap-peers", "", "The list of peers used to synchronize our blockchain with")

	return cmd
}
