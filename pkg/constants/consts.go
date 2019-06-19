package constants

const (
	AppLabel                     = "mongo-operator"
	MongoName                    = "mongo"
	MCRoleName                   = "mongo"
	GraceTime              		 = 15
	MongoStorageVolumeName       = "mongo-persistent-storage"
)

const (
	ExporterPort                 = 27017
	ExporterPortName             = "db"
	ExporterContainerName        = ""
	ExporterDefaultRequestCPU    = "1"
	ExporterDefaultLimitCPU      = "4"
	ExporterDefaultRequestMemory = "2Gi"
	ExporterDefaultLimitMemory   = "4Gi"
)
