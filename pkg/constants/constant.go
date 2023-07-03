

package constants

const (

	OPERATE_SUCCESS = 0
	UNKNOWN_ERROR = 9999

	LabelNodeResourceType    = "resource-type"
	LabelSshJobKindName      = "damon008.ai/ssh-svc"

	NoteTableName           = "note"
	UserTableName           = "user"
	SecretKey               = "secret key"
	IdentityKey             = "id"
	Total                   = "total"
	Notes                   = "notes"
	NoteID                  = "note_id"
	ApiServiceName          = "customer-service"
	ConsumerServiceName          = "consumer-service"
	OrderServiceName         = "order-service"
	UserServiceName         = "demouser"
	AnalyseServiceName         = "analyse-service"
	ProducerServiceName         = "producer-service"
	EnhanceServiceName         = "enhance-controller"
	HelloServiceName         = "hello-server"
	ElenaServiceName         = "demoElena"
	FileServiceName         = "demofile"
	TestServiceName         = "demotest"
	//MySQLDefaultDSN         = "root:qazwsxedc@tcp(192.168.6.56:3306)/crm1?charset=utf8&parseTime=True&loc=Local"
	MySQLDefaultDSN         = "root:qazwsxedc@tcp(121.37.173.206:3306)/crm?charset=utf8&parseTime=True&loc=Local"
	EtcdAddress             = "127.0.0.1:2379"
	CPURateLimit    float64 = 80.0
	DefaultLimit            = 10


	LicenceFileName = "/damon008/licence/LICENCE"
	PubKeyFileName="/damon008/secret/id_rsa.pub"
	PriKeyFileName="/damon008/secret/id_rsa"


	TIME_LAYOUT = "2006-01-02 15:04:05"


	ApiConfigPath     = "./server/cmd/handler/config.yaml"
	AuthConfigPath    = "./server/cmd/auth/config.yaml"
	BlobConfigPath    = "./server/cmd/blob/config.yaml"
	CarConfigPath     = "./server/cmd/car/config.yaml"
	ProfileConfigPath = "./server/cmd/profile/config.yaml"
	TripConfigPath    = "./server/cmd/trip/config.yaml"

	ApiGroup    = "API_GROUP"
	AuthGroup   = "AUTH_GROUP"
	BlobGroup   = "BLOB_GROUP"
	CarGroup    = "CAR_GROUP"
	RentalGroup = "RENTAL_GROUP"

	NacosLogDir   = "tmp/nacos/log"
	NacosCacheDir = "tmp/nacos/cache"
	NacosLogLevel = "debug"

	HlogFilePath = "./tmp/hlog/logs/"
	KlogFilePath = "./tmp/klog/logs/"

	MySqlDSN    = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	MongoURI    = "mongodb://%s:%s@%s:%d"
	RabbitMqURI = "amqp://%s:%s@%s:%d/"

	IPFlagName  = "ip"
	IPFlagValue = "0.0.0.0"
	IPFlagUsage = "address"

	PortFlagName  = "port"
	PortFlagUsage = "port"

	TCP = "tcp"

	FreePortAddress = "localhost:0"

	DefaultLicNumber = "100000000001"
	DefaultName      = "FreeCar"
	DefaultGender    = 1
	DefaultBirth     = 631152000000
)

var TerminalMode = []string{"sh", "-c"}

// 证书状态
var LicenceState error
// 登录状态
var AuthState = true
