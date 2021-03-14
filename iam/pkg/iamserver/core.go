//

package iamserver

import (
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alloyzeus/go-azfl/azfl/errors"
	lru "github.com/hashicorp/golang-lru"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/kadisoka/kadisoka-framework/foundation/pkg/app"
	mediastore "github.com/kadisoka/kadisoka-framework/foundation/pkg/media/store"
	"github.com/kadisoka/kadisoka-framework/foundation/pkg/realm"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iam"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/eav10n"
	"github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n"

	// SMS delivery service providers
	_ "github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n/nexmo"
	_ "github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n/telesign"
	_ "github.com/kadisoka/kadisoka-framework/iam/pkg/iamserver/pnv10n/twilio"

	// Media object storage modules
	_ "github.com/kadisoka/kadisoka-framework/foundation/pkg/media/store/local"
	_ "github.com/kadisoka/kadisoka-framework/foundation/pkg/media/store/minio"
	_ "github.com/kadisoka/kadisoka-framework/foundation/pkg/media/store/s3"
)

const secretFilesDir = "/run/secrets"

type Core struct {
	realmInfo realm.Info
	db        *sqlx.DB

	registeredUserInstanceIDCache *lru.ARCCache
	deletedUserInstanceIDCache    *lru.ARCCache

	iam.ServiceClient //TODO: not specifically client

	applicationDataProvider iam.ApplicationDataProvider
	mediaStore              *mediastore.Store

	eaVerifier *eav10n.Verifier
	pnVerifier *pnv10n.Verifier
}

// RealmInfo returns information about the realm this service is serving
// for.
func (core Core) RealmInfo() realm.Info { return core.realmInfo }

// RealmName returns the name of the realm of this service.
func (core Core) RealmName() string { return core.realmInfo.Name }

// NewCoreByConfig creates an instance of Core designed for use
// in identity provider services.
func NewCoreByConfig(
	coreCfg CoreConfig,
	appApp app.App,
	realmInfo realm.Info,
) (*Core, error) {
	realmName := realmInfo.Name

	iamDB, err := connectPostgres(coreCfg.DBURL)
	if err != nil {
		return nil, errors.Wrap("DB connection", err)
	}

	//TODO: get from secret storage (e.g., vault or AWS secret manager)
	jwtPrivateKeyFilenames := []string{
		filepath.Join(secretFilesDir, "jwt_ed25519.key"),
		filepath.Join(secretFilesDir, "jwt_rsa.key"),
		filepath.Join(secretFilesDir, "jwt.key"),
	}
	jwtKeyChain, err := iam.NewJWTKeyChainFromFiles(jwtPrivateKeyFilenames, "")
	if err != nil {
		return nil, errors.Wrap("JWT key chain loading", err)
	}

	// NOTE: We should store these data into a database instead of CSV file.
	clientDataCSVFilename := filepath.Join(secretFilesDir, "clients.csv")
	applicationDataProvider, err := newClientStaticDataProviderFromCSVFileByName(
		clientDataCSVFilename, 1)
	if err != nil {
		return nil, errors.Wrap("client data loading", err)
	}

	log.Info().Msg("Initializing media service...")
	log.Info().Msgf("Registered media object storage service integrations: %v",
		mediastore.ModuleNames())
	mediaStore, err := mediastore.New(coreCfg.Media, appApp)
	if err != nil {
		return nil, errors.Wrap("file service initialization", err)
	}

	log.Info().Msg("Initializing email-address verification services...")
	eaVerifier := eav10n.NewVerifier(realmInfo, iamDB, coreCfg.EAV)

	log.Info().Msg("Initializing phone-number verification service...")
	log.Info().Msgf("Registered SMS delivery service integrations: %v", pnv10n.ModuleNames())
	pnVerifier := pnv10n.NewVerifier(realmName, iamDB, coreCfg.PNV)

	registeredUserIDCache, err := lru.NewARC(65535)
	if err != nil {
		panic(err)
	}
	deletedUserAccountIDCache, err := lru.NewARC(65535)
	if err != nil {
		panic(err)
	}

	inst := &Core{
		realmInfo:                     realmInfo,
		db:                            iamDB,
		registeredUserInstanceIDCache: registeredUserIDCache,
		deletedUserInstanceIDCache:    deletedUserAccountIDCache,
		applicationDataProvider:       applicationDataProvider,
		mediaStore:                    mediaStore,
		eaVerifier:                    eaVerifier,
		pnVerifier:                    pnVerifier,
	}

	clientBase, err := iam.NewServiceClient(nil, jwtKeyChain, inst)
	if err != nil {
		panic(err)
	}

	inst.ServiceClient = clientBase

	return inst, nil
}

func (core *Core) isTestPhoneNumber(phoneNumber iam.PhoneNumber) bool {
	return phoneNumber.CountryCode() == 1 &&
		phoneNumber.NationalNumber() > 5550000 &&
		phoneNumber.NationalNumber() <= 5559999
}

func (core *Core) isTestEmailAddress(emailAddress iam.EmailAddress) bool {
	return false
}

func connectPostgres(dbURL string) (*sqlx.DB, error) {
	var db *sqlx.DB
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return nil, err
	}

	var maxIdleConns, maxOpenConns int64
	queryPart := parsedURL.Query()
	if maxIdleConnsStr := queryPart.Get("max_idle_conns"); maxIdleConnsStr != "" {
		queryPart.Del("max_idle_conns")
		maxIdleConns, err = strconv.ParseInt(maxIdleConnsStr, 10, 32)
		if err != nil {
			return nil, errors.Wrap("unable to parse max_idle_conns query parameter", err)
		}
	}
	if maxOpenConnsStr := queryPart.Get("max_open_conns"); maxOpenConnsStr != "" {
		queryPart.Del("max_open_conns")
		maxOpenConns, err = strconv.ParseInt(maxOpenConnsStr, 10, 32)
		if err != nil {
			return nil, errors.Wrap("unable to parse max_open_conns query parameter", err)
		}
	}
	if maxIdleConns == 0 {
		maxIdleConns = 2
	}
	if maxOpenConns == 0 {
		maxOpenConns = 8
	}

	parsedURL.RawQuery = queryPart.Encode()
	dbURL = parsedURL.String()
	for {
		db, err = sqlx.Connect("postgres", dbURL)
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "connect: connection refused") {
			return nil, err
		}
		const retryDuration = 5 * time.Second
		time.Sleep(retryDuration)
	}
	if db != nil {
		db.SetMaxIdleConns(int(maxIdleConns))
		db.SetMaxOpenConns(int(maxOpenConns))
	}
	return db, nil
}

type CoreConfig struct {
	DBURL string            `env:"DB_URL,required"`
	Media mediastore.Config `env:"MEDIA"`
	EAV   eav10n.Config     `env:"EAV"`
	PNV   pnv10n.Config     `env:"PNV"`
}

// CoreConfigSkeleton returns an instance of CoreConfig which has been
// configured to load config based on the internal system configuration.
// One kind of usages for a skeleton is to generate a template or documentations.
func CoreConfigSkeleton() CoreConfig {
	return CoreConfig{
		Media: mediastore.ConfigSkeleton(),
		PNV:   pnv10n.ConfigSkeleton(),
	}
}

func CoreConfigSkeletonPtr() *CoreConfig {
	cfg := CoreConfigSkeleton()
	return &cfg
}
