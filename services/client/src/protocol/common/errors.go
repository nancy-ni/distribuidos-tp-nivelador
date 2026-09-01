package common

const (
	InvalidAgencyIdError       = "AgencyID invalido"
	OpenInputFileError         = "Error abriendo el archivo input"
	OpenOutputFileError        = "Error abriendo el archivo output"
	InvalidBatchSizeError      = "Batch Size invalido"
	ScanFileError              = "Error escaneando linea del archivo"
	AssembleBetError           = "Error procesando apuesta: %s"
	PacketTooShortError        = "Longitud del paquete es corta"
	BetTooShort                = "Longitud de la apuesta es corta"
	AskWinnersTooShort         = "Longitud de la consulta de ganadores es corta"
	FinishTooShort             = "Longitud del mensaje final es corta"
	BatchTooShort              = "Longitud del bache es corta"
	UnexpectedMessage          = "Mensaje de codigo desconocido"
	DeserializeUint32Error     = "Error al deserializar Uint32"
	DeserializeUint16Error     = "Error al deserializar Uint16"
	DeserializePacketError     = "Error al deserializar el paquete"
	DeserializeBetError        = "Error al deserializar la apuesta"
	DeserializeAskWinnersError = "Error al deserializar la consulta de ganadores"
	DeserializeFinishError     = "Error al deserializar el mensaje final"
	DeserializeBatchError      = "Error al deserializar un bache"
)
